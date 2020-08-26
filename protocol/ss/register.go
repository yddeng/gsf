package ss

import (
	"github.com/yddeng/gsf/codec/pb"
	"github.com/yddeng/gsf/protocol/ss/ss"
)

const(
	Echo = 1

)

func init() {
	pb.Register("s2s",&ss.Echo{},1)

}
