package cluster

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/dnet/drpc"
	"github.com/yddeng/gsf/cluster/addr"
	"github.com/yddeng/gsf/util"
)

type rpcManager struct {
	rpcServer *drpc.Server
	rpcClient *drpc.Client
}

/*
type RPCChannel struct {
	session dnet.Session
}

func (this *RPCChannel) SendRequest(req *drpc.Request) error {
	if this.session == nil {
		return fmt.Errorf("drpc session is nil")
	}
	return this.session.Send(req)
}

func (this *RPCChannel) SendResponse(resp *drpc.Response) error {
	if this.session == nil {
		return fmt.Errorf("drpc session is nil")
	}
	return this.session.Send(resp)
}
*/

func AsynCall(logic addr.LogicAddr, data proto.Message, callback func(interface{}, error)) error {
	end := endpoints.getEndpointByLogic(logic)
	if end == nil {
		util.Logger().Errorf("%s is not found", logic.String())
		return fmt.Errorf("%s is not found", logic.String())
	}

	end.Lock()
	defer end.Unlock()
	return rpcMgr.rpcClient.AsynCall(end, proto.MessageName(data), data, rpcTimeout, callback)
}

func RegisterRPCMethod(rpcMsg proto.Message, h drpc.MethodHandler) {
	name := proto.MessageName(rpcMsg)
	util.Must(nil, rpcMgr.rpcServer.Register(name, h))
}
