#!/usr/bin/env python3
"""
Raft smart test‑client (Python ≤ 3.9 version)
---------------------------------------------
• Finds the current leader by probing the cluster.
• Sends subsequent client requests only to that leader.
• Re‑discovers the leader automatically if it changes.
"""

import time
from typing import Optional, Tuple

import grpc
import raft_pb2
import raft_pb2_grpc

# --------------------------------------------------------------------------- #
# Configuration
# --------------------------------------------------------------------------- #
NODES = [
    "node1:5001",
    "node2:5002",
    "node3:5003",
    "node4:5004",
    "node5:5005",
]

PROBE_OP = "__probe__"          # harmless operation to test leadership
SEND_INTERVAL = 5               # seconds between real client requests
DISCOVERY_TIMEOUT = 2           # seconds to wait for a probe reply
REQUEST_TIMEOUT = 5             # seconds to wait for a real reply


# --------------------------------------------------------------------------- #
# Helper functions
# --------------------------------------------------------------------------- #
def parse_leader_addr(result_str: str) -> Optional[str]:
    """
    If a follower’s reply embeds the leader’s address, extract and return it.
    """
    tokens = result_str.replace(",", " ").split()
    for tok in tokens:
        if ":" in tok and tok.startswith("node"):
            return tok
    return None


def probe_node(addr: str) -> Tuple[bool, Optional[str]]:
    """
    Try to submit a harmless request to *addr*.

    Returns (is_leader, leader_addr_if_known).
    """
    try:
        with grpc.insecure_channel(addr) as channel:
            stub = raft_pb2_grpc.RaftStub(channel)
            req = raft_pb2.ClientRequest(operation=PROBE_OP)
            resp = stub.SubmitClientRequest(req, timeout=DISCOVERY_TIMEOUT)

            if "is not the leader" not in resp.result:
                return True, addr

            hinted = parse_leader_addr(resp.result)
            return False, hinted
    except Exception:
        return False, None


def discover_leader() -> Optional[str]:
    """
    Ping every node until the leader is discovered or all probes fail.
    """
    print("[Client] Discovering leader …")
    for addr in NODES:
        is_leader, hinted = probe_node(addr)
        if is_leader:
            print(f"[Client] Leader is {addr}")
            return addr
        if hinted:
            print(f"[Client] Follower {addr} pointed us to {hinted}")
            return hinted

    # second pass in case only one follower knows the leader
    for addr in NODES:
        is_leader, _ = probe_node(addr)
        if is_leader:
            print(f"[Client] Leader is {addr}")
            return addr

    print("[Client] No leader reachable at the moment.")
    return None


# --------------------------------------------------------------------------- #
# Main loop
# --------------------------------------------------------------------------- #
def main() -> None:
    print("Waiting 10 s for the cluster to elect a leader …")
    time.sleep(10)

    leader: Optional[str] = discover_leader()
    op_index = 0

    while True:
        if leader is None:
            time.sleep(SEND_INTERVAL)
            leader = discover_leader()
            continue

        op_str = f"testKey{op_index}=testVal{op_index}"
        op_index += 1
        req = raft_pb2.ClientRequest(operation=op_str)

        print(f"[Client] → {leader}  SubmitClientRequest({op_str})")
        try:
            with grpc.insecure_channel(leader) as channel:
                stub = raft_pb2_grpc.RaftStub(channel)
                resp = stub.SubmitClientRequest(req, timeout=REQUEST_TIMEOUT)

            if "is not the leader" in resp.result:
                print(f"[Client] {leader} is no longer leader: {resp.result}")
                hinted = parse_leader_addr(resp.result)
                leader = hinted if hinted else None
            else:
                print(f"[Client] ✓  Response from leader {leader}: {resp.result}")

        except Exception as ex:
            print(f"[Client] ✗  Error talking to {leader}: {ex}")
            leader = None  # trigger re‑discovery

        time.sleep(SEND_INTERVAL)


if __name__ == "__main__":
    main()
