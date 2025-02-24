# import threading
# import time
# import os
# from failure_detector import FailureDetector
# from node import serve_fd
# import swim_pb2

# def run_fd_client():
#     membership_str = os.environ.get("MEMBERSHIP", "")
#     membership_list = []
#     if membership_str:
#         for entry in membership_str.split(","):
#             parts = entry.split(":")
#             if len(parts) == 3:
#                 node_id = int(parts[0])
#                 host = parts[1]
#                 port = int(parts[2])
#                 # Create NodeInfo using the proto message (for DC ports)
#                 membership_list.append(swim_pb2.NodeInfo(node_id=node_id, host=host, port=port))
#     node_id = int(os.environ.get("NODE_ID", "1"))
#     fd = FailureDetector(node_id, membership_list, T=5, k=3)
#     fd.run()

# if __name__ == '__main__':
#     t = threading.Thread(target=serve_fd)
#     t.daemon = True
#     t.start()
#     time.sleep(1)  # Give FD server time to start.
#     run_fd_client()


import threading
import time
import grpc
import os
from failure_detector import FailureDetector, subscribe_membership_updates
from node import serve_fd
import swim_pb2
import swim_pb2_grpc

def run_fd_client():
    membership_str = os.environ.get("MEMBERSHIP", "")
    membership_list = []
    if membership_str:
        for entry in membership_str.split(","):
            parts = entry.split(":")
            if len(parts) == 3:
                node_id = int(parts[0])
                host = parts[1]
                port = int(parts[2])
                membership_list.append(swim_pb2.NodeInfo(node_id=node_id, host=host, port=port))
    node_id = int(os.environ.get("NODE_ID", "1"))
    fd = FailureDetector(node_id, membership_list, T=5, k=3)
    # Start subscription thread to receive membership updates from local DC.
    local_dc_port = 60050 + node_id
    sub_thread = threading.Thread(target=subscribe_membership_updates, args=(local_dc_port, fd))
    sub_thread.daemon = True
    sub_thread.start()
    fd.run()

def bootstrap_new_node():
    bootstrap_addr = os.getenv("BOOTSTRAP_ADDRESS")
    node_id = int(os.getenv("NODE_ID"))
    
    channel = grpc.insecure_channel(bootstrap_addr)
    stub = swim_pb2_grpc.SwimServiceStub(channel)
    
    try:
        response = stub.Join(swim_pb2.JoinRequest(
            sender_id=node_id,
            host=f"node{node_id}",
            port=60050 + node_id
        ))
        print("\n=== Obtained Membership List ===")
        for member in response.membership_list:
            print(f"Node {member.node_id} - {member.host}:{member.port}")
        print("================================")
    except grpc.RpcError as e:
        print(f"Bootstrap failed: {e.code()}: {e.details()}")

if __name__ == '__main__':
    if os.getenv("BOOTSTRAP_NEEDED", "false").lower() == "true":
        bootstrap_new_node()
        exit(0)
    # Start FD server in a separate thread.
    t = threading.Thread(target=serve_fd)
    t.daemon = True
    t.start()
    # Give FD server time to start.
    time.sleep(1)
    run_fd_client()
