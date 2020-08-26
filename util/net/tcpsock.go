package net

import (
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const (
	started = 0x01 //0000 0001
	closed  = 0x02 //0000 0010
)

const sendBufChanSize = 1024

type TCPConn struct {
	flag         byte
	conn         *net.TCPConn
	ctx          interface{}   //用户数据
	readTimeout  time.Duration // 读超时
	writeTimeout time.Duration // 写超时

	codec       Codec       //编解码器
	sendBufChan chan []byte //发送队列

	msgCallback   func(interface{}, error) //消息回调
	closeCallback func(string)             //关闭连接回调
	closeReason   string                   //关闭原因

	lock sync.Mutex
}

func newTCPConn(conn *net.TCPConn) *TCPConn {
	return &TCPConn{
		conn:        conn,
		sendBufChan: make(chan []byte, sendBufChanSize),
	}
}

//读写超时
func (this *TCPConn) SetTimeout(readTimeout, writeTimeout time.Duration) {
	defer this.lock.Unlock()
	this.lock.Lock()

	this.readTimeout = readTimeout
	this.writeTimeout = writeTimeout
}

func (this *TCPConn) LocalAddr() net.Addr {
	return this.conn.LocalAddr()
}

func (this *TCPConn) NetConn() interface{} {
	return this.conn
}

//对端地址
func (this *TCPConn) RemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}

func (this *TCPConn) SetCodec(codec Codec) {
	this.lock.Lock()
	this.codec = codec
	this.lock.Unlock()
}

func (this *TCPConn) SetCloseCallBack(closeCallback func(reason string)) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.closeCallback = closeCallback
}

func (this *TCPConn) SetContext(ctx interface{}) {
	this.lock.Lock()
	this.ctx = ctx
	this.lock.Unlock()
}

func (this *TCPConn) Context() interface{} {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.ctx
}

//开启消息处理
func (this *TCPConn) Start(msgCb func(interface{}, error)) error {
	if msgCb == nil {
		return ErrNoMsgCallBack
	}

	this.lock.Lock()
	if this.flag == started {
		return ErrSessionStarted
	}
	this.flag = started

	if this.codec == nil {
		return ErrNoCodec
	}

	this.msgCallback = msgCb
	this.lock.Unlock()

	go this.receiveThread()
	go this.sendThread()

	return nil
}

func (this *TCPConn) isClose() bool {
	this.lock.Lock()
	defer this.lock.Unlock()
	return this.flag == closed
}

//接收线程
func (this *TCPConn) receiveThread() {
	for {
		if this.isClose() {
			return
		}

		if this.readTimeout > 0 {
			this.conn.SetReadDeadline(time.Now().Add(this.readTimeout))
		}

		msg, err := this.codec.Decode(this.conn)
		if this.isClose() {
			return
		}
		if err != nil {
			//if err == io.EOF {
			//	log.Println("read Close", err.Error())
			//} else {
			//	log.Println("read err: ", err.Error())
			//}
			//关闭连接
			//this.Close(err.Error())
			this.msgCallback(nil, err)
		} else {
			if msg != nil {
				this.msgCallback(msg, nil)
			}
		}
	}
}

//发送线程
func (this *TCPConn) sendThread() {
	defer this.close()
	for {
		data, isOpen := <-this.sendBufChan
		if !isOpen {
			break
		}
		if this.writeTimeout > 0 {
			this.conn.SetWriteDeadline(time.Now().Add(this.writeTimeout))
		}

		_, err := this.conn.Write(data)
		if err != nil {
			//log.Println("write err: ", err.Error())
			//this.Close(err.Error())
			this.msgCallback(nil, err)
		}

	}
}

func (this *TCPConn) Send(o interface{}) error {
	if o == nil {
		return ErrSendMsgNil
	}

	this.lock.Lock()
	if this.codec == nil {
		return ErrNoCodec
	}
	codec := this.codec
	this.lock.Unlock()

	data, err := codec.Encode(o)
	if err != nil {
		return err
	}

	return this.SendBytes(data)
}

func (this *TCPConn) SendBytes(data []byte) error {
	if len(data) == 0 {
		return ErrSendMsgNil
	}

	//非堵塞
	if len(this.sendBufChan) == sendBufChanSize {
		return ErrSendChanFull
	}

	this.lock.Lock()
	if this.flag == 0 {
		return ErrNotStarted
	}
	if this.flag == closed {
		return ErrSessionClosed
	}
	this.lock.Unlock()

	this.sendBufChan <- data
	return nil
}

/*
 主动关闭连接
 先关闭读，待写发送完毕关闭写
*/
func (this *TCPConn) Close(reason string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	if (this.flag & closed) > 0 {
		return
	}

	close(this.sendBufChan)
	this.closeReason = reason
	this.flag = closed
	this.conn.CloseRead()
}

func (this *TCPConn) close() {
	_ = this.conn.Close()
	this.lock.Lock()
	callback := this.closeCallback
	msg := this.closeReason
	this.lock.Unlock()
	if callback != nil {
		callback(msg)
	}
}

type TCPListener struct {
	listener *net.TCPListener
	started  int32
}

func NewTCPListener(network, addr string) (*TCPListener, error) {
	tcpAddr, err := net.ResolveTCPAddr(network, addr)
	if err != nil {
		return nil, err
	}

	listener, err := net.ListenTCP(tcpAddr.Network(), tcpAddr)
	return &TCPListener{listener: listener}, err
}

func (l *TCPListener) Listen(newClient func(session Session)) error {
	if newClient == nil {
		return ErrNewClientNil
	}

	if !atomic.CompareAndSwapInt32(&l.started, 0, 1) {
		return ErrSessionStarted
	}

	go func() {
		for {
			conn, err := l.listener.Accept()
			if err != nil {
				if ne, ok := err.(net.Error); ok && ne.Temporary() {
					continue
				} else {
					return
				}
			}
			newClient(newTCPConn(conn.(*net.TCPConn)))
		}
	}()

	return nil
}

func (l *TCPListener) Addr() net.Addr {
	return l.listener.Addr()
}

func (l *TCPListener) Close() {
	if atomic.CompareAndSwapInt32(&l.started, 1, 0) {
		_ = l.listener.Close()
	}

}

func DialTCP(network, addr string, timeout time.Duration) (Session, error) {
	tcpAddr, err := net.ResolveTCPAddr(network, addr)
	if err != nil {
		return nil, err
	}

	dialer := &net.Dialer{Timeout: timeout}
	conn, err := dialer.Dial(tcpAddr.Network(), addr)
	if err != nil {
		return nil, err
	}

	return newTCPConn(conn.(*net.TCPConn)), nil
}
