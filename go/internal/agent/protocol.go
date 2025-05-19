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

package agent

import (
	"encoding/binary"
	"fmt"
)

// Header describes the common header of protocol.
type Header struct {
	Magic        uint32
	ProtoType    uint16
	ProtoVersion uint16
	Sequence     uint64
	Length       uint32
	Reserved0    uint32
	Reserved1    uint32
}

const (
	// ProtoTypeKeepaliveReq defines the proto type of keepalive request.
	ProtoTypeKeepaliveReq = 0x9006

	// ProtoTypeKeepaliveResp defines the proto type of keepalive response.
	ProtoTypeKeepaliveResp = 0x9007

	// ProtoTypeDispatchMessage defines the proto type of dispatch message.
	ProtoTypeDispatchMessage = 0x9008

	// ProtoTypeRespondMessage defines the proto type of respond message.
	ProtoTypeRespondMessage = 0x900a
)

const (
	// magic number of header.
	magicNumber = 0xdeadbeef

	// version of protocol.
	protoVersion = 0x6

	lenUint8  = 1
	lenUint16 = 2
	lenUint32 = 4
	lenUint64 = 8
)

// HeaderLength returns the length of header.
func (h *Header) HeaderLength() uint32 {
	return lenUint32 +
		lenUint16 +
		lenUint16 +
		lenUint64 +
		lenUint32 +
		lenUint32 +
		lenUint32
}

// ReadBuffer reads the header from buffer.
func (h *Header) ReadBuffer(buf *Buffer) error {
	var err error

	if h.Magic, err = buf.DecodeUint32(); err != nil {
		return err
	}

	if h.Magic != magicNumber {
		return fmt.Errorf("got invalid magic number in header: 0x%x", h.Magic)
	}

	if h.ProtoType, err = buf.DecodeUint16(); err != nil {
		return err
	}

	if h.ProtoVersion, err = buf.DecodeUint16(); err != nil {
		return err
	}

	if h.Sequence, err = buf.DecodeUint64(); err != nil {
		return err
	}

	if h.Length, err = buf.DecodeUint32(); err != nil {
		return err
	}

	if h.Reserved0, err = buf.DecodeUint32(); err != nil {
		return err
	}

	if h.Reserved1, err = buf.DecodeUint32(); err != nil {
		return err
	}

	return nil
}

// EncodeBuffer encodes the header to buffer.
func (h *Header) EncodeBuffer() ([]byte, error) {
	buffer := make([]byte, h.HeaderLength())

	binary.BigEndian.PutUint32(buffer, magicNumber)
	binary.BigEndian.PutUint16(buffer[4:], h.ProtoType)
	binary.BigEndian.PutUint16(buffer[6:], h.ProtoVersion)
	binary.BigEndian.PutUint64(buffer[8:], h.Sequence)
	binary.BigEndian.PutUint32(buffer[16:], h.Length)
	binary.BigEndian.PutUint32(buffer[20:], h.Reserved0)
	binary.BigEndian.PutUint32(buffer[24:], h.Reserved1)

	return buffer, nil
}

// KeepaliveReq describes the keepalive request to gse agent.
type KeepaliveReq struct {
	PluginName string `json:"plugin_name"`
	Version    string `json:"version"`
	Pid        int    `json:"pid"`
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	Remark     string `json:"remark"`
}

// KeepaliveResp describes the keepalive response from gse agent.
type KeepaliveResp struct {
	AgentID    string `json:"agent_id"`
	Version    string `json:"version"`
	CloudID    int    `json:"cloud_id"`
	RunMode    int    `json:"run_mode"`
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
}

// SendMessage describes the message sent to gse agent.
type SendMessage struct {
	Name         string `json:"name"`
	TransmitType int    `json:"transmit_type"`
	Topic        string `json:"topic"`
	MessageID    string `json:"message_id"`
	SessionID    string `json:"session_id"`
}

// RecvMessage describes the message received from gse agent.
type RecvMessage struct {
	MessageID string `json:"message_id"`
	SessionID string `json:"session_id"`
}
