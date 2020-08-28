package cluster

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gsf/center/protocol"
	"github.com/yddeng/gsf/cluster/addr"
	"github.com/yddeng/gsf/codec/ss"
	"github.com/yddeng/gsf/util"
	net2 "github.com/yddeng/gsf/util/net"
	"github.com/yddeng/gsf/util/rpc"
	"net"
	"time"
)

type CenterPoint struct {
	centerAddr string
	self       *addr.Addr
	dialing    bool
	session    net2.Session
	handler    map[uint16]func(net2.Session, *ss.Message)
	rpcClient  *rpc.Client

	heartbeatTicker *time.Ticker
	heartbeat       *ss.Message
}

func connectCenter(centerAddr string, self *addr.Addr) {
	centerPoint = &CenterPoint{
		centerAddr: centerAddr,
		self:       self,
		dialing:    false,
		handler: map[uint16]func(net2.Session, *ss.Message){
			protocol.NotifyNodeInfoCmd: onNotifyNodeInfo,
			protocol.NodeLeaveCmd:      onNodeLeave,
			protocol.NodeChangeCmd:     onNodeChange,
		},
		rpcClient: rpc.NewClient(),
		heartbeat: ss.NewMessage(&protocol.Heartbeat{}),
	}

	centerPoint.dial()
}

func (this *CenterPoint) dial() {
	if this.dialing {
		return
	}

	this.dialing = true

	go func() {
		for {
			session, err := net2.DialTCP("tcp", this.centerAddr, time.Second*5)
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

func (this *CenterPoint) onConnected(session net2.Session) {
	eventQueue.Push(func() {
		this.dialing = false
		this.session = session

		session.SetCodec(ss.NewCodec(protocol.SS_SPACE, protocol.REQ_SPACE, protocol.RESP_SPACE))
		session.SetCloseCallBack(func(reason string) {
			eventQueue.Push(func() {
				this.session = nil
				if this.heartbeatTicker != nil {
					this.heartbeatTicker.Stop()
					this.heartbeatTicker = nil
				}
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

		// 请求登陆
		req := &protocol.LoginReq{
			Node: &protocol.NodeInfo{
				LogicAddr: uint32(this.self.Logic),
				NetAddr:   this.self.NetString(),
			},
		}
		err := this.rpcClient.AsynCall(&RPCChannel{session: session}, proto.MessageName(req), req, func(i interface{}, e error) {
			if e != nil {
				msg := fmt.Sprintf("loginResp failed, e %s", e.Error())
				util.Logger().Errorf(msg)
				session.Close(msg)
				return
			}
			resp := i.(*protocol.LoginResp)
			if !resp.GetOk() {
				msg := fmt.Sprintf("loginResp failed, msg %s", resp.GetMsg())
				util.Logger().Errorf(msg)
				session.Close(msg)
				return
			}
			util.Logger().Infoln("login center ok")
			// 在center上注册成功，心跳
			this.heartbeatTicker = util.StartLoopTask(time.Second, func() {
				eventQueue.Push(func() { this.send(this.heartbeat) })
			})
		})
		if err != nil {
			session.Close(err.Error())
		}
	})
}

func (this *CenterPoint) send(msg interface{}) {
	if this.session == nil {
		return
	}
	this.session.Send(msg)
}

func (this *CenterPoint) dispatchMsg(session net2.Session, msg *ss.Message) error {
	cmd := msg.GetCmd()
	if h, ok := this.handler[cmd]; ok {
		h(session, msg)
		return nil
	}
	return fmt.Errorf("dispatchMsg invailed cmd %d in nameSpace %s", cmd, protocol.SS_SPACE)
}

// 新节点上线，通知自己有哪些在线
func onNotifyNodeInfo(session net2.Session, msg *ss.Message) {
	req := msg.GetData().(*protocol.NotifyNodeInfo)
	util.Logger().Infof("onNotifyNodeInfo %v", req)
	for _, v := range req.GetNodes() {
		logicAddr := v.GetLogicAddr()
		end, ok := endPoints[logicAddr]
		if ok {
			util.Logger().Debugf("endpoint %s is exist", addr.LogicAddr(logicAddr).String())
			continue
		}
		netAddr, err := net.ResolveTCPAddr("tcp", v.GetNetAddr())
		if err != nil {
			util.Logger().Errorf("endpoint %s netAddr err: %s", addr.LogicAddr(logicAddr).String(), err)
			continue
		}
		end = &EndPoint{logic: &addr.Addr{
			Logic: addr.LogicAddr(logicAddr),
			Net:   netAddr,
		}}
		endPoints[logicAddr] = end
	}
}

func onNodeLeave(session net2.Session, msg *ss.Message) {
	req := msg.GetData().(*protocol.NodeLeave)
	util.Logger().Infof("onNodeLeave %v", req)
}

func onNodeChange(session net2.Session, msg *ss.Message) {
	req := msg.GetData().(*protocol.NodeChange)
	util.Logger().Infof("onNodeChange %v", req)
}
