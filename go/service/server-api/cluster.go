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

package serverapi

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/TencentBlueKing/bk-gse-sdk/go/internal/server"
	"github.com/TencentBlueKing/bk-gse-sdk/go/types"
)

// Cluster provides cluster handling methods.
type Cluster interface {
	// DispatchMessage dispatch message to agent through cluster.
	PluginDispatchMessage(ctx context.Context, messageID string, content []byte, agentIDList ...string) (*ClusterPluginDispatchMessageResp, error) // nolint:lll

	// EncoderDecoder provides encoder/decoder of cluster protocol.
	EncoderDecoder() ClusterEncoderDecoder
}

type clusterClient struct {
	*client
}

// EncoderDecoder provides encoder/decoder of cluster protocol.
func (c *clusterClient) EncoderDecoder() ClusterEncoderDecoder {
	return c
}

// PluginDispatchMessage dispatch message to agent through cluster.
func (c *clusterClient) PluginDispatchMessage(ctx context.Context, messageID string, content []byte, agentIDList ...string) (*ClusterPluginDispatchMessageResp, error) { // nolint:lll
	if c.conf.Token == "" {
		return nil, errors.Join(types.ErrInvalidConfig(), errors.New("cluster auth token is empty"))
	}

	resp, err := c.apiClient.Cluster().DispatchMessage(ctx,
		&server.ClusterDispatchMessageReq{
			SlotID:      c.conf.SlotID,
			Token:       c.conf.Token,
			MessageID:   messageID,
			AgentIDList: agentIDList,
			Content:     string(content),
		}, c.generateHeaders())

	if err != nil {
		return nil, err
	}

	result := &ClusterPluginDispatchMessageResp{
		Code:         resp.Code,
		Message:      resp.Message,
		AgentResults: make(map[string]*types.DispatchAgentResult, 0),
	}

	for _, item := range resp.Data.Results {
		result.AgentResults[item.AgentID] = &types.DispatchAgentResult{
			AgentID: item.AgentID,
			Code:    item.Code,
			Message: item.Message,
		}
	}

	return result, nil
}

// ClusterEncoderDecoder describes the encoder/decoder of cluster protocol.
type ClusterEncoderDecoder interface {
	// EncodePluginDispatchMessageRequest generates request body.
	EncodePluginDispatchMessageRequest(messageID string, content []byte, agentIDList ...string) ([]byte, error)

	// DecodePluginDispatchMessageResponse decodes the dispatch message response body.
	DecodePluginDispatchMessageResponse(body []byte) (*ClusterPluginDispatchMessageResp, error)

	// DecodePluginRespondMessageCallback decodes the respond message callback request body.
	DecodePluginRespondMessageCallback(body []byte) (*ClusterPluginRespondMessage, error)
}

// EncodePluginDispatchMessageRequest generates request body.
func (c *clusterClient) EncodePluginDispatchMessageRequest(messageID string, content []byte, agentIDList ...string) ([]byte, error) { // nolint:lll
	if c.conf.Token == "" {
		return nil, errors.Join(types.ErrInvalidConfig(), errors.New("cluster auth token is empty"))
	}

	return json.Marshal(&server.ClusterDispatchMessageReq{
		SlotID:      c.conf.SlotID,
		Token:       c.conf.Token,
		MessageID:   messageID,
		AgentIDList: agentIDList,
		Content:     string(content),
	})
}

// DecodePluginDispatchMessageResponse decodes the dispatch message response body.
func (c *clusterClient) DecodePluginDispatchMessageResponse(body []byte) (*ClusterPluginDispatchMessageResp, error) {
	resp := new(server.ClusterDispatchMessageResp)
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	result := &ClusterPluginDispatchMessageResp{
		Code:         resp.Code,
		Message:      resp.Message,
		AgentResults: make(map[string]*types.DispatchAgentResult, 0),
	}

	for _, item := range resp.Data.Results {
		result.AgentResults[item.AgentID] = &types.DispatchAgentResult{
			AgentID: item.AgentID,
			Code:    item.Code,
			Message: item.Message,
		}
	}

	return result, nil
}

// DecodePluginRespondMessageCallback decodes the respond message callback request body.
func (c *clusterClient) DecodePluginRespondMessageCallback(body []byte) (*ClusterPluginRespondMessage, error) {
	data := new(server.ClusterRespondMessage)
	if err := json.Unmarshal(body, data); err != nil {
		return nil, err
	}

	result := &ClusterPluginRespondMessage{
		MessageID: data.MessageID,
		AgentID:   data.AgentID,
		Content:   data.Content,
	}

	return result, nil
}
