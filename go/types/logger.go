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

import "log"

// Logger defines the logger.
type Logger interface {
	// Debug logs to DEBUG log.
	Debug(format string, args ...interface{})

	// Info logs to INFO log.
	Info(format string, args ...interface{})

	// Warn logs to WARNING log.
	Warn(format string, args ...interface{})

	// Error logs to ERROR log.
	Error(format string, args ...interface{})
}

// NewDefaultLogger creates a new DefaultLogger.
// level presents 0 ~ 3 for DEBUG ~ ERROR.
func NewDefaultLogger(level int) Logger {
	return &DefaultLogger{
		level: level,
	}
}

// DefaultLogger logs into stdout.
type DefaultLogger struct {
	// log level.
	level int
}

// Debug logs to stdout.
func (l DefaultLogger) Debug(format string, args ...interface{}) {
	if l.level <= 0 { // nolint:mnd
		log.Printf("[DEBUG]"+format, args...)
	}
}

// Info logs to stdout.
func (l DefaultLogger) Info(format string, args ...interface{}) {
	if l.level <= 1 { // nolint:mnd
		log.Printf("[INFO]"+format, args...)
	}
}

// Warn logs to stdout.
func (l DefaultLogger) Warn(format string, args ...interface{}) {
	if l.level <= 2 { // nolint:mnd
		log.Printf("[WARN]"+format, args...)
	}
}

// Error logs to stdout.
func (l DefaultLogger) Error(format string, args ...interface{}) {
	if l.level <= 3 { // nolint:mnd
		log.Printf("[ERROR]"+format, args...)
	}
}

// NewEmptyLogger creates a new logger which does nothing.
func NewEmptyLogger() Logger {
	return &EmptyLogger{}
}

// EmptyLogger does nothing.
type EmptyLogger struct {
}

// Debug does nothing.
func (l EmptyLogger) Debug(string, ...interface{}) {}

// Info does nothing.
func (l EmptyLogger) Info(string, ...interface{}) {}

// Warn does nothing.
func (l EmptyLogger) Warn(string, ...interface{}) {}

// Error does nothing.
func (l EmptyLogger) Error(string, ...interface{}) {}
