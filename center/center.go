package center

import (
	"fmt"
	"github.com/yddeng/dnet"
	"github.com/yddeng/dnet/drpc"
	"github.com/yddeng/dnet/dtcp"
	"github.com/yddeng/dutil/queue"
	"github.com/yddeng/gsf/center/protocol"
	"github.com/yddeng/gsf/codec/pb"
	"github.com/yddeng/gsf/codec/ss"
	"github.com/yddeng/gsf/util"
	"reflect"
	"time"
)

type Handler func(dnet.Session, *ss.Message)

var (
	msgHandlers      map[uint16]Handler
	eventQueue       *queue.EventQueue
	rpcServer        *drpc.Server
	nodes            map[uint32]*Node
	heartbeatTimeout = time.Second * 10
)

func registerHandler(cmd uint16, callback Handler) {
	_, ok := msgHandlers[cmd]
	if ok {
		return
	}

	msgHandlers[cmd] = callback
}

func dispatchMsg(session dnet.Session, msg *ss.Message) {
	if nil != msg {
		cmd := msg.GetCmd()
		handler, ok := msgHandlers[cmd]
		if ok {
			handler(session, msg)
		}
	}
}

func Launch(netAddr string) {
	l := util.Must(dtcp.NewTCPListener("tcp", netAddr)).(*dtcp.TCPListener)

	Init()

	util.Must(nil,
		l.Listen(func(session dnet.Session) {
			util.Logger().Infoln("new client", session.RemoteAddr().String())
			// 超时时间
			session.SetTimeout(heartbeatTimeout, 0)
			session.SetCodec(ss.NewCodec(protocol.SS_SPACE, protocol.REQ_SPACE, protocol.RESP_SPACE))
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
		}))

}

func Init() {
	msgHandlers = map[uint16]Handler{}
	nodes = map[uint32]*Node{}
	rpcServer = drpc.NewServer()
	eventQueue = queue.NewEventQueue(1024)
	eventQueue.Run(1)

	// ss
	registerHandler(protocol.HeartbeatCmd, onHeartbeat)

	// rpc
	rpcServer.Register(pb.GetNameById(protocol.REQ_SPACE, protocol.LoginReqCmd), onLogin)
}
