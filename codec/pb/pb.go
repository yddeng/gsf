package pb

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gsf/util/protocol"
	"github.com/yddeng/gsf/util/protocol/protobuf"
)

var nameSpace = map[string]*protocol.Protocol{}

func newProtocol() *protocol.Protocol {
	return protocol.NewProtoc(&protobuf.Protobuf{})
}

//根据名字注册实例(注意函数非线程安全，需要在初始化阶段完成所有消息的Register)
func Register(namespace string, msg proto.Message, id uint16) error {
	var ns *protocol.Protocol
	var ok bool

	if ns, ok = nameSpace[namespace]; !ok {
		ns = newProtocol()
		nameSpace[namespace] = ns
	}

	return ns.Register(id, msg)
}

func Marshal(namespace string, o interface{}) (uint16, []byte, error) {
	var ns *protocol.Protocol
	var ok bool
	if ns, ok = nameSpace[namespace]; !ok {
		return 0, nil, fmt.Errorf("invaild namespace:%s", namespace)
	}
	return ns.Marshal(o)
}

func Unmarshal(namespace string, id uint16, buff []byte) (interface{}, error) {
	var ns *protocol.Protocol
	var ok bool
	if ns, ok = nameSpace[namespace]; !ok {
		return nil, fmt.Errorf("invaild namespace:%s", namespace)
	}

	return ns.Unmarshal(id, buff)
}
