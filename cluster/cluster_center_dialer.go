package cluster

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/clugs/cluster/addr"
	"github.com/yddeng/clugs/cluster/clusterpb"
	"github.com/yddeng/clugs/codec/ss"
	"github.com/yddeng/clugs/logger"
	"github.com/yddeng/clugs/util"
	"github.com/yddeng/dnet"
	"github.com/yddeng/dnet/drpc"
	"reflect"
	"time"
)

type clusterCenterDialer struct {
	address string
	dialing bool
	session dnet.Session

	heartbeatTicker *time.Ticker
	heartbeat       *ss.Message
}

func dialCenter(centerAddr string) *clusterCenterDialer {
	dialer := &clusterCenterDialer{
		address:   centerAddr,
		dialing:   false,
		heartbeat: ss.NewMessage(&clusterpb.Heartbeat{}),
	}

	taskQueue.Push(func() { dialer.dial() })
	return dialer
}

func (this *clusterCenterDialer) dial() {
	if this.dialing {
		return
	}

	this.dialing = true

	go func() {
		for {
			if conn, err := dnet.DialTCP(this.address, time.Second*5); nil == err && conn != nil {
				this.onConnected(conn)
				return
			} else {
				time.Sleep(1 * time.Second)
			}
		}
	}()
}

func (this *clusterCenterDialer) onConnected(conn dnet.NetConn) {
	taskQueue.Push(func() {
		this.dialing = false

		this.session = dnet.NewTCPSession(conn,
			dnet.WithCodec(ss.NewCodec(clusterpb.SS_SPACE, clusterpb.REQ_SPACE, clusterpb.RESP_SPACE)),
			dnet.WithCloseCallback(func(session dnet.Session, reason error) {
				taskQueue.Push(func() {
					if this.heartbeatTicker != nil {
						this.heartbeatTicker.Stop()
						this.heartbeatTicker = nil
					}
					this.session = nil
					logger.Infof("onConnected session closed, reason: %s\n", reason)
					this.dial()
				})
			}),
			dnet.WithErrorCallback(func(session dnet.Session, err error) {
				logger.Error("onConnected session error:", err)
				session.Close(err)
			}),
			dnet.WithMessageCallback(func(session dnet.Session, message interface{}) {
				taskQueue.Push(func() {
					var err error
					switch message.(type) {
					case *ss.Message:
						err = this.dispatchMsg(session, message.(*ss.Message))
					//case *drpc.Request:
					case *drpc.Response:
						err = rpcMgr.rpcClient.OnRPCResponse(message.(*drpc.Response))
					default:
						err = fmt.Errorf("invalid type:%s", reflect.TypeOf(message).String())
					}
					if err != nil {
						logger.Errorf("onConnected dispatch error: %s. \n", err.Error())
					}
				})
			}),
		)

		// 注册身份信息
		req := &clusterpb.LoginReq{
			Node: &clusterpb.NodeInfo{
				LogicAddr: LocalAddr.Logic.Uint32(),
				NetAddr:   LocalAddr.NetString(),
			},
		}
		if err := rpcMgr.rpcClient.Go(this, proto.MessageName(req), req, rpcTimeout, func(i interface{}, e error) {
			if e != nil || !i.(*clusterpb.LoginResp).GetOk() {
				var msg string
				if e != nil {
					msg = fmt.Sprintf("onConnected loginResp failed, error %s", e.Error())
				} else {
					msg = fmt.Sprintf("onConnected loginResp false, msg %s", i.(*clusterpb.LoginResp).GetMsg())
				}
				logger.Error(msg)
				panic(msg)
				return
			}

			logger.Info("onConnected login center ok")
			// 在center上注册成功，心跳
			this.heartbeatTicker = util.LoopTask(time.Second, func() {
				_ = this.send(this.heartbeat)
			})
		}); err != nil {

		}
	})
}

func (this *clusterCenterDialer) send(msg interface{}) error {
	if this.session == nil {
		return fmt.Errorf("session is nil")
	}
	return this.session.Send(msg)
}

func (this *clusterCenterDialer) SendRequest(req *drpc.Request) error {
	return this.send(req)
}

func (this *clusterCenterDialer) SendResponse(resp *drpc.Response) error {
	return this.send(resp)
}

func (this *clusterCenterDialer) dispatchMsg(session dnet.Session, msg *ss.Message) error {
	cmd := msg.GetCmd()
	switch cmd {
	case clusterpb.NotifyNodeInfoCmd:
		this.onNotifyNodeInfo(session, msg)
	case clusterpb.NodeLeaveCmd:
		this.onNodeLeave(session, msg)
	case clusterpb.NodeEnterCmd:
		this.onNodeEnter(session, msg)
	case clusterpb.HeartbeatCmd:

	default:
		return fmt.Errorf("dispatchMsg invailed cmd %d in nameSpace %s", cmd, clusterpb.SS_SPACE)
	}
	return nil
}

// 通知自己有哪些在线(新增节点，删除节点)
func (this *clusterCenterDialer) onNotifyNodeInfo(session dnet.Session, msg *ss.Message) {
	req := msg.GetData().(*clusterpb.NotifyNodeInfo)
	logger.Infof("onNotifyNodeInfo %v", req)

	existNode := map[addr.LogicAddr]struct{}{} // center 上现存的节点

	// 添加或修改 在 center 上存在的节点
	for _, v := range req.GetNodes() {
		logicAddr, err := addr.MakeAddr(addr.LogicAddr(v.GetLogicAddr()).String(), v.GetNetAddr())
		if err != nil {
			logger.Errorf("onNotifyNodeInfo error :%s. ", err.Error())
			continue
		}

		existNode[logicAddr.Logic] = struct{}{}
		end := endGroup.getEndpoint(logicAddr.Logic)
		if end != nil {
			// 已存在节点，继续验证地址
			end.Lock()
			// 新节点上来，替换原有连接
			if end.logic.NetString() != logicAddr.NetString() {
				if end.session != nil {
					// todo 直接断开吗
					end.session.Close(fmt.Errorf("logicAddr %s replace neatAddr %s to %s", end.logic.Logic.String(), end.logic.NetString(), logicAddr.NetString()))
				}
				end.logic.Net = logicAddr.Net
			}
			end.Unlock()
		} else {
			// 不存在节点，新增
			endGroup.addEndpoint(logicAddr)
		}
	}

	// 移除本地在 center 上不存在的节点
	needRmNode := map[addr.LogicAddr]struct{}{}
	endGroup.each(func(end *endpoint) bool {
		logic := end.logic.Logic
		if _, ok := existNode[logic]; !ok {
			needRmNode[logic] = struct{}{}
		}
		return true
	})
	for logicAddr := range needRmNode {
		endGroup.delEndpoint(logicAddr)
	}
}

func (this *clusterCenterDialer) onNodeLeave(session dnet.Session, msg *ss.Message) {
	req := msg.GetData().(*clusterpb.NodeLeave)
	logger.Infof("onNodeLeave %v", req)
	endGroup.delEndpoint(addr.LogicAddr(req.GetLogicAddr()))
}

func (this *clusterCenterDialer) onNodeEnter(session dnet.Session, msg *ss.Message) {
	req := msg.GetData().(*clusterpb.NodeEnter)
	logger.Infof("onNodeEnter %v", req)

	logicAddr, err := addr.MakeAddr(addr.LogicAddr(req.GetNode().GetLogicAddr()).String(), req.GetNode().GetNetAddr())
	if err != nil {
		logger.Errorf("onNodeEnter error :%s. ", err.Error())
		return
	}

	endGroup.addEndpoint(logicAddr)
}
