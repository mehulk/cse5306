import grpc
from concurrent import futures
import swim_pb2_grpc
import threading
import os
import time

from node import Node
from failure_detector import FailureDetector

class NodeInfo:
    def __init__(self, node_id, host, port):
        self.node_id = node_id
        self.host = host
        self.port = port
        self.last_heard_from = 0

def serve(node_id, host, port, membership_list):
    # Start a gRPC server for the Python FD + stubs for dissemination
    node = Node(node_id, membership_list)
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    swim_pb2_grpc.add_SwimServiceServicer_to_server(node, server)
    server.add_insecure_port(f'{host}:{port}')
    server.start()
    
    # Start the Failure Detector in a background thread
    failure_detector = FailureDetector(node_id, membership_list)
    threading.Thread(target=failure_detector.run, daemon=True).start()
    
    server.wait_for_termination()

if __name__ == '__main__':
    # Hard-code or environment for membership
    membership_list = [
        NodeInfo(1, 'node1', 50051),
        NodeInfo(2, 'node2', 50052),
        NodeInfo(3, 'node3', 50053),
        NodeInfo(4, 'node4', 50054),
        NodeInfo(5, 'node5', 50055),
    ]
    
    node_id = int(os.environ.get('NODE_ID', 1))
    node_info = next(node for node in membership_list if node.node_id == node_id)

    # Start the Python FD server
    # (In practice, we might want the Python container to use the same container name, or 
    #  rely on docker-compose's DNS. Adjust as needed.)
    serve(node_id, node_info.host, node_info.port, membership_list)
