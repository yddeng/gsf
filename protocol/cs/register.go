package cs

import (
	"github.com/yddeng/gsf/codec/pb"
	"github.com/yddeng/gsf/protocol/cs/cspb"
)

const (
	CS_SPACE = "c2s"
	SC_SPACE = "s2c"
) 

const(
	Heartbeat = 1
	Echo = 2

)

func init() {
	//toS
	pb.Register(CS_SPACE,&cspb.HeartbeatToS{},1)
	pb.Register(CS_SPACE,&cspb.EchoToS{},2)

	//toC
	pb.Register(SC_SPACE,&cspb.HeartbeatToC{},1)
	pb.Register(SC_SPACE,&cspb.EchoToC{},2)

}
