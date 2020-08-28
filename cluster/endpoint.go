package cluster

import (
	"fmt"
	"github.com/yddeng/gsf/cluster/addr"
	"github.com/yddeng/gsf/util/net"
	"github.com/yddeng/gsf/util/rpc"
)

type EndPoint struct {
	logic   *addr.Addr
	session net.Session
	dialing bool
}

// 传入 *ss.Message, *rpc.Request, *rpc.Response
func (this *EndPoint) Send(o interface{}) error {
	if this.session == nil {
		return fmt.Errorf("%s session is nil", this.logic.Logic.String())
	}
	return this.session.Send(o)
}

func (this *EndPoint) SendRequest(req *rpc.Request) error {
	return this.Send(req)
}

func (this *EndPoint) SendResponse(resp *rpc.Response) error {
	return this.Send(resp)
}
