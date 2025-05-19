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

// Package types provides the types and constants used throughout the sdk.
package types

// AgentStatus describes the agent status.
type AgentStatus int

const (
	// AgentStatusUnknown means the status is unknown.
	AgentStatusUnknown AgentStatus = -1

	// AgentStatusInit means the agent is initializing.
	AgentStatusInit AgentStatus = 0

	// AgentStatusStarting means the agent is starting.
	AgentStatusStarting AgentStatus = 1

	// AgentStatusRunning means the agent is running.
	AgentStatusRunning AgentStatus = 2

	// AgentStatusDamage means the agent is damage.
	AgentStatusDamage AgentStatus = 3

	// AgentStatusBusy means the agent is busy.
	AgentStatusBusy AgentStatus = 4

	// AgentStatusUpgrading means the agent is upgrading.
	AgentStatusUpgrading AgentStatus = 5

	// AgentStatusStopping means the agent is stopping.
	AgentStatusStopping AgentStatus = 6

	// AgentStatusUninit means the agent is uninitialized.
	AgentStatusUninit AgentStatus = 7
)

// AgentInfo describes the agent info.
type AgentInfo struct {
	AgentID    string
	Version    string
	CloudID    int
	RunMode    int
	StatusCode AgentStatus
	Status     string
}

// IsRunning returns true if the agent is running.
func (info *AgentInfo) IsRunning() bool {
	return info.StatusCode == AgentStatusRunning
}

// IsProxy returns true if the agent is proxy.
func (info *AgentInfo) IsProxy() bool {
	return info.RunMode == 0
}

// DispatchAgentResult describes the dispatch result of single agent.
type DispatchAgentResult struct {
	AgentID string
	Code    int
	Message string
}
