import grpc, time, itertools
import raft_pb2, raft_pb2_grpc

NODES = [f"node{i}:{5000+i}" for i in range(1,6)]
print("Starting test client – will probe cluster every 10 s …")

def submit(op: str):
    for ep in itertools.cycle(NODES):
        try:
            stub = raft_pb2_grpc.RaftStub(grpc.insecure_channel(ep))
            resp = stub.SubmitClientRequest(raft_pb2.ClientRequest(operation=op), timeout=2)
            print(f"Client sends RPC SubmitClientRequest to Node Cluster Leader.")
            if "committed" in resp.result:
                print(f"[CLIENT] Leader {ep} accepted → {resp.result}")
                return
            else:
                print(f"[CLIENT] {ep} said: {resp.result}")
        except Exception as ex:
            print(f"[CLIENT] {ep} unreachable: {ex}")

if __name__ == "__main__":
    time.sleep(10)  # allow cluster to stabilise
    counter = 0
    while True:
        counter += 1
        submit(f"operation‑{counter}")
        time.sleep(10)