package ss

import (
	"github.com/yddeng/clugs/codec/pb"
	"github.com/yddeng/clugs/protocol/ss/sspb"
)

const (
	SS_SPACE = "s2s"
)

const(
	Echo = 1

)

func init() {
	pb.Register(SS_SPACE,&sspb.Echo{},1)

}
