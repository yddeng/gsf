package cluster

import (
	"encoding/binary"
	"fmt"
	"github.com/yddeng/dnet"
	"github.com/yddeng/dnet/drpc"
	"github.com/yddeng/dnet/dtcp"
	"github.com/yddeng/dutil/buffer"
	"github.com/yddeng/gsf/cluster/addr"
	"github.com/yddeng/gsf/codec/ss"
	protorpc "github.com/yddeng/gsf/protocol/rpc"
	protoss "github.com/yddeng/gsf/protocol/ss"
	"github.com/yddeng/gsf/util"
	"io"
	"net"
	"reflect"
	"time"
)

// 节点间建立连接时，用于第一次消息编码
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
			data := encode(selfPoint.logic.Logic, selfPoint.logic.NetString())
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
				session := dtcp.NewTCPConn(conn)
				connectOk(end, session)
			}

		} else {
			end.Unlock()
			util.Logger().Errorf("dial endpoint %s netAddr %s error:%s \n", end.logic.Logic.String(), end.logic.NetString(), err)
			dialFailed(end, err)
		}
	}()
}

func dialFailed(end *endpoint, err error) {
	isSame := endpoints.getEndpointByLogic(end.logic.Logic) == end

	end.Lock()
	defer end.Unlock()
	end.dialing = false

	now := time.Now()
	if end.session == nil {
		util.Logger().Errorf("connectFailed err %s \n", err)
		if isSame && now.Before(end.dialTimeout) {
			time.Sleep(time.Second)
			dial(end)
			return
		} else {
			end.ssMsg = end.ssMsg[0:0]
			reqMsg := end.reqMsg
			end.reqMsg = end.reqMsg[0:0]
			logicAddr := end.logic.Logic.String()

			eventQueue.Push(func() {
				for _, req := range reqMsg {
					_ = rpcMgr.rpcClient.OnRPCResponse(&drpc.Response{
						SeqNo: req.SeqNo,
						Data:  nil,
						Err:   fmt.Errorf("connect logicAddr %s failed", logicAddr),
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

func accept(conn *net.TCPConn) {
	conn.SetReadDeadline(time.Now().Add(heartbeatTime))
	buff := make([]byte, 64)
	_, err := io.ReadFull(conn, buff)
	if err != nil {
		conn.Write(codeFailed)
		util.Logger().Errorf("accept conn, read err %s", err)
		return
	}
	conn.SetReadDeadline(time.Time{})

	logic, netStr := decode(buff)
	end := endpoints.getEndpointByLogic(logic)
	if end == nil {
		conn.Write(codeFailed)
		util.Logger().Errorf("accept conn, logic %s is nil", logic.String())
		return
	}

	end.Lock()
	if end.logic.NetString() != netStr {
		end.Unlock()
		conn.Write(codeFailed)
		util.Logger().Errorf("accept conn, logic %s netAddr not equal %s != %s", logic.String(), end.logic.NetString(), netStr)
		return
	}
	end.Unlock()

	conn.Write(codeOk)

	// 连接成功
	session := dtcp.NewTCPConn(conn)
	connectOk(end, session)
}

func connectOk(end *endpoint, session dnet.Session) {

	session.SetTimeout(heartbeatTime, 0)
	session.SetCodec(ss.NewCodec(protoss.SS_SPACE, protorpc.REQ_SPACE, protorpc.RESP_SPACE))
	session.SetCloseCallBack(func(session dnet.Session, reason string) {
		end.Lock()
		defer end.Unlock()
		end.session = nil
		session.SetContext(nil)
		util.Logger().Infof("endpoint %s session closed, reason: %s\n", end.logic.Logic.String(), reason)
	})

	end.Lock()
	defer end.Unlock()

	end.dialing = false
	end.dialTimeout = time.Time{}

	if end.session != nil {
		util.Logger().Infof("endpoint %s already connect", end.logic.Logic.String())
		session.Close(fmt.Sprintf("endpoint %s already connect", end.logic.Logic.String()))
		return
	}

	util.Logger().Infof("endpoint connect %s <-> %s", selfPoint.logic.Logic.String(), end.logic.Logic.String())

	end.session = session
	session.SetContext(end)

	session.Start(func(data interface{}, err error) {
		if err != nil {
			session.Close(err.Error())
		} else {
			eventQueue.Push(func() {
				var err error
				switch data.(type) {
				case *ss.Message:
					end.Lock()
					err = dispatchSS(end.logic.Logic, data.(*ss.Message))
					end.Unlock()
				case *drpc.Request:
					err = rpcMgr.rpcServer.OnRPCRequest(end, data.(*drpc.Request))
				case *drpc.Response:
					err = rpcMgr.rpcClient.OnRPCResponse(data.(*drpc.Response))
				default:
					err = fmt.Errorf("invalid type:%s", reflect.TypeOf(data).String())
				}
				if err != nil {
					util.Logger().Errorf(err.Error())
				}
			})
		}
	})

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
