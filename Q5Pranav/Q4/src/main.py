import os
import time
import grpc
from concurrent import futures
import threading

import swim_pb2_grpc
from node import Node
from failure_detector import FailureDetector

class NodeInfo:
    def __init__(self, node_id, host, port):
        self.node_id = node_id
        self.host = host
        self.port = port
        self.last_heard_from = 0

def serve(node_id, membership_list):
    """
    Starts the gRPC server for the Python Failure Detector (Q4).
    Binds on port 50050 + node_id.
    """
    # 1. Create Node object and attach to gRPC server
    node = Node(node_id, membership_list)
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    swim_pb2_grpc.add_SwimServiceServicer_to_server(node, server)

    # 2. Bind to 0.0.0.0 on port (50050 + node_id)
    fd_port = 50050 + node_id
    server.add_insecure_port(f'0.0.0.0:{fd_port}')
    server.start()

    # 3. Start Failure Detector in the background
    failure_detector = FailureDetector(node_id, membership_list)
    threading.Thread(target=failure_detector.run, daemon=True).start()

    # 4. Keep server alive
    server.wait_for_termination()

if __name__ == '__main__':
    # 1. Define membership list for all nodes
    membership_list = [
        NodeInfo(1, 'localhost', 50051),
        NodeInfo(2, 'localhost', 50052),
        NodeInfo(3, 'localhost', 50053),
        NodeInfo(4, 'localhost', 50054),
        NodeInfo(5, 'localhost', 50055)
    ]

    # 2. Figure out this node’s ID from environment
    node_id = int(os.environ.get('NODE_ID', 1))

    # 3. If this node is not the bootstrap node, wait 10 seconds
    #    so node 1 has time to fully start
    if node_id != 1:
        print(f"Node {node_id} sleeping 10s to let bootstrap node 1 come up...")
        time.sleep(10)

    # 4. Start the server on (50050 + node_id)
    print(f"Node {node_id} is starting Python Failure Detector on port {50050 + node_id}...")
    serve(node_id, membership_list)
