package cluster

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gsf/cluster/addr"
	"github.com/yddeng/gsf/codec/ss"
	dnet "github.com/yddeng/gsf/util/net"
	"sync"
	"time"
)

type endpoint struct {
	logic       *addr.Addr
	session     dnet.Session
	dialing     bool
	dialTimeout time.Time
	postMsg     []proto.Message
	callMsg     []*call
	*sync.Mutex
}

type call struct {
	msg      proto.Message
	callback func(interface{}, error)
	deadline time.Time
	to       addr.LogicAddr
}

func newEndpoint(logic *addr.Addr) *endpoint {
	return &endpoint{
		logic:       logic,
		dialTimeout: time.Time{},
		postMsg:     make([]proto.Message, 0, 4),
		callMsg:     make([]*call, 0, 4),
		Mutex:       new(sync.Mutex),
	}
}

// 传入 *ss.Message
func (this *endpoint) send(msg *ss.Message) error {
	if this.session == nil {
		return fmt.Errorf("%s session is nil", this.logic.Logic.String())
	}
	return this.session.Send(msg)
}

//func (this *endpoint) SendRequest(req *rpc.Request) error {
//	return this.Send(req)
//}
//
//func (this *endpoint) SendResponse(resp *rpc.Response) error {
//	return this.Send(resp)
//}

type endpointGroup struct {
	logic2End map[addr.LogicAddr]*endpoint
	type2End  map[uint32]*endpoint
	*sync.Mutex
}

func (this *endpointGroup) addEndpoint(logic *addr.Addr) {
	this.Lock()
	defer this.Unlock()
	_, ok := this.logic2End[logic.Logic]
	if !ok {
		this.logic2End[logic.Logic] = newEndpoint(logic)
	}
}

func (this *endpointGroup) removeEndpoint(logic addr.LogicAddr) {
	this.Lock()
	defer this.Unlock()

	_, ok := this.logic2End[logic]
	if ok {
		delete(this.logic2End, logic)
	}
}

func (this *endpointGroup) getEndpointByLogic(logic addr.LogicAddr) *endpoint {
	this.Lock()
	defer this.Unlock()

	end, ok := this.logic2End[logic]
	if ok {
		return end
	}
	return nil
}
