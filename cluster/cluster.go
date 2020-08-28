package cluster

import (
	"github.com/yddeng/gsf/cluster/addr"
	"github.com/yddeng/gsf/codec/ss"
	protorpc "github.com/yddeng/gsf/protocol/rpc"
	protoss "github.com/yddeng/gsf/protocol/ss"
	"github.com/yddeng/gsf/util"
	"github.com/yddeng/gsf/util/net"
	"github.com/yddeng/gsf/util/queue"
	"github.com/yddeng/gsf/util/rpc"
	"time"
)

var (
	eventQueue  *queue.EventQueue
	endPoints   map[uint32]*EndPoint
	selfPoint   *EndPoint
	centerPoint *CenterPoint
)

func onClose(session net.Session, reason string) {

}

func Launcher(centerAddr string, self *addr.Addr) error {
	l, err := net.NewTCPListener("tcp", self.Net.String())
	if err != nil {
		return err
	}

	Init()

	connectCenter(centerAddr, self)
	selfPoint = &EndPoint{logic: self}

	// 集群内通信
	l.Listen(func(session net.Session) {
		util.Logger().Infoln("new client", session.RemoteAddr().String())
		// 超时时间
		session.SetTimeout(10*time.Second, 0)
		session.SetCodec(ss.NewCodec(protoss.SS_SPACE, protorpc.REQ_SPACE, protorpc.RESP_SPACE))
		session.SetCloseCallBack(func(reason string) {
			onClose(session, reason)
		})

		err := session.Start(func(data interface{}, err error) {
			if err != nil {
				session.Close(err.Error())
			} else {
				eventQueue.Push(func() {
					var err error
					switch data.(type) {
					case *ss.Message:
						//dispatchMsg(session, data.(*ss.Message))
					case *rpc.Request:
						//err = rpcServer.OnRPCRequest(&Node{session: session}, data.(*rpc.Request))
					case *rpc.Response:
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
	return nil
}

func Init() {
	eventQueue = queue.NewEventQueue(1024)
	eventQueue.Run(1)
	endPoints = map[uint32]*EndPoint{}
}
