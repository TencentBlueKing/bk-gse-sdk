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
	"net/http"

	"github.com/TencentBlueKing/bk-gse-sdk/go/internal"
	"github.com/TencentBlueKing/bk-gse-sdk/go/types"
)

// OptionFn defines the function type for setting options.
type OptionFn func(*Config)

// WithBaseHeader sets the base header to add to requests.
func WithBaseHeader(header http.Header) OptionFn {
	return func(c *Config) {
		c.BaseHeader = internal.CopyHeaders(header)
	}
}

// WithBaseURL sets the base URL of the server.
func WithBaseURL(url string) OptionFn {
	return func(c *Config) {
		c.BaseURL = url
	}
}

// WithClient sets the HTTP client to use for requests.
func WithClient(client *http.Client) OptionFn {
	return func(c *Config) {
		c.Client = client
	}
}

// WithClusterAuth sets the cluster auth.
func WithClusterAuth(slotID int, token string) OptionFn {
	return func(c *Config) {
		c.SlotID = slotID
		c.Token = token
	}
}

// WithAPIGwAuth sets the api gw auth.
func WithAPIGwAuth(appCode, appSecret string) OptionFn {
	return func(c *Config) {
		c.AppCode = appCode
		c.AppSecret = appSecret
	}
}

// WithLogger sets the logger to use for requests.
func WithLogger(logger types.Logger) OptionFn {
	return func(c *Config) {
		c.Logger = logger
	}
}
