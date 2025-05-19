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

package server

// ClusterDispatchMessageReq defines the http request body of cluster dispatch message.
type ClusterDispatchMessageReq struct {
	SlotID      int      `json:"slot_id"`
	Token       string   `json:"token"`
	MessageID   string   `json:"message_id"`
	AgentIDList []string `json:"agent_id_list"`
	Content     string   `json:"content"`
}

// ClusterDispatchMessageResp defines the http response body of cluster dispatch message.
type ClusterDispatchMessageResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Results []*ClusterAgentResult `json:"results"`
	} `json:"data"`
}

// ClusterDispatchMultiMessageReq defines the http request body of cluster dispatch multi message.
type ClusterDispatchMultiMessageReq struct {
	SlotID           int                    `json:"slot_id"`
	Token            string                 `json:"token"`
	MessageID        string                 `json:"message_id"`
	AgentMessageList []*ClusterAgentMessage `json:"agent_message_list"`
	FrontContent     string                 `json:"front_content"`
	BackContent      string                 `json:"back_content"`
}

// ClusterDispatchMultiMessageResp defines the http response body of cluster dispatch multi message.
type ClusterDispatchMultiMessageResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Results []*ClusterAgentResult `json:"results"`
	} `json:"data"`
}

// ClusterAgentResult defines the result of each agent from cluster.
type ClusterAgentResult struct {
	AgentID string `json:"bk_agent_id"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// ClusterAgentMessage defines the message of each agent from cluster.
type ClusterAgentMessage struct {
	AgentID string `json:"bk_agent_id"`
	Content string `json:"content"`
}

// ClusterRespondMessage describes the body of respond message from agent.
type ClusterRespondMessage struct {
	MessageID string `json:"message_id"`
	AgentID   string `json:"bk_agent_id"`
	Content   string `json:"content"`
}
