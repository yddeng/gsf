package cluster

import (
	"encoding/binary"
	"fmt"
	"github.com/yddeng/clugs/cluster/addr"
	"github.com/yddeng/clugs/cluster/clusterpb"
	"github.com/yddeng/clugs/codec/ss"
	"github.com/yddeng/clugs/logger"
	"github.com/yddeng/dnet"
	"github.com/yddeng/dnet/drpc"
	"github.com/yddeng/dutil/buffer"
	"io"
	"net"
	"reflect"
	"time"
)

// 节点间建立连接时，用于第一次消息编码
// 将自己的逻辑地址打包
// 4 + 2 + n, 逻辑地址、net地址长度、net地址
func encode(logic addr.LogicAddr, netAddr string) []byte {
	buff := buffer.NewBufferWithCap(64)
	buff.WriteUint32BE(logic.Uint32())
	netAddrSize := len(netAddr)
	buff.WriteUint16BE(uint16(netAddrSize))
	buff.WriteString(netAddr)
	placeholder := make([]byte, 64-4-2-netAddrSize)
	buff.Write(placeholder)
	return buff.Bytes()
}

// 节点间建立连接时，用于第一次消息解码
// 拆包逻辑地址
func decode(data []byte) (logic addr.LogicAddr, netAddr string) {
	buff := buffer.NewBuffer(data)
	logicId, _ := buff.ReadUint32BE()
	logic = addr.LogicAddr(logicId)
	netAddrSize, _ := buff.ReadUint16BE()
	netAddr, _ = buff.ReadString(int(netAddrSize))
	return
}

// 建立连接
func dial(end *endpoint) {
	if end.dialing {
		return
	}

	end.dialing = true
	if end.dialTimeout.IsZero() {
		end.dialTimeout = time.Now().Add(rpcTimeout / 2)
	}

	go func() {
		end.Lock()
		conn, err := net.DialTCP("tcp", nil, end.logic.Net)
		if nil == err {
			data := encode(LocalAddr.Logic, LocalAddr.NetString())
			end.Unlock()

			conn.SetWriteDeadline(time.Now().Add(heartbeatTime))
			_, e := conn.Write(data)
			if e != nil {
				conn.Close()
				dialFailed(end, e)
				return
			}
			conn.SetWriteDeadline(time.Time{})

			buff := make([]byte, 4)
			_, e = io.ReadFull(conn, buff)
			if e != nil {
				conn.Close()
				dialFailed(end, e)
				return
			}

			code := binary.BigEndian.Uint32(buff)
			if code != 1 {
				conn.Close()
				dialFailed(end, fmt.Errorf("code = %d", code))
				return
			} else {
				connectOk(end, conn)
			}

		} else {
			end.Unlock()
			logger.Errorf("cluster:dial endpoint %s netAddr %s error:%s \n", end.logic.Logic.String(), end.logic.NetString(), err)
			dialFailed(end, err)
		}
	}()
}

func dialFailed(end *endpoint, err error) {
	isSame := endGroup.getEndpoint(end.logic.Logic) == end

	end.Lock()
	defer end.Unlock()
	end.dialing = false

	now := time.Now()
	if end.session == nil {
		logger.Errorf("cluster:dialFailed error %s \n", err)
		if isSame && now.Before(end.dialTimeout) {
			time.Sleep(time.Millisecond * 100)
			dial(end)
			return
		} else {
			end.ssMsg = end.ssMsg[0:0]
			reqMsg := end.reqMsg
			end.reqMsg = end.reqMsg[0:0]
			logicAddr := end.logic.Logic.String()

			taskQueue.Push(func() {
				for _, req := range reqMsg {
					_ = rpcMgr.rpcClient.OnRPCResponse(&drpc.Response{
						Seq:   req.Seq,
						Data:  nil,
						Error: fmt.Sprintf("connect logicAddr %s failed", logicAddr),
					})
				}
			})
		}
	}
	end.dialTimeout = time.Time{}
}

var (
	codeOk     []byte
	codeFailed []byte
)

func init() {
	codeOk = make([]byte, 4)
	codeFailed = make([]byte, 4)
	binary.BigEndian.PutUint32(codeOk, 1)
	binary.BigEndian.PutUint32(codeFailed, 0)
}

func acceptConn(conn *net.TCPConn) {
	conn.SetReadDeadline(time.Now().Add(heartbeatTime))
	buff := make([]byte, 64)
	_, err := io.ReadFull(conn, buff)
	if err != nil {
		conn.Write(codeFailed)
		logger.Errorf("cluster:acceptConn read error %s. ", err)
		return
	}
	conn.SetReadDeadline(time.Time{})

	logic, netStr := decode(buff)
	end := endGroup.getEndpoint(logic)
	if end == nil {
		conn.Write(codeFailed)
		logger.Errorf("cluster:acceptConn logic %s is nil", logic.String())
		return
	}

	end.Lock()
	if end.logic.NetString() != netStr {
		end.Unlock()
		conn.Write(codeFailed)
		logger.Errorf("cluster:acceptConn logic %s netAddr not equal %s != %s", logic.String(), end.logic.NetString(), netStr)
		return
	}
	end.Unlock()

	conn.Write(codeOk)

	// 连接成功
	connectOk(end, conn)
}

func connectOk(end *endpoint, conn dnet.NetConn) {
	session := dnet.NewTCPSession(conn,
		dnet.WithTimeout(heartbeatTime, 0),
		dnet.WithCodec(ss.NewCodec(clusterpb.SS_SPACE, clusterpb.REQ_SPACE, clusterpb.RESP_SPACE)),
		dnet.WithCloseCallback(func(session dnet.Session, reason error) {
			end.Lock()
			defer end.Unlock()
			end.session = nil
			session.SetContext(nil)
			logger.Infof("cluster:connectOK endpoint %s session closed, reason: %s\n", end.logic.Logic.String(), reason)
		}),
		dnet.WithErrorCallback(func(session dnet.Session, err error) {
			logger.Error("cluster:connectOK session error:", err)
			session.Close(err)
		}),
		dnet.WithMessageCallback(func(session dnet.Session, message interface{}) {
			taskQueue.Push(func() {
				var err error
				switch message.(type) {
				case *ss.Message:
					end.Lock()
					err = dispatchSS(end.logic.Logic, message.(*ss.Message))
					end.Unlock()
				case *drpc.Request:
					err = rpcMgr.rpcServer.OnRPCRequest(end, message.(*drpc.Request))
				case *drpc.Response:
					err = rpcMgr.rpcClient.OnRPCResponse(message.(*drpc.Response))
				default:
					err = fmt.Errorf("invalid type:%s", reflect.TypeOf(message).String())
				}
				if err != nil {
					logger.Errorf("cluster:connectOK dispatch error: %s. \n", err.Error())
				}
			})
		}),
	)

	end.Lock()
	defer end.Unlock()

	end.dialing = false
	end.dialTimeout = time.Time{}

	if end.session != nil {
		logger.Infof("cluster:connectOK endpoint %s already connect", end.logic.Logic.String())
		session.Close(fmt.Errorf("cluster:connectOK endpoint %s already connect", end.logic.Logic.String()))
		return
	}

	logger.Infof("cluster:connectOK endpoint connection %s <-> %s", LocalAddr.Logic.String(), end.logic.Logic.String())

	end.session = session
	session.SetContext(end)

	// 将消息发送出去
	for _, msg := range end.ssMsg {
		_ = end.send(msg)
	}
	for _, req := range end.reqMsg {
		_ = end.send(req)
	}
	end.ssMsg = end.ssMsg[0:0]
	end.reqMsg = end.reqMsg[0:0]

}
