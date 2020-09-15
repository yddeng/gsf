package net

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

type WSConn struct {
	flag         byte
	conn         *websocket.Conn
	ctx          interface{}   //用户数据
	readTimeout  time.Duration // 读超时
	writeTimeout time.Duration // 写超时

	sendBufChan chan []byte //发送队列

	msgCallback   func(interface{}, error)             //消息回调
	closeCallback func(session Session, reason string) //关闭连接回调
	closeReason   string                               //关闭原因

	lock sync.Mutex
}

func NewWSConn(conn *websocket.Conn) *WSConn {
	return &WSConn{
		conn:        conn,
		sendBufChan: make(chan []byte, sendBufChanSize),
	}
}

//读写超时
func (this *WSConn) SetTimeout(readTimeout, writeTimeout time.Duration) {
	defer this.lock.Unlock()
	this.lock.Lock()

	this.readTimeout = readTimeout
	this.writeTimeout = writeTimeout
}

func (this *WSConn) LocalAddr() net.Addr {
	return this.conn.LocalAddr()
}

func (this *WSConn) NetConn() interface{} {
	return this.conn
}

//对端地址
func (this *WSConn) RemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}

func (this *WSConn) SetCodec(codec Codec) {}

func (this *WSConn) SetCloseCallBack(closeCallback func(session Session, reason string)) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.closeCallback = closeCallback
}

func (this *WSConn) SetContext(ctx interface{}) {
	this.lock.Lock()
	this.ctx = ctx
	this.lock.Unlock()
}

func (this *WSConn) Context() interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.ctx
}

//开启消息处理
func (this *WSConn) Start(msgCb func(interface{}, error)) error {
	if msgCb == nil {
		return ErrNoMsgCallBack
	}

	this.lock.Lock()
	if this.flag == started {
		return ErrStateFailed
	}
	this.flag = started
	this.msgCallback = msgCb
	this.lock.Unlock()

	go this.receiveThread()
	go this.sendThread()

	return nil
}

func (this *WSConn) isClose() bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.flag == closed
}

//接收线程
func (this *WSConn) receiveThread() {
	for {
		if this.isClose() {
			return
		}
		if this.readTimeout > 0 {
			_ = this.conn.SetReadDeadline(time.Now().Add(this.readTimeout))
		}
		_, msg, err := this.conn.ReadMessage()
		if this.isClose() {
			return
		}
		if err != nil {
			this.msgCallback(nil, err)
		} else {
			this.msgCallback(msg, err)
		}
	}
}

//发送线程
func (this *WSConn) sendThread() {
	defer this.close()
	for {
		data, isOpen := <-this.sendBufChan
		if !isOpen {
			break
		}
		if this.writeTimeout > 0 {
			_ = this.conn.SetWriteDeadline(time.Now().Add(this.writeTimeout))
		}

		err := this.conn.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			this.msgCallback(nil, err)
		}

	}
}

func (this *WSConn) Send(o interface{}) error {
	if o == nil {
		return ErrSendMsgNil
	}

	data, ok := o.([]byte)
	if !ok {
		return fmt.Errorf("interface {} is %s,need []byte or use SendMsg(data []byte)", reflect.TypeOf(o).String())
	}

	return this.SendBytes(data)
}

func (this *WSConn) SendBytes(data []byte) error {
	if len(data) == 0 {
		return ErrSendMsgNil
	}

	//非堵塞
	if len(this.sendBufChan) == sendBufChanSize {
		return ErrSendChanFull
	}

	this.lock.Lock()
	if this.flag == 0 {
		return ErrStateFailed
	}
	if this.flag == closed {
		return ErrStateFailed
	}
	this.lock.Unlock()

	this.sendBufChan <- data
	return nil
}

/*
 主动关闭连接
 先关闭读，待写发送完毕关闭写
*/
func (this *WSConn) Close(reason string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if (this.flag & closed) > 0 {
		return
	}

	close(this.sendBufChan)
	this.closeReason = reason
	this.flag = closed
}

func (this *WSConn) close() {
	_ = this.conn.Close()
	this.lock.Lock()
	callback := this.closeCallback
	msg := this.closeReason
	this.lock.Unlock()
	if callback != nil {
		callback(this, msg)
	}
}

type WSListener struct {
	listener *net.TCPListener
	upgrader *websocket.Upgrader
	origin   string
	started  int32
}

func NewWSListener(network, addr, origin string, upgrader ...*websocket.Upgrader) (*WSListener, error) {
	tcpAddr, err := net.ResolveTCPAddr(network, addr)
	if err != nil {
		return nil, err
	}
	listener, err := net.ListenTCP(tcpAddr.Network(), tcpAddr)
	if err != nil {
		return nil, err
	}

	l := &WSListener{
		listener: listener,
		origin:   origin,
	}

	if len(upgrader) > 0 {
		l.upgrader = upgrader[0]
	} else {
		l.upgrader = &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				// allow all connections by default
				return true
			},
		}
	}
	return l, nil
}

func (this *WSListener) Close() {
	if !atomic.CompareAndSwapInt32(&this.started, 1, 0) {
		this.listener.Close()
	}
}

func (this *WSListener) Listen(newClient func(Session)) error {

	if newClient == nil {
		return ErrNewClientNil
	}

	if !atomic.CompareAndSwapInt32(&this.started, 0, 1) {
		return ErrStateFailed
	}

	http.HandleFunc(this.origin, func(w http.ResponseWriter, r *http.Request) {
		c, err := this.upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("wssocket Upgrade failed:%s\n", err.Error())
			return
		}
		newClient(NewWSConn(c))
	})

	go func() {
		err := http.Serve(this.listener, nil)
		if err != nil {
			log.Printf("http.Serve() failed:%s\n", err.Error())
		}

		_ = this.listener.Close()
	}()

	return nil
}

func DialWS(addr, path string, timeout time.Duration) (Session, error) {
	u := url.URL{Scheme: "ws", Host: addr, Path: path}
	websocket.DefaultDialer.HandshakeTimeout = timeout
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}
	return NewWSConn(conn), nil
}
