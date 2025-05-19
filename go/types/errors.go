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

package types

import "errors"

var (
	errAlreadyLaunched   = errors.New("already launched")
	errAlreadyTerminated = errors.New("already terminated")
	errNotLaunched       = errors.New("not launched")
	errNotConnected      = errors.New("not connected")
	errContextDone       = errors.New("context done")
	errInvalidProtocol   = errors.New("invalid protocol")
	errNotAthorized      = errors.New("not authorized")
	errInvalidConfig     = errors.New("invalid config")
)

// ErrAlreadyLaunched defines the error when client already launched.
func ErrAlreadyLaunched() error {
	return errAlreadyLaunched
}

// ErrAlreadyTerminated defines the error when client already terminated.
func ErrAlreadyTerminated() error {
	return errAlreadyTerminated
}

// ErrNotLaunched defines the error when client not launched.
func ErrNotLaunched() error {
	return errNotLaunched
}

// NotConnected defines the error when client not connected.
func NotConnected() error {
	return errNotConnected
}

// ErrContextDone defines the error when context done.
func ErrContextDone() error {
	return errContextDone
}

// ErrInvalidProtocol defines the error when protocol invalid.
func ErrInvalidProtocol() error {
	return errInvalidProtocol
}

// ErrNotAthorized defines the error when not authorized.
func ErrNotAthorized() error {
	return errNotAthorized
}

// ErrInvalidConfig defines the error when config invalid.
func ErrInvalidConfig() error {
	return errInvalidConfig
}
