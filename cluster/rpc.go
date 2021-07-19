package cluster

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/clugs/cluster/addr"
	"github.com/yddeng/clugs/logger"
	"github.com/yddeng/dnet/drpc"
)

type rpcManager struct {
	rpcServer *drpc.Server
	rpcClient *drpc.Client
}

var rpcMgr = &rpcManager{
	rpcServer: drpc.NewServer(),
	rpcClient: drpc.NewClient(),
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

func AsyncCall(logic addr.LogicAddr, data proto.Message, callback func(interface{}, error)) error {
	end := endGroup.getEndpoint(logic)
	if end == nil {
		logger.Errorf("%s is not found", logic.String())
		return fmt.Errorf("%s is not found", logic.String())
	}

	end.Lock()
	defer end.Unlock()
	return rpcMgr.rpcClient.Go(end, proto.MessageName(data), data, drpc.DefaultRPCTimeout, callback)
}

func RegisterRPCMethod(name string, h drpc.MethodHandler) {
	rpcMgr.rpcServer.Register(name, h)
}
