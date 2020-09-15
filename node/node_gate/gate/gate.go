package gate

import (
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/dnet"
	"github.com/yddeng/dnet/dtcp"
	"github.com/yddeng/gsf/cluster"
	"github.com/yddeng/gsf/cluster/addr"
	"github.com/yddeng/gsf/protocol/ss"
	"github.com/yddeng/gsf/util"
)

func Launch(externalAddr string) error {
	listener, err := dtcp.NewTCPListener("tcp", externalAddr)
	if err != nil {
		util.Logger().Errorln(err)
		return err
	}

	_ = listener.Listen(func(session dnet.Session) {

	})

	return nil
}

func init() {
	cluster.RegisterSSMethod(ss.Echo, func(from addr.LogicAddr, msg proto.Message) {
		util.Logger().Debugf("ss echo from %s msg %v", from.String(), msg)
	})
}
