// Code generated by protoc-gen-go. DO NOT EDIT.
// source: swap_resolver.proto

/*
Package swapresolver is a generated protocol buffer package.

It is generated from these files:
	swap_resolver.proto

It has these top-level messages:
	ResolveReq
	ResolveResp
	TakeOrderReq
	TakeOrderResp
	SuggestDealReq
	SuggestDealResp
	SwapReq
	SwapResp
*/
package swapresolver

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type CoinType int32

const (
	CoinType_BTC CoinType = 0
	CoinType_LTC CoinType = 1
)

var CoinType_name = map[int32]string{
	0: "BTC",
	1: "LTC",
}
var CoinType_value = map[string]int32{
	"BTC": 0,
	"LTC": 1,
}

func (x CoinType) String() string {
	return proto.EnumName(CoinType_name, int32(x))
}
func (CoinType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// 6
type ResolveReq struct {
	Hash string `protobuf:"bytes,1,opt,name=hash" json:"hash,omitempty"`
}

func (m *ResolveReq) Reset()                    { *m = ResolveReq{} }
func (m *ResolveReq) String() string            { return proto.CompactTextString(m) }
func (*ResolveReq) ProtoMessage()               {}
func (*ResolveReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ResolveReq) GetHash() string {
	if m != nil {
		return m.Hash
	}
	return ""
}

type ResolveResp struct {
	Preimage string `protobuf:"bytes,1,opt,name=preimage" json:"preimage,omitempty"`
}

func (m *ResolveResp) Reset()                    { *m = ResolveResp{} }
func (m *ResolveResp) String() string            { return proto.CompactTextString(m) }
func (*ResolveResp) ProtoMessage()               {}
func (*ResolveResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *ResolveResp) GetPreimage() string {
	if m != nil {
		return m.Preimage
	}
	return ""
}

// 1 - CLI to taker
type TakeOrderReq struct {
	Orderid     string   `protobuf:"bytes,1,opt,name=orderid" json:"orderid,omitempty"`
	TakerAmount int64    `protobuf:"varint,2,opt,name=taker_amount,json=takerAmount" json:"taker_amount,omitempty"`
	TakerCoin   CoinType `protobuf:"varint,3,opt,name=taker_coin,json=takerCoin,enum=swapresolver.CoinType" json:"taker_coin,omitempty"`
	MakerAmount int64    `protobuf:"varint,4,opt,name=maker_amount,json=makerAmount" json:"maker_amount,omitempty"`
	MakerCoin   CoinType `protobuf:"varint,5,opt,name=maker_coin,json=makerCoin,enum=swapresolver.CoinType" json:"maker_coin,omitempty"`
}

func (m *TakeOrderReq) Reset()                    { *m = TakeOrderReq{} }
func (m *TakeOrderReq) String() string            { return proto.CompactTextString(m) }
func (*TakeOrderReq) ProtoMessage()               {}
func (*TakeOrderReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *TakeOrderReq) GetOrderid() string {
	if m != nil {
		return m.Orderid
	}
	return ""
}

func (m *TakeOrderReq) GetTakerAmount() int64 {
	if m != nil {
		return m.TakerAmount
	}
	return 0
}

func (m *TakeOrderReq) GetTakerCoin() CoinType {
	if m != nil {
		return m.TakerCoin
	}
	return CoinType_BTC
}

func (m *TakeOrderReq) GetMakerAmount() int64 {
	if m != nil {
		return m.MakerAmount
	}
	return 0
}

func (m *TakeOrderReq) GetMakerCoin() CoinType {
	if m != nil {
		return m.MakerCoin
	}
	return CoinType_BTC
}

// 14 taker to CLI
type TakeOrderResp struct {
	RPreimage []byte `protobuf:"bytes,1,opt,name=r_preimage,json=rPreimage,proto3" json:"r_preimage,omitempty"`
}

func (m *TakeOrderResp) Reset()                    { *m = TakeOrderResp{} }
func (m *TakeOrderResp) String() string            { return proto.CompactTextString(m) }
func (*TakeOrderResp) ProtoMessage()               {}
func (*TakeOrderResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *TakeOrderResp) GetRPreimage() []byte {
	if m != nil {
		return m.RPreimage
	}
	return nil
}

// 2 - from taker to maker
type SuggestDealReq struct {
	Orderid     string   `protobuf:"bytes,1,opt,name=orderid" json:"orderid,omitempty"`
	TakerDealId string   `protobuf:"bytes,2,opt,name=taker_deal_id,json=takerDealId" json:"taker_deal_id,omitempty"`
	TakerAmount int64    `protobuf:"varint,3,opt,name=taker_amount,json=takerAmount" json:"taker_amount,omitempty"`
	TakerCoin   CoinType `protobuf:"varint,4,opt,name=taker_coin,json=takerCoin,enum=swapresolver.CoinType" json:"taker_coin,omitempty"`
	MakerAmount int64    `protobuf:"varint,5,opt,name=maker_amount,json=makerAmount" json:"maker_amount,omitempty"`
	MakerCoin   CoinType `protobuf:"varint,6,opt,name=maker_coin,json=makerCoin,enum=swapresolver.CoinType" json:"maker_coin,omitempty"`
	TakerPubkey string   `protobuf:"bytes,7,opt,name=taker_pubkey,json=takerPubkey" json:"taker_pubkey,omitempty"`
}

func (m *SuggestDealReq) Reset()                    { *m = SuggestDealReq{} }
func (m *SuggestDealReq) String() string            { return proto.CompactTextString(m) }
func (*SuggestDealReq) ProtoMessage()               {}
func (*SuggestDealReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *SuggestDealReq) GetOrderid() string {
	if m != nil {
		return m.Orderid
	}
	return ""
}

func (m *SuggestDealReq) GetTakerDealId() string {
	if m != nil {
		return m.TakerDealId
	}
	return ""
}

func (m *SuggestDealReq) GetTakerAmount() int64 {
	if m != nil {
		return m.TakerAmount
	}
	return 0
}

func (m *SuggestDealReq) GetTakerCoin() CoinType {
	if m != nil {
		return m.TakerCoin
	}
	return CoinType_BTC
}

func (m *SuggestDealReq) GetMakerAmount() int64 {
	if m != nil {
		return m.MakerAmount
	}
	return 0
}

func (m *SuggestDealReq) GetMakerCoin() CoinType {
	if m != nil {
		return m.MakerCoin
	}
	return CoinType_BTC
}

func (m *SuggestDealReq) GetTakerPubkey() string {
	if m != nil {
		return m.TakerPubkey
	}
	return ""
}

// 3 from maker back to taker
type SuggestDealResp struct {
	Orderid     string `protobuf:"bytes,1,opt,name=orderid" json:"orderid,omitempty"`
	RHash       []byte `protobuf:"bytes,2,opt,name=r_hash,json=rHash,proto3" json:"r_hash,omitempty"`
	MakerDealId string `protobuf:"bytes,3,opt,name=maker_deal_id,json=makerDealId" json:"maker_deal_id,omitempty"`
	MakerPubkey string `protobuf:"bytes,4,opt,name=maker_pubkey,json=makerPubkey" json:"maker_pubkey,omitempty"`
}

func (m *SuggestDealResp) Reset()                    { *m = SuggestDealResp{} }
func (m *SuggestDealResp) String() string            { return proto.CompactTextString(m) }
func (*SuggestDealResp) ProtoMessage()               {}
func (*SuggestDealResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *SuggestDealResp) GetOrderid() string {
	if m != nil {
		return m.Orderid
	}
	return ""
}

func (m *SuggestDealResp) GetRHash() []byte {
	if m != nil {
		return m.RHash
	}
	return nil
}

func (m *SuggestDealResp) GetMakerDealId() string {
	if m != nil {
		return m.MakerDealId
	}
	return ""
}

func (m *SuggestDealResp) GetMakerPubkey() string {
	if m != nil {
		return m.MakerPubkey
	}
	return ""
}

// 4 from taker to maker
type SwapReq struct {
	MakerDealId string `protobuf:"bytes,1,opt,name=maker_deal_id,json=makerDealId" json:"maker_deal_id,omitempty"`
}

func (m *SwapReq) Reset()                    { *m = SwapReq{} }
func (m *SwapReq) String() string            { return proto.CompactTextString(m) }
func (*SwapReq) ProtoMessage()               {}
func (*SwapReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *SwapReq) GetMakerDealId() string {
	if m != nil {
		return m.MakerDealId
	}
	return ""
}

// 13 maker to taker
type SwapResp struct {
	RPreimage []byte `protobuf:"bytes,1,opt,name=r_preimage,json=rPreimage,proto3" json:"r_preimage,omitempty"`
}

func (m *SwapResp) Reset()                    { *m = SwapResp{} }
func (m *SwapResp) String() string            { return proto.CompactTextString(m) }
func (*SwapResp) ProtoMessage()               {}
func (*SwapResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *SwapResp) GetRPreimage() []byte {
	if m != nil {
		return m.RPreimage
	}
	return nil
}

func init() {
	proto.RegisterType((*ResolveReq)(nil), "swapresolver.ResolveReq")
	proto.RegisterType((*ResolveResp)(nil), "swapresolver.ResolveResp")
	proto.RegisterType((*TakeOrderReq)(nil), "swapresolver.TakeOrderReq")
	proto.RegisterType((*TakeOrderResp)(nil), "swapresolver.TakeOrderResp")
	proto.RegisterType((*SuggestDealReq)(nil), "swapresolver.SuggestDealReq")
	proto.RegisterType((*SuggestDealResp)(nil), "swapresolver.SuggestDealResp")
	proto.RegisterType((*SwapReq)(nil), "swapresolver.SwapReq")
	proto.RegisterType((*SwapResp)(nil), "swapresolver.SwapResp")
	proto.RegisterEnum("swapresolver.CoinType", CoinType_name, CoinType_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for SwapResolver service

type SwapResolverClient interface {
	// ResolveHash is used by LND to request translation of Rhash to a pre-image.
	// the resolver may return the preimage and error indicating that there is no
	// such hash/deal
	ResolveHash(ctx context.Context, in *ResolveReq, opts ...grpc.CallOption) (*ResolveResp, error)
}

type swapResolverClient struct {
	cc *grpc.ClientConn
}

func NewSwapResolverClient(cc *grpc.ClientConn) SwapResolverClient {
	return &swapResolverClient{cc}
}

func (c *swapResolverClient) ResolveHash(ctx context.Context, in *ResolveReq, opts ...grpc.CallOption) (*ResolveResp, error) {
	out := new(ResolveResp)
	err := grpc.Invoke(ctx, "/swapresolver.SwapResolver/ResolveHash", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for SwapResolver service

type SwapResolverServer interface {
	// ResolveHash is used by LND to request translation of Rhash to a pre-image.
	// the resolver may return the preimage and error indicating that there is no
	// such hash/deal
	ResolveHash(context.Context, *ResolveReq) (*ResolveResp, error)
}

func RegisterSwapResolverServer(s *grpc.Server, srv SwapResolverServer) {
	s.RegisterService(&_SwapResolver_serviceDesc, srv)
}

func _SwapResolver_ResolveHash_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResolveReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SwapResolverServer).ResolveHash(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swapresolver.SwapResolver/ResolveHash",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SwapResolverServer).ResolveHash(ctx, req.(*ResolveReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _SwapResolver_serviceDesc = grpc.ServiceDesc{
	ServiceName: "swapresolver.SwapResolver",
	HandlerType: (*SwapResolverServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ResolveHash",
			Handler:    _SwapResolver_ResolveHash_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "swap_resolver.proto",
}

// Client API for P2P service

type P2PClient interface {
	// TakeOrder is called to initiate a swap between maker and taker
	// it is a temporary service needed until the integration with XUD
	// intended to be called from CLI to simulate order taking by taker
	TakeOrder(ctx context.Context, in *TakeOrderReq, opts ...grpc.CallOption) (*TakeOrderResp, error)
	// SuggestDeal is called by the taker to inform the maker that he
	// would like to execute a swap. The maker may reject the request
	// for now, the maker can only accept/reject and can't rediscuss the
	// deal or suggest partial amount. If accepted the maker should respond
	// with a hash that would be used for teh swap.
	SuggestDeal(ctx context.Context, in *SuggestDealReq, opts ...grpc.CallOption) (*SuggestDealResp, error)
	// Swap initiates the swap. It is called by the taker to confirm that
	// he has the hash and confirm the deal.
	Swap(ctx context.Context, in *SwapReq, opts ...grpc.CallOption) (*SwapResp, error)
}

type p2PClient struct {
	cc *grpc.ClientConn
}

func NewP2PClient(cc *grpc.ClientConn) P2PClient {
	return &p2PClient{cc}
}

func (c *p2PClient) TakeOrder(ctx context.Context, in *TakeOrderReq, opts ...grpc.CallOption) (*TakeOrderResp, error) {
	out := new(TakeOrderResp)
	err := grpc.Invoke(ctx, "/swapresolver.P2P/TakeOrder", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *p2PClient) SuggestDeal(ctx context.Context, in *SuggestDealReq, opts ...grpc.CallOption) (*SuggestDealResp, error) {
	out := new(SuggestDealResp)
	err := grpc.Invoke(ctx, "/swapresolver.P2P/SuggestDeal", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *p2PClient) Swap(ctx context.Context, in *SwapReq, opts ...grpc.CallOption) (*SwapResp, error) {
	out := new(SwapResp)
	err := grpc.Invoke(ctx, "/swapresolver.P2P/Swap", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for P2P service

type P2PServer interface {
	// TakeOrder is called to initiate a swap between maker and taker
	// it is a temporary service needed until the integration with XUD
	// intended to be called from CLI to simulate order taking by taker
	TakeOrder(context.Context, *TakeOrderReq) (*TakeOrderResp, error)
	// SuggestDeal is called by the taker to inform the maker that he
	// would like to execute a swap. The maker may reject the request
	// for now, the maker can only accept/reject and can't rediscuss the
	// deal or suggest partial amount. If accepted the maker should respond
	// with a hash that would be used for teh swap.
	SuggestDeal(context.Context, *SuggestDealReq) (*SuggestDealResp, error)
	// Swap initiates the swap. It is called by the taker to confirm that
	// he has the hash and confirm the deal.
	Swap(context.Context, *SwapReq) (*SwapResp, error)
}

func RegisterP2PServer(s *grpc.Server, srv P2PServer) {
	s.RegisterService(&_P2P_serviceDesc, srv)
}

func _P2P_TakeOrder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TakeOrderReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(P2PServer).TakeOrder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swapresolver.P2P/TakeOrder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(P2PServer).TakeOrder(ctx, req.(*TakeOrderReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _P2P_SuggestDeal_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SuggestDealReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(P2PServer).SuggestDeal(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swapresolver.P2P/SuggestDeal",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(P2PServer).SuggestDeal(ctx, req.(*SuggestDealReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _P2P_Swap_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SwapReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(P2PServer).Swap(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/swapresolver.P2P/Swap",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(P2PServer).Swap(ctx, req.(*SwapReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _P2P_serviceDesc = grpc.ServiceDesc{
	ServiceName: "swapresolver.P2P",
	HandlerType: (*P2PServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "TakeOrder",
			Handler:    _P2P_TakeOrder_Handler,
		},
		{
			MethodName: "SuggestDeal",
			Handler:    _P2P_SuggestDeal_Handler,
		},
		{
			MethodName: "Swap",
			Handler:    _P2P_Swap_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "swap_resolver.proto",
}

func init() { proto.RegisterFile("swap_resolver.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 480 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x94, 0xdf, 0x6e, 0xd3, 0x30,
	0x14, 0xc6, 0xe7, 0x26, 0xfd, 0x93, 0xd3, 0x6c, 0x4c, 0x46, 0x9b, 0x42, 0xd8, 0xa4, 0xe2, 0xab,
	0x0e, 0x89, 0x5e, 0x14, 0x21, 0xae, 0xa1, 0x13, 0x02, 0x69, 0x12, 0x55, 0x96, 0xfb, 0xc8, 0xa3,
	0x56, 0x1b, 0xb5, 0x6e, 0x8c, 0xdd, 0x32, 0xed, 0x11, 0x78, 0x10, 0xde, 0x88, 0x4b, 0x1e, 0x06,
	0xd9, 0x49, 0x5a, 0x9b, 0xb2, 0x16, 0x71, 0x77, 0x7c, 0xfc, 0xf9, 0xcb, 0x39, 0x3f, 0x1f, 0x07,
	0x9e, 0xaa, 0x7b, 0x2a, 0x32, 0xc9, 0x54, 0xb1, 0xf8, 0xc6, 0xe4, 0x40, 0xc8, 0x62, 0x55, 0xe0,
	0x50, 0x27, 0xeb, 0x1c, 0xe9, 0x01, 0x24, 0x65, 0x9c, 0xb0, 0xaf, 0x18, 0x83, 0x3f, 0xa3, 0x6a,
	0x16, 0xa1, 0x1e, 0xea, 0x07, 0x89, 0x89, 0xc9, 0x15, 0x74, 0x37, 0x0a, 0x25, 0x70, 0x0c, 0x1d,
	0x21, 0x59, 0xce, 0xe9, 0x94, 0x55, 0xb2, 0xcd, 0x9a, 0xfc, 0x42, 0x10, 0xa6, 0x74, 0xce, 0x3e,
	0xcb, 0x09, 0x93, 0xda, 0x2f, 0x82, 0x76, 0xa1, 0xe3, 0x7c, 0x52, 0x69, 0xeb, 0x25, 0x7e, 0x01,
	0xe1, 0x8a, 0xce, 0x99, 0xcc, 0x28, 0x2f, 0xd6, 0xcb, 0x55, 0xd4, 0xe8, 0xa1, 0xbe, 0x97, 0x74,
	0x4d, 0xee, 0x9d, 0x49, 0xe1, 0x37, 0x00, 0xa5, 0xe4, 0x4b, 0x91, 0x2f, 0x23, 0xaf, 0x87, 0xfa,
	0x27, 0xc3, 0xf3, 0x81, 0x5d, 0xfd, 0x60, 0x54, 0xe4, 0xcb, 0xf4, 0x41, 0xb0, 0x24, 0x30, 0x4a,
	0xbd, 0xd4, 0xce, 0xdc, 0x76, 0xf6, 0x4b, 0x67, 0xee, 0x3a, 0xf3, 0xad, 0x73, 0x73, 0xbf, 0x33,
	0xaf, 0x9d, 0xc9, 0x00, 0x8e, 0xad, 0xee, 0x94, 0xc0, 0x97, 0x00, 0x32, 0x73, 0x68, 0x84, 0x49,
	0x20, 0xc7, 0x35, 0x8e, 0x1f, 0x0d, 0x38, 0xb9, 0x5d, 0x4f, 0xa7, 0x4c, 0xad, 0xae, 0x19, 0x5d,
	0xec, 0x07, 0x42, 0xe0, 0xb8, 0xec, 0x76, 0xc2, 0xe8, 0x22, 0xcb, 0x27, 0x86, 0x48, 0x50, 0x11,
	0xd1, 0xc7, 0x3f, 0xed, 0x42, 0xf3, 0x0e, 0x41, 0xf3, 0xff, 0x17, 0x5a, 0xf3, 0x10, 0xb4, 0xd6,
	0x3f, 0x42, 0xdb, 0xd6, 0x2c, 0xd6, 0x77, 0x73, 0xf6, 0x10, 0xb5, 0xad, 0xb6, 0xc6, 0x26, 0x45,
	0xbe, 0x23, 0x78, 0xe2, 0x70, 0x52, 0x62, 0x0f, 0xa8, 0x33, 0x68, 0xc9, 0xcc, 0x4c, 0x69, 0xc3,
	0x00, 0x6f, 0xca, 0x8f, 0x54, 0xcd, 0x34, 0x3f, 0xee, 0xf0, 0xf3, 0xca, 0x0f, 0x71, 0x97, 0x1f,
	0xb7, 0x6b, 0xf1, 0x2d, 0x49, 0x55, 0xcb, 0x2b, 0x68, 0xdf, 0xde, 0x53, 0xa1, 0xef, 0x6a, 0xc7,
	0x11, 0xed, 0x38, 0x92, 0x2b, 0xe8, 0x94, 0xf2, 0x83, 0xd3, 0xf0, 0xf2, 0x02, 0x3a, 0x35, 0x1f,
	0xdc, 0x06, 0xef, 0x7d, 0x3a, 0x3a, 0x3d, 0xd2, 0xc1, 0x4d, 0x3a, 0x3a, 0x45, 0xc3, 0x14, 0xc2,
	0xca, 0xc8, 0xa0, 0xc4, 0xd7, 0x9b, 0x57, 0x67, 0xba, 0x8b, 0x5c, 0xd0, 0xdb, 0x27, 0x1b, 0x3f,
	0x7b, 0x64, 0x47, 0x09, 0x72, 0x34, 0xfc, 0x89, 0xc0, 0x1b, 0x0f, 0xc7, 0xf8, 0x03, 0x04, 0x9b,
	0xc9, 0xc5, 0xb1, 0x7b, 0xc2, 0x7e, 0xb0, 0xf1, 0xf3, 0x47, 0xf7, 0xb4, 0x1f, 0xbe, 0x81, 0xae,
	0x75, 0x51, 0xf8, 0xc2, 0x55, 0xbb, 0xb3, 0x1e, 0x5f, 0xee, 0xd9, 0x35, 0x6e, 0x6f, 0xc1, 0xd7,
	0x3d, 0xe3, 0xb3, 0x3f, 0x84, 0x25, 0xff, 0xf8, 0xfc, 0x6f, 0x69, 0x7d, 0xf0, 0xae, 0x65, 0xfe,
	0x64, 0xaf, 0x7f, 0x07, 0x00, 0x00, 0xff, 0xff, 0xe5, 0x2f, 0x29, 0x85, 0xe0, 0x04, 0x00, 0x00,
}
