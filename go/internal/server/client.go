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

// Package server provides gse server API processing.
package server

import "github.com/TencentBlueKing/bk-gse-sdk/go/internal"

// Client provides all handling methods in gse server.
type Client interface {
	// Cluster provides cluster handling methods.
	Cluster() Cluster
}

// New creates a new client.
func New(conf Config) Client {
	return &client{
		conf: conf,
	}
}

type client struct {
	conf Config
}

// Cluster provides cluster handling methods.
func (c *client) Cluster() Cluster {
	return c
}

// Post returns a new request with method POST.
func (c *client) post() *Request {
	return &Request{
		client:   c.conf.Client,
		method:   "POST",
		headers:  internal.CopyHeaders(c.conf.BaseHeader),
		basePath: c.conf.BaseURL,
	}
}
