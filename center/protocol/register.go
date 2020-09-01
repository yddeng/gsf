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
	pb.Register("center_req", &LoginReq{}, 1)
	pb.Register("center_resp", &LoginResp{}, 2)
	pb.Register("center_ss", &Heartbeat{}, 3)
	pb.Register("center_ss", &NotifyNodeInfo{}, 4)
	pb.Register("center_ss", &NodeLeave{}, 5)
	pb.Register("center_ss", &NodeEnter{}, 6)
}
