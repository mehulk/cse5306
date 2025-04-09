import grpc
import raft_pb2
import raft_pb2_grpc
import time
import threading
import argparse
import random
from concurrent import futures

# For Q4 log replication
import raft_pb2 as go_raft_pb2
import raft_pb2_grpc as go_raft_pb2_grpc

def parse_peer_id(peer_str):
    """
    E.g. "node2:5002" -> 2
    """
    host_part = peer_str.split(":")[0]
    return int(host_part.replace("node", ""))


class LogEntry:
    def __init__(self, operation, term, index):
        self.operation = operation
        self.term = term
        self.index = index


class RaftNode(raft_pb2_grpc.RaftServicer):
    def __init__(self, node_id, port, peers):
        self.node_id = node_id
        self.state = "follower"
        self.port = port
        self.peers = peers  # e.g. ["node2:5002", "node3:5003", ...]
        self.term = 0
        self.commit_index = 0
        self.log = []  # list of LogEntry
        self.lock = threading.Lock()
        self.grpc_channels = {}
        self.stop_event = threading.Event()
        self.voted_for = None
        self.current_leader_id = None

        # For majority acknowledgment
        self.pending_acks = {}  # map entry_index -> set of node_ids that have acked

        self.election_timeout = random.uniform(1.5, 3.0)
        self.heartbeat_interval = 1.0
        self.last_heartbeat = time.time()

        print(f"Node {self.node_id} starting as FOLLOWER (timeout={self.election_timeout:.2f}s)")
        threading.Thread(target=self.monitor_election_timeout, daemon=True).start()

    # ---------------------------------------------------
    # gRPC channel helper (caches stubs)
    # ---------------------------------------------------
    def get_stub(self, peer_addr):
        """Return a cached RaftStub for the given peer address."""
        with self.lock:
            if peer_addr not in self.grpc_channels:
                channel = grpc.insecure_channel(peer_addr)
                self.grpc_channels[peer_addr] = raft_pb2_grpc.RaftStub(channel)
            return self.grpc_channels[peer_addr]

    # ---------------------------------------------------
    # Helper: replicate entire log to local Go
    # ---------------------------------------------------
    def replicate_log_in_go(self):
        """
        Called unconditionally at the end of every RPC to replicate this node's entire log to its corresponding Go node.
        """
        go_port = 6000 + self.node_id
        target = f"q4node{self.node_id}:{go_port}"
        print(f"Node {self.node_id} sends RPC ReplicateLog to Node q4node{self.node_id} (unconditional)")

        try:
            channel = grpc.insecure_channel(target)
            stub = go_raft_pb2_grpc.LogReplicatorStub(channel)

            with self.lock:
                proto_entries = [
                    go_raft_pb2.LogEntry(operation=e.operation, term=e.term, index=e.index)
                    for e in self.log
                ]
                req = go_raft_pb2.AppendEntriesRequest(
                    term=self.term,
                    leaderId=self.node_id,
                    log=proto_entries,
                    commitIndex=self.commit_index,
                )

            response = stub.ReplicateLog(req)
            print(
                f"Node {self.node_id} got ReplicateLog response from Node q4node{self.node_id}: success={response.success}"
            )
        except Exception as ex:
            print(f"Node {self.node_id} failed to replicate log in Go: {ex}")

    # ---------------------------------------------------
    # Raft RPC methods
    # ---------------------------------------------------
    def AppendEntries(self, request, context):
        # SERVER-SIDE print
        print(
            f"Node {self.node_id} runs RPC AppendEntries called by Node {request.leaderId}"
        )

        success = False
        with self.lock:
            if request.term >= self.term:
                self.current_leader_id = request.leaderId
                self.last_heartbeat = time.time()

                if self.state != "follower":
                    print(
                        f"Node {self.node_id} steps down to FOLLOWER (Leader={request.leaderId})"
                    )
                    self.state = "follower"
                    self.voted_for = None

                self.term = request.term
                print(
                    f"Node {self.node_id} receives HEARTBEAT from Leader {request.leaderId}"
                )

                # Merge new entries
                self.merge_log_entries(request.log)
                if request.commitIndex > self.commit_index:
                    self.commit_index = request.commitIndex
                    print(
                        f"Node {self.node_id} updated commit index to {self.commit_index}"
                    )

                success = True
            else:
                print(
                    f"Node {self.node_id} rejects AppendEntries from Node {request.leaderId} (term {request.term} < {self.term})"
                )
                success = False

        # Unconditional replication to Go after the RPC
        self.replicate_log_in_go()

        return raft_pb2.AppendEntriesResponse(success=success, nodeId=self.node_id)

    def RequestVote(self, request, context):
        print(
            f"Node {self.node_id} runs RPC RequestVote called by Node {request.candidateId}"
        )

        granted = False
        with self.lock:
            if request.term > self.term:
                self.term = request.term
                self.state = "follower"
                self.voted_for = None

            if request.term == self.term and (
                self.voted_for is None or self.voted_for == request.candidateId
            ):
                self.voted_for = request.candidateId
                self.last_heartbeat = time.time()  # reset timer when voting
                print(f"Node {self.node_id} grants vote to Node {request.candidateId}")
                granted = True
            else:
                print(
                    f"Node {self.node_id} rejects vote for Node {request.candidateId}"
                )

        self.replicate_log_in_go()
        return raft_pb2.VoteResponse(voteGranted=granted)

    def SubmitClientRequest(self, request, context):
        print(f"Node {self.node_id} runs RPC SubmitClientRequest called by a client.")

        result_str = f"Node {self.node_id} is not the leader. Current state={self.state}"

        with self.lock:
            if self.state == "leader":
                new_index = len(self.log) + 1
                entry = LogEntry(request.operation, self.term, new_index)
                self.log.append(entry)
                print(
                    f"Node {self.node_id} appended operation '{request.operation}' at index {new_index}"
                )

                success = self.broadcast_new_entry(new_index)
                if success:
                    self.commit_index = new_index
                    result_str = (
                        f"Operation '{request.operation}' committed by leader {self.node_id}"
                    )
                    print(
                        f"Node {self.node_id} updated commit_index to {self.commit_index} after majority ack."
                    )
                else:
                    result_str = "Replication failed; no majority ack"
            else:
                if self.current_leader_id and self.current_leader_id != self.node_id:
                    leader_addr = f"node{self.current_leader_id}:{5000 + self.current_leader_id}"
                    print(
                        f"Node {self.node_id} sends RPC SubmitClientRequest to Node {self.current_leader_id}"
                    )
                    try:
                        channel = grpc.insecure_channel(leader_addr)
                        stub = raft_pb2_grpc.RaftStub(channel)
                        resp = stub.SubmitClientRequest(request)
                        result_str = resp.result
                    except Exception as ex:
                        result_str = (
                            f"Node {self.node_id} failed to forward to leader {self.current_leader_id}: {ex}"
                        )
                        print(result_str)

        self.replicate_log_in_go()
        return raft_pb2.ClientResponse(result=result_str)

    # ---------------------------------------------------
    # Single entry broadcast for majority ack
    # ---------------------------------------------------
    def broadcast_new_entry(self, entry_index):
        with self.lock:
            self.pending_acks[entry_index] = {self.node_id}
            new_entry = self.log[-1]
            single_entry = raft_pb2.LogEntry(
                operation=new_entry.operation, term=new_entry.term, index=new_entry.index
            )
            request = raft_pb2.AppendEntriesRequest(
                term=self.term,
                leaderId=self.node_id,
                log=[single_entry],
                commitIndex=self.commit_index,
            )

        threads = []
        for peer in self.peers:
            peer_id = parse_peer_id(peer)
            t = threading.Thread(
                target=self.send_append_one_entry,
                args=(peer, peer_id, request, entry_index),
            )
            t.start()
            threads.append(t)

        for t in threads:
            t.join(timeout=2.0)

        with self.lock:
            ack_count = len(self.pending_acks[entry_index])
            majority = len(self.peers) // 2 + 1
            ok = ack_count >= majority
            if ok:
                print(
                    f"Node {self.node_id} got majority ack_count={ack_count} for entry index={entry_index}"
                )
            else:
                print(
                    f"Node {self.node_id} failed to get majority ack (only {ack_count}) for entry index={entry_index}"
                )
            del self.pending_acks[entry_index]
            return ok

    def send_append_one_entry(self, peer, peer_id, request, entry_index):
        print(
            f"Node {self.node_id} sends RPC AppendEntries (single entry) to Node {peer_id}"
        )
        try:
            stub = self.get_stub(peer)
            response = stub.AppendEntries(request)
            if response.success:
                with self.lock:
                    self.pending_acks[entry_index].add(peer_id)
        except Exception as ex:
            print(
                f"Node {self.node_id} error sending single-entry AppendEntries to Node {peer_id}: {ex}"
            )

    # ---------------------------------------------------
    # Merging log entries on follower
    # ---------------------------------------------------
    def merge_log_entries(self, incoming_entries):
        local_last_index = self.log[-1].index if self.log else 0
        for entry in incoming_entries:
            if entry.index > local_last_index:
                self.log.append(LogEntry(entry.operation, entry.term, entry.index))
                local_last_index = entry.index

    # ---------------------------------------------------
    # Heartbeats & Elections
    # ---------------------------------------------------
    def send_heartbeat_once(self):
        """Send a single heartbeat immediately after becoming leader."""
        with self.lock:
            entries = [
                raft_pb2.LogEntry(operation=e.operation, term=e.term, index=e.index)
                for e in self.log
            ]
            request = raft_pb2.AppendEntriesRequest(
                term=self.term,
                leaderId=self.node_id,
                log=entries,
                commitIndex=self.commit_index,
            )

        for peer in self.peers:
            peer_id = parse_peer_id(peer)
            print(
                f"Node {self.node_id} sends RPC AppendEntries (immediate HB) to Node {peer_id}"
            )
            try:
                stub = self.get_stub(peer)
                stub.AppendEntries(request)
            except Exception as ex:
                print(f"Heartbeat error to Node {peer_id}: {ex}")

    def broadcast_heartbeat(self):
        while not self.stop_event.is_set():
            time.sleep(self.heartbeat_interval)
            with self.lock:
                if self.state != "leader":
                    continue
                entries = [
                    raft_pb2.LogEntry(operation=e.operation, term=e.term, index=e.index)
                    for e in self.log
                ]
                request = raft_pb2.AppendEntriesRequest(
                    term=self.term,
                    leaderId=self.node_id,
                    log=entries,
                    commitIndex=self.commit_index,
                )
            print(f"Leader Node {self.node_id} sends HEARTBEAT to all nodes")
            for peer in self.peers:
                peer_id = parse_peer_id(peer)
                print(f"Node {self.node_id} sends RPC AppendEntries to Node {peer_id}")
                try:
                    stub = self.get_stub(peer)
                    resp = stub.AppendEntries(request)
                    print(
                        f"Node {self.node_id} got response from Node {peer_id}: success={resp.success}"
                    )
                except Exception as ex:
                    print(f"Error contacting Node {peer_id}: {ex}")

    def start_heartbeat(self):
        threading.Thread(target=self.broadcast_heartbeat, daemon=True).start()

    def monitor_election_timeout(self):
        while not self.stop_event.is_set():
            time.sleep(0.1)
            with self.lock:
                if self.state != "follower":
                    continue
                if time.time() - self.last_heartbeat >= self.election_timeout:
                    print(
                        f"Node {self.node_id} election timeout expired; starting election"
                    )
                    threading.Thread(target=self.start_election, daemon=True).start()

    def start_election(self):
        with self.lock:
            self.state = "candidate"
            self.term += 1
            self.last_heartbeat = time.time()
            print(f"Node {self.node_id} becomes CANDIDATE for term {self.term}")

        vote_request = raft_pb2.VoteRequest(candidateId=self.node_id, term=self.term)
        votes_received = 1  # vote for self
        for peer in self.peers:
            peer_id = parse_peer_id(peer)
            print(f"Node {self.node_id} sends RPC RequestVote to Node {peer_id}")
            try:
                stub = self.get_stub(peer)
                response = stub.RequestVote(vote_request)
                if response.voteGranted:
                    votes_received += 1
                    print(f"Node {peer_id} granted vote to Node {self.node_id}")
            except Exception as ex:
                print(
                    f"Node {self.node_id} error requesting vote from Node {peer_id}: {ex}"
                )

        majority = len(self.peers) // 2 + 1
        with self.lock:
            elected = votes_received >= majority and self.state == "candidate"
            if elected:
                self.state = "leader"
                print(f"Node {self.node_id} is elected LEADER for term {self.term}")
            else:
                self.state = "follower"
                print(
                    f"Node {self.node_id} failed to become leader; reverting to FOLLOWER"
                )

        if elected:
            self.send_heartbeat_once()  # immediate heartbeat
            self.start_heartbeat()

    def stop(self):
        print(f"Stopping Node {self.node_id}")
        with self.lock:
            self.stop_event.set()


def serve(node):
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    raft_pb2_grpc.add_RaftServicer_to_server(node, server)
    server.add_insecure_port(f"[::]:{node.port}")
    server.start()
    try:
        while True:
            time.sleep(86400)
    except KeyboardInterrupt:
        server.stop(0)


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("-id", type=int, required=True)
    parser.add_argument("-port", type=int, required=True)
    parser.add_argument("-peers", type=str, required=True)
    args = parser.parse_args()

    node_instance = RaftNode(args.id, args.port, args.peers.split(","))
    serve(node_instance)


if __name__ == "__main__":
    main()
