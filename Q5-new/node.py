import grpc
from concurrent import futures
import swim_pb2
import swim_pb2_grpc
import time
import os

FD_BASE = 50050  # FD port = FD_BASE + node_id

class FDNode(swim_pb2_grpc.SwimServiceServicer):
    def __init__(self, node_id):
        self.node_id = node_id

    def Ping(self, request, context):
        print(f"Component FailureDetector of Node {self.node_id} runs RPC Ping called by Node {request.sender_id}", flush=True)
        return swim_pb2.PingResponse(sender_id=self.node_id)

    def IndirectPing(self, request, context):
        print(f"Component FailureDetector of Node {self.node_id} runs RPC IndirectPing called by Node {request.sender_id} for target {request.target_id}", flush=True)
        # Here, for simplicity, try to ping the target node using its FD port.
        target_fd_port = FD_BASE + request.target_id
        target_host = f"node{request.target_id}"  # Using Docker DNS naming
        try:
            with grpc.insecure_channel(f"{target_host}:{target_fd_port}") as channel:
                stub = swim_pb2_grpc.SwimServiceStub(channel)
                stub.Ping(swim_pb2.PingRequest(sender_id=self.node_id))
            return swim_pb2.IndirectPingResponse(sender_id=self.node_id, success=True)
        except Exception as e:
            return swim_pb2.IndirectPingResponse(sender_id=self.node_id, success=False)

def serve_fd():
    node_id = int(os.environ.get("NODE_ID", "1"))
    fd_port = FD_BASE + node_id
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    swim_pb2_grpc.add_SwimServiceServicer_to_server(FDNode(node_id), server)
    server.add_insecure_port(f"0.0.0.0:{fd_port}")
    print(f"Failure Detector (FD) for Node {node_id} listening on port {fd_port}", flush=True)
    server.start()
    while True:
        time.sleep(3600)

if __name__ == '__main__':
    serve_fd()
