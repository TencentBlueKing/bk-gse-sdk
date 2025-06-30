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

import "encoding/binary"

const (
	/*
	 * Data Plugin Protocol
	 */

	// ProtoTypeDataPluginSyncConfigReq defines the proto type of data plugin sync config request.
	ProtoTypeDataPluginSyncConfigReq = 0x0A

	// ProtoTypeDataPluginSyncConfigResp defines the proto type of data plugin sync config response.
	ProtoTypeDataPluginSyncConfigResp = ProtoTypeDataPluginSyncConfigReq // same as request message-number.

	// ProtoTypeDataPluginReportReq defines the proto type of data plugin report request.
	ProtoTypeDataPluginReportReq = 0xc01
)

// NewDataUpHeader creates a new header.
func NewDataUpHeader() *DataUpHeader {
	return &DataUpHeader{}
}

// DataUpHeader describes the data up header of protocol.
type DataUpHeader struct {
	ProtoType  uint32
	DataID     uint32
	UTCTime    uint32
	BodyLength uint32
	Reserved0  uint32
	Reserved1  uint32
}

// NewHeader creates a new header.
func (h *DataUpHeader) NewHeader() IHeader {
	return NewDataUpHeader()
}

// HeaderLength returns the length of header.
func (h *DataUpHeader) HeaderLength() uint32 {
	return lenUint32 +
		lenUint32 +
		lenUint32 +
		lenUint32 +
		lenUint32 +
		lenUint32
}

// TotalLength returns the length of total protocol, header + body.
func (h *DataUpHeader) TotalLength() uint32 {
	return h.HeaderLength() + h.BodyLength
}

// ReadBuffer reads the header from buffer.
func (h *DataUpHeader) ReadBuffer(buf *Buffer) error {
	var err error

	if h.ProtoType, err = buf.DecodeUint32(); err != nil {
		return err
	}

	if h.DataID, err = buf.DecodeUint32(); err != nil {
		return err
	}

	if h.UTCTime, err = buf.DecodeUint32(); err != nil {
		return err
	}

	if h.BodyLength, err = buf.DecodeUint32(); err != nil {
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
func (h *DataUpHeader) EncodeBuffer() ([]byte, error) {
	buffer := make([]byte, h.HeaderLength())

	binary.BigEndian.PutUint32(buffer, h.ProtoType)
	binary.BigEndian.PutUint32(buffer[4:], h.DataID)
	binary.BigEndian.PutUint32(buffer[8:], h.UTCTime)
	binary.BigEndian.PutUint32(buffer[12:], h.BodyLength)
	binary.BigEndian.PutUint32(buffer[16:], h.Reserved0)
	binary.BigEndian.PutUint32(buffer[20:], h.Reserved1)

	return buffer, nil
}

// NewDataDownHeader creates a new header.
func NewDataDownHeader() *DataDownHeader {
	return &DataDownHeader{}
}

// DataDownHeader describes the data down header of protocol.
type DataDownHeader struct {
	ProtoType  uint32
	BodyLength uint32
}

// NewHeader creates a new header.
func (h *DataDownHeader) NewHeader() IHeader {
	return NewDataDownHeader()
}

// HeaderLength returns the length of header.
func (h *DataDownHeader) HeaderLength() uint32 {
	return lenUint32 +
		lenUint32
}

// TotalLength returns the length of total protocol, header + body.
func (h *DataDownHeader) TotalLength() uint32 {
	return h.HeaderLength() + h.BodyLength
}

// ReadBuffer reads the header from buffer.
func (h *DataDownHeader) ReadBuffer(buf *Buffer) error {
	var err error

	if h.ProtoType, err = buf.DecodeUint32(); err != nil {
		return err
	}

	if h.BodyLength, err = buf.DecodeUint32(); err != nil {
		return err
	}

	return nil
}

// EncodeBuffer encodes the header to buffer.
func (h *DataDownHeader) EncodeBuffer() ([]byte, error) {
	buffer := make([]byte, h.HeaderLength())

	binary.BigEndian.PutUint32(buffer, h.ProtoType)
	binary.BigEndian.PutUint32(buffer[4:], h.BodyLength)

	return buffer, nil
}

// DataPluginSyncConfigReq describes the data plugin sync config request.
type DataPluginSyncConfigReq struct {
}

// DataPluginSyncConfigResp describes the data plugin sync config response.
type DataPluginSyncConfigResp struct {
	CloudID int    `json:"cloud_id"`
	AgentID string `json:"bk_agent_id"`
}
