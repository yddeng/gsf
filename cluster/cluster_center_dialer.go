package cluster

import (
	"fmt"
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
	clu     *cluster
	address string
	dialing bool
	session dnet.Session

	heartbeatTicker *time.Ticker
	heartbeat       *ss.Message
}

func dialCenter(centerAddr string, clu *cluster) *clusterCenterDialer {
	dialer := &clusterCenterDialer{
		clu:       clu,
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
					logger.Infof("cluster.center_dialer:onConnected session closed, reason: %s\n", reason)
					this.dial()
				})
			}),
			dnet.WithErrorCallback(func(session dnet.Session, err error) {
				logger.Error("cluster.center_dialer:onConnected session error:", err)
				session.Close(err)
			}),
			dnet.WithMessageCallback(func(session dnet.Session, message interface{}) {
				var err error
				switch message.(type) {
				case *ss.Message:
					err = this.dispatchMsg(session, message.(*ss.Message))
				//case *drpc.Request:
				case *drpc.Response:
					err = this.rpcClient.OnRPCResponse(message.(*drpc.Response))
				default:
					err = fmt.Errorf("invalid type:%s", reflect.TypeOf(message).String())
				}
				if err != nil {
					logger.Errorf("cluster.center_dialer:onConnected dispatch error: %s. \n", err.Error())
				}
			}),
		)

		// 注册身份信息
		req := &clusterpb.LoginReq{
			Node: &clusterpb.NodeInfo{
				LogicAddr: LocalAddr.Logic.Uint32(),
				NetAddr:   LocalAddr.NetString(),
			},
		}
		if err := RPCGo(this, req, func(i interface{}, e error) {
			if e != nil {
				msg := fmt.Sprintf("cluster.center_dialer:onConnected loginResp failed, e %s", e.Error())
				logger.Error(msg)
				panic(msg)
				return
			}
			resp := i.(*clusterpb.LoginResp)
			if !resp.GetOk() {
				msg := fmt.Sprintf("cluster.center_dialer:onConnected loginResp failed, msg %s", resp.GetMsg())
				logger.Error(msg)
				panic(msg)
				return
			}
			logger.Info("cluster.center_dialer:onConnected login center ok")
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
