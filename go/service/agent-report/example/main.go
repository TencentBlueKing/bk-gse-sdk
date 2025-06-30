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

// Package main show how to use agent report.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	agentreport "github.com/TencentBlueKing/bk-gse-sdk/go/service/agent-report"
	"github.com/TencentBlueKing/bk-gse-sdk/go/types"
)

const (
	configFilePath = "agent_report_example.json"
)

// Config defines the config.
type Config struct {
	// DomainSocketPath is the domain socket path.
	DomainSocketPath string `json:"domain_socket_path"`

	// DataID is the data-id(channel-id) to report.
	DataID uint32 `json:"data_id"`
}

// global config.
// nolint:gochecknoglobals
var config Config

func run() {
	// get a new client with options.
	client, err := agentreport.New(
		agentreport.WithDomainSocketPath(config.DomainSocketPath),
		agentreport.WithLogger(types.NewDefaultLogger(1)),
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
	fmt.Println("agent cloud-id: ", agentInfo.CloudID)

	// after launch successfully, you can report data to agent.
	// the message will be report to the data-id(channel-id) which you set.
	for i := 0; ; time.Sleep(time.Second) {
		if err = client.ReportData(
			context.Background(),
			config.DataID,
			[]byte(fmt.Sprintf("[this is the data reported from client %d]", i)),
		); err != nil {
			panic(err)
		}

		i++
		fmt.Printf("[%d] reported data to server\n", i)
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
