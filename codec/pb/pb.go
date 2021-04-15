package pb

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/clugs/codec/pb/protobuf"
)

type SpaceProtocol struct {
	protoMap map[string]*Protocol
	cmdMap   map[string]*spaceCmd
}

type spaceCmd struct {
	cmd2Name map[uint16]string
	name2Cmd map[string]uint16
}

func NewSpaceProtocol() *SpaceProtocol {
	return &SpaceProtocol{
		protoMap: map[string]*Protocol{},
		cmdMap:   map[string]*spaceCmd{},
	}
}

//根据名字注册实例(注意函数非线程安全，需要在初始化阶段完成所有消息的Register)
func (this *SpaceProtocol) Register(namespace string, msg proto.Message, id uint16) {
	var (
		ns *Protocol
		sc *spaceCmd
		ok bool
	)

	if ns, ok = this.protoMap[namespace]; !ok {
		ns = NewProtocol(new(protobuf.Protobuf))
		this.protoMap[namespace] = ns
	}

	if sc, ok = this.cmdMap[namespace]; !ok {
		sc = &spaceCmd{
			cmd2Name: map[uint16]string{},
			name2Cmd: map[string]uint16{},
		}
		this.cmdMap[namespace] = sc
	}

	if _, ok := sc.cmd2Name[id]; ok {
		panic(fmt.Sprintf("id %d id areadly register", id))
	}

	name := proto.MessageName(msg)
	sc.cmd2Name[id] = name
	sc.name2Cmd[name] = id

	ns.Register(id, msg)
}

func (this *SpaceProtocol) Marshal(namespace string, o interface{}) (uint16, []byte, error) {
	var ns *Protocol
	var ok bool
	if ns, ok = this.protoMap[namespace]; !ok {
		return 0, nil, fmt.Errorf("invaild namespace:%s", namespace)
	}
	return ns.Marshal(o)
}

func (this *SpaceProtocol) Unmarshal(namespace string, id uint16, buff []byte) (interface{}, error) {
	var ns *Protocol
	var ok bool
	if ns, ok = this.protoMap[namespace]; !ok {
		return nil, fmt.Errorf("invaild namespace:%s", namespace)
	}

	return ns.Unmarshal(id, buff)
}

func (this *SpaceProtocol) GetNameById(namespace string, id uint16) string {
	sc, ok := this.cmdMap[namespace]
	if !ok {
		return ""
	}
	return sc.cmd2Name[id]
}

func (this *SpaceProtocol) GetIdByName(namespace string, name string) uint16 {
	sc, ok := this.cmdMap[namespace]
	if !ok {
		return 0
	}
	return sc.name2Cmd[name]
}

var defP = NewSpaceProtocol()

func Register(namespace string, msg proto.Message, id uint16) {
	defP.Register(namespace, msg, id)
}

func Marshal(namespace string, o interface{}) (uint16, []byte, error) {
	return defP.Marshal(namespace, o)
}

func Unmarshal(namespace string, id uint16, buff []byte) (interface{}, error) {
	return defP.Unmarshal(namespace, id, buff)
}

func GetNameById(namespace string, id uint16) string {
	return defP.GetNameById(namespace, id)
}

func GetIdByName(namespace string, name string) uint16 {
	return defP.GetIdByName(namespace, name)
}
