/*
 * Copyright 2023 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package server

import (
	"net/http"

	"github.com/cloudwego/hertz/pkg/network"
	"github.com/cloudwego/hertz/pkg/protocol"
)

type hertzWriter struct {
	sentHeaders bool
	statusCode  int
	baseRw      *protocol.Response
	headers     http.Header
	buf         []byte
}

func newWriter(response *protocol.Response) *hertzWriter {
	return &hertzWriter{
		sentHeaders: false,
		statusCode:  http.StatusOK,
		baseRw:      response,
		headers:     http.Header{},
		buf:         []byte{},
	}
}

var (
	_ http.ResponseWriter = (*hertzWriter)(nil)
	_ network.ExtWriter   = (*hertzWriter)(nil)
)

func (h *hertzWriter) Header() http.Header {
	return h.headers
}

func (h *hertzWriter) Write(b []byte) (int, error) {
	h.syncHeaders()
	h.baseRw.SetStatusCode(h.statusCode)
	h.buf = append(h.buf, b...)
	return len(b), nil
}

func (h *hertzWriter) syncHeaders() {
	if !h.sentHeaders {
		for hname, hval := range h.headers {
			h.baseRw.Header.Set(hname, hval[0])
		}
	}
}

func (h *hertzWriter) WriteHeader(code int) {
	h.syncHeaders()
	h.statusCode = code
}

func (h *hertzWriter) Flush() error {
	h.syncHeaders()
	return nil
}

func (h *hertzWriter) Finalize() error {
	return nil
}

type hijackerWriter struct {
	http.ResponseWriter
}

func newHijackWriter(w http.ResponseWriter) network.ExtWriter {
	return &hijackerWriter{
		ResponseWriter: w,
	}
}

func (*hijackerWriter) Finalize() error {
	return nil
}

func (*hijackerWriter) Flush() error {
	return nil
}
