package cluster

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/clugs/cluster/addr"
	"github.com/yddeng/clugs/logger"
	"github.com/yddeng/clugs/protocol/rpc/rpcpb"
	"github.com/yddeng/dnet/drpc"
	"testing"
	"time"
)

func init() {
	RegisterRPCMethod(proto.MessageName(&rpcpb.EchoReq{}), func(replier *drpc.Replier, req interface{}) {
		m := req.(*rpcpb.EchoReq)
		fmt.Println("rpc handler", m.GetMsg())

		replier.Reply(&rpcpb.EchoResp{Msg: "yes, I'm 1.1.1"}, nil)
	})
	logge := logger.New("log", "cluster_test")
	logger.InitLogger(logge)
}

func launch(centerAddr, logicAddr, localNetAddr string) {
	logic, err := addr.MakeAddr(logicAddr, localNetAddr)
	if err != nil {
		panic(err)
	}
	Launch(centerAddr, logic)
}

func TestCall(t *testing.T) {
	launch("127.0.0.1:9874", "1.1.1", "127.0.0.1:6547")
	// 自调用
	time.Sleep(time.Second)
	logic, _ := addr.MakeLogicAddr("1.1.1")
	AsyncCall(logic, &rpcpb.EchoReq{Msg: "I'm 1.1.1"}, func(i interface{}, e error) {
		fmt.Println(i, e)
	})

	time.Sleep(time.Second * 20)
}

func TestCall2(t *testing.T) {
	launch("127.0.0.1:9874", "1.1.2", "127.0.0.1:6548")
	//
	time.Sleep(time.Second)
	logic, _ := addr.MakeLogicAddr("1.1.1")
	AsyncCall(logic, &rpcpb.EchoReq{Msg: "I'm 1.1.2"}, func(i interface{}, e error) {
		fmt.Println(i, e)
	})
	time.Sleep(time.Second * 20)
}
