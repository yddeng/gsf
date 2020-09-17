package gate

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/dnet"
	"github.com/yddeng/dnet/drpc"
	"github.com/yddeng/dnet/dtcp"
	"github.com/yddeng/gsf/center/protocol"
	"github.com/yddeng/gsf/cluster"
	"github.com/yddeng/gsf/cluster/addr"
	"github.com/yddeng/gsf/codec/cs"
	protocs "github.com/yddeng/gsf/protocol/cs"
	"github.com/yddeng/gsf/protocol/ss"
	"github.com/yddeng/gsf/util"
	"reflect"
)

func Launch(externalAddr string) error {
	listener, err := dtcp.NewTCPListener("tcp", externalAddr)
	if err != nil {
		util.Logger().Errorln(err)
		return err
	}

	err = listener.Listen(func(session dnet.Session) {
		util.Logger().Infoln("new client", session.RemoteAddr().String())
		// 超时时间
		session.SetTimeout(0, 0)
		session.SetCodec(cs.NewCodec(protocs.SC_SPACE, protocs.CS_SPACE))
		session.SetCloseCallBack(onClose)

		err := session.Start(func(data interface{}, err error) {
			if err != nil {
				session.Close(err.Error())
			} else {
				eventQueue.Push(func() {
					var err error
					switch data.(type) {
					case *ss.Message:
						dispatchMsg(session, data.(*ss.Message))
					case *drpc.Request:
						err = rpcServer.OnRPCRequest(&Node{session: session}, data.(*drpc.Request))
					//case *drpc.Response:
					default:
						err = fmt.Errorf("invalid type:%s", reflect.TypeOf(data).String())
					}
					if err != nil {
						util.Logger().Errorf(err.Error())
					}
				})
			}
		})
		if err != nil {
			util.Logger().Errorf("%s start session err: %s", session.RemoteAddr().String(), err)
		}
	})
	return err
}

func init() {
	cluster.RegisterSSMethod(ss.Echo, func(from addr.LogicAddr, msg proto.Message) {
		util.Logger().Debugf("ss echo from %s msg %v", from.String(), msg)
	})
}
