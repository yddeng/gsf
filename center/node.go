package center

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gsf/center/protocol"
	"github.com/yddeng/gsf/cluster/addr"
	"github.com/yddeng/gsf/codec/ss"
	"github.com/yddeng/gsf/util"
	"github.com/yddeng/gsf/util/net"
	"github.com/yddeng/gsf/util/rpc"
	net2 "net"
)

type Node struct {
	LogicAddr *addr.Addr
	session   net.Session
}

func broadcast(msg proto.Message, except ...addr.LogicAddr) {
	excep := addr.LogicAddr(0)
	if len(except) > 0 {
		excep = except[0]
	}
	for _, n := range nodes {
		if n.LogicAddr.Logic != excep {
			n.Send(ss.NewMessage(msg))
		}
	}
}

// 传入 *ss.Message, *rpc.Request, *rpc.Response
func (this *Node) Send(o interface{}) error {
	if this.session == nil {
		return fmt.Errorf("%s session is nil", this.LogicAddr.Logic.String())
	}
	return this.session.Send(o)
}

func (this *Node) SendRequest(req *rpc.Request) error {
	return this.Send(req)
}

func (this *Node) SendResponse(resp *rpc.Response) error {
	return this.Send(resp)
}

func onClose(session net.Session, reason string) {
	util.Logger().Infoln("onClose", reason)
	ctx := session.Context()
	if ctx != nil {
		n := ctx.(*Node)
		n.session = nil
		session.SetContext(nil)
	}
}

func onLogin(replyer *rpc.Replyer, arg interface{}) {
	req := arg.(*protocol.LoginReq)
	logic := req.GetNode().GetLogicAddr()
	netStr := req.GetNode().GetNetAddr()
	session := replyer.Channel.(*Node).session
	resp := &protocol.LoginResp{}
	util.Logger().Infof("onLogin %v", req)

	netAddr, err := net2.ResolveTCPAddr("tcp", netStr)
	if err != nil {
		util.Logger().Errorf(err.Error())
		resp.Msg = err.Error()
		_ = replyer.Reply(resp, nil)
		return
	}

	n, ok := nodes[logic]
	if !ok {
		n = &Node{
			LogicAddr: &addr.Addr{
				Logic: addr.LogicAddr(logic),
				Net:   netAddr,
			},
			session: session,
		}
		nodes[logic] = n
		session.SetContext(n)

		resp.Ok = true
		util.Logger().Infof("add node %d \n", n.LogicAddr.Logic)
		_ = replyer.Reply(resp, nil)

		notify := &protocol.NotifyNodeInfo{
			Nodes: []*protocol.NodeInfo{{
				LogicAddr: logic,
				NetAddr:   netAddr.String(),
			}},
		}
		// 通知所有节点，新节点上线,除了自己
		broadcast(notify, n.LogicAddr.Logic)

		for _, node := range nodes {
			if node.LogicAddr.Logic != n.LogicAddr.Logic {
				notify.Nodes = append(notify.Nodes, &protocol.NodeInfo{
					LogicAddr: uint32(node.LogicAddr.Logic),
					NetAddr:   node.LogicAddr.Net.String(),
				})
			}
		}
		// 通知自己，有哪些节点在线
		_ = n.Send(ss.NewMessage(notify))

	} else {
		if n.session != nil {
			resp.Msg = "session is already connect"
			util.Logger().Infof("node %d %s\n", n.LogicAddr.Logic, resp.GetMsg())
			_ = replyer.Reply(resp, nil)
			return
		}
		n.session = session
		session.SetContext(n)

		util.Logger().Infof("reLogin node %d \n", n.LogicAddr.Logic)
		resp.Ok = true
		_ = replyer.Reply(resp, nil)

		// 换了新地址
		if n.LogicAddr.Net.String() != netAddr.String() {
			change := &protocol.NodeChange{Nodes: []*protocol.NodeInfo{
				{
					LogicAddr: logic,
					NetAddr:   netAddr.String(),
				},
			}}
			broadcast(change)
		}
	}
}

func onHeartbeat(session net.Session, msg *ss.Message) {}
