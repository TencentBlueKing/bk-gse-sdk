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

package agentmessage

import (
	"time"

	"github.com/TencentBlueKing/bk-gse-sdk/go/types"
)

// OptionFn defines the function type for setting options.
type OptionFn func(*Config)

// WithDomainSocketPath sets the domain socket path.
func WithDomainSocketPath(path string) OptionFn {
	return func(c *Config) {
		c.DomainSocketPath = path
	}
}

// WithLocalSocketPort sets the socket port on windows.
func WithLocalSocketPort(port uint) OptionFn {
	return func(c *Config) {
		c.LocalSocketPort = port
	}
}

// WithPluginName sets the plugin name.
func WithPluginName(name string) OptionFn {
	return func(c *Config) {
		c.PluginName = name
	}
}

// WithPluginVersion sets the plugin version.
func WithPluginVersion(version string) OptionFn {
	return func(c *Config) {
		c.PluginVersion = version
	}
}

// WithReconnectInterval sets the reconnect interval.
func WithReconnectInterval(interval time.Duration) OptionFn {
	return func(c *Config) {
		c.ReconnectInterval = interval
	}
}

// WithKeepaliveInterval sets the keepalive interval.
func WithKeepaliveInterval(interval time.Duration) OptionFn {
	return func(c *Config) {
		c.KeepaliveInterval = interval
	}
}

// WithMaxMessageSizeBytes sets the max message size in bytes.
func WithMaxMessageSizeBytes(size uint32) OptionFn {
	return func(c *Config) {
		c.MaxMessageSizeBytes = size
	}
}

// WithRecvCallback sets the callback function for receiving message.
func WithRecvCallback(callback Callback) OptionFn {
	return func(c *Config) {
		c.RecvCallback = callback
	}
}

// WithLogger sets the logger.
func WithLogger(logger types.Logger) OptionFn {
	return func(c *Config) {
		c.Logger = logger
	}
}

// DisableLogger disables the logger.
func DisableLogger() OptionFn {
	return func(c *Config) {
		c.Logger = types.NewEmptyLogger()
	}
}
