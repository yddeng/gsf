// Code generated by protoc-gen-go. DO NOT EDIT.
// source: center.proto

package center

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type NodeInfo struct {
	LogicAddr            uint32   `protobuf:"varint,1,opt,name=logicAddr,proto3" json:"logicAddr,omitempty"`
	NetAddr              string   `protobuf:"bytes,2,opt,name=netAddr,proto3" json:"netAddr,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NodeInfo) Reset()         { *m = NodeInfo{} }
func (m *NodeInfo) String() string { return proto.CompactTextString(m) }
func (*NodeInfo) ProtoMessage()    {}
func (*NodeInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_1de517c49d537f4b, []int{0}
}

func (m *NodeInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodeInfo.Unmarshal(m, b)
}
func (m *NodeInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodeInfo.Marshal(b, m, deterministic)
}
func (m *NodeInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeInfo.Merge(m, src)
}
func (m *NodeInfo) XXX_Size() int {
	return xxx_messageInfo_NodeInfo.Size(m)
}
func (m *NodeInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeInfo.DiscardUnknown(m)
}

var xxx_messageInfo_NodeInfo proto.InternalMessageInfo

func (m *NodeInfo) GetLogicAddr() uint32 {
	if m != nil {
		return m.LogicAddr
	}
	return 0
}

func (m *NodeInfo) GetNetAddr() string {
	if m != nil {
		return m.NetAddr
	}
	return ""
}

type LoginReq struct {
	Node                 *NodeInfo `protobuf:"bytes,1,opt,name=node,proto3" json:"node,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *LoginReq) Reset()         { *m = LoginReq{} }
func (m *LoginReq) String() string { return proto.CompactTextString(m) }
func (*LoginReq) ProtoMessage()    {}
func (*LoginReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_1de517c49d537f4b, []int{1}
}

func (m *LoginReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginReq.Unmarshal(m, b)
}
func (m *LoginReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginReq.Marshal(b, m, deterministic)
}
func (m *LoginReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginReq.Merge(m, src)
}
func (m *LoginReq) XXX_Size() int {
	return xxx_messageInfo_LoginReq.Size(m)
}
func (m *LoginReq) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginReq.DiscardUnknown(m)
}

var xxx_messageInfo_LoginReq proto.InternalMessageInfo

func (m *LoginReq) GetNode() *NodeInfo {
	if m != nil {
		return m.Node
	}
	return nil
}

type LoginResp struct {
	ErrCode              int32    `protobuf:"varint,1,opt,name=errCode,proto3" json:"errCode,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginResp) Reset()         { *m = LoginResp{} }
func (m *LoginResp) String() string { return proto.CompactTextString(m) }
func (*LoginResp) ProtoMessage()    {}
func (*LoginResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_1de517c49d537f4b, []int{2}
}

func (m *LoginResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginResp.Unmarshal(m, b)
}
func (m *LoginResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginResp.Marshal(b, m, deterministic)
}
func (m *LoginResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginResp.Merge(m, src)
}
func (m *LoginResp) XXX_Size() int {
	return xxx_messageInfo_LoginResp.Size(m)
}
func (m *LoginResp) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginResp.DiscardUnknown(m)
}

var xxx_messageInfo_LoginResp proto.InternalMessageInfo

func (m *LoginResp) GetErrCode() int32 {
	if m != nil {
		return m.ErrCode
	}
	return 0
}

func (m *LoginResp) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

type Heartbeat struct {
	Timestamp            int64    `protobuf:"varint,1,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Heartbeat) Reset()         { *m = Heartbeat{} }
func (m *Heartbeat) String() string { return proto.CompactTextString(m) }
func (*Heartbeat) ProtoMessage()    {}
func (*Heartbeat) Descriptor() ([]byte, []int) {
	return fileDescriptor_1de517c49d537f4b, []int{3}
}

func (m *Heartbeat) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Heartbeat.Unmarshal(m, b)
}
func (m *Heartbeat) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Heartbeat.Marshal(b, m, deterministic)
}
func (m *Heartbeat) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Heartbeat.Merge(m, src)
}
func (m *Heartbeat) XXX_Size() int {
	return xxx_messageInfo_Heartbeat.Size(m)
}
func (m *Heartbeat) XXX_DiscardUnknown() {
	xxx_messageInfo_Heartbeat.DiscardUnknown(m)
}

var xxx_messageInfo_Heartbeat proto.InternalMessageInfo

func (m *Heartbeat) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

type NotifyNodeInfo struct {
	Nodes                []*NodeInfo `protobuf:"bytes,1,rep,name=nodes,proto3" json:"nodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *NotifyNodeInfo) Reset()         { *m = NotifyNodeInfo{} }
func (m *NotifyNodeInfo) String() string { return proto.CompactTextString(m) }
func (*NotifyNodeInfo) ProtoMessage()    {}
func (*NotifyNodeInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_1de517c49d537f4b, []int{4}
}

func (m *NotifyNodeInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NotifyNodeInfo.Unmarshal(m, b)
}
func (m *NotifyNodeInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NotifyNodeInfo.Marshal(b, m, deterministic)
}
func (m *NotifyNodeInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NotifyNodeInfo.Merge(m, src)
}
func (m *NotifyNodeInfo) XXX_Size() int {
	return xxx_messageInfo_NotifyNodeInfo.Size(m)
}
func (m *NotifyNodeInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_NotifyNodeInfo.DiscardUnknown(m)
}

var xxx_messageInfo_NotifyNodeInfo proto.InternalMessageInfo

func (m *NotifyNodeInfo) GetNodes() []*NodeInfo {
	if m != nil {
		return m.Nodes
	}
	return nil
}

type NodeLeave struct {
	Nodes                []uint32 `protobuf:"varint,1,rep,packed,name=nodes,proto3" json:"nodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NodeLeave) Reset()         { *m = NodeLeave{} }
func (m *NodeLeave) String() string { return proto.CompactTextString(m) }
func (*NodeLeave) ProtoMessage()    {}
func (*NodeLeave) Descriptor() ([]byte, []int) {
	return fileDescriptor_1de517c49d537f4b, []int{5}
}

func (m *NodeLeave) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodeLeave.Unmarshal(m, b)
}
func (m *NodeLeave) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodeLeave.Marshal(b, m, deterministic)
}
func (m *NodeLeave) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeLeave.Merge(m, src)
}
func (m *NodeLeave) XXX_Size() int {
	return xxx_messageInfo_NodeLeave.Size(m)
}
func (m *NodeLeave) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeLeave.DiscardUnknown(m)
}

var xxx_messageInfo_NodeLeave proto.InternalMessageInfo

func (m *NodeLeave) GetNodes() []uint32 {
	if m != nil {
		return m.Nodes
	}
	return nil
}

type NodeChange struct {
	Nodes                []*NodeInfo `protobuf:"bytes,1,rep,name=nodes,proto3" json:"nodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *NodeChange) Reset()         { *m = NodeChange{} }
func (m *NodeChange) String() string { return proto.CompactTextString(m) }
func (*NodeChange) ProtoMessage()    {}
func (*NodeChange) Descriptor() ([]byte, []int) {
	return fileDescriptor_1de517c49d537f4b, []int{6}
}

func (m *NodeChange) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodeChange.Unmarshal(m, b)
}
func (m *NodeChange) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodeChange.Marshal(b, m, deterministic)
}
func (m *NodeChange) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeChange.Merge(m, src)
}
func (m *NodeChange) XXX_Size() int {
	return xxx_messageInfo_NodeChange.Size(m)
}
func (m *NodeChange) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeChange.DiscardUnknown(m)
}

var xxx_messageInfo_NodeChange proto.InternalMessageInfo

func (m *NodeChange) GetNodes() []*NodeInfo {
	if m != nil {
		return m.Nodes
	}
	return nil
}

func init() {
	proto.RegisterType((*NodeInfo)(nil), "nodeInfo")
	proto.RegisterType((*LoginReq)(nil), "loginReq")
	proto.RegisterType((*LoginResp)(nil), "loginResp")
	proto.RegisterType((*Heartbeat)(nil), "heartbeat")
	proto.RegisterType((*NotifyNodeInfo)(nil), "notifyNodeInfo")
	proto.RegisterType((*NodeLeave)(nil), "nodeLeave")
	proto.RegisterType((*NodeChange)(nil), "nodeChange")
}

func init() { proto.RegisterFile("center.proto", fileDescriptor_1de517c49d537f4b) }

var fileDescriptor_1de517c49d537f4b = []byte{
	// 253 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x50, 0xb1, 0x4e, 0xc3, 0x30,
	0x10, 0x55, 0x08, 0x85, 0xfa, 0xa0, 0x08, 0x59, 0x0c, 0x11, 0x2a, 0x22, 0x78, 0x6a, 0x07, 0x5a,
	0x01, 0x03, 0x33, 0xed, 0x84, 0x84, 0x18, 0x3c, 0xb2, 0xb9, 0xc9, 0xc5, 0x8d, 0xd4, 0xd8, 0xc1,
	0xb6, 0x90, 0xf8, 0x7b, 0x74, 0x4e, 0xad, 0xc2, 0xc4, 0xe4, 0x77, 0xcf, 0xef, 0xde, 0xbb, 0x3b,
	0x38, 0xaf, 0xd0, 0x04, 0x74, 0x8b, 0xde, 0xd9, 0x60, 0xc5, 0x0a, 0xc6, 0xc6, 0xd6, 0xf8, 0x6a,
	0x1a, 0xcb, 0xa7, 0xc0, 0x76, 0x56, 0xb7, 0xd5, 0x4b, 0x5d, 0xbb, 0x22, 0x2b, 0xb3, 0xd9, 0x44,
	0x1e, 0x08, 0x5e, 0xc0, 0xa9, 0xc1, 0x10, 0xff, 0x8e, 0xca, 0x6c, 0xc6, 0x64, 0x2a, 0xc5, 0x1c,
	0xc6, 0x24, 0x33, 0x12, 0x3f, 0xf9, 0x0d, 0x1c, 0x93, 0x5f, 0x6c, 0x3f, 0x7b, 0x64, 0x8b, 0x64,
	0x2e, 0x23, 0x2d, 0x9e, 0x87, 0x08, 0x23, 0xd1, 0xf7, 0xe4, 0x88, 0xce, 0xad, 0x93, 0x7c, 0x24,
	0x53, 0xc9, 0x2f, 0x21, 0xef, 0xbc, 0xde, 0xe7, 0x10, 0x14, 0x73, 0x60, 0x5b, 0x54, 0x2e, 0x6c,
	0x50, 0x05, 0x1a, 0x34, 0xb4, 0x1d, 0xfa, 0xa0, 0xba, 0x3e, 0xb6, 0xe6, 0xf2, 0x40, 0x88, 0x07,
	0xb8, 0x30, 0x36, 0xb4, 0xcd, 0xf7, 0x7b, 0x5a, 0xec, 0x16, 0x46, 0x94, 0xee, 0x8b, 0xac, 0xcc,
	0xff, 0x4e, 0x35, 0xf0, 0xe2, 0x0e, 0x18, 0x81, 0x37, 0x54, 0x5f, 0xc8, 0xaf, 0x7e, 0xab, 0x27,
	0x49, 0x72, 0x0f, 0x40, 0x60, 0xbd, 0x55, 0x46, 0xe3, 0xbf, 0x8e, 0xab, 0xe9, 0xc7, 0xb5, 0xf6,
	0xcd, 0x32, 0x1e, 0xb9, 0xb2, 0xbb, 0xe5, 0x70, 0xf4, 0xfd, 0xb3, 0x39, 0x89, 0xfc, 0xd3, 0x4f,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x4b, 0x56, 0xc4, 0xaa, 0x8c, 0x01, 0x00, 0x00,
}
