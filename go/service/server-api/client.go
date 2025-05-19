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

// Package serverapi provides gse server methods handling.
package serverapi

import (
	"fmt"
	"net/http"

	"git.woa.com/bk-gse/bk-gse-sdk/internal/server"
)

// Client provides all handling methods in gse server.
type Client interface {
	// Cluster provides cluster handling methods.
	Cluster() Cluster
}

// New creates a new client.
func New(opts ...OptionFn) (Client, error) {
	conf := NewDefaultConfig()

	for _, opt := range opts {
		opt(conf)
	}

	if err := conf.Validate(); err != nil {
		return nil, err
	}

	c := &client{conf: conf}
	c.apiClient = server.New(server.Config{
		BaseHeader: conf.BaseHeader,
		BaseURL:    conf.BaseURL,
		Client:     conf.Client,
		Logger:     conf.Logger,
	})

	return c, nil
}

type client struct {
	conf *Config

	apiClient server.Client
}

// Cluster provides cluster handling methods.
func (c *client) Cluster() Cluster {
	return &clusterClient{client: c}
}

const (
	headerAPIGwAuthKey = "X-Bkapi-Authorization"
)

func (c *client) generateHeaders() http.Header {
	header := http.Header{}

	// api-gateway auth headers.
	if c.conf.AppCode != "" {
		header.Set(headerAPIGwAuthKey,
			fmt.Sprintf("{\"bk_app_code\": \"%s\", \"bk_app_secret\": \"%s\"}", c.conf.AppCode, c.conf.AppSecret))
	}

	return header
}
