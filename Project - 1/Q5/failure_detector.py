import random
import time
import grpc
import swim_pb2
import swim_pb2_grpc
import os
import threading
from google.protobuf.empty_pb2 import Empty

FD_BASE = 50050  # FD port = FD_BASE + node_id (e.g., node1: 50051)
DC_BASE = 60050  # DC port = DC_BASE + node_id (e.g., node1: 60051)

class FailureDetector:
    def __init__(self, node_id, membership_list, T=5, k=3):
        """
        :param node_id: ID of this node.
        :param membership_list: List of NodeInfo objects (with DC ports) for all nodes.
        :param T: Ping interval in seconds.
        :param k: Number of nodes to contact indirectly.
        """
        self.node_id = node_id
        self.membership_list = membership_list
        self.T = T
        self.k = k

    def run(self):
        while True:
            time.sleep(self.T)
            if len(self.membership_list) <= 1:
                continue
            target_node = random.choice(self.membership_list)
            if target_node.node_id != self.node_id:
                self.ping_node(target_node)

    def ping_node(self, target_node):
        target_fd_port = FD_BASE + target_node.node_id
        print(f"[FD] Node {self.node_id}: Sending direct Ping to Node {target_node.node_id} (FD port {target_fd_port})", flush=True)
        print(f"Component FailureDetector of Node {self.node_id} sends RPC Ping to Component FailureDetector of Node {target_node.node_id}", flush=True)

        try:
            with grpc.insecure_channel(f"{target_node.host}:{target_fd_port}") as channel:
                stub = swim_pb2_grpc.SwimServiceStub(channel)
                stub.Ping(swim_pb2.PingRequest(sender_id=self.node_id))
            print(f"[FD] Direct Ping successful: Node {target_node.node_id} is healthy.", flush=True)
        except grpc.RpcError as e:
            print(f"[FD] Direct Ping to Node {target_node.node_id} failed: {e.details()}", flush=True)
            self.indirect_ping(target_node)

    def indirect_ping(self, target_node):
        other_nodes = [node for node in self.membership_list if node.node_id not in (self.node_id, target_node.node_id)]
        if not other_nodes:
            print(f"[FD] No helper nodes available for indirect ping to Node {target_node.node_id}.", flush=True)
            self.mark_node_failed(target_node)
            return
        ping_nodes = random.sample(other_nodes, min(self.k, len(other_nodes)))
        success = False
        for helper_node in ping_nodes:
            helper_fd_port = FD_BASE + helper_node.node_id
            print(f"[FD] Node {self.node_id}: Sending Indirect Ping via Node {helper_node.node_id} to check Node {target_node.node_id}", flush=True)
            print(f"Component FailureDetector of Node {self.node_id} sends RPC IndirectPing to Component FailureDetector of Node {helper_node.node_id} for target {target_node.node_id}", flush=True)

            try:
                with grpc.insecure_channel(f"{helper_node.host}:{helper_fd_port}") as channel:
                    stub = swim_pb2_grpc.SwimServiceStub(channel)
                    response = stub.IndirectPing(swim_pb2.IndirectPingRequest(
                        sender_id=self.node_id,
                        target_id=target_node.node_id
                    ))
                    if response.success:
                        print(f"[FD] Indirect Ping confirmed: Node {target_node.node_id} is healthy (via Node {helper_node.node_id}).", flush=True)
                        success = True
                        break
            except grpc.RpcError as e:
                print(f"[FD] Indirect Ping via Node {helper_node.node_id} failed for Node {target_node.node_id}: {e.details()}", flush=True)
                continue
        if not success:
            self.mark_node_failed(target_node)

    def mark_node_failed(self, target_node):
        print(f"[FD] Node {self.node_id}: Detected failure of Node {target_node.node_id}. Removing it from membership.", flush=True)
        if target_node in self.membership_list:
            self.membership_list.remove(target_node)
        print(f"[FD] Updated membership list: {[node.node_id for node in self.membership_list]}", flush=True)
        local_dc_port = DC_BASE + self.node_id
        try:
            with grpc.insecure_channel(f"localhost:{local_dc_port}") as channel:
                stub = swim_pb2_grpc.SwimServiceStub(channel)
                print(f"[FD] Node {self.node_id}: Reporting failure of Node {target_node.node_id} to local DC on port {local_dc_port}.", flush=True)
                stub.BroadcastFailure(swim_pb2.BroadcastFailureRequest(
                    sender_id=self.node_id,
                    failed_node_id=target_node.node_id
                ))
        except grpc.RpcError as e:
            print(f"[FD] Error reporting failure to local DC: {e.details()}", flush=True)

def subscribe_membership_updates(local_dc_port, fd_instance):
    channel = grpc.insecure_channel(f"localhost:{local_dc_port}")
    stub = swim_pb2_grpc.SwimServiceStub(channel)
    try:
        for update in stub.StreamMembership(Empty()):
            fd_instance.membership_list = update.membership_list
            print(f"[FD] Membership updated from DC: {[node.node_id for node in update.membership_list]}", flush=True)
    except grpc.RpcError as e:
        print(f"[FD] StreamMembership error: {e.details()}", flush=True)
