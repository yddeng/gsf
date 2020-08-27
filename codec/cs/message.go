package cs

import (
	"github.com/golang/protobuf/proto"
)

type Message struct {
	seqNo   uint32
	data    proto.Message
	cmd     uint16
	errCode uint16 // 错误码
}

func NewMessage(seqNo uint32, data proto.Message) *Message {
	return &Message{seqNo: seqNo, data: data}
}

func ErrMessage(seqNo uint32, cmd uint16, errCode uint16) *Message {
	return &Message{seqNo: seqNo, cmd: cmd, errCode: errCode}
}

func (this *Message) GetData() proto.Message {
	return this.data
}

func (this *Message) GetSeqNo() uint32 {
	return this.seqNo
}

func (this *Message) GetErrCode() uint16 {
	return this.errCode
}

func (this *Message) GetCmd() uint16 {
	return this.cmd
}
