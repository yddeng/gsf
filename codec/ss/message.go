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

func (this *Message) GetData() interface{} {
	return this.data
}

func (this *Message) GetCmd() uint16 {
	return this.cmd
}
