// Copyright 2018 The go-zeromq Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package zmq4

import (
	"bytes"
	"fmt"
	"io"
)

type MsgType byte

const (
	UsrMsg MsgType = 0
	CmdMsg MsgType = 1
)

// Msg is a ZMTP message, possibly composed of multiple frames.
type Msg struct {
	Frames    [][]byte
	err       error
	Type      MsgType
	multipart bool
}

// NewMsg creates a message and sets the content to frame.
//
// The underlying storage of frame is now shared with the message!
func NewMsg(frame []byte) Msg {
	return Msg{Frames: [][]byte{frame}}
}

// NewMsgFrom creates a message and sets the content to the given frames.
//
// The underlying storage of the frames is now shared with the message!
func NewMsgFrom(frames ...[]byte) Msg {
	return Msg{Frames: frames}
}

func NewMsgString(frame string) Msg {
	return NewMsg([]byte(frame))
}

func NewMsgFromString(frames []string) Msg {
	msg := Msg{Frames: make([][]byte, len(frames))}
	for i, frame := range frames {
		msg.Frames[i] = append(msg.Frames[i], []byte(frame)...)
	}
	return msg
}

func (msg Msg) isCmd() bool {
	return msg.Type == CmdMsg
}

func (msg Msg) Err() error {
	return msg.err
}

// Bytes returns the concatenated content of all its frames.
func (msg Msg) Bytes() []byte {
	buf := make([]byte, 0, msg.size())
	for _, frame := range msg.Frames {
		buf = append(buf, frame...)
	}
	return buf
}

func (msg Msg) size() int {
	n := 0
	for _, frame := range msg.Frames {
		n += len(frame)
	}
	return n
}

func (msg Msg) Size() int {
	return msg.size()
}

func (msg Msg) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString("Msg{Frames:{")
	for i, frame := range msg.Frames {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(buf, "%q", frame)
	}
	buf.WriteString("}}")
	return buf.String()
}

func (msg Msg) Clone() Msg {
	o := Msg{Frames: make([][]byte, len(msg.Frames))}
	for i, frame := range msg.Frames {
		o.Frames[i] = make([]byte, len(frame))
		copy(o.Frames[i], frame)
	}
	return o
}

// The Reset method is used to clear a message.
//
// The messages Size() will be 0 afterwards and
// len(Frames) will also be 0.
// Reset does not touch the underlying storage. It just
// sets the len of the slices to 0.
func (m *Msg) Reset() {
	for i := range m.Frames {
		f := m.Frames[i]
		allocator.free(&f)
		m.Frames[i] = nil
	}
	fr := m.Frames[:0]
	m.Frames = fr
}

// Cmd is a ZMTP Cmd as per:
//  https://rfc.zeromq.org/spec:23/ZMTP/#formal-grammar
type Cmd struct {
	Name string
	Body []byte
}

func (cmd *Cmd) unmarshalZMTP(data []byte) error {
	if len(data) == 0 {
		return io.ErrUnexpectedEOF
	}
	n := int(data[0])
	if n > len(data)-1 {
		return ErrBadCmd
	}
	cmd.Name = string(data[1 : n+1])
	cmd.Body = data[n+1:]
	return nil
}

func (cmd *Cmd) marshalZMTP() ([]byte, error) {
	n := len(cmd.Name)
	if n > 255 {
		return nil, ErrBadCmd
	}

	buf := make([]byte, 0, 1+n+len(cmd.Body))
	buf = append(buf, byte(n))
	buf = append(buf, []byte(cmd.Name)...)
	buf = append(buf, cmd.Body...)
	return buf, nil
}

// ZMTP commands as per:
//  https://rfc.zeromq.org/spec:23/ZMTP/#commands
const (
	CmdCancel      = "CANCEL"
	CmdError       = "ERROR"
	CmdHello       = "HELLO"
	CmdInitiate    = "INITIATE"
	CmdPing        = "PING"
	CmdPong        = "PONG"
	CmdReady       = "READY"
	CmdSubscribe   = "SUBSCRIBE"
	CmdUnsubscribe = "UNSUBSCRIBE"
	CmdWelcome     = "WELCOME"
)
