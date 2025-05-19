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

package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Request handle the http request things.
type Request struct {
	client *http.Client

	method  string
	headers http.Header
	body    []byte

	basePath string
	subPath  string

	processErr error
}

// SubResourcef set sub reousrce path and args.
func (r *Request) SubResourcef(subPath string, args ...interface{}) *Request {
	if len(args) == 0 {
		r.subPath = subPath

		return r
	}

	r.subPath = fmt.Sprint(subPath, args)

	return r
}

// Body set request body.
func (r *Request) Body(body any) *Request {
	if r.processErr != nil {
		return r
	}

	buf, err := json.Marshal(body)
	if err != nil {
		r.processFail(err)

		return r
	}

	r.body = buf

	return r
}

// Headers set request header.
func (r *Request) Headers(header http.Header) *Request {
	if r.headers == nil {
		r.headers = header

		return r
	}

	for key, values := range header {
		for _, v := range values {
			r.headers.Add(key, v)
		}
	}

	return r
}

// Do send request.
func (r *Request) Do() *Result {
	if r.processErr != nil {
		return &Result{processErr: r.processErr}
	}

	req, err := http.NewRequest(r.method, r.getURL(), bytes.NewReader(r.body))
	if err != nil {
		r.processFail(err)

		return &Result{processErr: r.processErr}
	}

	req.Header = r.headers

	resp, err := r.client.Do(req)
	if err != nil {
		r.processFail(err)

		return &Result{processErr: r.processErr}
	}

	if resp.StatusCode != http.StatusOK {
		return &Result{processErr: fmt.Errorf("request to %s failed, status code: %d", r.getURL(), resp.StatusCode)}
	}

	return &Result{resp: resp}
}

func (r *Request) getURL() string {
	return strings.TrimRight(r.basePath, "/") + "/" + strings.TrimLeft(r.subPath, "/")
}

func (r *Request) processFail(err error) {
	if err == nil {
		return
	}

	if r.processErr == nil {
		r.processErr = err

		return
	}

	r.processErr = errors.Join(r.processErr, err)
}

// Result represents the result of a request.
type Result struct {
	processErr error

	resp *http.Response
}

// Into parses the response body into v.
func (r *Result) Into(responseBody any) error {
	if r.processErr != nil {
		return r.processErr
	}

	if r.resp == nil {
		return errors.New("no response")
	}

	data, err := io.ReadAll(r.resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, responseBody)
}
