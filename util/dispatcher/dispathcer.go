package dispatcher

import (
	"github.com/golang/protobuf/proto"
	"github.com/yddeng/gsf/util/net"
)

/*
 事件分发器
 每一条协议注册一个处理函数
*/

type Handler func(net.Session, proto.Message)

type Dispatcher struct {
	handlers map[uint16]Handler
}

func NewDispatch() *Dispatcher {
	return &Dispatcher{
		handlers: map[uint16]Handler{},
	}
}

func (this *Dispatcher) Registerk(cmd uint16, callback Handler) {
	_, ok := this.handlers[cmd]
	if ok {
		return
	}

	this.handlers[cmd] = callback
}

func (this *Dispatcher) Dispatch(session net.Session, cmd uint16, msg proto.Message) {
	if nil != msg {
		handler, ok := this.handlers[cmd]
		if ok {
			handler(session, msg)
		}
	}
}
