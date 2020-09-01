package center

import (
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gsf/center/protocol"
	"github.com/yddeng/gsf/codec/ss"
	"github.com/yddeng/gsf/util"
	dnet "github.com/yddeng/gsf/util/net"
	"github.com/yddeng/gsf/util/queue"
	"github.com/yddeng/gsf/util/rpc"
	"time"
)

type Handler func(dnet.Session, *ss.Message)

var (
	msgHandlers      map[uint16]Handler
	eventQueue       *queue.EventQueue
	rpcServer        *rpc.Server
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

func Launcher(netAddr string) {
	l := util.Must(dnet.NewTCPListener("tcp", netAddr)).(*dnet.TCPListener)

	Init()

	util.Must(nil,
		l.Listen(func(session dnet.Session) {
			util.Logger().Infoln("new client", session.RemoteAddr().String())
			// 超时时间
			session.SetTimeout(heartbeatTimeout, 0)
			session.SetCodec(ss.NewCodec("center_ss", "center_req", "center_resp"))
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
							dispatchMsg(session, data.(*ss.Message))
						case *rpc.Request:
							err = rpcServer.OnRPCRequest(&Node{session: session}, data.(*rpc.Request))
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
		}))

}

func Init() {
	msgHandlers = map[uint16]Handler{}
	nodes = map[uint32]*Node{}
	rpcServer = rpc.NewServer()
	eventQueue = queue.NewEventQueue(1024)
	eventQueue.Run(1)

	// ss
	registerHandler(protocol.HeartbeatCmd, onHeartbeat)

	// rpc
	util.Must(nil, rpcServer.Register(proto.MessageName(&protocol.LoginReq{}), onLogin))
}
