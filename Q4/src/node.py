import grpc
from concurrent import futures
import swim_pb2
import swim_pb2_grpc

class Node(swim_pb2_grpc.SwimServiceServicer):
    def __init__(self, node_id, membership_list):
        self.node_id = node_id
        self.membership_list = membership_list

    def Ping(self, request, context):
        print(f"Component FailureDetector of Node {self.node_id} runs RPC Ping called by Component FailureDetector of Node {request.sender_id}",flush=True)
        return swim_pb2.PingResponse(sender_id=self.node_id)

    def IndirectPing(self, request, context):
        print(f"Component FailureDetector of Node {self.node_id} runs RPC IndirectPing called by Component FailureDetector of Node {request.sender_id}",flush=True)
        target_node = next((node for node in self.membership_list if node.node_id == request.target_id), None)
        if target_node:
            try:
                with grpc.insecure_channel(f'{target_node.host}:{target_node.port}') as channel:
                    stub = swim_pb2_grpc.SwimServiceStub(channel)
                    response = stub.Ping(swim_pb2.PingRequest(sender_id=self.node_id))
                    return swim_pb2.IndirectPingResponse(sender_id=self.node_id, success=True)
            except grpc.RpcError:
                return swim_pb2.IndirectPingResponse(sender_id=self.node_id, success=False)
        return swim_pb2.IndirectPingResponse(sender_id=self.node_id, success=False)
