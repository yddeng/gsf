package cluster

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/yddeng/clugs/cluster/addr"
	"github.com/yddeng/clugs/cluster/clusterpb"
	"github.com/yddeng/clugs/codec/pb"
	"github.com/yddeng/clugs/codec/ss"
	"github.com/yddeng/clugs/logger"
	"github.com/yddeng/dnet"
	"github.com/yddeng/dnet/drpc"
	"reflect"
)

type clusterCenter struct {
	ends      map[addr.LogicAddr]*endpoint
	handlers  map[uint16]func(dnet.Session, *ss.Message)
	rpcServer *drpc.Server
}

var center *clusterCenter

func (this *clusterCenter) broadcast(msg proto.Message, excepts ...addr.LogicAddr) {
	except := addr.LogicAddr(0)
	if len(excepts) > 0 {
		except = excepts[0]
	}
	for logic, end := range this.ends {
		if logic != except {
			end.send(ss.NewMessage(msg))
		}
	}
}

func (this *clusterCenter) dispatch(session dnet.Session, message interface{}) {
	taskQueue.Push(func() {
		var err error
		switch message.(type) {
		case *ss.Message:
			ssMsg := message.(*ss.Message)
			if handler, ok := this.handlers[ssMsg.GetCmd()]; ok {
				handler(session, ssMsg)
			}

		case *drpc.Request:
			err = this.rpcServer.OnRPCRequest(&endpoint{session: session}, message.(*drpc.Request))

		//case *drpc.Response:

		default:
			err = fmt.Errorf("invalid type:%s", reflect.TypeOf(message).String())
		}
		if err != nil {
			logger.Errorf("cluster.center:dispatch error: %s. \n", err.Error())
		}
	})
}

func (this *clusterCenter) onHeartbeat(session dnet.Session, msg *ss.Message) {
	session.Send(ss.NewMessage(msg.GetData()))
}

func (this *clusterCenter) onLogin(replyer *drpc.Replier, arg interface{}) {

	req := arg.(*clusterpb.LoginReq)
	nodeInfo := req.GetNode()
	session := replyer.Channel.(*endpoint).session

	resp := &clusterpb.LoginResp{}
	logger.Infof("cluster.center:onLogin %v", req)

	logicAddr, err := addr.MakeAddr(addr.LogicAddr(nodeInfo.GetLogicAddr()).String(), nodeInfo.GetNetAddr())
	if err != nil {
		logger.Errorf("cluster.center:onLogin error :%s. ", err.Error())
		resp.Msg = err.Error()
		_ = replyer.Reply(resp)
		session.Close(err)
		return
	}

	end, ok := this.ends[logicAddr.Logic]
	if !ok {
		end = &endpoint{
			logic:   logicAddr,
			session: session,
		}
		this.ends[logicAddr.Logic] = end
		session.SetContext(end)
		logger.Infof("cluster.center:onLogin add endpoint [%s:%s] \n", logicAddr.Logic.String(), logicAddr.NetString())

		resp.Ok = true
		_ = replyer.Reply(resp)

	} else {
		// 已经有节点在该逻辑地址上启动。
		// 可能出现情况：该逻辑地址已被占用，但新节点上线时原有节点网络闪断，导致这条请求合法。
		if end.session != nil {
			resp.Msg = fmt.Sprintf("logicAddr %s is already register,address %s. \n", end.logic.Logic.String(), end.logic.NetString())
			_ = replyer.Reply(resp)
			logger.Infof("cluster.center:onLogin %s. \n", resp.GetMsg())
			session.Close(errors.New(resp.GetMsg()))
			return
		}
		logger.Infof("cluster.center:onLogin reconnect endpoint %s, net address from %s to %s. \n", logicAddr.Logic.String(), end.logic.NetString(), logicAddr.NetString())

		end.session = session
		end.logic.Net = logicAddr.Net
		session.SetContext(end)

		resp.Ok = true
		_ = replyer.Reply(resp)
	}

	enter := &clusterpb.NodeEnter{
		Node: &clusterpb.NodeInfo{
			LogicAddr: logicAddr.Logic.Uint32(),
			NetAddr:   logicAddr.NetString(),
		},
	}
	// 通知所有节点，新节点上线,除了自己
	this.broadcast(enter, logicAddr.Logic)

	notify := &clusterpb.NotifyNodeInfo{
		Nodes: make([]*clusterpb.NodeInfo, 0, len(this.ends)),
	}
	for _, e := range this.ends {
		notify.Nodes = append(notify.Nodes, &clusterpb.NodeInfo{
			LogicAddr: e.logic.Logic.Uint32(),
			NetAddr:   e.logic.NetString(),
		})
	}
	// 通知自己，有哪些节点在线,包括自己
	_ = end.send(ss.NewMessage(notify))

}

func (this *clusterCenter) onClose(session dnet.Session, err error) {
	logger.Infof("cluster.center:onClose %s. \n", err.Error())
	ctx := session.Context()
	if ctx != nil {
		end := ctx.(*endpoint)
		end.session = nil
		session.SetContext(nil)

		logger.Infof("cluster.center:onClose endpoint %s onClose %s. \n", end.logic.Logic.String(), err.Error())
		delete(this.ends, end.logic.Logic)

		// 通知所有节点，节点离线
		this.broadcast(&clusterpb.NodeLeave{
			LogicAddr: end.logic.Logic.Uint32(),
		})
	}
}

func (this *clusterCenter) Stop() {

}

func LunchCenter(netAddr string) *clusterCenter {
	center = &clusterCenter{
		ends:      map[addr.LogicAddr]*endpoint{},
		handlers:  map[uint16]func(dnet.Session, *ss.Message){},
		rpcServer: drpc.NewServer(),
	}

	center.rpcServer.Register(pb.GetNameById(clusterpb.REQ_SPACE, clusterpb.LoginReqCmd), center.onLogin)

	logger.Infof("cluster.center:LunchCenter serveTCP %s. \n", netAddr)

	go func() {
		if err := dnet.ServeTCP(netAddr, dnet.HandleFunc(func(conn dnet.NetConn) {
			logger.Infof("cluster.center:LunchCenter remote address %s. \n", conn.RemoteAddr().String())

			dnet.NewTCPSession(conn,
				dnet.WithTimeout(heartbeatTime, 0),
				dnet.WithCodec(ss.NewCodec(clusterpb.SS_SPACE, clusterpb.REQ_SPACE, clusterpb.RESP_SPACE)),
				dnet.WithCloseCallback(center.onClose),
				dnet.WithErrorCallback(func(session dnet.Session, err error) {
					logger.Errorf("cluster.center:LunchCenter session error:%s. \n", err.Error())
					session.Close(err)
				}),
				dnet.WithMessageCallback(center.dispatch),
			)

		})); err != nil {
			panic("cluster.center:LunchCenter " + err.Error())
		}
	}()

	return center
}
