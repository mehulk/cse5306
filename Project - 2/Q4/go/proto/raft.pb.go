// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: raft.proto

package raft

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type LogEntry struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Operation     string                 `protobuf:"bytes,1,opt,name=operation,proto3" json:"operation,omitempty"`
	Term          int32                  `protobuf:"varint,2,opt,name=term,proto3" json:"term,omitempty"`
	Index         int32                  `protobuf:"varint,3,opt,name=index,proto3" json:"index,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *LogEntry) Reset() {
	*x = LogEntry{}
	mi := &file_raft_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *LogEntry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogEntry) ProtoMessage() {}

func (x *LogEntry) ProtoReflect() protoreflect.Message {
	mi := &file_raft_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogEntry.ProtoReflect.Descriptor instead.
func (*LogEntry) Descriptor() ([]byte, []int) {
	return file_raft_proto_rawDescGZIP(), []int{0}
}

func (x *LogEntry) GetOperation() string {
	if x != nil {
		return x.Operation
	}
	return ""
}

func (x *LogEntry) GetTerm() int32 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *LogEntry) GetIndex() int32 {
	if x != nil {
		return x.Index
	}
	return 0
}

type AppendEntriesRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Term          int32                  `protobuf:"varint,1,opt,name=term,proto3" json:"term,omitempty"`               // Current term of the leader
	LeaderId      int32                  `protobuf:"varint,2,opt,name=leaderId,proto3" json:"leaderId,omitempty"`       // Leader's ID
	Log           []*LogEntry            `protobuf:"bytes,3,rep,name=log,proto3" json:"log,omitempty"`                  // Log entries to replicate
	CommitIndex   int32                  `protobuf:"varint,4,opt,name=commitIndex,proto3" json:"commitIndex,omitempty"` // Index of highest log entry known to be committed
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AppendEntriesRequest) Reset() {
	*x = AppendEntriesRequest{}
	mi := &file_raft_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AppendEntriesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppendEntriesRequest) ProtoMessage() {}

func (x *AppendEntriesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_raft_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppendEntriesRequest.ProtoReflect.Descriptor instead.
func (*AppendEntriesRequest) Descriptor() ([]byte, []int) {
	return file_raft_proto_rawDescGZIP(), []int{1}
}

func (x *AppendEntriesRequest) GetTerm() int32 {
	if x != nil {
		return x.Term
	}
	return 0
}

func (x *AppendEntriesRequest) GetLeaderId() int32 {
	if x != nil {
		return x.LeaderId
	}
	return 0
}

func (x *AppendEntriesRequest) GetLog() []*LogEntry {
	if x != nil {
		return x.Log
	}
	return nil
}

func (x *AppendEntriesRequest) GetCommitIndex() int32 {
	if x != nil {
		return x.CommitIndex
	}
	return 0
}

type AppendEntriesResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Success       bool                   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	NodeId        int32                  `protobuf:"varint,2,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AppendEntriesResponse) Reset() {
	*x = AppendEntriesResponse{}
	mi := &file_raft_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AppendEntriesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AppendEntriesResponse) ProtoMessage() {}

func (x *AppendEntriesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_raft_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AppendEntriesResponse.ProtoReflect.Descriptor instead.
func (*AppendEntriesResponse) Descriptor() ([]byte, []int) {
	return file_raft_proto_rawDescGZIP(), []int{2}
}

func (x *AppendEntriesResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *AppendEntriesResponse) GetNodeId() int32 {
	if x != nil {
		return x.NodeId
	}
	return 0
}

type VoteRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CandidateId   int32                  `protobuf:"varint,1,opt,name=candidateId,proto3" json:"candidateId,omitempty"`
	Term          int32                  `protobuf:"varint,2,opt,name=term,proto3" json:"term,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *VoteRequest) Reset() {
	*x = VoteRequest{}
	mi := &file_raft_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VoteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteRequest) ProtoMessage() {}

func (x *VoteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_raft_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteRequest.ProtoReflect.Descriptor instead.
func (*VoteRequest) Descriptor() ([]byte, []int) {
	return file_raft_proto_rawDescGZIP(), []int{3}
}

func (x *VoteRequest) GetCandidateId() int32 {
	if x != nil {
		return x.CandidateId
	}
	return 0
}

func (x *VoteRequest) GetTerm() int32 {
	if x != nil {
		return x.Term
	}
	return 0
}

type VoteResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	VoteGranted   bool                   `protobuf:"varint,1,opt,name=voteGranted,proto3" json:"voteGranted,omitempty"`
	Term          int32                  `protobuf:"varint,2,opt,name=term,proto3" json:"term,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *VoteResponse) Reset() {
	*x = VoteResponse{}
	mi := &file_raft_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *VoteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteResponse) ProtoMessage() {}

func (x *VoteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_raft_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteResponse.ProtoReflect.Descriptor instead.
func (*VoteResponse) Descriptor() ([]byte, []int) {
	return file_raft_proto_rawDescGZIP(), []int{4}
}

func (x *VoteResponse) GetVoteGranted() bool {
	if x != nil {
		return x.VoteGranted
	}
	return false
}

func (x *VoteResponse) GetTerm() int32 {
	if x != nil {
		return x.Term
	}
	return 0
}

type ClientRequest struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// The operation requested by the client
	Operation     string `protobuf:"bytes,1,opt,name=operation,proto3" json:"operation,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ClientRequest) Reset() {
	*x = ClientRequest{}
	mi := &file_raft_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ClientRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientRequest) ProtoMessage() {}

func (x *ClientRequest) ProtoReflect() protoreflect.Message {
	mi := &file_raft_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientRequest.ProtoReflect.Descriptor instead.
func (*ClientRequest) Descriptor() ([]byte, []int) {
	return file_raft_proto_rawDescGZIP(), []int{5}
}

func (x *ClientRequest) GetOperation() string {
	if x != nil {
		return x.Operation
	}
	return ""
}

type ClientResponse struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Result of the operation
	Result        string `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ClientResponse) Reset() {
	*x = ClientResponse{}
	mi := &file_raft_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ClientResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientResponse) ProtoMessage() {}

func (x *ClientResponse) ProtoReflect() protoreflect.Message {
	mi := &file_raft_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientResponse.ProtoReflect.Descriptor instead.
func (*ClientResponse) Descriptor() ([]byte, []int) {
	return file_raft_proto_rawDescGZIP(), []int{6}
}

func (x *ClientResponse) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

var File_raft_proto protoreflect.FileDescriptor

var file_raft_proto_rawDesc = string([]byte{
	0x0a, 0x0a, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x72, 0x61,
	0x66, 0x74, 0x22, 0x52, 0x0a, 0x08, 0x4c, 0x6f, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x1c,
	0x0a, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x65, 0x72, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x74, 0x65, 0x72, 0x6d,
	0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52,
	0x05, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x22, 0x8a, 0x01, 0x0a, 0x14, 0x41, 0x70, 0x70, 0x65, 0x6e,
	0x64, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x74,
	0x65, 0x72, 0x6d, 0x12, 0x1a, 0x0a, 0x08, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x49, 0x64, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x6c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x20, 0x0a, 0x03, 0x6c, 0x6f, 0x67, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x72,
	0x61, 0x66, 0x74, 0x2e, 0x4c, 0x6f, 0x67, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x03, 0x6c, 0x6f,
	0x67, 0x12, 0x20, 0x0a, 0x0b, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x49, 0x6e, 0x64, 0x65, 0x78,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0b, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x49, 0x6e,
	0x64, 0x65, 0x78, 0x22, 0x49, 0x0a, 0x15, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74,
	0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73,
	0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x22, 0x43,
	0x0a, 0x0b, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x20, 0x0a,
	0x0b, 0x63, 0x61, 0x6e, 0x64, 0x69, 0x64, 0x61, 0x74, 0x65, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x0b, 0x63, 0x61, 0x6e, 0x64, 0x69, 0x64, 0x61, 0x74, 0x65, 0x49, 0x64, 0x12,
	0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x74,
	0x65, 0x72, 0x6d, 0x22, 0x44, 0x0a, 0x0c, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x76, 0x6f, 0x74, 0x65, 0x47, 0x72, 0x61, 0x6e, 0x74,
	0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0b, 0x76, 0x6f, 0x74, 0x65, 0x47, 0x72,
	0x61, 0x6e, 0x74, 0x65, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x04, 0x74, 0x65, 0x72, 0x6d, 0x22, 0x2d, 0x0a, 0x0d, 0x43, 0x6c, 0x69,
	0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1c, 0x0a, 0x09, 0x6f, 0x70,
	0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6f,
	0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x28, 0x0a, 0x0e, 0x43, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x32, 0xc8, 0x01, 0x0a, 0x04, 0x52, 0x61, 0x66, 0x74, 0x12, 0x48, 0x0a, 0x0d, 0x41,
	0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x12, 0x1a, 0x2e, 0x72,
	0x61, 0x66, 0x74, 0x2e, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x2e,
	0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x34, 0x0a, 0x0b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x56, 0x6f, 0x74, 0x65, 0x12, 0x11, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x56, 0x6f, 0x74, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x12, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x56,
	0x6f, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x40, 0x0a, 0x13, 0x53,
	0x75, 0x62, 0x6d, 0x69, 0x74, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x13, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x43,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0x9d, 0x01,
	0x0a, 0x0d, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x12,
	0x47, 0x0a, 0x0c, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x65, 0x4c, 0x6f, 0x67, 0x12,
	0x1a, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74,
	0x72, 0x69, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x72, 0x61,
	0x66, 0x74, 0x2e, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x43, 0x0a, 0x0e, 0x41, 0x63, 0x6b, 0x52,
	0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x2e, 0x72, 0x61, 0x66,
	0x74, 0x2e, 0x41, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x45, 0x6e, 0x74, 0x72, 0x69, 0x65, 0x73, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x1a, 0x14, 0x2e, 0x72, 0x61, 0x66, 0x74, 0x2e, 0x43,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x11, 0x5a,
	0x0f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x61, 0x66, 0x74, 0x3b, 0x72, 0x61, 0x66, 0x74,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_raft_proto_rawDescOnce sync.Once
	file_raft_proto_rawDescData []byte
)

func file_raft_proto_rawDescGZIP() []byte {
	file_raft_proto_rawDescOnce.Do(func() {
		file_raft_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_raft_proto_rawDesc), len(file_raft_proto_rawDesc)))
	})
	return file_raft_proto_rawDescData
}

var file_raft_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_raft_proto_goTypes = []any{
	(*LogEntry)(nil),              // 0: raft.LogEntry
	(*AppendEntriesRequest)(nil),  // 1: raft.AppendEntriesRequest
	(*AppendEntriesResponse)(nil), // 2: raft.AppendEntriesResponse
	(*VoteRequest)(nil),           // 3: raft.VoteRequest
	(*VoteResponse)(nil),          // 4: raft.VoteResponse
	(*ClientRequest)(nil),         // 5: raft.ClientRequest
	(*ClientResponse)(nil),        // 6: raft.ClientResponse
}
var file_raft_proto_depIdxs = []int32{
	0, // 0: raft.AppendEntriesRequest.log:type_name -> raft.LogEntry
	1, // 1: raft.Raft.AppendEntries:input_type -> raft.AppendEntriesRequest
	3, // 2: raft.Raft.RequestVote:input_type -> raft.VoteRequest
	5, // 3: raft.Raft.SubmitClientRequest:input_type -> raft.ClientRequest
	1, // 4: raft.LogReplicator.ReplicateLog:input_type -> raft.AppendEntriesRequest
	2, // 5: raft.LogReplicator.AckReplication:input_type -> raft.AppendEntriesResponse
	2, // 6: raft.Raft.AppendEntries:output_type -> raft.AppendEntriesResponse
	4, // 7: raft.Raft.RequestVote:output_type -> raft.VoteResponse
	6, // 8: raft.Raft.SubmitClientRequest:output_type -> raft.ClientResponse
	2, // 9: raft.LogReplicator.ReplicateLog:output_type -> raft.AppendEntriesResponse
	6, // 10: raft.LogReplicator.AckReplication:output_type -> raft.ClientResponse
	6, // [6:11] is the sub-list for method output_type
	1, // [1:6] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_raft_proto_init() }
func file_raft_proto_init() {
	if File_raft_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_raft_proto_rawDesc), len(file_raft_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_raft_proto_goTypes,
		DependencyIndexes: file_raft_proto_depIdxs,
		MessageInfos:      file_raft_proto_msgTypes,
	}.Build()
	File_raft_proto = out.File
	file_raft_proto_goTypes = nil
	file_raft_proto_depIdxs = nil
}
