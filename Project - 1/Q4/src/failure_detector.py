import random
import time
import grpc
import swim_pb2
import swim_pb2_grpc

class FailureDetector:
    def __init__(self, node_id, membership_list, T=5, k=3):
        self.node_id = node_id
        self.membership_list = membership_list
        self.T = T
        self.k = k

    def run(self):
        while True:
            time.sleep(self.T)
            target_node = random.choice(self.membership_list)
            if target_node.node_id != self.node_id:
                self.ping_node(target_node)

    def ping_node(self, target_node):
        print(f"Component FailureDetector of Node {self.node_id} sends RPC Ping to Component FailureDetector of Node {target_node.node_id}",flush=True)
        try:
            with grpc.insecure_channel(f'{target_node.host}:{target_node.port}') as channel:
                stub = swim_pb2_grpc.SwimServiceStub(channel)
                response = stub.Ping(swim_pb2.PingRequest(sender_id=self.node_id))
                target_node.last_heard_from = time.time()
            print(f"Node {target_node.node_id} is healthy",flush=True)
        except grpc.RpcError:
            self.indirect_ping(target_node)

    def indirect_ping(self, target_node):
        other_nodes = [node for node in self.membership_list if node.node_id not in [self.node_id, target_node.node_id]]
        ping_nodes = random.sample(other_nodes, min(self.k, len(other_nodes)))
        
        for ping_node in ping_nodes:
            print(f"Component FailureDetector of Node {self.node_id} sends RPC IndirectPing to Component FailureDetector of Node {ping_node.node_id}",flush=True)
            try:
                with grpc.insecure_channel(f'{ping_node.host}:{ping_node.port}') as channel:
                    stub = swim_pb2_grpc.SwimServiceStub(channel)
                    response = stub.IndirectPing(swim_pb2.IndirectPingRequest(sender_id=self.node_id, target_id=target_node.node_id))
                    if response.success:
                        target_node.last_heard_from = time.time()
                        print(f"Node {target_node.node_id} is healthy (indirect ping)",flush=True)
                        return
            except grpc.RpcError:
                continue

        print(f"Node {target_node.node_id} has failed",flush=True)
        self.membership_list.remove(target_node)
