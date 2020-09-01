package cluster

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gsf/cluster/addr"
	"github.com/yddeng/gsf/util"
	dnet "github.com/yddeng/gsf/util/net"
	"github.com/yddeng/gsf/util/rpc"
	"time"
)

type rpcManager struct {
	rpcServer *rpc.Server
	rpcClient *rpc.Client
}

type RPCChannel struct {
	session dnet.Session
}

func (this *RPCChannel) SendRequest(req *rpc.Request) error {
	if this.session == nil {
		return fmt.Errorf("rpc session is nil")
	}
	return this.session.Send(req)
}

func (this *RPCChannel) SendResponse(resp *rpc.Response) error {
	if this.session == nil {
		return fmt.Errorf("rpc session is nil")
	}
	return this.session.Send(resp)
}

func (this *rpcManager) asynCall(end *endpoint, data proto.Message, callback func(interface{}, error)) error {
	end.Lock()
	defer end.Unlock()
	// 给自己的消息
	if end.logic.Logic == selfPoint.logic.Logic {

	}
	if end.session == nil {
		end.callMsg = append(end.callMsg, &call{
			msg:      data,
			callback: callback,
			deadline: time.Now().Add(rpcTimeout),
			to:       end.logic.Logic,
		})
		dial(end)
		return nil
	}
	return this.rpcClient.AsynCall(&RPCChannel{session: end.session}, proto.MessageName(data), data, rpcTimeout, callback)
}

func AsynCall(logic addr.LogicAddr, data proto.Message, callback func(interface{}, error)) error {
	end := endpoints.getEndpointByLogic(logic)
	if end == nil {
		return fmt.Errorf("%s is not found", logic.String())
	}
	return rpcMgr.asynCall(end, data, callback)
}

func RegisterRPCMethod(rpcMsg proto.Message, h rpc.MethodHandler) {
	name := proto.MessageName(rpcMsg)
	util.Must(nil, rpcMgr.rpcServer.Register(name, h))
}
