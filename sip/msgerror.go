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

package sip

import (
	"bytes"
	"fmt"
)

// Code is result of ParseMsg.
type Code int

const (
	IncompleteHeader Code = iota
	IncompletePayload
	ParseError
)

//type MsgIncompleteError struct {
//	Msg []byte
//}
//
//func (err MsgIncompleteError) Error() string {
//	return fmt.Sprintf("Incomplete SIP message:\n%s", err.Msg)
//}

type MsgParseError struct {
	Code   Code
	Msg    []byte
	Offset int
}

func (err MsgParseError) Error() string {
	switch err.Code {
	case IncompleteHeader:
		return "Incomplete SIP Header"
	case IncompletePayload:
		return string(err.Msg)
	case ParseError:
		lines := bytes.Split(err.Msg, []byte("\r\n"))
		var b bytes.Buffer
		o := 0
		line := 0
		lineOffset := 0
		for i := 0; i < len(lines); i++ {
			if o <= err.Offset && err.Offset < o+len(lines[i]) {
				b.Write(lines[i])
				line = i + 1
				lineOffset = err.Offset - o
				b.WriteByte('\n')
				for j := 0; j < lineOffset; j++ {
					b.WriteByte(' ')
				}
				b.WriteByte('^')
				b.WriteByte('\n')
			}
			o += len(lines[i]) + 2
		}
		return fmt.Sprintf("Error at line %d offset %d:\n%s", line, lineOffset, b.String())
	}

	return ""
}
