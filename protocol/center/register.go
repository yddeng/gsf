package ss

import (
	"github.com/yddeng/gsf/codec/pb"
	"github.com/yddeng/gsf/protocol/center/center"
)

func init() {
	pb.Register("center_req", &center.LoginReq{}, 1)
	pb.Register("center_resp", &center.LoginResp{}, 2)
	pb.Register("center_ss", &center.Heartbeat{}, 3)
	pb.Register("center_ss", &center.NotifyNodeInfo{}, 4)
	pb.Register("center_ss", &center.NodeLeave{}, 5)
	pb.Register("center_ss", &center.NodeChange{}, 6)
}
