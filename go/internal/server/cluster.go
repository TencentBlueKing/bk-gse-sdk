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

import (
	"context"
	"net/http"
)

// Cluster provides all cluster methods in gse server.
type Cluster interface {
	// DispatchMessage dispatch message to agents through cluster.
	DispatchMessage(ctx context.Context, request *ClusterDispatchMessageReq, header http.Header) (*ClusterDispatchMessageResp, error) // nolint: lll

	// DispatchMultiMessage dispatch multi message to agents through cluster.
	DispatchMultiMessage(ctx context.Context, request *ClusterDispatchMultiMessageReq, header http.Header) (*ClusterDispatchMultiMessageResp, error) // nolint: lll
}

const (
	clusterPathPrefix = "/api/v2/cluster"

	clusterPathDispatchMessage      = clusterPathPrefix + "/dispatch_message"
	clusterPathDispatchMultiMessage = clusterPathPrefix + "/dispatch_multi_message"
)

// DispatchMessage dispatch message to agents through cluster.
func (c *client) DispatchMessage(_ context.Context, request *ClusterDispatchMessageReq, header http.Header) (
	*ClusterDispatchMessageResp, error) {

	resp := new(ClusterDispatchMessageResp)

	err := c.post().
		SubResourcef(clusterPathDispatchMessage).
		Headers(header).
		Body(request).
		Do().
		Into(resp)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

// DispatchMultiMessage dispatch multi message to agents through cluster.
func (c *client) DispatchMultiMessage(_ context.Context, request *ClusterDispatchMultiMessageReq, header http.Header) (
	*ClusterDispatchMultiMessageResp, error) {

	resp := new(ClusterDispatchMultiMessageResp)

	err := c.post().
		SubResourcef(clusterPathDispatchMultiMessage).
		Headers(header).
		Body(request).
		Do().
		Into(resp)

	if err != nil {
		return nil, err
	}

	return resp, nil
}
