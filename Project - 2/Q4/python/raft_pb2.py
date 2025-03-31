# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# NO CHECKED-IN PROTOBUF GENCODE
# source: raft.proto
# Protobuf Python Version: 5.29.0
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import descriptor_pool as _descriptor_pool
from google.protobuf import runtime_version as _runtime_version
from google.protobuf import symbol_database as _symbol_database
from google.protobuf.internal import builder as _builder
_runtime_version.ValidateProtobufRuntimeVersion(
    _runtime_version.Domain.PUBLIC,
    5,
    29,
    0,
    '',
    'raft.proto'
)
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor_pool.Default().AddSerializedFile(b'\n\nraft.proto\x12\x04raft\":\n\x08LogEntry\x12\x11\n\toperation\x18\x01 \x01(\t\x12\x0c\n\x04term\x18\x02 \x01(\x05\x12\r\n\x05index\x18\x03 \x01(\x05\"h\n\x14\x41ppendEntriesRequest\x12\x0c\n\x04term\x18\x01 \x01(\x05\x12\x10\n\x08leaderId\x18\x02 \x01(\x05\x12\x1b\n\x03log\x18\x03 \x03(\x0b\x32\x0e.raft.LogEntry\x12\x13\n\x0b\x63ommitIndex\x18\x04 \x01(\x05\"8\n\x15\x41ppendEntriesResponse\x12\x0f\n\x07success\x18\x01 \x01(\x08\x12\x0e\n\x06nodeId\x18\x02 \x01(\x05\"0\n\x0bVoteRequest\x12\x13\n\x0b\x63\x61ndidateId\x18\x01 \x01(\x05\x12\x0c\n\x04term\x18\x02 \x01(\x05\"1\n\x0cVoteResponse\x12\x13\n\x0bvoteGranted\x18\x01 \x01(\x08\x12\x0c\n\x04term\x18\x02 \x01(\x05\"\"\n\rClientRequest\x12\x11\n\toperation\x18\x01 \x01(\t\" \n\x0e\x43lientResponse\x12\x0e\n\x06result\x18\x01 \x01(\t2\xc8\x01\n\x04Raft\x12H\n\rAppendEntries\x12\x1a.raft.AppendEntriesRequest\x1a\x1b.raft.AppendEntriesResponse\x12\x34\n\x0bRequestVote\x12\x11.raft.VoteRequest\x1a\x12.raft.VoteResponse\x12@\n\x13SubmitClientRequest\x12\x13.raft.ClientRequest\x1a\x14.raft.ClientResponse2\x9d\x01\n\rLogReplicator\x12G\n\x0cReplicateLog\x12\x1a.raft.AppendEntriesRequest\x1a\x1b.raft.AppendEntriesResponse\x12\x43\n\x0e\x41\x63kReplication\x12\x1b.raft.AppendEntriesResponse\x1a\x14.raft.ClientResponseB\x11Z\x0fproto/raft;raftb\x06proto3')

_globals = globals()
_builder.BuildMessageAndEnumDescriptors(DESCRIPTOR, _globals)
_builder.BuildTopDescriptorsAndMessages(DESCRIPTOR, 'raft_pb2', _globals)
if not _descriptor._USE_C_DESCRIPTORS:
  _globals['DESCRIPTOR']._loaded_options = None
  _globals['DESCRIPTOR']._serialized_options = b'Z\017proto/raft;raft'
  _globals['_LOGENTRY']._serialized_start=20
  _globals['_LOGENTRY']._serialized_end=78
  _globals['_APPENDENTRIESREQUEST']._serialized_start=80
  _globals['_APPENDENTRIESREQUEST']._serialized_end=184
  _globals['_APPENDENTRIESRESPONSE']._serialized_start=186
  _globals['_APPENDENTRIESRESPONSE']._serialized_end=242
  _globals['_VOTEREQUEST']._serialized_start=244
  _globals['_VOTEREQUEST']._serialized_end=292
  _globals['_VOTERESPONSE']._serialized_start=294
  _globals['_VOTERESPONSE']._serialized_end=343
  _globals['_CLIENTREQUEST']._serialized_start=345
  _globals['_CLIENTREQUEST']._serialized_end=379
  _globals['_CLIENTRESPONSE']._serialized_start=381
  _globals['_CLIENTRESPONSE']._serialized_end=413
  _globals['_RAFT']._serialized_start=416
  _globals['_RAFT']._serialized_end=616
  _globals['_LOGREPLICATOR']._serialized_start=619
  _globals['_LOGREPLICATOR']._serialized_end=776
# @@protoc_insertion_point(module_scope)
