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

// Package agentmessage provides agent message handling.
package agentmessage

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"sync/atomic"
	"time"

	"git.woa.com/bk-gse/bk-gse-sdk/internal"
	"git.woa.com/bk-gse/bk-gse-sdk/internal/agent"
	"git.woa.com/bk-gse/bk-gse-sdk/types"
)

// Client provides all handling methods in agent message.
type Client interface {
	// Launch starts connecting to an agent and holding, wait until it's connected or the context is done.
	Launch(ctx context.Context) error

	// Terminate terminates the connection holding from an agent.
	Terminate(ctx context.Context) error

	// SendMessage sends a message respond to server though agent.
	SendMessage(ctx context.Context, messageID string, content []byte) error

	// GetAgentInfo returns agent info.
	GetAgentInfo() (types.AgentInfo, error)
}

// Callback defines a callback function for client to call when receive a message.
type Callback func(messageID string, content []byte)

// New creates a new agent-message client.
func New(opts ...OptionFn) (Client, error) {
	conf := NewDefaultConfig()

	for _, opt := range opts {
		opt(conf)
	}

	if err := conf.Validate(); err != nil {
		return nil, err
	}

	c := &client{conf: conf}
	c.client = agent.New(agent.Config{
		DomainSocketPath:    conf.DomainSocketPath,
		LocalSocketPort:     conf.LocalSocketPort,
		PluginName:          conf.PluginName,
		ReconnectInterval:   conf.ReconnectInterval,
		MaxMessageSizeBytes: conf.MaxMessageSizeBytes,
		RecvCallback: func(header agent.Header, content []byte) {
			c.handleReceive(header, content)
		},
		Logger: conf.Logger,
	})

	return c, nil
}

type client struct {
	conf *Config

	authorized atomic.Bool

	done chan struct{}

	client agent.Client

	// agentInfo describes the agent newest info from keepalive response.
	agentInfo types.AgentInfo
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

// SendMessage sends a message respond to server though agent.
func (c *client) SendMessage(ctx context.Context, messageID string, content []byte) error {
	return c.sendMessage(ctx, messageID, content)
}

// GetAgentInfo returns agent info.
func (c *client) GetAgentInfo() (types.AgentInfo, error) {
	if !c.authorized.Load() {
		return types.AgentInfo{}, types.ErrNotAthorized()
	}

	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.agentInfo, nil
}

func (c *client) sendMessage(ctx context.Context, messageID string, content []byte) error {
	request := agent.SendMessage{
		Name:         c.conf.PluginName,
		TransmitType: 0,
		Topic:        "",
		MessageID:    messageID,
		SessionID:    "",
	}
	info, err := json.Marshal(request)
	if err != nil {
		c.conf.Logger.Error("marshal send message info failed. message-id: %s, err: %v", messageID, err)
		return err
	}

	header := agent.Header{
		ProtoType: agent.ProtoTypeRespondMessage,
		Sequence:  internal.GenerateSequence(),
		Reserved0: uint32(len(info)),
		Reserved1: uint32(len(content)),
	}

	buffer := make([]byte, len(info)+len(content))
	copy(buffer, info)
	copy(buffer[len(info):], content)

	if err = c.client.SendMessage(ctx, header, buffer); err != nil {
		c.conf.Logger.Error("send message to agent failed. message-id: %s, err: %v", messageID, err)
		return err
	}

	c.conf.Logger.Debug("sent message to agent. message-id: %s, content: %s", messageID, string(content))

	return nil
}

func (c *client) handleReceive(header agent.Header, content []byte) {
	switch header.ProtoType {
	case agent.ProtoTypeKeepaliveResp:
		go c.handleKeepaliveResp(header, content)

	case agent.ProtoTypeDispatchMessage:
		c.handleDispatchMessage(header, content)

	default:
		c.conf.Logger.Warn("received unknown message type: 0x%x", header.ProtoType)
	}
}

func (c *client) handleKeepaliveResp(_ agent.Header, content []byte) {
	var resp agent.KeepaliveResp
	if err := json.Unmarshal(content, &resp); err != nil {
		c.conf.Logger.Warn("unmarshal keepalive response failed: %v", err)
		return
	}

	c.mutex.Lock()
	c.agentInfo = types.AgentInfo{
		AgentID:    resp.AgentID,
		Version:    resp.Version,
		CloudID:    resp.CloudID,
		RunMode:    resp.RunMode,
		StatusCode: types.AgentStatus(resp.StatusCode),
		Status:     resp.Status,
	}
	c.mutex.Unlock()

	c.authorized.Store(true)

	c.conf.Logger.Debug("received keepalive response: %v", resp)
}

func (c *client) handleDispatchMessage(header agent.Header, content []byte) {
	infoLen := header.Reserved0
	dataLen := header.Reserved1

	if infoLen == 0 || infoLen+dataLen != uint32(len(content)) || infoLen+dataLen+header.HeaderLength() != header.Length {
		c.conf.Logger.Error("recv message from agent with invalid length: %v", header)
		return
	}

	var resp agent.RecvMessage
	if err := json.Unmarshal(content[:infoLen], &resp); err != nil {
		c.conf.Logger.Error("unmarshal recv message failed: %v", err)
		return
	}

	c.conf.Logger.Debug("received dispatch message: %v", resp)

	c.conf.RecvCallback(resp.MessageID, content[infoLen:])
}

func (c *client) holdKeepalive() {
	c.conf.Logger.Info("start sending keepalive with interval %s", c.conf.KeepaliveInterval.String())

	for ; ; time.Sleep(c.conf.KeepaliveInterval) {
		select {
		case <-c.done:
			c.conf.Logger.Info("stop sending keepalive")
			return

		default:
			request := agent.KeepaliveReq{
				PluginName: c.conf.PluginName,
				Version:    c.conf.PluginVersion,
				Pid:        os.Getpid(),
				StatusCode: 0,
				Status:     "ok",
				Remark:     "",
			}

			buf, err := json.Marshal(&request)
			if err != nil {
				c.conf.Logger.Warn("marshal keepalive request failed: %v", err)
				continue
			}

			header := agent.Header{
				ProtoType: agent.ProtoTypeKeepaliveReq,
				Sequence:  internal.GenerateSequence(),
			}

			if err = c.client.SendMessage(context.Background(), header, buf); err != nil {
				c.conf.Logger.Warn("send keepalive request failed: %v", err)
				continue
			}

			c.conf.Logger.Debug("send keepalive request succeed")
		}
	}
}
