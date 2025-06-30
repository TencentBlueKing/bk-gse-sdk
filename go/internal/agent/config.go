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

package agent

import (
	"time"

	"github.com/TencentBlueKing/bk-gse-sdk/go/types"
)

// Config describes the agent domainsocket communication configuration.
type Config struct {
	// DomainSocketPath describes the agent message domain socket path on unix machine.
	DomainSocketPath string

	// LocalSocketPort describes the agent message socket port on windows machine.
	LocalSocketPort uint

	// PluginName describes the plugin's name who is using this SDK.
	// plugin name will be used to do authority and identify which server to send message to.
	PluginName string

	// ReconnectInterval describes the reconnect interval when connection lost.
	ReconnectInterval time.Duration

	// MaxMessageSizeBytes describes the max message size in bytes.
	MaxMessageSizeBytes uint32

	// RecvCallback describes the callback function for agent message service to call when receive a message.
	RecvCallback Callback

	// RecvHeader describes the header for agent message service to call when receive a message.
	RecvHeader IHeader

	// Logger describes the logger for this service.
	// default logger will prints to stdout.
	Logger types.Logger
}

// Callback describes the callback function for agent message service to call when receive a message.
type Callback func(header IHeader, content []byte)
