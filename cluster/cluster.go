package cluster

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/clugs/cluster/addr"
	"github.com/yddeng/dutil/task"
	"net"
	"sync"
	"time"
)

var (
	taskQueue = task.NewTaskQueue(1024)

	heartbeatTime = time.Second * 30 // 集群节点间超时时间间隔
	rpcTimeout    = time.Second * 8  // rpc 请求超时时间间隔

	LocalAddr *addr.Addr
	clu *cluster = &cluster{
		ssHandler:  map[uint16]func(from addr.LogicAddr, msg proto.Message),
		Mutex: new(sync.Mutex),
	}
)

func init() {
	taskQueue.Run()
}


type cluster struct {
	ssHandler  map[uint16]func(from addr.LogicAddr, msg proto.Message)
	endGroup *endpointGroup
	centerDialer *clusterCenterDialer
	*sync.Mutex
}

func Launch(centerAddr string, localAddr *addr.Addr)  {
	l, err := net.ListenTCP("tcp", localAddr.Net)
	if err != nil {
		panic(err)
	}

	LocalAddr = localAddr
	clu = &cluster{
		endGroup: &endpointGroup{
			logic2End: map[addr.LogicAddr]*endpoint{},
			type2End:  map[uint32]*endpoint{},
			Mutex:     new(sync.Mutex),
		},
		Mutex: nil,
	}

	clu.centerDialer = dialCenter(centerAddr, clu)


	go func() {
		for {
			conn, err := l.AcceptTCP()
			if err != nil {
				if ne, ok := err.(net.Error); ok && ne.Temporary() {
					continue
				} else {
					util.Logger().Errorln(err)
					return
				}
			}

			// 新连接验证
			go accept(conn)
		}
	}()

}

func Post(logic addr.LogicAddr, msg proto.Message) error {
	end := endpoints.getEndpointByLogic(logic)
	if end == nil {
		util.Logger().Errorf("%s is not found", logic.String())
		return fmt.Errorf("%s is not found", logic.String())
	}

	end.Lock()
	defer end.Unlock()
	return end.send(ss.NewMessage(msg))
}

func RegisterSSMethod(cmd uint16, h func(from addr.LogicAddr, msg proto.Message)) {
	_, ok := ssHandler[cmd]
	if ok {
		panic(fmt.Sprintf("register ss method cmd %d already registed", cmd))
	}
	ssHandler[cmd] = h
}

func dispatchSS(from addr.LogicAddr, msg *ss.Message) error {
	cmd := msg.GetCmd()
	h, ok := ssHandler[cmd]
	if ok {
		h(from, msg.GetData())
		return nil
	}
	return fmt.Errorf("dispatchSS invailed cmd %d ", cmd)
}
