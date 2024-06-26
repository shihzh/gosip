// Copyright 2020 Justine Alexandra Roberts Tunney
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Addr / Via Parameter Library

package sip

import (
	"bytes"
	"strings"
)

// XHeader is a linked list storing an unrecognized SIP headers.
type XHeader struct {
	Name  string // tokenc
	Value []byte // UTF8, never nil
	Next  *XHeader
}

// Get returns an entry in O(n) time.
func (h *XHeader) Get(name string) *XHeader {
	if h == nil {
		return nil
	}
	if strings.EqualFold(h.Name, name) {
		return h
	}
	return h.Next.Get(name)
}

// String turns XHeader into a string.
func (h *XHeader) String() string {
	var b bytes.Buffer
	h.Append(&b)
	return b.String()
}

// Append serializes headers in insertion order.
func (h *XHeader) Append(b Writer) {
	if h == nil {
		return
	}
	h.Next.Append(b)
	appendSanitized(b, []byte(h.Name), tokenc)
	b.WriteString(": ")
	b.Write(h.Value)
	b.WriteString("\r\n")
}
