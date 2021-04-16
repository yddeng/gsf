package rpc

import (
	"github.com/yddeng/clugs/codec/pb"
	"github.com/yddeng/clugs/protocol/rpc/rpcpb"
)

const (
	REQ_SPACE  = "rpc_req"
	RESP_SPACE = "rpc_resp"
)

const(
	Echo = 1

)

func init() {
	// rpc_req
	pb.Register(REQ_SPACE,&rpcpb.EchoReq{},1)

	// rpc_resp
	pb.Register(RESP_SPACE,&rpcpb.EchoResp{},1)

}
