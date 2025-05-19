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

// Package main show how to use server api.
package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"git.woa.com/bk-gse/bk-gse-sdk/internal"
	serverapi "git.woa.com/bk-gse/bk-gse-sdk/service/server-api"
	"git.woa.com/bk-gse/bk-gse-sdk/types"
)

const (
	configFilePath = "server_example.json"
)

// Config defines the config.
type Config struct {
	ListenAddr string `json:"listen_addr"`
	SlotID     int    `json:"slot_id"`
	Token      string `json:"token"`

	// AppCode defines the app code for apigw.
	AppCode string `json:"app_code"`

	// AppSecret defines the app secret for apigw.
	AppSecret string `json:"app_secret"`

	// GSEBaseURL defines the gse base url.
	GSEBaseURL string `json:"gse_base_url"`

	// AgentID describes the agent-id of agent where the plugin installed.
	AgentIDList []string `json:"agent_id_list"`
}

// global config.
// nolint:gochecknoglobals
var config Config

// Server brings up a http server and serves the request.
// send message to agent through cluster.
type Server struct {
	client serverapi.Client
}

// HTTPCallbackHandler handle the respond message callback from agent through cluster.
func (s *Server) HTTPCallbackHandler(resp http.ResponseWriter, req *http.Request) {
	defer func() {
		_ = req.Body.Close()
	}()

	data, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	message, err := s.client.Cluster().EncoderDecoder().DecodePluginRespondMessageCallback(data)
	if err != nil {
		panic(err)
	}

	fmt.Printf("[%s] received message from agent(%s): %s\n",
		message.MessageID, message.AgentID, message.Content)

	resp.WriteHeader(http.StatusOK)
}

// ServeHTTPServer serve http server.
func (s *Server) ServeHTTPServer() {
	http.HandleFunc("/callback", s.HTTPCallbackHandler)

	if err := http.ListenAndServe(config.ListenAddr, nil); err != nil { // nolint:gosec
		panic(err)
	}
}

// GenerateMessageID generate a message id.
func GenerateMessageID() string {
	return fmt.Sprintf("example-server-message-id: %d", internal.GenerateSequence())
}

func run() {
	client, err := serverapi.New(
		serverapi.WithBaseURL(config.GSEBaseURL),
		serverapi.WithClient(http.DefaultClient),
		serverapi.WithClusterAuth(config.SlotID, config.Token),
		serverapi.WithAPIGwAuth(config.AppCode, config.AppSecret),
		serverapi.WithLogger(types.NewDefaultLogger(1)),
	)
	if err != nil {
		panic(err)
	}

	svr := &Server{client: client}
	go svr.ServeHTTPServer()

	for ; ; time.Sleep(time.Second) {
		messageID := GenerateMessageID()
		resp, err := client.Cluster().PluginDispatchMessage(
			context.Background(),
			messageID,
			[]byte("[this is the message from example server]"),
			config.AgentIDList...,
		)
		if err != nil {
			panic(err)
		}

		if resp.Code != 0 {
			panic(fmt.Errorf("[%s] send message to agent failed, code: %d, message: %s", messageID, resp.Code, resp.Message))
		}

		for agentID, agentResult := range resp.AgentResults {
			err = errors.Join(err,
				fmt.Errorf("[%s] send message to agent failed, agent-id: %s, agent-code: %d, agent-message: %s",
					messageID, agentID, agentResult.Code, agentResult.Message))
		}
		if err != nil {
			panic(err)
		}

		fmt.Printf("[%s] sent message to agents(%v)\n", messageID, config.AgentIDList)
	}
}

func main() {
	configContent, err := os.ReadFile(configFilePath)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(configContent, &config); err != nil {
		panic(err)
	}

	run()
}
