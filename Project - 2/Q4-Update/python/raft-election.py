import grpc
import raft_pb2
import raft_pb2_grpc
import time
import threading
import argparse
import random
from concurrent import futures

# A simple class to hold a log entry: operation, term, index.
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
        self.peers = peers  # e.g. ["node2:5002","node3:5003",...]
        self.term = 0
        self.commit_index = 0
        self.log = []  # list of LogEntry
        self.lock = threading.Lock()
        self.grpc_channels = {}
        self.stop_event = threading.Event()
        self.voted_for = None
        self.votes_received = 0

        # For Q4: replicate from Python to Go nodes. Example:
        #   "q4node<id>:<6000 + id>"
        # We'll do that after we commit a new entry or on heartbeat.
        # We keep a separate channel cache for Go if desired, or reuse self.grpc_channels with different keys.
        self.go_channels = {}  # for q4

        # Timeouts
        self.election_timeout = random.uniform(1.5, 3.0)  # 1.5–3s
        self.heartbeat_interval = 1.0  # 1s
        self.last_heartbeat = time.time()

        print(f"Node {self.node_id} starting as FOLLOWER with election timeout {self.election_timeout:.2f}s")

        # Start the background thread to monitor election timeouts
        threading.Thread(target=self.monitor_election_timeout, daemon=True).start()

    # --------------------- RPCs (Raft) ---------------------

    def AppendEntries(self, request, context):
        print(f"Node {self.node_id} runs RPC AppendEntries called by Node {request.leaderId}")

        with self.lock:
            if request.term >= self.term:
                # reset election timer
                self.last_heartbeat = time.time()
                if self.state != "follower":
                    print(f"Node {self.node_id} steps down to FOLLOWER on receiving HEARTBEAT from Leader {request.leaderId}")
                    self.state = "follower"
                    self.voted_for = None

                self.term = request.term
                print(f"Node {self.node_id} receives HEARTBEAT from Leader {request.leaderId}")

                # update commitIndex if leader has advanced it
                if request.commitIndex > self.commit_index:
                    self.commit_index = request.commitIndex
                    print(f"Node {self.node_id} updated commit index to {self.commit_index}")

                return raft_pb2.AppendEntriesResponse(success=True)
            else:
                print(f"Node {self.node_id} rejects AppendEntries from Node {request.leaderId} (term {request.term} < {self.term})")
                return raft_pb2.AppendEntriesResponse(success=False)

    def RequestVote(self, request, context):
        print(f"Node {self.node_id} runs RPC RequestVote called by Node {request.candidateId}")

        with self.lock:
            if request.term > self.term:
                self.term = request.term
                self.state = "follower"
                self.voted_for = None

            if request.term == self.term and (self.voted_for is None or self.voted_for == request.candidateId):
                self.voted_for = request.candidateId
                print(f"Node {self.node_id} grants vote to Node {request.candidateId}")
                return raft_pb2.VoteResponse(voteGranted=True, term=self.term)
            else:
                print(f"Node {self.node_id} rejects vote for Node {request.candidateId}")
                return raft_pb2.VoteResponse(voteGranted=False, term=self.term)

    def SubmitClientRequest(self, request, context):
        # SERVER side log
        print(f"Node {self.node_id} runs RPC SubmitClientRequest called by a client (or node).")
        with self.lock:
            if self.state == "leader":
                # Append new operation to local log
                new_index = len(self.log) + 1
                entry = LogEntry(operation=request.operation, term=self.term, index=new_index)
                self.log.append(entry)
                print(f"Node {self.node_id} appended new client op '{request.operation}' at log index {new_index}.")

                # Attempt to replicate the single new entry to all Python followers
                majority_ok = self._broadcast_new_entry(new_index)
                if majority_ok:
                    # once majority, we can commit
                    self.commit_index = new_index
                    print(f"Node {self.node_id} commits index={new_index} after majority ACK.")

                    # replicate the entire log to the Go node (Q4)
                    self._replicate_log_to_go()
                    return raft_pb2.ClientResponse(result=f"Operation '{request.operation}' committed at index {new_index}")
                else:
                    return raft_pb2.ClientResponse(result="Replication failed: no majority")
            else:
                # Not the leader => forward
                # if we have a known leader, forward
                if hasattr(self, 'current_leader_id') and self.current_leader_id and self.current_leader_id != self.node_id:
                    leader_port = 5000 + self.current_leader_id
                    leader_addr = f"node{self.current_leader_id}:{leader_port}"

                    # CLIENT side log from the perspective of "this node" forwarding
                    print(f"Node {self.node_id} sends RPC SubmitClientRequest to Node {self.current_leader_id}")

                    try:
                        channel = grpc.insecure_channel(leader_addr)
                        stub = raft_pb2_grpc.RaftStub(channel)
                        return stub.SubmitClientRequest(request)
                    except Exception as ex:
                        msg = f"Node {self.node_id} failed to forward to leader {self.current_leader_id}: {ex}"
                        print(msg)
                        return raft_pb2.ClientResponse(result=msg)
                else:
                    msg = f"Node {self.node_id} is not the leader and no known leader."
                    print(msg)
                    return raft_pb2.ClientResponse(result=msg)

    # --------------------- Python-level Replication of New Entry ---------------------

    def _broadcast_new_entry(self, entry_index):
        """
        Send the single newly-appended entry to each Python follower to get an ACK.
        If majority acknowledges, return True.
        """
        # We'll need a concurrency approach or a quick blocking approach
        # For simplicity, do a blocking approach:
        ack_count = 1  # leader acks itself
        new_entry = self.log[-1]  # the appended entry
        single_append_request = raft_pb2.AppendEntriesRequest(
            term=self.term,
            leaderId=self.node_id,
            # We'll send just 1 new entry
            log=[raft_pb2.LogEntry(operation=new_entry.operation, term=new_entry.term, index=new_entry.index)],
            commitIndex=self.commit_index
        )

        def send_append(peer):
            nonlocal ack_count
            print(f"Node {self.node_id} sends RPC AppendEntries (new entry) to Node {peer}")
            try:
                if peer not in self.grpc_channels:
                    self.grpc_channels[peer] = raft_pb2_grpc.RaftStub(grpc.insecure_channel(peer))
                stub = self.grpc_channels[peer]
                resp = stub.AppendEntries(single_append_request)
                if resp.success:
                    with self.lock:
                        ack_count += 1
            except Exception as ex:
                print(f"Node {self.node_id} error sending new entry to {peer}: {ex}")

        threads = []
        for peer in self.peers:
            t = threading.Thread(target=send_append, args=(peer,))
            t.start()
            threads.append(t)

        for t in threads:
            t.join(timeout=2.0)

        with self.lock:
            majority_needed = len(self.peers)//2 + 1
            print(f"Node {self.node_id} got ack_count={ack_count}; majority={majority_needed}")
            return (ack_count >= majority_needed)

    # --------------------- Q4: replicate the entire log to Go node(s) ---------------------

    def _replicate_log_to_go(self):
        """
        Example: replicate to just this node's "twin" Go container, e.g. q4nodeX:600X
        If you want to replicate to all 5 Go containers, you can do that in a loop.
        """

        # replicate to exactly "q4node<node_id> : (6000 + node_id)"
        target_host = f"q4node{self.node_id}"
        target_port = 6000 + self.node_id
        target_addr = f"{target_host}:{target_port}"

        print(f"Node {self.node_id} sends RPC ReplicateLog to Node q4node{self.node_id}")
        try:
            if target_addr not in self.go_channels:
                ch = grpc.insecure_channel(target_addr)
                self.go_channels[target_addr] = pb.LogReplicatorStub(ch)

            stub = self.go_channels[target_addr]
            # build entire log
            proto_entries = []
            for entry in self.log:
                proto_entries.append(
                    pb.LogEntry(operation=entry.operation, term=entry.term, index=entry.index)
                )

            req = pb.AppendEntriesRequest(
                term=self.term,
                leaderId=self.node_id,
                log=proto_entries,
                commitIndex=self.commit_index
            )
            resp = stub.ReplicateLog(req)
            print(f"Node {self.node_id} got ReplicateLog ack from Go node: success={resp.success}")
        except Exception as ex:
            print(f"Node {self.node_id} error replicating to Go node q4node{self.node_id}: {ex}")


    # --------------------- Heartbeat + Elections ---------------------

    def broadcast_heartbeat(self):
        while not self.stop_event.is_set():
            if self.state == "leader":
                with self.lock:
                    request = raft_pb2.AppendEntriesRequest(
                        term=self.term,
                        leaderId=self.node_id,
                        commitIndex=self.commit_index,
                    )
                    print(f"Leader Node {self.node_id} sends HEARTBEAT to all nodes")

                    for peer in self.peers:
                        try:
                            print(f"Node {self.node_id} sends RPC AppendEntries to Node {peer}")
                            if peer not in self.grpc_channels:
                                channel = grpc.insecure_channel(peer)
                                self.grpc_channels[peer] = raft_pb2_grpc.RaftStub(channel)
                            stub = self.grpc_channels[peer]
                            response = stub.AppendEntries(request)
                            print(f"Node {self.node_id} received response from Node {peer}: success={response.success}")
                        except Exception as e:
                            print(f"Error contacting Node {peer}: {e}")

                time.sleep(self.heartbeat_interval)

    def start_heartbeat(self):
        threading.Thread(target=self.broadcast_heartbeat, daemon=True).start()

    def monitor_election_timeout(self):
        while not self.stop_event.is_set():
            time.sleep(0.1)
            with self.lock:
                elapsed_time = time.time() - self.last_heartbeat
                if elapsed_time >= self.election_timeout and self.state == "follower":
                    print(f"Node {self.node_id} election timeout expired; starting election")
                    threading.Thread(target=self.start_election, daemon=True).start()

    def start_election(self):
        with self.lock:
            self.state = "candidate"
            self.term += 1
            self.voted_for = self.node_id
            print(f"Node {self.node_id} becomes CANDIDATE for term {self.term}")

        vote_req = raft_pb2.VoteRequest(candidateId=self.node_id, term=self.term)

        votes_received = 1  # vote for self
        for peer in self.peers:
            try:
                print(f"Node {self.node_id} sends RPC RequestVote to Node {peer}")
                if peer not in self.grpc_channels:
                    channel = grpc.insecure_channel(peer)
                    self.grpc_channels[peer] = raft_pb2_grpc.RaftStub(channel)
                stub = self.grpc_channels[peer]
                resp = stub.RequestVote(vote_req)
                if resp.voteGranted:
                    with self.lock:
                        votes_received += 1
                        print(f"Node {peer} granted vote to Node {self.node_id}")
            except Exception as e:
                print(f"Node {self.node_id} error requesting vote from {peer}: {e}")

        with self.lock:
            majority_needed = len(self.peers)//2 + 1
            if votes_received >= majority_needed and self.state == "candidate":
                print(f"Node {self.node_id} is elected LEADER for term {self.term}")
                self.state = "leader"
                # you might track self.current_leader_id = self.node_id
                threading.Thread(target=self.start_heartbeat, daemon=True).start()
            else:
                print(f"Node {self.node_id} failed to become leader; reverting to FOLLOWER")
                self.state = "follower"

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
    import sys
    parser = argparse.ArgumentParser()
    parser.add_argument("-id", type=int, required=True)
    parser.add_argument("-port", type=int, required=True)
    parser.add_argument("-peers", type=str, required=True)
    args = parser.parse_args()

    peers_list = args.peers.split(",")

    node_instance = RaftNode(args.id, args.port, peers_list)
    serve(node_instance)

if __name__ == "__main__":
    main()
