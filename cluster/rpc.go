package cluster

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gsf/cluster/addr"
	"github.com/yddeng/gsf/util"
	"github.com/yddeng/gsf/util/rpc"
)

type rpcManager struct {
	rpcServer *rpc.Server
	rpcClient *rpc.Client
}

/*
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

func RegisterRPCMethod(rpcMsg proto.Message, h rpc.MethodHandler) {
	name := proto.MessageName(rpcMsg)
	util.Must(nil, rpcMgr.rpcServer.Register(name, h))
}
