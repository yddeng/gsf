package cluster

import (
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gsf/cluster/addr"
	"github.com/yddeng/gsf/util/queue"
	"github.com/yddeng/gsf/util/rpc"
	"sync"
	"time"
)

var (
	eventQueue = queue.NewEventQueue(1024)

	selfPoint *endpoint
	centerP   *centerPoint

	ssHandler = map[uint16]func(from addr.LogicAddr, msg proto.Message){}
	endpoints = &endpointGroup{
		logic2End: map[addr.LogicAddr]*endpoint{},
		type2End:  map[uint32]*endpoint{},
		Mutex:     new(sync.Mutex),
	}
	rpcMgr = &rpcManager{
		rpcServer: rpc.NewServer(),
		rpcClient: rpc.NewClient(),
	}

	heartbeatTime = time.Second * 30 // 集群节点间超时时间间隔
	rpcTimeout    = time.Second * 8  // rpc 请求超时时间间隔
)

func init() {
	eventQueue.Run(1)
}
