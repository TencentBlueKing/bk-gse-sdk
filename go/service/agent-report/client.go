/*
 * TencentBlueKing is pleased to support the open source community by making
 * 蓝鲸智云 - 管控平台(BlueKing - General Service Engine) available.
 * Copyright (C) 2025 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed
 * under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
 * CONDITIONS OF ANY KIND, either express or implied. See the License for the specific
 * language governing permissions and limitations under the License.We undertake not
 * to change the open source license (MIT license) applicable to the current version
 * of the project delivered to anyone in the future.
 */

// Package agentreport provides agent data report handling.
package agentreport

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/TencentBlueKing/bk-gse-sdk/go/internal/agent"
	"github.com/TencentBlueKing/bk-gse-sdk/go/types"
)

// Client provides all handling methods in agent data report.
type Client interface {
	// Launch starts connecting to an agent and holding, wait until it's connected or the context is done.
	Launch(ctx context.Context) error

	// Terminate terminates the connection holding from an agent.
	Terminate(ctx context.Context) error

	// ReportData sends a data report to server though agent.
	ReportData(ctx context.Context, dataID uint32, content []byte) error

	// GetAgentInfo returns agent info.
	GetAgentInfo() (types.AgentSimpleInfo, error)
}

// New creates a new agent-report client.
func New(opts ...OptionFn) (Client, error) {
	conf := NewDefaultConfig()

	for _, opt := range opts {
		opt(conf)
	}

	if err := conf.Validate(); err != nil {
		return nil, err
	}

	c := &client{
		conf: conf,
		done: make(chan struct{}),
	}
	c.client = agent.New(agent.Config{
		DomainSocketPath:    conf.DomainSocketPath,
		LocalSocketPort:     conf.LocalSocketPort,
		ReconnectInterval:   conf.ReconnectInterval,
		MaxMessageSizeBytes: conf.MaxMessageSizeBytes,
		RecvCallback: func(header agent.IHeader, content []byte) {
			c.handleReceive(header, content)
		},
		RecvHeader: agent.NewDataDownHeader(),
		Logger:     conf.Logger,
	})

	return c, nil
}

type client struct {
	conf *Config

	done chan struct{}

	client agent.Client

	// agentInfo describes the agent newest info from keepalive response.
	agentInfo types.AgentSimpleInfo
	mutex     sync.RWMutex
}

// Launch starts connecting to an agent and holding, wait until it's connected or the context is done.
func (c *client) Launch(ctx context.Context) error {
	if err := c.client.Launch(ctx); err != nil {
		return err
	}

	go c.holdKeepalive() // nolint:contextcheck

	return nil
}

// Terminate terminates the connection holding from an agent.
func (c *client) Terminate(ctx context.Context) error {
	if err := c.client.Terminate(ctx); err != nil {
		return err
	}

	c.done <- struct{}{}

	return nil
}

// ReportData sends a data report to server though agent.
func (c *client) ReportData(ctx context.Context, dataID uint32, content []byte) error {
	return c.reportData(ctx, dataID, content)
}

// GetAgentInfo returns agent info.
func (c *client) GetAgentInfo() (types.AgentSimpleInfo, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.agentInfo, nil
}

func (c *client) reportData(ctx context.Context, dataID uint32, content []byte) error {
	header := agent.NewDataUpHeader()
	header.ProtoType = agent.ProtoTypeDataPluginReportReq
	header.DataID = dataID
	header.UTCTime = uint32(time.Now().UTC().UnixMilli())
	header.BodyLength = uint32(len(content))

	return c.client.SendMessage(ctx, header, content)
}

func (c *client) handleReceive(recvHeader agent.IHeader, content []byte) {
	header, ok := recvHeader.(*agent.DataDownHeader)
	if !ok {
		c.conf.Logger.Warn("received unknown header: %v", recvHeader)
		return
	}

	switch header.ProtoType {
	case agent.ProtoTypeDataPluginSyncConfigResp:
		go c.handleKeepaliveResp(header, content)

	default:
		c.conf.Logger.Warn("received unknown message type: 0x%x", header.ProtoType)
	}
}

func (c *client) handleKeepaliveResp(_ *agent.DataDownHeader, content []byte) {
	var resp agent.DataPluginSyncConfigResp
	if err := json.Unmarshal(content, &resp); err != nil {
		c.conf.Logger.Warn("unmarshal keepalive(sync config) response failed: %v", err)
		return
	}

	c.mutex.Lock()
	c.agentInfo = types.AgentSimpleInfo{
		AgentID: resp.AgentID,
		CloudID: resp.CloudID,
	}
	c.mutex.Unlock()

	c.conf.Logger.Debug("received keepalive(sync config) response: %v", resp)
}

func (c *client) holdKeepalive() {
	c.conf.Logger.Info("start sending keepalive(sync config) with interval %s", c.conf.KeepaliveInterval.String())

	for ; ; time.Sleep(c.conf.KeepaliveInterval) {
		select {
		case <-c.done:
			c.conf.Logger.Info("stop sending keepalive(sync config)")
			return

		default:
			header := agent.NewDataUpHeader()
			header.ProtoType = agent.ProtoTypeDataPluginSyncConfigReq
			header.BodyLength = 0

			if err := c.client.SendMessage(context.Background(), header, nil); err != nil {
				c.conf.Logger.Warn("send keepalive(sync config) request failed: %v", err)
				continue
			}

			c.conf.Logger.Debug("send keepalive(sync config) request succeed")
		}
	}
}
