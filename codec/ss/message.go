package ss

import (
	"github.com/golang/protobuf/proto"
)

type Message struct {
	data proto.Message
	cmd  uint16
}

func NewMessage(data proto.Message) *Message {
	return &Message{data: data}
}

func (this *Message) GetData() proto.Message {
	return this.data
}

func (this *Message) GetCmd() uint16 {
	return this.cmd
}

func (this *Message) SetCmd(cmd uint16) {
	this.cmd = cmd
}
