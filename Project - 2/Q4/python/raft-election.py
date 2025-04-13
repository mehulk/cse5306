#!/usr/bin/env python3
import grpc, time, threading, argparse, random
from concurrent import futures

import raft_pb2, raft_pb2_grpc
# Go stub names identical (generated inside python dir too)
import raft_pb2 as go_pb2
import raft_pb2_grpc as go_grpc

# ---------------- utility ----------------

def parse_peer_id(peer: str) -> int:
    return int(peer.split(":")[0].replace("node", ""))

class LogEntry:
    def __init__(self, op: str, term: int, idx: int):
        self.operation, self.term, self.index = op, term, idx

# ---------------- Raft node ----------------
class RaftNode(raft_pb2_grpc.RaftServicer):
    def __init__(self, node_id: int, port: int, peers: list[str]):
        self.node_id, self.port, self.peers = node_id, port, peers
        self.state = "follower"
        self.term = 0
        self.voted_for = None
        self.current_leader_id = None

        self.log: list[LogEntry] = []
        self.commit_index = 0
        self.pending_acks: dict[int,set[int]] = {}

        self.heartbeat_interval = 1.0
        self.election_timeout = random.uniform(1.5, 3.0)
        self.last_heartbeat = time.time()

        self.lock = threading.RLock()
        self.grpc_channels: dict[str, raft_pb2_grpc.RaftStub] = {}
        self.stop_event = threading.Event()

        print(f"Node {self.node_id} starting as FOLLOWER (timeout={self.election_timeout:.2f}s)")
        threading.Thread(target=self._monitor_election_timeout, daemon=True).start()

    # ---------- gRPC helpers ----------
    def _stub(self, target: str):
        if target not in self.grpc_channels:
            channel = grpc.insecure_channel(target)
            self.grpc_channels[target] = raft_pb2_grpc.RaftStub(channel)
        return self.grpc_channels[target]

    # ---------- Go replication ----------
    def _replicate_to_go(self):
        """Send *current* log to paired Go node (always).
        Each side keeps the longest log, so overwrites are safe."""
        target = f"q4node{self.node_id}:600{self.node_id}"
        stub = go_grpc.LogReplicatorStub(grpc.insecure_channel(target))
        with self.lock:
            entries = [go_pb2.LogEntry(operation=e.operation, term=e.term, index=e.index)
                       for e in self.log]
            req = go_pb2.AppendEntriesRequest(term=self.term, leaderId=self.node_id,
                                              log=entries, commitIndex=self.commit_index)
        try:
            print(f"Node {self.node_id} sends RPC ReplicateLog to Node q4node{self.node_id}")
            stub.ReplicateLog(req, timeout=1)
        except Exception as ex:
            print(f"Node {self.node_id} replicate‑to‑Go error: {ex}")

    # ---------- Raft RPCs ----------
    def AppendEntries(self, request, context):
        print(f"Node {self.node_id} runs RPC AppendEntries called by Node {request.leaderId}")
        with self.lock:
            if request.term < self.term:
                return raft_pb2.AppendEntriesResponse(success=False, nodeId=self.node_id)

            # step‑down if necessary
            self.state = "follower"
            self.term = request.term
            self.current_leader_id = request.leaderId
            self.last_heartbeat = time.time()

            # merge log
            if len(request.log) > len(self.log):
                self.log = [LogEntry(e.operation, e.term, e.index) for e in request.log]
            if request.commitIndex > self.commit_index:
                self.commit_index = request.commitIndex

        # forward to Go side
        self._replicate_to_go()
        return raft_pb2.AppendEntriesResponse(success=True, nodeId=self.node_id)

    def RequestVote(self, request, context):
        print(f"Node {self.node_id} runs RPC RequestVote called by Node {request.candidateId}")
        with self.lock:
            if request.term > self.term:
                self.term = request.term
                self.voted_for = None
                self.state = "follower"

            grant = (request.term == self.term and
                      (self.voted_for in (None, request.candidateId)))
            if grant:
                self.voted_for = request.candidateId
        return raft_pb2.VoteResponse(voteGranted=grant, term=self.term)

    def SubmitClientRequest(self, request, context):
        print(f"Node {self.node_id} runs RPC SubmitClientRequest called by a client.")
        with self.lock:
            if self.state != "leader":
                if self.current_leader_id:
                    # forward
                    leader = f"node{self.current_leader_id}:{5000+self.current_leader_id}"
                    print(f"Node {self.node_id} forwards client req to Leader {self.current_leader_id}")
                    try:
                        resp = self._stub(leader).SubmitClientRequest(request, timeout=2)
                        return resp
                    except Exception as ex:
                        return raft_pb2.ClientResponse(result=f"forward failed: {ex}")
                return raft_pb2.ClientResponse(result="Node is not the leader yet")

            # we *are* leader
            idx = len(self.log) + 1
            self.log.append(LogEntry(request.operation, self.term, idx))
            self.pending_acks[idx] = {self.node_id}

        # broadcast the single new entry for ACKs
        self._broadcast_append_one(idx)

        with self.lock:
            maj = len(self.peers)//2 + 1
            committed = len(self.pending_acks[idx]) >= maj
            if committed:
                self.commit_index = idx
                del self.pending_acks[idx]
                # notify Go immediately
                self._replicate_to_go()
                return raft_pb2.ClientResponse(result=f"Operation '{request.operation}' committed at index {idx}")
            else:
                return raft_pb2.ClientResponse(result="Failed to reach majority")

    # ---------- internal helpers ----------
    def _broadcast_append_one(self, entry_idx: int):
        with self.lock:
            entry = self.log[entry_idx-1]
            req = raft_pb2.AppendEntriesRequest(term=self.term, leaderId=self.node_id,
                                                log=[raft_pb2.LogEntry(operation=entry.operation,
                                                                       term=entry.term, index=entry.index)],
                                                commitIndex=self.commit_index)
        def _worker(peer: str):
            peer_id = parse_peer_id(peer)
            print(f"Node {self.node_id} sends RPC AppendEntries(single) to Node {peer_id}")
            try:
                resp = self._stub(peer).AppendEntries(req, timeout=1)
                if resp.success:
                    with self.lock:
                        self.pending_acks[entry_idx].add(peer_id)
            except Exception:
                pass
        threads=[threading.Thread(target=_worker,args=(p,)) for p in self.peers]
        for t in threads: t.start()
        for t in threads: t.join(1)

    # ---------- heartbeat / election ----------
    def _broadcast_heartbeat_loop(self):
        while not self.stop_event.is_set():
            time.sleep(self.heartbeat_interval)
            if self.state != "leader":
                continue
            with self.lock:
                req = raft_pb2.AppendEntriesRequest(term=self.term, leaderId=self.node_id,
                                                    log=[raft_pb2.LogEntry(operation=e.operation, term=e.term, index=e.index) for e in self.log],
                                                    commitIndex=self.commit_index)
            for peer in self.peers:
                peer_id = parse_peer_id(peer)
                print(f"Leader Node {self.node_id} sends AppendEntries (Heartbeat) to Node {peer_id}")
                try:
                    self._stub(peer).AppendEntries(req, timeout=1)
                except Exception:
                    pass
            # also replicate to Go
            self._replicate_to_go()

    def _start_heartbeat(self):
        threading.Thread(target=self._broadcast_heartbeat_loop, daemon=True).start()

    def _monitor_election_timeout(self):
        while not self.stop_event.is_set():
            time.sleep(0.1)
            if self.state == "leader":
                continue
            if time.time() - self.last_heartbeat >= self.election_timeout:
                self._start_election()

    def _start_election(self):
        with self.lock:
            self.state = "candidate"
            self.term += 1
            self.last_heartbeat = time.time()
            votes = 1
            req = raft_pb2.VoteRequest(candidateId=self.node_id, term=self.term)
        for peer in self.peers:
            peer_id = parse_peer_id(peer)
            print(f"Node {self.node_id} sends RPC RequestVote to Node {peer_id}")
            try:
                resp = self._stub(peer).RequestVote(req, timeout=1)
                if resp.voteGranted:
                    votes += 1
            except Exception:
                pass
        if votes >= len(self.peers)//2 + 1:
            with self.lock:
                self.state = "leader"
                self.current_leader_id = self.node_id
            print(f"Node {self.node_id} elected LEADER for term {self.term}")
            self._start_heartbeat()
        else:
            with self.lock:
                self.state = "follower"

# ---------- gRPC server wrapper ----------

def serve(node: RaftNode):
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    raft_pb2_grpc.add_RaftServicer_to_server(node, server)
    server.add_insecure_port(f"[::]:{node.port}")
    server.start()
    server.wait_for_termination()

# ---------- entry ----------
if __name__ == "__main__":
    ap = argparse.ArgumentParser()
    ap.add_argument("-id", type=int, required=True)
    ap.add_argument("-port", type=int, required=True)
    ap.add_argument("-peers", required=True)
    args = ap.parse_args()
    node = RaftNode(args.id, args.port, args.peers.split(","))
    serve(node)