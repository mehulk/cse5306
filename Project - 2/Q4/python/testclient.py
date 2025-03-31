import time
import grpc
import raft_pb2
import raft_pb2_grpc

def run_client():
    print("Starting the test client. Will send requests every 10 seconds.")
    endpoints = [
        "node1:5001",
        "node2:5002",
        "node3:5003",
        "node4:5004",
        "node5:5005"
    ]

    # Wait for cluster to settle
    time.sleep(10)

    operation_count = 0
    while True:
        operation_count += 1
        operation_str = f"Test Operation #{operation_count}"
        request = raft_pb2.ClientRequest(operation=operation_str)

        # We'll try each endpoint until we find the leader
        leader_found = False
        for ep in endpoints:
            print(f"[TestClient] sends RPC SubmitClientRequest to {ep} (operation={operation_str})")
            try:
                channel = grpc.insecure_channel(ep)
                stub = raft_pb2_grpc.RaftStub(channel)
                response = stub.SubmitClientRequest(request)

                if "is not the leader" not in response.result:
                    print(f"[TestClient] Leader found at {ep}. Response: {response.result}")
                    leader_found = True
                    break
                else:
                    print(f"[TestClient] {ep} is not leader: {response.result}")
            except Exception as e:
                print(f"[TestClient] Error connecting to {ep}: {e}")

        if not leader_found:
            print("[TestClient] No leader found or request failed this round.")

        time.sleep(10)  # Wait 10 seconds before sending next request

if __name__ == "__main__":
    run_client()
