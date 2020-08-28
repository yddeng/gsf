package cluster

import (
	"github.com/yddeng/gsf/util/net"
	"github.com/yddeng/gsf/util/rpc"
)

type RPCChannel struct {
	session net.Session
}

func (this *RPCChannel) SendRequest(req *rpc.Request) error {
	return this.session.Send(req)
}

func (this *RPCChannel) SendResponse(resp *rpc.Response) error {
	return this.session.Send(resp)
}
