// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: gozab.proto

package __

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// FollowerLeaderClient is the client API for FollowerLeader service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FollowerLeaderClient interface {
	Broadcast(ctx context.Context, in *PropTxn, opts ...grpc.CallOption) (*AckTxn, error)
	Commit(ctx context.Context, in *CommitTxn, opts ...grpc.CallOption) (*Empty, error)
	HeartBeat(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
}

type followerLeaderClient struct {
	cc grpc.ClientConnInterface
}

func NewFollowerLeaderClient(cc grpc.ClientConnInterface) FollowerLeaderClient {
	return &followerLeaderClient{cc}
}

func (c *followerLeaderClient) Broadcast(ctx context.Context, in *PropTxn, opts ...grpc.CallOption) (*AckTxn, error) {
	out := new(AckTxn)
	err := c.cc.Invoke(ctx, "/gozab.FollowerLeader/Broadcast", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *followerLeaderClient) Commit(ctx context.Context, in *CommitTxn, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/gozab.FollowerLeader/Commit", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *followerLeaderClient) HeartBeat(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/gozab.FollowerLeader/HeartBeat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FollowerLeaderServer is the server API for FollowerLeader service.
// All implementations must embed UnimplementedFollowerLeaderServer
// for forward compatibility
type FollowerLeaderServer interface {
	Broadcast(context.Context, *PropTxn) (*AckTxn, error)
	Commit(context.Context, *CommitTxn) (*Empty, error)
	HeartBeat(context.Context, *Empty) (*Empty, error)
	mustEmbedUnimplementedFollowerLeaderServer()
}

// UnimplementedFollowerLeaderServer must be embedded to have forward compatible implementations.
type UnimplementedFollowerLeaderServer struct {
}

func (UnimplementedFollowerLeaderServer) Broadcast(context.Context, *PropTxn) (*AckTxn, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Broadcast not implemented")
}
func (UnimplementedFollowerLeaderServer) Commit(context.Context, *CommitTxn) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Commit not implemented")
}
func (UnimplementedFollowerLeaderServer) HeartBeat(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HeartBeat not implemented")
}
func (UnimplementedFollowerLeaderServer) mustEmbedUnimplementedFollowerLeaderServer() {}

// UnsafeFollowerLeaderServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FollowerLeaderServer will
// result in compilation errors.
type UnsafeFollowerLeaderServer interface {
	mustEmbedUnimplementedFollowerLeaderServer()
}

func RegisterFollowerLeaderServer(s grpc.ServiceRegistrar, srv FollowerLeaderServer) {
	s.RegisterService(&FollowerLeader_ServiceDesc, srv)
}

func _FollowerLeader_Broadcast_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PropTxn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FollowerLeaderServer).Broadcast(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gozab.FollowerLeader/Broadcast",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FollowerLeaderServer).Broadcast(ctx, req.(*PropTxn))
	}
	return interceptor(ctx, in, info, handler)
}

func _FollowerLeader_Commit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommitTxn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FollowerLeaderServer).Commit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gozab.FollowerLeader/Commit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FollowerLeaderServer).Commit(ctx, req.(*CommitTxn))
	}
	return interceptor(ctx, in, info, handler)
}

func _FollowerLeader_HeartBeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FollowerLeaderServer).HeartBeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gozab.FollowerLeader/HeartBeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FollowerLeaderServer).HeartBeat(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// FollowerLeader_ServiceDesc is the grpc.ServiceDesc for FollowerLeader service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FollowerLeader_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gozab.FollowerLeader",
	HandlerType: (*FollowerLeaderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Broadcast",
			Handler:    _FollowerLeader_Broadcast_Handler,
		},
		{
			MethodName: "Commit",
			Handler:    _FollowerLeader_Commit_Handler,
		},
		{
			MethodName: "HeartBeat",
			Handler:    _FollowerLeader_HeartBeat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gozab.proto",
}

// LeaderUserClient is the client API for LeaderUser service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LeaderUserClient interface {
	Store(ctx context.Context, in *Vec, opts ...grpc.CallOption) (*Empty, error)
	Retrieve(ctx context.Context, in *GetTxn, opts ...grpc.CallOption) (*ResultTxn, error)
	Identify(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Vote, error)
}

type leaderUserClient struct {
	cc grpc.ClientConnInterface
}

func NewLeaderUserClient(cc grpc.ClientConnInterface) LeaderUserClient {
	return &leaderUserClient{cc}
}

func (c *leaderUserClient) Store(ctx context.Context, in *Vec, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/gozab.LeaderUser/Store", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *leaderUserClient) Retrieve(ctx context.Context, in *GetTxn, opts ...grpc.CallOption) (*ResultTxn, error) {
	out := new(ResultTxn)
	err := c.cc.Invoke(ctx, "/gozab.LeaderUser/Retrieve", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *leaderUserClient) Identify(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Vote, error) {
	out := new(Vote)
	err := c.cc.Invoke(ctx, "/gozab.LeaderUser/Identify", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LeaderUserServer is the server API for LeaderUser service.
// All implementations must embed UnimplementedLeaderUserServer
// for forward compatibility
type LeaderUserServer interface {
	Store(context.Context, *Vec) (*Empty, error)
	Retrieve(context.Context, *GetTxn) (*ResultTxn, error)
	Identify(context.Context, *Empty) (*Vote, error)
	mustEmbedUnimplementedLeaderUserServer()
}

// UnimplementedLeaderUserServer must be embedded to have forward compatible implementations.
type UnimplementedLeaderUserServer struct {
}

func (UnimplementedLeaderUserServer) Store(context.Context, *Vec) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Store not implemented")
}
func (UnimplementedLeaderUserServer) Retrieve(context.Context, *GetTxn) (*ResultTxn, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Retrieve not implemented")
}
func (UnimplementedLeaderUserServer) Identify(context.Context, *Empty) (*Vote, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Identify not implemented")
}
func (UnimplementedLeaderUserServer) mustEmbedUnimplementedLeaderUserServer() {}

// UnsafeLeaderUserServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LeaderUserServer will
// result in compilation errors.
type UnsafeLeaderUserServer interface {
	mustEmbedUnimplementedLeaderUserServer()
}

func RegisterLeaderUserServer(s grpc.ServiceRegistrar, srv LeaderUserServer) {
	s.RegisterService(&LeaderUser_ServiceDesc, srv)
}

func _LeaderUser_Store_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Vec)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LeaderUserServer).Store(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gozab.LeaderUser/Store",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LeaderUserServer).Store(ctx, req.(*Vec))
	}
	return interceptor(ctx, in, info, handler)
}

func _LeaderUser_Retrieve_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTxn)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LeaderUserServer).Retrieve(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gozab.LeaderUser/Retrieve",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LeaderUserServer).Retrieve(ctx, req.(*GetTxn))
	}
	return interceptor(ctx, in, info, handler)
}

func _LeaderUser_Identify_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LeaderUserServer).Identify(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gozab.LeaderUser/Identify",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LeaderUserServer).Identify(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// LeaderUser_ServiceDesc is the grpc.ServiceDesc for LeaderUser service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LeaderUser_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gozab.LeaderUser",
	HandlerType: (*LeaderUserServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Store",
			Handler:    _LeaderUser_Store_Handler,
		},
		{
			MethodName: "Retrieve",
			Handler:    _LeaderUser_Retrieve_Handler,
		},
		{
			MethodName: "Identify",
			Handler:    _LeaderUser_Identify_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gozab.proto",
}

// VoterCandidateClient is the client API for VoterCandidate service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type VoterCandidateClient interface {
	AskVote(ctx context.Context, in *Epoch, opts ...grpc.CallOption) (*Vote, error)
	NewEpoch(ctx context.Context, in *Epoch, opts ...grpc.CallOption) (*EpochHist, error)
	NewLeader(ctx context.Context, in *EpochHist, opts ...grpc.CallOption) (*Vote, error)
	CommitNewLeader(ctx context.Context, in *Epoch, opts ...grpc.CallOption) (*Empty, error)
}

type voterCandidateClient struct {
	cc grpc.ClientConnInterface
}

func NewVoterCandidateClient(cc grpc.ClientConnInterface) VoterCandidateClient {
	return &voterCandidateClient{cc}
}

func (c *voterCandidateClient) AskVote(ctx context.Context, in *Epoch, opts ...grpc.CallOption) (*Vote, error) {
	out := new(Vote)
	err := c.cc.Invoke(ctx, "/gozab.VoterCandidate/AskVote", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *voterCandidateClient) NewEpoch(ctx context.Context, in *Epoch, opts ...grpc.CallOption) (*EpochHist, error) {
	out := new(EpochHist)
	err := c.cc.Invoke(ctx, "/gozab.VoterCandidate/NewEpoch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *voterCandidateClient) NewLeader(ctx context.Context, in *EpochHist, opts ...grpc.CallOption) (*Vote, error) {
	out := new(Vote)
	err := c.cc.Invoke(ctx, "/gozab.VoterCandidate/NewLeader", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *voterCandidateClient) CommitNewLeader(ctx context.Context, in *Epoch, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/gozab.VoterCandidate/CommitNewLeader", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// VoterCandidateServer is the server API for VoterCandidate service.
// All implementations must embed UnimplementedVoterCandidateServer
// for forward compatibility
type VoterCandidateServer interface {
	AskVote(context.Context, *Epoch) (*Vote, error)
	NewEpoch(context.Context, *Epoch) (*EpochHist, error)
	NewLeader(context.Context, *EpochHist) (*Vote, error)
	CommitNewLeader(context.Context, *Epoch) (*Empty, error)
	mustEmbedUnimplementedVoterCandidateServer()
}

// UnimplementedVoterCandidateServer must be embedded to have forward compatible implementations.
type UnimplementedVoterCandidateServer struct {
}

func (UnimplementedVoterCandidateServer) AskVote(context.Context, *Epoch) (*Vote, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AskVote not implemented")
}
func (UnimplementedVoterCandidateServer) NewEpoch(context.Context, *Epoch) (*EpochHist, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewEpoch not implemented")
}
func (UnimplementedVoterCandidateServer) NewLeader(context.Context, *EpochHist) (*Vote, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewLeader not implemented")
}
func (UnimplementedVoterCandidateServer) CommitNewLeader(context.Context, *Epoch) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CommitNewLeader not implemented")
}
func (UnimplementedVoterCandidateServer) mustEmbedUnimplementedVoterCandidateServer() {}

// UnsafeVoterCandidateServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to VoterCandidateServer will
// result in compilation errors.
type UnsafeVoterCandidateServer interface {
	mustEmbedUnimplementedVoterCandidateServer()
}

func RegisterVoterCandidateServer(s grpc.ServiceRegistrar, srv VoterCandidateServer) {
	s.RegisterService(&VoterCandidate_ServiceDesc, srv)
}

func _VoterCandidate_AskVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Epoch)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VoterCandidateServer).AskVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gozab.VoterCandidate/AskVote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VoterCandidateServer).AskVote(ctx, req.(*Epoch))
	}
	return interceptor(ctx, in, info, handler)
}

func _VoterCandidate_NewEpoch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Epoch)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VoterCandidateServer).NewEpoch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gozab.VoterCandidate/NewEpoch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VoterCandidateServer).NewEpoch(ctx, req.(*Epoch))
	}
	return interceptor(ctx, in, info, handler)
}

func _VoterCandidate_NewLeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EpochHist)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VoterCandidateServer).NewLeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gozab.VoterCandidate/NewLeader",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VoterCandidateServer).NewLeader(ctx, req.(*EpochHist))
	}
	return interceptor(ctx, in, info, handler)
}

func _VoterCandidate_CommitNewLeader_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Epoch)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(VoterCandidateServer).CommitNewLeader(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/gozab.VoterCandidate/CommitNewLeader",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(VoterCandidateServer).CommitNewLeader(ctx, req.(*Epoch))
	}
	return interceptor(ctx, in, info, handler)
}

// VoterCandidate_ServiceDesc is the grpc.ServiceDesc for VoterCandidate service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var VoterCandidate_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gozab.VoterCandidate",
	HandlerType: (*VoterCandidateServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AskVote",
			Handler:    _VoterCandidate_AskVote_Handler,
		},
		{
			MethodName: "NewEpoch",
			Handler:    _VoterCandidate_NewEpoch_Handler,
		},
		{
			MethodName: "NewLeader",
			Handler:    _VoterCandidate_NewLeader_Handler,
		},
		{
			MethodName: "CommitNewLeader",
			Handler:    _VoterCandidate_CommitNewLeader_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "gozab.proto",
}
