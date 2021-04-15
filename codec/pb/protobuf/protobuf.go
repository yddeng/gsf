package protobuf

import "github.com/golang/protobuf/proto"

/*
 protobuf 协议解析
*/

type Protobuf struct{}

func (_ Protobuf) Unmarshal(data []byte, o interface{}) (err error) {
	return proto.Unmarshal(data, o.(proto.Message))
}

func (_ Protobuf) Marshal(data interface{}) ([]byte, error) {
	return proto.Marshal(data.(proto.Message))
}
