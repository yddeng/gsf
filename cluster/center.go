package cluster

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gsf/center/protocol"
	"github.com/yddeng/gsf/cluster/addr"
	"github.com/yddeng/gsf/codec/ss"
	"github.com/yddeng/gsf/util"
	dnet "github.com/yddeng/gsf/util/net"
	"github.com/yddeng/gsf/util/rpc"
	"net"
	"time"
)

type centerPoint struct {
	centerAddr string
	self       *addr.Addr
	dialing    bool
	session    dnet.Session
	handler    map[uint16]func(dnet.Session, *ss.Message)
	rpcClient  *rpc.Client

	heartbeatTicker *time.Ticker
	heartbeat       *ss.Message
}

func connectCenter(centerAddr string, self *addr.Addr) {
	centerP = &centerPoint{
		centerAddr: centerAddr,
		self:       self,
		dialing:    false,
		handler: map[uint16]func(dnet.Session, *ss.Message){
			protocol.NotifyNodeInfoCmd: onNotifyNodeInfo,
			protocol.NodeLeaveCmd:      onNodeLeave,
			protocol.NodeEnterCmd:      onNodeEnter,
		},
		rpcClient: rpc.NewClient(),
		heartbeat: ss.NewMessage(&protocol.Heartbeat{}),
	}

	eventQueue.Push(func() { centerP.dial() })
}

func (this *centerPoint) dial() {
	if this.dialing {
		return
	}

	this.dialing = true

	go func() {
		for {
			session, err := dnet.DialTCP("tcp", this.centerAddr, time.Second*5)
			if nil == err && session != nil {
				this.onConnected(session)
				return
			} else {
				//util.Logger().Errorf("dial center %s error:%s \n", this.centerAddr, err)
				time.Sleep(1 * time.Second)
			}
		}
	}()
}

func (this *centerPoint) onConnected(session dnet.Session) {
	eventQueue.Push(func() {
		this.dialing = false
		this.session = session

		session.SetCodec(ss.NewCodec(protocol.SS_SPACE, protocol.REQ_SPACE, protocol.RESP_SPACE))
		session.SetCloseCallBack(func(reason string) {
			eventQueue.Push(func() {
				if this.heartbeatTicker != nil {
					this.heartbeatTicker.Stop()
					this.heartbeatTicker = nil
				}
				this.session = nil
				util.Logger().Infof("center session closed, reason: %s\n", reason)
				this.dial()
			})
		})

		session.Start(func(data interface{}, err error) {
			if err != nil {
				session.Close(err.Error())
			} else {
				eventQueue.Push(func() {
					var err error
					switch data.(type) {
					case *ss.Message:
						err = this.dispatchMsg(session, data.(*ss.Message))
					case *rpc.Request:
					case *rpc.Response:
						err = this.rpcClient.OnRPCResponse(data.(*rpc.Response))
					}
					if err != nil {
						util.Logger().Errorf(err.Error())
					}
				})
			}
		})

		// 注册身份信息
		req := &protocol.LoginReq{
			Node: &protocol.NodeInfo{
				LogicAddr: uint32(this.self.Logic),
				NetAddr:   this.self.NetString(),
			},
		}
		err := this.rpcClient.AsynCall(this, proto.MessageName(req), req, rpcTimeout, func(i interface{}, e error) {
			if e != nil {
				msg := fmt.Sprintf("loginResp failed, e %s", e.Error())
				util.Logger().Errorf(msg)
				panic(msg)
				return
			}
			resp := i.(*protocol.LoginResp)
			if !resp.GetOk() {
				msg := fmt.Sprintf("loginResp failed, msg %s", resp.GetMsg())
				util.Logger().Errorf(msg)
				panic(msg)
				return
			}
			util.Logger().Infoln("login center ok")
			// 在center上注册成功，心跳
			this.heartbeatTicker = util.StartLoopTask(time.Second, func() {
				_ = this.send(this.heartbeat)
			})
		})
		if err != nil {
			session.Close(err.Error())
		}
	})
}

func (this *centerPoint) send(msg interface{}) error {
	if this.session == nil {
		return fmt.Errorf("session is nil")
	}
	return this.session.Send(msg)
}

func (this *centerPoint) SendRequest(req *rpc.Request) error {
	return this.send(req)
}

func (this *centerPoint) SendResponse(resp *rpc.Response) error {
	return this.send(resp)
}

func (this *centerPoint) dispatchMsg(session dnet.Session, msg *ss.Message) error {
	cmd := msg.GetCmd()
	if h, ok := this.handler[cmd]; ok {
		h(session, msg)
		return nil
	}
	return fmt.Errorf("dispatchMsg invailed cmd %d in nameSpace %s", cmd, protocol.SS_SPACE)
}

// 通知自己有哪些在线
func onNotifyNodeInfo(session dnet.Session, msg *ss.Message) {
	req := msg.GetData().(*protocol.NotifyNodeInfo)
	util.Logger().Infof("onNotifyNodeInfo %v", req)
	for _, v := range req.GetNodes() {
		logicAddr := addr.LogicAddr(v.GetLogicAddr())
		netAddr, err := net.ResolveTCPAddr("tcp", v.GetNetAddr())
		if err != nil {
			util.Logger().Errorf("endpoint %s netAddr err: %s", logicAddr.String(), err)
			continue
		}
		end := endpoints.getEndpointByLogic(logicAddr)
		if end != nil {
			end.Lock()
			// 新节点上来，替换原有连接
			if end.logic.NetString() != netAddr.String() {
				if end.session != nil {
					end.session.Close(fmt.Sprintf("logicAddr %s replace neatAddr %s -> %s", end.logic.Logic.String(), end.logic.NetString(), netAddr.String()))
				}
				end.logic.Net = netAddr
			}
			end.Unlock()
		} else {
			endpoints.addEndpoint(&addr.Addr{
				Logic: logicAddr,
				Net:   netAddr,
			})
		}
	}

}

func onNodeLeave(session dnet.Session, msg *ss.Message) {
	req := msg.GetData().(*protocol.NodeLeave)
	util.Logger().Infof("onNodeLeave %v", req)
	endpoints.removeEndpoint(addr.LogicAddr(req.GetLogicAddr()))
}

func onNodeEnter(session dnet.Session, msg *ss.Message) {
	req := msg.GetData().(*protocol.NodeEnter)
	util.Logger().Infof("onNodeEnter %v", req)
	logicAddr := addr.LogicAddr(req.GetNode().GetLogicAddr())
	netAddr, err := net.ResolveTCPAddr("tcp", req.GetNode().GetNetAddr())
	if err != nil {
		util.Logger().Errorf("endpoint %s netAddr err: %s", logicAddr.String(), err)
		return
	}
	endpoints.addEndpoint(&addr.Addr{
		Logic: logicAddr,
		Net:   netAddr,
	})
}
