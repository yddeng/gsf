package protocol

import (
	"github.com/yddeng/gsf/codec/pb"
)

const (
	SS_SPACE   = "center_ss"
	REQ_SPACE  = "center_req"
	RESP_SPACE = "center_resp"
)

const (
	HeartbeatCmd      = 3
	NotifyNodeInfoCmd = 4
	NodeLeaveCmd      = 5
	NodeEnterCmd      = 6
)

func init() {
	pb.Register(REQ_SPACE, &LoginReq{}, 1)
	pb.Register(RESP_SPACE, &LoginResp{}, 2)
	pb.Register(SS_SPACE, &Heartbeat{}, 3)
	pb.Register(SS_SPACE, &NotifyNodeInfo{}, 4)
	pb.Register(SS_SPACE, &NodeLeave{}, 5)
	pb.Register(SS_SPACE, &NodeEnter{}, 6)
}
