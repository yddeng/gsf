package center

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
)

type Node struct {
	LogicAddr *addr.Addr
	session   dnet.Session
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

// todo 网络闪断
func onClose(session dnet.Session, reason string) {
	util.Logger().Infoln("onClose", reason)
	ctx := session.Context()
	if ctx != nil {
		n := ctx.(*Node)
		n.session = nil
		session.SetContext(nil)

		util.Logger().Infoln(n.LogicAddr.Logic.String(), "onClose", reason)
		tId := n.LogicAddr.Logic.Uint32()
		delete(nodes, tId)

		// 通知所有节点，节点离线
		broadcast(&protocol.NodeLeave{
			LogicAddr: tId,
		})
	}
}

func onLogin(replyer *rpc.Replyer, arg interface{}) {
	req := arg.(*protocol.LoginReq)
	logic := req.GetNode().GetLogicAddr()
	netStr := req.GetNode().GetNetAddr()
	session := replyer.Channel.(*Node).session
	resp := &protocol.LoginResp{}
	util.Logger().Infof("onLogin %v", req)

	netAddr, err := net.ResolveTCPAddr("tcp", netStr)
	if err != nil {
		resp.Msg = err.Error()
		_ = replyer.Reply(resp, nil)
		util.Logger().Errorf(err.Error())
		session.Close(err.Error())
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
		_ = replyer.Reply(resp, nil)
		util.Logger().Infof("add node %d \n", n.LogicAddr.Logic)

	} else {
		// 已经有节点在该逻辑地址上启动。
		// 可能出现情况：该逻辑地址已被占用，但新节点上线时原有节点网络闪断，导致这条请求合法。
		if n.session != nil {
			resp.Msg = fmt.Sprintf("logicAddr %s netAddr %s session is already connect\n", n.LogicAddr.Logic.String(), n.LogicAddr.NetString())
			_ = replyer.Reply(resp, nil)
			util.Logger().Infof(resp.GetMsg())
			session.Close(resp.GetMsg())
			return
		}

		n.session = session
		n.LogicAddr.Net = netAddr
		session.SetContext(n)

		util.Logger().Infof("reconnect node %d \n", n.LogicAddr.Logic)
		resp.Ok = true
		_ = replyer.Reply(resp, nil)

	}

	enter := &protocol.NodeEnter{
		Node: &protocol.NodeInfo{
			LogicAddr: logic,
			NetAddr:   netAddr.String(),
		},
	}
	// 通知所有节点，新节点上线,除了自己
	broadcast(enter, n.LogicAddr.Logic)

	notify := &protocol.NotifyNodeInfo{
		Nodes: make([]*protocol.NodeInfo, 0, len(nodes)),
	}
	for _, node := range nodes {
		notify.Nodes = append(notify.Nodes, &protocol.NodeInfo{
			LogicAddr: uint32(node.LogicAddr.Logic),
			NetAddr:   node.LogicAddr.NetString(),
		})
	}
	// 通知自己，有哪些节点在线,包括自己
	_ = n.Send(ss.NewMessage(notify))
}

func onHeartbeat(session dnet.Session, msg *ss.Message) {}
