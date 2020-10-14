package cluster

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/dnet/drpc"
	"github.com/yddeng/gsf/cluster/addr"
	"github.com/yddeng/gsf/codec/pb"
	"github.com/yddeng/gsf/protocol/rpc"
	"github.com/yddeng/gsf/protocol/rpc/rpcpb"
	"github.com/yddeng/gsf/protocol/ss"
	"github.com/yddeng/gsf/protocol/ss/sspb"
	"github.com/yddeng/gsf/util"
	"testing"
	"time"
)

func TestLauncher(t *testing.T) {
	util.InitLogger("log", "test")
	logic, err := addr.MakeAddr("1.1.1", "127.0.0.1:6547")
	if err != nil {
		fmt.Println(1, err)
		return
	}

	err = Launch("127.0.0.1:9874", logic)
	if err != nil {
		fmt.Println(2, err)
		return
	}
	select {}
}

func TestAsynCall1(t *testing.T) {
	util.InitLogger("log", "asynCall1")
	logic1, _ := addr.MakeAddr("1.1.1", "127.0.0.1:6547")
	logic2, _ := addr.MakeAddr("1.1.2", "127.0.0.1:6548")

	RegisterRPCMethod(pb.GetNameById(rpc.REQ_SPACE, rpc.Echo), func(replyer *drpc.Replyer, req interface{}) {
		msg := req.(*rpcpb.EchoReq)
		fmt.Println(msg)

		resp := &rpcpb.EchoResp{
			Msg: "reply " + msg.GetMsg(),
		}
		replyer.Reply(resp, nil)
	})
	RegisterSSMethod(ss.Echo, func(from addr.LogicAddr, msg proto.Message) {
		req := msg.(*sspb.Echo)
		fmt.Println(req)
	})

	err := Launch("127.0.0.1:9874", logic1)
	if err != nil {
		fmt.Println(2, err)
		return
	}

	time.Sleep(time.Second)
	AsynCall(logic1.Logic, &rpcpb.EchoReq{Msg: "hello call1"}, func(i interface{}, e error) {
		if e != nil {
			fmt.Println(e)
			return
		}
		msg := i.(*rpcpb.EchoResp)
		fmt.Println(msg)
	})
	Post(logic1.Logic, &sspb.Echo{Msg: "hello post1"})

	// 超时断开连接，后重连
	time.Sleep(time.Second * 40)
	AsynCall(logic2.Logic, &rpcpb.EchoReq{Msg: "hello call1"}, func(i interface{}, e error) {
		if e != nil {
			fmt.Println(e)
			return
		}
		msg := i.(*rpcpb.EchoResp)
		fmt.Println(msg)
	})
	Post(logic2.Logic, &sspb.Echo{Msg: "hello post1"})
	select {}
}

func TestAsynCall2(t *testing.T) {
	util.InitLogger("log", "asynCall2")
	logic1, _ := addr.MakeAddr("1.1.1", "127.0.0.1:6547")
	logic2, _ := addr.MakeAddr("1.1.2", "127.0.0.1:6548")

	RegisterRPCMethod(pb.GetNameById(rpc.REQ_SPACE, rpc.Echo), func(replyer *drpc.Replyer, req interface{}) {
		msg := req.(*rpcpb.EchoReq)
		fmt.Println(msg)

		resp := &rpcpb.EchoResp{
			Msg: "reply " + msg.GetMsg(),
		}
		replyer.Reply(resp, nil)
	})
	RegisterSSMethod(ss.Echo, func(from addr.LogicAddr, msg proto.Message) {
		req := msg.(*sspb.Echo)
		fmt.Println(req)
	})

	err := Launch("127.0.0.1:9874", logic2)
	if err != nil {
		fmt.Println(2, err)
		return
	}

	time.Sleep(time.Second)
	AsynCall(logic1.Logic, &rpcpb.EchoReq{Msg: "hello call2"}, func(i interface{}, e error) {
		if e != nil {
			fmt.Println(e)
			return
		}
		msg := i.(*rpcpb.EchoResp)
		fmt.Println(msg)
	})
	Post(logic1.Logic, &sspb.Echo{Msg: "hello post2"})

	select {}
}
