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
	"errors"
	"net/http"

	"github.com/TencentBlueKing/bk-gse-sdk/go/types"
)

// NewDefaultConfig creates a default configuration for server service.
func NewDefaultConfig() *Config {
	return &Config{
		BaseHeader: make(http.Header),
		BaseURL:    "",
		Client:     &http.Client{},
		Logger:     types.NewDefaultLogger(defaultLoggerLevel),
	}
}

const (
	defaultLoggerLevel = 1 // INFO
)

// Config describes the server configuration.
type Config struct {
	// BaseHeader is the base HTTP header to add to requests.
	BaseHeader http.Header

	// BaseURL is the base URL of the server.
	BaseURL string

	// Client is the HTTP client to use for requests.
	Client *http.Client

	// SlotID, Token provides the cluster dispatching authentication.
	SlotID int
	Token  string

	// AppCode, AppSecret provides the APIGateway authentication.
	AppCode   string
	AppSecret string

	// Logger is the logger to use for requests.
	Logger types.Logger
}

// Validate validates the configuration.
func (c Config) Validate() error {
	if c.Client == nil {
		return errors.Join(types.ErrInvalidConfig(), errors.New("http client is nil"))
	}

	if c.Logger == nil {
		return errors.Join(types.ErrInvalidConfig(), errors.New("logger is empty"))
	}

	return nil
}
