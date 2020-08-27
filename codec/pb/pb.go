package pb

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gsf/util/protocol"
	"github.com/yddeng/gsf/util/protocol/protobuf"
)

var nameSpace = map[string]*protocol.Protocol{}

type nameCmd struct {
	cmd2Name map[uint16]string
	name2Cmd map[string]uint16
}

var spaceCmd map[string]*nameCmd

func newProtocol() *protocol.Protocol {
	return protocol.NewProtoc(&protobuf.Protobuf{})
}

//根据名字注册实例(注意函数非线程安全，需要在初始化阶段完成所有消息的Register)
func Register(namespace string, msg proto.Message, id uint16) error {
	var ns *protocol.Protocol
	var sc *nameCmd
	var ok bool

	if ns, ok = nameSpace[namespace]; !ok {
		ns = newProtocol()
		nameSpace[namespace] = ns
	}

	if sc, ok = spaceCmd[namespace]; !ok {
		sc = &nameCmd{
			cmd2Name: map[uint16]string{},
			name2Cmd: map[string]uint16{},
		}
		spaceCmd[namespace] = sc
	}
	name := proto.MessageName(msg)
	sc.cmd2Name[id] = name
	sc.name2Cmd[name] = id

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

func GetNameById(namespace string, id uint16) string {
	sc, ok := spaceCmd[namespace]
	if !ok {
		return ""
	}
	return sc.cmd2Name[id]
}

func GetIdByName(namespace string, name string) uint16 {
	sc, ok := spaceCmd[namespace]
	if !ok {
		return 0
	}
	return sc.name2Cmd[name]
}
