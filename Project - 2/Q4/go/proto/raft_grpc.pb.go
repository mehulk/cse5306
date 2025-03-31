// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: raft.proto

package raft

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Raft_AppendEntries_FullMethodName       = "/raft.Raft/AppendEntries"
	Raft_RequestVote_FullMethodName         = "/raft.Raft/RequestVote"
	Raft_SubmitClientRequest_FullMethodName = "/raft.Raft/SubmitClientRequest"
)

// RaftClient is the client API for Raft service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// For Q3’s existing service
type RaftClient interface {
	AppendEntries(ctx context.Context, in *AppendEntriesRequest, opts ...grpc.CallOption) (*AppendEntriesResponse, error)
	RequestVote(ctx context.Context, in *VoteRequest, opts ...grpc.CallOption) (*VoteResponse, error)
	SubmitClientRequest(ctx context.Context, in *ClientRequest, opts ...grpc.CallOption) (*ClientResponse, error)
}

type raftClient struct {
	cc grpc.ClientConnInterface
}

func NewRaftClient(cc grpc.ClientConnInterface) RaftClient {
	return &raftClient{cc}
}

func (c *raftClient) AppendEntries(ctx context.Context, in *AppendEntriesRequest, opts ...grpc.CallOption) (*AppendEntriesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AppendEntriesResponse)
	err := c.cc.Invoke(ctx, Raft_AppendEntries_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftClient) RequestVote(ctx context.Context, in *VoteRequest, opts ...grpc.CallOption) (*VoteResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(VoteResponse)
	err := c.cc.Invoke(ctx, Raft_RequestVote_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *raftClient) SubmitClientRequest(ctx context.Context, in *ClientRequest, opts ...grpc.CallOption) (*ClientResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ClientResponse)
	err := c.cc.Invoke(ctx, Raft_SubmitClientRequest_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RaftServer is the server API for Raft service.
// All implementations must embed UnimplementedRaftServer
// for forward compatibility.
//
// For Q3’s existing service
type RaftServer interface {
	AppendEntries(context.Context, *AppendEntriesRequest) (*AppendEntriesResponse, error)
	RequestVote(context.Context, *VoteRequest) (*VoteResponse, error)
	SubmitClientRequest(context.Context, *ClientRequest) (*ClientResponse, error)
	mustEmbedUnimplementedRaftServer()
}

// UnimplementedRaftServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedRaftServer struct{}

func (UnimplementedRaftServer) AppendEntries(context.Context, *AppendEntriesRequest) (*AppendEntriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AppendEntries not implemented")
}
func (UnimplementedRaftServer) RequestVote(context.Context, *VoteRequest) (*VoteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RequestVote not implemented")
}
func (UnimplementedRaftServer) SubmitClientRequest(context.Context, *ClientRequest) (*ClientResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubmitClientRequest not implemented")
}
func (UnimplementedRaftServer) mustEmbedUnimplementedRaftServer() {}
func (UnimplementedRaftServer) testEmbeddedByValue()              {}

// UnsafeRaftServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RaftServer will
// result in compilation errors.
type UnsafeRaftServer interface {
	mustEmbedUnimplementedRaftServer()
}

func RegisterRaftServer(s grpc.ServiceRegistrar, srv RaftServer) {
	// If the following call pancis, it indicates UnimplementedRaftServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Raft_ServiceDesc, srv)
}

func _Raft_AppendEntries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppendEntriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftServer).AppendEntries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Raft_AppendEntries_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftServer).AppendEntries(ctx, req.(*AppendEntriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Raft_RequestVote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftServer).RequestVote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Raft_RequestVote_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftServer).RequestVote(ctx, req.(*VoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Raft_SubmitClientRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClientRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RaftServer).SubmitClientRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Raft_SubmitClientRequest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RaftServer).SubmitClientRequest(ctx, req.(*ClientRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Raft_ServiceDesc is the grpc.ServiceDesc for Raft service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Raft_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "raft.Raft",
	HandlerType: (*RaftServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AppendEntries",
			Handler:    _Raft_AppendEntries_Handler,
		},
		{
			MethodName: "RequestVote",
			Handler:    _Raft_RequestVote_Handler,
		},
		{
			MethodName: "SubmitClientRequest",
			Handler:    _Raft_SubmitClientRequest_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "raft.proto",
}

const (
	LogReplicator_ReplicateLog_FullMethodName   = "/raft.LogReplicator/ReplicateLog"
	LogReplicator_AckReplication_FullMethodName = "/raft.LogReplicator/AckReplication"
)

// LogReplicatorClient is the client API for LogReplicator service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// For Q4’s log replication in Go
// Alternatively, you can reuse the 'Raft' service itself, or add
// specialized RPC calls for replication. Example below:
type LogReplicatorClient interface {
	// Node -> LogReplicator server to store an incoming log entry
	ReplicateLog(ctx context.Context, in *AppendEntriesRequest, opts ...grpc.CallOption) (*AppendEntriesResponse, error)
	// Possibly a separate RPC for acknowledging or retrieving commit states, etc.
	AckReplication(ctx context.Context, in *AppendEntriesResponse, opts ...grpc.CallOption) (*ClientResponse, error)
}

type logReplicatorClient struct {
	cc grpc.ClientConnInterface
}

func NewLogReplicatorClient(cc grpc.ClientConnInterface) LogReplicatorClient {
	return &logReplicatorClient{cc}
}

func (c *logReplicatorClient) ReplicateLog(ctx context.Context, in *AppendEntriesRequest, opts ...grpc.CallOption) (*AppendEntriesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AppendEntriesResponse)
	err := c.cc.Invoke(ctx, LogReplicator_ReplicateLog_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logReplicatorClient) AckReplication(ctx context.Context, in *AppendEntriesResponse, opts ...grpc.CallOption) (*ClientResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ClientResponse)
	err := c.cc.Invoke(ctx, LogReplicator_AckReplication_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LogReplicatorServer is the server API for LogReplicator service.
// All implementations must embed UnimplementedLogReplicatorServer
// for forward compatibility.
//
// For Q4’s log replication in Go
// Alternatively, you can reuse the 'Raft' service itself, or add
// specialized RPC calls for replication. Example below:
type LogReplicatorServer interface {
	// Node -> LogReplicator server to store an incoming log entry
	ReplicateLog(context.Context, *AppendEntriesRequest) (*AppendEntriesResponse, error)
	// Possibly a separate RPC for acknowledging or retrieving commit states, etc.
	AckReplication(context.Context, *AppendEntriesResponse) (*ClientResponse, error)
	mustEmbedUnimplementedLogReplicatorServer()
}

// UnimplementedLogReplicatorServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedLogReplicatorServer struct{}

func (UnimplementedLogReplicatorServer) ReplicateLog(context.Context, *AppendEntriesRequest) (*AppendEntriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReplicateLog not implemented")
}
func (UnimplementedLogReplicatorServer) AckReplication(context.Context, *AppendEntriesResponse) (*ClientResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AckReplication not implemented")
}
func (UnimplementedLogReplicatorServer) mustEmbedUnimplementedLogReplicatorServer() {}
func (UnimplementedLogReplicatorServer) testEmbeddedByValue()                       {}

// UnsafeLogReplicatorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LogReplicatorServer will
// result in compilation errors.
type UnsafeLogReplicatorServer interface {
	mustEmbedUnimplementedLogReplicatorServer()
}

func RegisterLogReplicatorServer(s grpc.ServiceRegistrar, srv LogReplicatorServer) {
	// If the following call pancis, it indicates UnimplementedLogReplicatorServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&LogReplicator_ServiceDesc, srv)
}

func _LogReplicator_ReplicateLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppendEntriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogReplicatorServer).ReplicateLog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LogReplicator_ReplicateLog_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogReplicatorServer).ReplicateLog(ctx, req.(*AppendEntriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogReplicator_AckReplication_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AppendEntriesResponse)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogReplicatorServer).AckReplication(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LogReplicator_AckReplication_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogReplicatorServer).AckReplication(ctx, req.(*AppendEntriesResponse))
	}
	return interceptor(ctx, in, info, handler)
}

// LogReplicator_ServiceDesc is the grpc.ServiceDesc for LogReplicator service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LogReplicator_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "raft.LogReplicator",
	HandlerType: (*LogReplicatorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ReplicateLog",
			Handler:    _LogReplicator_ReplicateLog_Handler,
		},
		{
			MethodName: "AckReplication",
			Handler:    _LogReplicator_AckReplication_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "raft.proto",
}
