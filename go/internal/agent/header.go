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

// IHeader describes the header interface.
type IHeader interface {
	// NewHeader creates a new header with default params.
	NewHeader() IHeader

	// HeaderLength returns the length of header.
	HeaderLength() uint32

	// TotalLength returns the length of whole protocol, header + body.
	TotalLength() uint32

	// ReadBuffer reads the header from buffer.
	ReadBuffer(buf *Buffer) error

	// EncodeBuffer encodes the header to buffer.
	EncodeBuffer() ([]byte, error)
}
