package rpc

import (
	"github.com/yddeng/gsf/codec/pb"
	"github.com/yddeng/gsf/protocol/rpc/rpc"
)

const(
	Echo = 1

)

func init() {
	// rpc_req
	pb.Register("rpc_req",&rpc.EchoReq{},1)

	// rpc_resp
	pb.Register("rpc_resp",&rpc.EchoResp{},1)

}
