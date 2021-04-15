package clusterpb

import (
	"github.com/yddeng/clugs/codec/pb"
)

const (
	SS_SPACE   = "center_ss"
	REQ_SPACE  = "center_req"
	RESP_SPACE = "center_resp"
)

const (
	LoginReqCmd       = 1
	LoginRespCmd      = 2
	HeartbeatCmd      = 3
	NotifyNodeInfoCmd = 4
	NodeLeaveCmd      = 5
	NodeEnterCmd      = 6
)

func init() {
	pb.Register(REQ_SPACE, &LoginReq{}, LoginReqCmd)
	pb.Register(RESP_SPACE, &LoginResp{}, LoginRespCmd)
	pb.Register(SS_SPACE, &Heartbeat{}, HeartbeatCmd)
	pb.Register(SS_SPACE, &NotifyNodeInfo{}, NotifyNodeInfoCmd)
	pb.Register(SS_SPACE, &NodeLeave{}, NodeLeaveCmd)
	pb.Register(SS_SPACE, &NodeEnter{}, NodeEnterCmd)
}
