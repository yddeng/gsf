package codec

import (
	"github.com/golang/protobuf/proto"
)

type Message struct {
	seqNo uint64
	data  proto.Message
	cmd   uint16
}

func NewMessage(seqNo uint64, data proto.Message) *Message {
	return &Message{seqNo: seqNo, data: data}
}

func (this *Message) GetData() interface{} {
	return this.data
}

func (this *Message) GetCmd() uint16 {
	return this.cmd
}

func (this *Message) GetSeqNo() uint64 {
	return this.seqNo
}
