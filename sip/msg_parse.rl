// -*-go-*-
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
	"fmt"
)

%% machine msg;
%% include sip "sip.rl";
%% write data;

// ParseMsg turns a SIP message byte slice into a data structure.
func ParseMsg(data []byte) (msg *Msg, pos int, err error) {
	if data == nil {
		return nil, 0, nil
	}
	msg = new(Msg)
	viap := &msg.Via
	cs := 0
	p := 0
	pe := len(data)
	eof := len(data)
	buf := make([]byte, len(data))
	amt := 0
	mark := 0
	clen := 0
	ctype := ""
	var name string
	var hex byte
	var value *string
	var via *Via
	var addrp **Addr
	var addr *Addr

	%% main := Message;
	%% write init;
	%% write exec;

	if cs < msg_first_final {
		if p == pe {
			return nil, p, MsgParseError{Code: IncompleteHeader, Msg: []uint8{}, Offset: 0}
		} else {
			return nil, p, MsgParseError{Code: ParseError, Msg: data, Offset: p}
		}
	}

	if clen > 0 {
		if clen > len(data) - p {
			return nil, p, MsgParseError{
			    Code: IncompletePayload,
			    Msg: []byte(fmt.Sprintf("Content-Length: %d", clen)),
			    Offset: len(data) - p}
		}
		msg.Payload = &MiscPayload{T: ctype, D: data[p:p+clen]}
	}
	return msg, p+clen, nil
}

func lookAheadWSP(data []byte, p, pe int) bool {
	return p + 2 < pe && (data[p+2] == ' ' || data[p+2] == '\t')
}

func lastAddr(addrp **Addr) **Addr {
	if *addrp == nil {
		return addrp
	} else {
		return lastAddr(&(*addrp).Next)
	}
}
