# import random
# import time
# import grpc
# import swim_pb2
# import swim_pb2_grpc

# class FailureDetector:
#     def __init__(self, node_id, membership_list, T=5, k=3):
#         """
#         :param node_id: ID of this node.
#         :param membership_list: List of NodeInfo objects describing all known nodes.
#         :param T: Ping interval (seconds).
#         :param k: Number of nodes to contact indirectly if direct ping fails.
#         """
#         self.node_id = node_id
#         self.membership_list = membership_list
#         self.T = T
#         self.k = k

#     def run(self):
#         """
#         Main loop for the Failure Detector. Repeatedly selects a random target node
#         to Ping, and sleeps T seconds between iterations.
#         """
#         while True:
#             time.sleep(self.T)

#             # If there's only 1 or 0 nodes in membership, no need to ping
#             if len(self.membership_list) <= 1:
#                 continue

#             target_node = random.choice(self.membership_list)
#             # Avoid self-ping
#             if target_node.node_id != self.node_id:
#                 self.ping_node(target_node)

#     def ping_node(self, target_node):
#         """
#         Attempt a direct Ping to the given target node.
#         If that fails, try indirect pings through up to k other nodes.
#         """
#         print(f"Component FailureDetector of Node {self.node_id} sends RPC Ping to "
#               f"Component FailureDetector of Node {target_node.node_id}", flush=True)

#         # Attempt direct ping
#         try:
#             with grpc.insecure_channel(f"{target_node.host}:{target_node.port}") as channel:
#                 stub = swim_pb2_grpc.SwimServiceStub(channel)
#                 stub.Ping(swim_pb2.PingRequest(sender_id=self.node_id))

#             # If successful:
#             print(f"Node {target_node.node_id} is healthy", flush=True)

#         except grpc.RpcError:
#             # If direct ping fails, attempt indirect
#             self.indirect_ping(target_node)

#     def indirect_ping(self, target_node):
#         """
#         Attempt an IndirectPing to the target_node via up to k other nodes in membership.
#         If all indirect pings fail, declare the target node as failed.
#         """
#         # Filter out ourselves and the target node from the potential 'helpers'
#         other_nodes = [node for node in self.membership_list
#                        if node.node_id not in (self.node_id, target_node.node_id)]

#         # Randomly select up to k of those for indirect ping
#         ping_nodes = random.sample(other_nodes, min(self.k, len(other_nodes)))

#         for helper_node in ping_nodes:
#             print(f"Component FailureDetector of Node {self.node_id} sends RPC IndirectPing to "
#                   f"Component FailureDetector of Node {helper_node.node_id}", flush=True)
#             try:
#                 with grpc.insecure_channel(f"{helper_node.host}:{helper_node.port}") as channel:
#                     stub = swim_pb2_grpc.SwimServiceStub(channel)
#                     response = stub.IndirectPing(swim_pb2.IndirectPingRequest(
#                         sender_id=self.node_id,
#                         target_id=target_node.node_id
#                     ))
#                     if response.success:
#                         print(f"Node {target_node.node_id} is healthy (indirect ping)", flush=True)
#                         return  # The node is alive, so stop here

#             except grpc.RpcError:
#                 # If we can’t even contact the helper, continue to the next
#                 pass

#         # If we exhausted all helpers without success, the target is likely failed
#         self.mark_node_failed(target_node)

#     def mark_node_failed(self, target_node):
#         """
#         Remove the failed node from membership and inform the local Go Dissemination
#         component that the node has failed.
#         """
#         print(f"Node {target_node.node_id} has failed", flush=True)
#         if target_node in self.membership_list:
#             self.membership_list.remove(target_node)

#         # Notify local Dissemination service so it can broadcast the failure
#         try:
#             # Assuming your Go Dissemination is listening on localhost:6000 in the same container
#             with grpc.insecure_channel("localhost:6000") as channel:
#                 stub = swim_pb2_grpc.SwimServiceStub(channel)
#                 print(f"Component FailureDetector of Node {self.node_id} sends RPC BroadcastFailure "
#                       f"to Component Dissemination of Node {self.node_id}", flush=True)
#                 stub.BroadcastFailure(
#                     swim_pb2.BroadcastFailureRequest(
#                         sender_id=self.node_id,
#                         failed_node_id=target_node.node_id
#                     )
#                 )
#         except grpc.RpcError as e:
#             print(f"Error calling BroadcastFailure on local Go server: {e}", flush=True)
import random
import time
import grpc
import swim_pb2
import swim_pb2_grpc

class FailureDetector:
    def __init__(self, node_id, membership_list, T=5, k=3):
        """
        :param node_id: ID of this node.
        :param membership_list: List of NodeInfo objects describing all known nodes.
        :param T: Ping interval (seconds).
        :param k: Number of nodes to contact indirectly if direct ping fails.
        """
        self.node_id = node_id
        self.membership_list = membership_list
        self.T = T
        self.k = k

    def run(self):
        """
        Main loop for the Failure Detector. Repeatedly selects a random target node
        to Ping, and sleeps T seconds between iterations.
        """
        while True:
            time.sleep(self.T)

            # If there's only 1 or 0 nodes in membership, no need to ping
            if len(self.membership_list) <= 1:
                continue

            target_node = random.choice(self.membership_list)
            # Avoid self-ping
            if target_node.node_id != self.node_id:
                self.ping_node(target_node)

    def ping_node(self, target_node):
        """
        Attempt a direct Ping to the given target node.
        If that fails, try indirect pings through up to k other nodes.
        """
        print(f"[PING] Node {self.node_id} → Node {target_node.node_id}", flush=True)

        # Attempt direct ping
        try:
            with grpc.insecure_channel(f"{target_node.host}:{target_node.port}") as channel:
                stub = swim_pb2_grpc.SwimServiceStub(channel)
                stub.Ping(swim_pb2.PingRequest(sender_id=self.node_id))

            # If successful:
            print(f"[HEALTHY] Node {target_node.node_id} responded", flush=True)

        except grpc.RpcError as e:
            print(f"[WARNING] Node {target_node.node_id} did not respond! Attempting indirect ping...", flush=True)
            # If direct ping fails, attempt indirect
            self.indirect_ping(target_node)

    def indirect_ping(self, target_node):
        """
        Attempt an IndirectPing to the target_node via up to k other nodes in membership.
        If all indirect pings fail, declare the target node as failed.
        """
        # Filter out ourselves and the target node from the potential 'helpers'
        other_nodes = [node for node in self.membership_list if node.node_id not in (self.node_id, target_node.node_id)]

        # Randomly select up to k of those for indirect ping
        ping_nodes = random.sample(other_nodes, min(self.k, len(other_nodes)))

        for helper_node in ping_nodes:
            print(f"[INDIRECT PING] Node {self.node_id} → Node {helper_node.node_id} → Node {target_node.node_id}", flush=True)
            try:
                with grpc.insecure_channel(f"{helper_node.host}:{helper_node.port}") as channel:
                    stub = swim_pb2_grpc.SwimServiceStub(channel)
                    response = stub.IndirectPing(swim_pb2.IndirectPingRequest(
                        sender_id=self.node_id,
                        target_id=target_node.node_id
                    ))
                    if response.success:
                        print(f"[INDIRECT SUCCESS] Node {target_node.node_id} is healthy", flush=True)
                        return  # The node is alive, so stop here

            except grpc.RpcError as e:
                print(f"[ERROR] Indirect ping via Node {helper_node.node_id} failed: {e}", flush=True)

        # If we exhausted all helpers without success, the target is likely failed
        self.mark_node_failed(target_node)

    def mark_node_failed(self, target_node):
        """
        Remove the failed node from membership and inform the local Go Dissemination
        component that the node has failed.
        """
        print(f"[FAILURE DETECTED] Node {target_node.node_id} has failed", flush=True)

        # Remove from local membership
        if target_node in self.membership_list:
            self.membership_list.remove(target_node)

        # Notify local Dissemination component to broadcast the failure
        try:
            with grpc.insecure_channel("localhost:6000") as channel:
                stub = swim_pb2_grpc.SwimServiceStub(channel)
                print(f"[BROADCAST] Node {self.node_id} → Dissemination: BroadcastFailure for Node {target_node.node_id}", flush=True)
                response = stub.BroadcastFailure(
                    swim_pb2.BroadcastFailureRequest(
                        sender_id=self.node_id,
                        failed_node_id=target_node.node_id
                    )
                )
                if response.success:
                    print(f"[BROADCAST SUCCESS] Node {target_node.node_id} failure broadcasted", flush=True)
                else:
                    print(f"[BROADCAST FAILED] Dissemination did not acknowledge failure", flush=True)

        except grpc.RpcError as e:
            print(f"[ERROR] BroadcastFailure RPC failed: {e}", flush=True)
