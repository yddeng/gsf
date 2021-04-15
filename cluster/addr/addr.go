package addr

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

const (
	GroupMask  uint32 = 0xFFFC0000 //高14 (1,16383)
	TypeMask   uint32 = 0x0003FC00 //中8 (1,254)
	ServerMask uint32 = 0x000003FF //低10 (1,1023)
)

var (
	ErrInvalidAddrFmt = fmt.Errorf("invalid addr format")
	ErrInvalidGroup   = fmt.Errorf("group should between(1,16383)")
	ErrInvalidType    = fmt.Errorf("type should between(1,254)")
	ErrInvalidServer  = fmt.Errorf("server should between(1,1023)")
)

type LogicAddr uint32

type Addr struct {
	Logic LogicAddr
	Net   *net.TCPAddr
}

func (this *Addr) NetString() string {
	return this.Net.String()
}

func MakeAddr(logic string, tcpAddr string) (*Addr, error) {
	logicAddr, err := MakeLogicAddr(logic)
	if nil != err {
		return nil, err
	}

	netAddr, err := net.ResolveTCPAddr("tcp", tcpAddr)
	if nil != err {
		return nil, err
	}

	return &Addr{
		Logic: logicAddr,
		Net:   netAddr,
	}, nil
}

func (this LogicAddr) Uint32() uint32 {
	return uint32(this)
}

func (this LogicAddr) Group() uint32 {
	return (uint32(this) & GroupMask) >> 18
}

func (this LogicAddr) Type() uint32 {
	return (uint32(this) & TypeMask) >> 10
}

func (this LogicAddr) Server() uint32 {
	return uint32(this) & ServerMask
}

func (this LogicAddr) String() string {
	return fmt.Sprintf("%d.%d.%d", this.Group(), this.Type(), this.Server())
}

func (this LogicAddr) Empty() bool {
	return uint32(this) == 0
}

func (this *LogicAddr) Clear() {
	(*this) = 0
}

func MakeLogicAddr(addr string) (LogicAddr, error) {
	var err error
	v := strings.Split(addr, ".")
	if len(v) != 3 {
		return LogicAddr(0), ErrInvalidAddrFmt
	}

	group, err := strconv.Atoi(v[0])
	if nil != err {
		return LogicAddr(0), ErrInvalidGroup
	}

	if 0 == group || uint32(group) > (GroupMask>>18) {
		return LogicAddr(0), ErrInvalidGroup
	}

	tt, err := strconv.Atoi(v[1])
	if nil != err {
		return LogicAddr(0), ErrInvalidType
	}

	if 0 == tt || uint32(tt) > ((TypeMask>>10)-1) {
		return LogicAddr(0), ErrInvalidType
	}

	server, err := strconv.Atoi(v[2])
	if nil != err {
		return LogicAddr(0), ErrInvalidServer
	}

	if server == 0 || uint32(server) > ServerMask {
		return LogicAddr(0), ErrInvalidServer
	}

	return LogicAddr(0 | (uint32(tt) << 10) | (uint32(group) << 18) | (uint32(server))), nil
}
