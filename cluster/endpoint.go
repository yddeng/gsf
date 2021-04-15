package cluster

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/clugs/cluster/addr"
	"github.com/yddeng/clugs/codec/ss"
	"github.com/yddeng/dnet"
	"github.com/yddeng/dnet/drpc"
	"reflect"
	"sync"
	"time"
)

type endpoint struct {
	logic       *addr.Addr
	session     dnet.Session
	dialing     bool
	dialTimeout time.Time

	ssMsg  []*ss.Message
	reqMsg []*drpc.Request
	mtx    *sync.Mutex
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
		ssMsg:       make([]*ss.Message, 0, 4),
		reqMsg:      make([]*drpc.Request, 0, 4),
	}
}

// 传入 *ss.Message, *drpc.Request, *drpc.Response
func (this *endpoint) send(msg interface{}) error {
	// 发送给自己的消息，直接处理
	if this.logic.Logic == selfPoint.logic.Logic {
		eventQueue.Push(func() {
			var err error
			switch msg.(type) {
			case *ss.Message:
				req := msg.(*ss.Message)
				req.SetCmd(pb.GetIdByName(protoss.SS_SPACE, proto.MessageName(req.GetData())))
				err = dispatchSS(selfPoint.logic.Logic, msg.(*ss.Message))
			case *drpc.Request:
				err = rpcMgr.rpcServer.OnRPCRequest(selfPoint, msg.(*drpc.Request))
			case *drpc.Response:
				err = rpcMgr.rpcClient.OnRPCResponse(msg.(*drpc.Response))
			default:
				err = fmt.Errorf("invalid type:%s", reflect.TypeOf(msg).String())
			}
			if err != nil {
				util.Logger().Errorf(err.Error())
			}
		})
		return nil
	}

	// 与对端木有建立连接。先暂存消息，建立连接
	if this.session == nil {
		switch msg.(type) {
		case *ss.Message:
			this.ssMsg = append(this.ssMsg, msg.(*ss.Message))
		case *drpc.Request:
			this.reqMsg = append(this.reqMsg, msg.(*drpc.Request))
		default:
			util.Logger().Debugf("pending invalid type:%s", reflect.TypeOf(msg).String())
		}
		dial(this)
		return nil //fmt.Errorf("%s session is nil", this.logic.Logic.String())
	}
	return this.session.Send(msg)
}

func (this *endpoint) SendRequest(req *drpc.Request) error {
	return this.send(req)
}

func (this *endpoint) SendResponse(resp *drpc.Response) error {
	return this.send(resp)
}

type endpointGroup struct {
	logic2End map[addr.LogicAddr]*endpoint
	type2End  map[uint32]*endpoint
	*sync.Mutex
}

func (this *endpointGroup) rangeEach(f func(end *endpoint) bool) {
	this.Lock()
	defer this.Unlock()
	for _, end := range this.logic2End {
		if !f(end) {
			break
		}
	}
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
