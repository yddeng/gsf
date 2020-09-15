package cluster

import (
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/dnet/drpc"
	"github.com/yddeng/dutil/queue"
	"github.com/yddeng/gsf/cluster/addr"
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
		rpcServer: drpc.NewServer(),
		rpcClient: drpc.NewClient(),
	}

	heartbeatTime = time.Second * 30 // 集群节点间超时时间间隔
	rpcTimeout    = time.Second * 8  // rpc 请求超时时间间隔
)

func init() {
	eventQueue.Run(1)
}
