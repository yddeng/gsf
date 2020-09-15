package gate

import (
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gsf/cluster"
	"github.com/yddeng/gsf/cluster/addr"
	"github.com/yddeng/gsf/protocol/ss"
	"github.com/yddeng/gsf/util"
	"github.com/yddeng/gsf/util/net"
)

func Launch(externalAddr string) error {
	listener, err := net.NewTCPListener("tcp", externalAddr)
	if err != nil {
		util.Logger().Errorln(err)
		return err
	}

	_ = listener.Listen(func(session net.Session) {

	})

	return nil
}

func init() {
	cluster.RegisterSSMethod(ss.Echo, func(from addr.LogicAddr, msg proto.Message) {
		util.Logger().Debugf("ss echo from %s msg %v", from.String(), msg)
	})
}
