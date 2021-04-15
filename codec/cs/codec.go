package cs

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/clugs/codec/pb"
	_ "github.com/yddeng/clugs/protocol/cs"
	"github.com/yddeng/dutil/buffer"
	"io"
	"reflect"
)

// 编解码器
// 消息 -- 格式: 消息头(序列号＋协议号+flag+错误码+消息体长度), 消息体

const (
	seqNoSize   = 4                                                       // 消息序列号
	cmdSize     = 2                                                       // 消息ID（消息体的编码ID，对应的反序列化结构）
	flagSize    = 1                                                       // 标记（消息体是否加密，压缩）
	errCodeSize = 2                                                       // 消息错误码（请求失败的消息直接返回错误码，没有消息体）
	bodySize    = 2                                                       // 消息长度（消息体的长度）
	headSize    = seqNoSize + cmdSize + flagSize + errCodeSize + bodySize // 消息头长度
	buffSize    = 65535 - headSize                                        // 消息体最大长度
)

type Codec struct {
	decode, encode string
	*Decoder
}

func NewCodec(decode, encode string) *Codec {
	return &Codec{
		decode: decode,
		encode: encode,
		Decoder: &Decoder{
			readBuf: buffer.NewBufferWithCap(65535),
		},
	}
}

type Decoder struct {
	readBuf  *buffer.Buffer
	readHead bool
	seqNo    uint32
	cmd      uint16
	flag     byte
	errCode  uint16
	bodyLen  uint16
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

func (decoder *Codec) unPack() (*Message, error) {

	if !decoder.readHead {
		if decoder.readBuf.Len() < headSize {
			return nil, nil
		}

		decoder.seqNo, _ = decoder.readBuf.ReadUint32BE()
		decoder.cmd, _ = decoder.readBuf.ReadUint16BE()
		decoder.flag, _ = decoder.readBuf.ReadByte()
		decoder.errCode, _ = decoder.readBuf.ReadUint16BE()
		decoder.bodyLen, _ = decoder.readBuf.ReadUint16BE()
		decoder.readHead = true
	}

	var msg proto.Message
	if decoder.bodyLen != 0 {
		if decoder.readBuf.Len() < int(decoder.bodyLen) {
			return nil, nil
		}

		data, _ := decoder.readBuf.ReadBytes(int(decoder.bodyLen))
		i, err := pb.Unmarshal(decoder.decode, decoder.cmd, data)
		if err != nil {
			return nil, err
		}
		msg = i.(proto.Message)
	}

	decoder.readHead = false
	return &Message{
		seqNo:   decoder.seqNo,
		data:    msg,
		cmd:     decoder.cmd,
		errCode: decoder.errCode,
	}, nil
}

//编码
func (encoder *Codec) Encode(o interface{}) ([]byte, error) {
	msg, ok := o.(*Message)
	if !ok {
		return nil, fmt.Errorf(" o'type %s is't *cs.Message", reflect.TypeOf(o).String())
	}

	var dataLen int
	var data []byte
	var cmd = msg.GetCmd()
	var err error
	if msg.errCode == 0 {
		cmd, data, err = pb.Marshal(encoder.encode, msg.GetData())
		if err != nil {
			return nil, err
		}

		dataLen = len(data)
		if dataLen > buffSize {
			return nil, fmt.Errorf("encode dataLen is too large,len: %d", dataLen)
		}
	}

	msgLen := dataLen + headSize
	buff := buffer.NewBufferWithCap(msgLen)

	//写入seqNo
	buff.WriteUint32BE(msg.GetSeqNo())
	//写入cmd
	buff.WriteUint16BE(cmd)
	//写入flag
	buff.WriteByte(uint8(0))
	// errCode
	buff.WriteUint16BE(msg.GetErrCode())
	//写入data长度
	buff.WriteUint16BE(uint16(dataLen))
	if dataLen != 0 {
		//data数据
		buff.WriteBytes(data)
	}
	return buff.Bytes(), nil
}
