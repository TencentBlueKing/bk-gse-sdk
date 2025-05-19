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
	"io"
	"net"
)

// Buffer provides buffer management for binary data.
type Buffer struct {
	conn                 net.Conn
	buf                  []byte
	pos, limit, capacity uint32
}

// NewBuffer creates a new Buffer instance with capacity base on tcp connection.
func NewBuffer(conn net.Conn, capacity uint32) *Buffer {
	return &Buffer{
		conn:     conn,
		buf:      make([]byte, capacity),
		capacity: capacity,
	}
}

// Read reads target num bytes data from connection to the buf.
func (b *Buffer) Read(num uint32) error {
	if b.limit+num > b.capacity {
		return fmt.Errorf("override buffer capacity %d", b.limit+num)
	}

	// read message data from tcp connection.
	if _, err := io.ReadFull(b.conn, b.buf[b.limit:b.limit+num]); err != nil {
		return err
	}

	// count read num.
	b.limit += num

	return nil
}

// DecodeUint8 decodes from buffer and returns an uint8 num.
func (b *Buffer) DecodeUint8() (uint8, error) {
	if b.pos+1 > b.limit {
		// not enough, read more.
		if err := b.Read(lenUint8 - (b.limit - b.pos)); err != nil {
			return 0, err
		}
	}

	// decode uint8 data.
	x := b.buf[b.pos : b.pos+1]
	b.pos++

	return x[0], nil
}

// DecodeUint16 decodes from buffer and returns an uint16 num.
func (b *Buffer) DecodeUint16() (uint16, error) {
	if b.pos+2 > b.limit {
		// not enough, read more.
		if err := b.Read(lenUint16 - (b.limit - b.pos)); err != nil {
			return 0, err
		}
	}

	// decode uint16 data.
	x := binary.BigEndian.Uint16(b.buf[b.pos : b.pos+2])
	b.pos += lenUint16

	return x, nil
}

// DecodeUint32 decodes from buffer and returns an uint32 num.
func (b *Buffer) DecodeUint32() (uint32, error) {
	if b.pos+4 > b.limit {
		// not enough, read more.
		if err := b.Read(lenUint32 - (b.limit - b.pos)); err != nil {
			return 0, err
		}
	}

	// decode uint32 data.
	x := binary.BigEndian.Uint32(b.buf[b.pos : b.pos+4])
	b.pos += lenUint32

	return x, nil
}

// DecodeUint64 decodes from buffer and returns an uint64 num.
func (b *Buffer) DecodeUint64() (uint64, error) {
	if b.pos+8 > b.limit {
		// not enough, read more.
		if err := b.Read(lenUint64 - (b.limit - b.pos)); err != nil {
			return 0, err
		}
	}

	// decode uint64 data.
	x := binary.BigEndian.Uint64(b.buf[b.pos : b.pos+8])
	b.pos += lenUint64

	return x, nil
}

// DecodeBytes decodes from buffer and returns raw bytes.
func (b *Buffer) DecodeBytes(length uint32) ([]byte, error) {
	if b.pos+length > b.limit {
		// not enough, read more.
		if err := b.Read(length - (b.limit - b.pos)); err != nil {
			return nil, err
		}
	}

	x := b.buf[b.pos : b.pos+length]
	b.pos += length

	return x, nil
}
