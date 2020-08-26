package cs

import (
	"github.com/yddeng/gsf/codec/pb"
	"github.com/yddeng/gsf/protocol/cs/cs"
)

const(
	Heartbeat = 1
	Echo = 2

)

func init() {
	//toS
	pb.Register("c2s",&cs.HeartbeatToS{},1)
	pb.Register("c2s",&cs.EchoToS{},2)

	//toC
	pb.Register("s2c",&cs.HeartbeatToC{},1)
	pb.Register("s2c",&cs.EchoToC{},2)

}
