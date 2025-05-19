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

// Package main show how to use agent message.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/TencentBlueKing/bk-gse-sdk/go/internal"
	agentmessage "github.com/TencentBlueKing/bk-gse-sdk/go/service/agent-message"
	"github.com/TencentBlueKing/bk-gse-sdk/go/types"
)

const (
	configFilePath = "agent_message_example.json"

	pidFilePerm = 0x644
)

// Config defines the config.
type Config struct {
	// PluginName is the plugin name. this should be matched with the name registered in server.
	PluginName string `json:"plugin_name"`

	// PluginVersion is the plugin version.
	PluginVersion string `json:"plugin_version"`

	// DomainSocketPath is the domain socket path.
	DomainSocketPath string `json:"domain_socket_path"`

	// PidFile is the pid file of this process.
	PidFile string `json:"pid_file"`
}

// global config.
// nolint:gochecknoglobals
var config Config

// MessageCallback handle the message received from agent.
func MessageCallback(messageID string, content []byte) {
	fmt.Printf("[%s] received message from server: %s\n", messageID, string(content))
}

// GenerateMessageID generate a message id.
func GenerateMessageID() string {
	return fmt.Sprintf("example-client-message-id: %d", internal.GenerateSequence())
}

func run() {
	// get a new client with options.
	client, err := agentmessage.New(
		agentmessage.WithPluginName(config.PluginName),
		agentmessage.WithPluginVersion(config.PluginVersion),
		agentmessage.WithDomainSocketPath(config.DomainSocketPath),
		agentmessage.WithRecvCallback(MessageCallback),
		agentmessage.WithLogger(types.NewDefaultLogger(1)),
	)
	if err != nil {
		panic(err)
	}

	// launch client, it will try to connect to agent and keep the connection.
	if err = client.Launch(context.Background()); err != nil {
		panic(err)
	}

	// wait for a while to receive keepalive response which provides the agent info.
	time.Sleep(3 * time.Second) // nolint:mnd

	// get agent info.
	agentInfo, err := client.GetAgentInfo()
	if err != nil {
		panic(err)
	}

	fmt.Println("agent-id: ", agentInfo.AgentID)
	fmt.Println("agent version: ", agentInfo.Version)
	fmt.Println("agent cloud-id: ", agentInfo.CloudID)
	fmt.Println("agent is running: ", agentInfo.IsRunning())
	fmt.Println("agent is a proxy: ", agentInfo.IsProxy())

	// after launch successfully, you can send message to agent.
	// the message will be respond to your registered API in server.
	for ; ; time.Sleep(time.Second) {
		messageID := GenerateMessageID()
		if err = client.SendMessage(
			context.Background(),
			messageID,
			[]byte("[this is the message from example client]"),
		); err != nil {
			panic(err)
		}

		fmt.Printf("[%s] sent message to server\n", messageID)
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

	if err = os.WriteFile(config.PidFile, []byte(strconv.Itoa(os.Getpid())), pidFilePerm); err != nil {
		panic(err)
	}

	run()
}
