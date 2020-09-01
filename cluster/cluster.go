package cluster

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gsf/cluster/addr"
	"github.com/yddeng/gsf/codec/ss"
	"github.com/yddeng/gsf/util"
	"net"
)

func Launcher(centerAddr string, self *addr.Addr) error {
	l, err := net.ListenTCP("tcp", self.Net)
	if err != nil {
		return err
	}

	connectCenter(centerAddr, self)
	selfPoint = &endpoint{logic: self}

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

	return nil
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
