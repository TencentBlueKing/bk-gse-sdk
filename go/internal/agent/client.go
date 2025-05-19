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

// Package agent provides communication management with gse agent.
package agent

import (
	"context"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"git.woa.com/bk-gse/bk-gse-sdk/types"
)

// Client provides communication management with gse agent.
type Client interface {
	// Launch starts connecting to an agent and holding, wait until it's connected or the context is done.
	Launch(ctx context.Context) error

	// Terminate terminates the connection holding from an agent.
	Terminate(ctx context.Context) error

	// IsConnected returns whether it's connected to an agent.
	IsConnected() bool

	// SendMessage sends a message respond to server though agent.
	SendMessage(ctx context.Context, header Header, content []byte) error
}

// New creates a new client.
func New(conf Config) Client {
	return &client{
		conf: conf,
		done: make(chan struct{}),
	}
}

type client struct {
	conf Config

	launched  atomic.Bool
	connected atomic.Bool

	conn  net.Conn
	mutex sync.Mutex

	done chan struct{}
}

// Launch starts connecting to an agent and holding, wait until it's connected or the context is done.
func (c *client) Launch(ctx context.Context) error {
	if c.launched.Load() {
		return types.ErrAlreadyLaunched()
	}

	c.launched.Store(true)

	notifyConnectedOnce := make(chan struct{})
	defer close(notifyConnectedOnce)

	go c.holdConnection(notifyConnectedOnce)

	for {
		select {
		case <-ctx.Done():
			c.launched.Store(false)
			c.done <- struct{}{}

			return types.ErrContextDone()

		case <-notifyConnectedOnce:
			return nil
		}
	}
}

// Terminate terminates the connection holding from an agent.
func (c *client) Terminate(_ context.Context) error {
	if !c.launched.Load() {
		return types.ErrNotLaunched()
	}

	c.launched.Store(false)
	c.done <- struct{}{}

	return nil
}

// IsConnected returns whether it's connected to an agent.
func (c *client) IsConnected() bool {
	return c.connected.Load()
}

// SendMessage sends a message respond to server though agent.
func (c *client) SendMessage(_ context.Context, header Header, content []byte) error {
	if !c.launched.Load() {
		return types.ErrNotLaunched()
	}

	if !c.connected.Load() {
		return types.NotConnected()
	}

	c.mutex.Lock()
	conn := c.conn
	c.mutex.Unlock()

	if conn == nil {
		return types.NotConnected()
	}

	header.ProtoVersion = protoVersion
	header.Length = uint32(len(content)) + header.HeaderLength()
	headerBuf, err := header.EncodeBuffer()
	if err != nil {
		return err
	}

	buffer := make([]byte, len(headerBuf)+len(content))
	copy(buffer, headerBuf)
	copy(buffer[len(headerBuf):], content)

	if _, err := conn.Write(buffer); err != nil {
		return err
	}

	return nil
}

func (c *client) holdConnection(notifyConnectedOnce chan<- struct{}) {
	for {
		select {
		case <-c.done:
			return

		default:
			conn, err := c.Dial()
			if err != nil {
				c.conf.Logger.Warn("connect to socket failed: %v. retrying in %s", err, c.conf.ReconnectInterval.String())

				// interval before retry.
				time.Sleep(c.conf.ReconnectInterval)

				continue
			}

			c.conf.Logger.Info("connected to socket: %s", conn.RemoteAddr().String())

			c.connectionConnect(conn)

			// notify connected once
			if notifyConnectedOnce != nil {
				notifyConnectedOnce <- struct{}{}
				notifyConnectedOnce = nil
			}

			// brings up receive handler.
			receiveErr := make(chan error)
			go func() { receiveErr <- c.handleReceive(conn) }()

			select {
			case <-c.done:
				c.conf.Logger.Info("forcely disconnected from socket: %s", conn.RemoteAddr().String())
				c.connectionDisconnect()

				return

			case <-receiveErr:
				c.conf.Logger.Warn("lost connection from socket: %s", conn.RemoteAddr().String())
				c.connectionDisconnect()
			}
		}
	}
}

func (c *client) connectionConnect(conn net.Conn) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.conn != nil {
		_ = c.conn.Close()
	}

	c.conn = conn
	c.connected.Store(true)
}

func (c *client) connectionDisconnect() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.conn != nil {
		_ = c.conn.Close()
	}

	c.conn = nil
	c.connected.Store(false)
}

func (c *client) handleReceive(conn net.Conn) error {
	var err error

	for {
		buffer := NewBuffer(conn, c.conf.MaxMessageSizeBytes)

		var header Header
		if err := header.ReadBuffer(buffer); err != nil {
			return err
		}

		if header.Length < header.HeaderLength() {
			return types.ErrInvalidProtocol()
		}

		var raw []byte
		if raw, err = buffer.DecodeBytes(header.Length - header.HeaderLength()); err != nil {
			return err
		}

		content := make([]byte, len(raw))
		copy(content, raw)

		c.conf.RecvCallback(header, content)
	}
}
