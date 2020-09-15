package ss

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gsf/codec/pb"
	_ "github.com/yddeng/gsf/protocol/rpc"
	_ "github.com/yddeng/gsf/protocol/ss"
	"github.com/yddeng/gsf/util/buffer"
	"github.com/yddeng/gsf/util/rpc"
	"io"
	"reflect"
)

// 编解码器
// 消息 -- 格式: 消息头(消息类型+序列号+协议号+消息体长度), 消息体

const (
	ttSize    = 1                                       // 消息类型
	seqNoSize = 8                                       // 消息序列号
	cmdSize   = 2                                       // 协议号（消息体的编码ID，对应的反序列化结构）
	bodySize  = 2                                       // 消息体长度
	headSize  = ttSize + seqNoSize + cmdSize + bodySize // 消息头长度
	buffSize  = 65535 - headSize                        // 消息体最大长度
)

const (
	SS_Message     = 0x01 // ss普通消息
	RPC_Request    = 0x02 // rpc请求
	RPC_Response   = 0x04 // rpc回复
	RPC_Resp_Error = 0x12 // rpc回复错误
)

type Codec struct {
	ss, req, resp string
	*Decoder
}

func NewCodec(ss, req, resp string) *Codec {
	return &Codec{
		ss:   ss,
		req:  req,
		resp: resp,
		Decoder: &Decoder{
			readBuf: buffer.NewBufferWithCap(65535),
		},
	}
}

type Decoder struct {
	readBuf  *buffer.Buffer
	readHead bool
	tt       byte
	seqNo    uint64
	cmd      uint16
	bodyLen  uint16
}

func checkTT(tt byte) byte {
	if tt&SS_Message != 0 {
		return SS_Message
	}
	if tt&RPC_Request != 0 {
		return RPC_Request
	}
	if tt&RPC_Response != 0 {
		return RPC_Response
	}
	return 0
}

//解码
func (decoder *Codec) Decode(reader io.Reader) (interface{}, error) {
	for {
		msg, err := decoder.unPack()

		if msg != nil {
			return msg, nil

		} else if err == nil {
			_, err1 := decoder.readBuf.ReadFrom(reader)
			if err1 != nil {
				return nil, err1
			}
		} else {
			return nil, err
		}
	}
}

func (decoder *Codec) unPack() (interface{}, error) {
	if !decoder.readHead {
		if decoder.readBuf.Len() < headSize {
			return nil, nil
		}

		decoder.tt, _ = decoder.readBuf.ReadByte()
		decoder.seqNo, _ = decoder.readBuf.ReadUint64BE()
		decoder.cmd, _ = decoder.readBuf.ReadUint16BE()
		decoder.bodyLen, _ = decoder.readBuf.ReadUint16BE()
		decoder.readHead = true
	}

	if decoder.readBuf.Len() < int(decoder.bodyLen) {
		return nil, nil
	}

	data, err := decoder.readBuf.ReadBytes(int(decoder.bodyLen))
	if err != nil {
		return nil, err
	}
	var msg interface{}

	switch checkTT(decoder.tt) {
	case SS_Message:
		m, err := pb.Unmarshal(decoder.ss, decoder.cmd, data)
		if err != nil {
			return nil, err
		}
		msg = &Message{
			data: m.(proto.Message),
			cmd:  decoder.cmd,
		}

	case RPC_Request:
		m, err := pb.Unmarshal(decoder.req, decoder.cmd, data)
		if err != nil {
			return nil, err
		}
		msg = &rpc.Request{
			SeqNo:    decoder.seqNo,
			Method:   pb.GetNameById(decoder.req, decoder.cmd),
			Data:     m,
			NeedResp: true,
		}
	case RPC_Response:
		resp := &rpc.Response{SeqNo: decoder.seqNo}
		if decoder.tt == RPC_Resp_Error {
			resp.Err = fmt.Errorf(string(data))

		} else {
			m, err := pb.Unmarshal(decoder.resp, decoder.cmd, data)
			if err != nil {
				return nil, err
			}
			resp.Data = m
		}
		msg = resp
	default:
		err = fmt.Errorf("unPack err: tt is %d", decoder.tt)
	}

	decoder.readHead = false
	return msg, err
}

//编码
func (encoder *Codec) Encode(o interface{}) ([]byte, error) {
	var tt byte
	var seqNo uint64
	var cmd uint16
	var data []byte
	var bodyLen int
	var err error

	switch o.(type) {
	case *Message:
		msg := o.(*Message)
		tt = SS_Message
		cmd, data, err = pb.Marshal(encoder.ss, msg.GetData())
		if err != nil {
			return nil, err
		}
	case *rpc.Request:
		msg := o.(*rpc.Request)
		tt = RPC_Request
		seqNo = msg.SeqNo
		cmd, data, err = pb.Marshal(encoder.req, msg.Data)
		if err != nil {
			return nil, err
		}

	case *rpc.Response:
		msg := o.(*rpc.Response)
		if msg.Err != nil {
			tt = RPC_Resp_Error
			data = []byte(msg.Err.Error())
		} else {
			tt = RPC_Response
			seqNo = msg.SeqNo
			cmd, data, err = pb.Marshal(encoder.resp, msg.Data)
			if err != nil {
				return nil, err
			}
		}
	default:
		return nil, fmt.Errorf("invailed type:%s", reflect.TypeOf(o).String())
	}

	bodyLen = len(data)
	if bodyLen > buffSize {
		return nil, fmt.Errorf("encode dataLen is too large,len: %d", bodyLen)
	}

	totalLen := headSize + bodyLen
	buff := buffer.NewBufferWithCap(totalLen)
	//tt
	buff.WriteUint8BE(tt)
	//seq
	buff.WriteUint64BE(seqNo)
	//cmd
	buff.WriteUint16BE(cmd)
	//bodylen
	buff.WriteUint16BE(uint16(bodyLen))
	//body
	buff.WriteBytes(data)

	return buff.Bytes(), nil
}
