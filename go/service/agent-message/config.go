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
	"errors"
	"time"

	"github.com/TencentBlueKing/bk-gse-sdk/go/types"
)

// NewDefaultConfig creates a default configuration for agent-message service.
func NewDefaultConfig() *Config {
	return &Config{
		DomainSocketPath:    "",
		LocalSocketPort:     0,
		PluginName:          "",
		ReconnectInterval:   defaultReconnectInterval,
		KeepaliveInterval:   defaultKeepaliveInterval,
		MaxMessageSizeBytes: defaultMaxMessageSizeBytes,
		RecvCallback:        func(string, []byte) {},
		Logger:              types.NewDefaultLogger(defaultLoggerLevel),
	}
}

const (
	defaultReconnectInterval   = 1 * time.Second
	defaultKeepaliveInterval   = 3 * time.Second
	defaultMaxMessageSizeBytes = 1024 * 1024 * 10
	defaultLoggerLevel         = 1 // INFO
)

// Config defines the configuration for agent-message service.
type Config struct {
	// DomainSocketPath describes the agent message domain socket path on unix machine.
	DomainSocketPath string

	// LocalSocketPort describes the agent message socket port on windows machine.
	LocalSocketPort uint

	// PluginName describes the plugin's name who is using this SDK.
	// plugin name will be used to do authority and identify which server to send message to.
	PluginName string

	// PluginVersion describes the plugin's version who is using this SDK.
	PluginVersion string

	// ReconnectInterval describes the reconnect interval when connection lost.
	ReconnectInterval time.Duration

	// KeepaliveInterval describes the keepalive interval when connection is alive.
	KeepaliveInterval time.Duration

	// MaxMessageSizeBytes describes the max message size in bytes.
	MaxMessageSizeBytes uint32

	// RecvCallback describes the callback function for agent message service to call when receive a message.
	RecvCallback Callback

	// Logger describes the logger for this service.
	// default logger will
	Logger types.Logger
}

// Validate validates the configuration.
func (c Config) Validate() error {
	if c.PluginName == "" {
		return errors.Join(types.ErrInvalidConfig(), errors.New("plugin name is empty"))
	}

	if c.PluginVersion == "" {
		return errors.Join(types.ErrInvalidConfig(), errors.New("plugin version is empty"))
	}

	if c.ReconnectInterval == 0 {
		return errors.Join(types.ErrInvalidConfig(), errors.New("reconnect interval is 0"))
	}

	if c.KeepaliveInterval == 0 {
		return errors.Join(types.ErrInvalidConfig(), errors.New("keepalive interval is 0"))
	}

	if c.RecvCallback == nil {
		return errors.Join(types.ErrInvalidConfig(), errors.New("recv callback function is empty"))
	}

	if c.Logger == nil {
		return errors.Join(types.ErrInvalidConfig(), errors.New("logger is empty"))
	}

	return nil
}
