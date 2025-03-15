import grpc
import raft_pb2
import raft_pb2_grpc
import time
import threading
import argparse
import random
from concurrent import futures

# Simple class to hold a log entry.
class LogEntry:
    def __init__(self, operation, term, index):
        self.operation = operation
        self.term = term
        self.index = index

# Raft node implementation.
class RaftNode(raft_pb2_grpc.RaftServicer):
    def __init__(self, node_id, port, peers):
        self.node_id = node_id
        self.state = "follower"  # All nodes start as followers.
        self.port = port
        self.peers = peers  # List of peer addresses (e.g., "node2:5002")
        self.term = 0
        self.commit_index = 0
        self.log = []  # List of LogEntry objects.
        self.lock = threading.Lock()
        self.grpc_channels = {}  # Cache for gRPC channels to peers.
        self.stop_event = threading.Event()
        self.voted_for = None
        self.votes_received = 0

        # Election timeout: randomized between 1.5 and 3 seconds.
        self.election_timeout = random.uniform(1.5, 3.0)
        # Heartbeat interval: fixed at 1 second.
        self.heartbeat_interval = 1.0

        # Last heartbeat received time.
        self.last_heartbeat = time.time()

        print(f"Node {self.node_id} starting as FOLLOWER with election timeout {self.election_timeout:.2f}s")
        threading.Thread(target=self.monitor_election_timeout, daemon=True).start()

    # --- gRPC RPC Methods ---
    def AppendEntries(self, request, context):
        print(f"Node {self.node_id} runs RPC AppendEntries called by Node {request.leaderId}")
        
        with self.lock:
            if request.term >= self.term:
                # Reset election timer since heartbeat was received.
                self.last_heartbeat = time.time()
                
                if self.state != "follower":
                    print(f"Node {self.node_id} steps down to FOLLOWER on receiving HEARTBEAT from Leader {request.leaderId}")
                    self.state = "follower"
                    self.voted_for = None
                
                self.term = request.term  # Update our term.
                print(f"Node {self.node_id} receives HEARTBEAT from Leader {request.leaderId}")
                
                if request.commitIndex > self.commit_index:
                    self.commit_index = request.commitIndex
                    print(f"Node {self.node_id} updated commit index to {self.commit_index}")
                
                return raft_pb2.AppendEntriesResponse(success=True)
            else:
                print(f"Node {self.node_id} rejected AppendEntries from Node {request.leaderId} (term {request.term} < {self.term})")
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
                return raft_pb2.VoteResponse(voteGranted=True)
            else:
                print(f"Node {self.node_id} rejects vote for Node {request.candidateId}")
                return raft_pb2.VoteResponse(voteGranted=False)

    # --- Heartbeat and Election Functions ---
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
                            if peer not in self.grpc_channels:
                                channel = grpc.insecure_channel(peer)
                                stub = raft_pb2_grpc.RaftStub(channel)
                                self.grpc_channels[peer] = stub
                            else:
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
            print(f"Node {self.node_id} becomes CANDIDATE for term {self.term}")
        
        vote_request = raft_pb2.VoteRequest(candidateId=self.node_id, term=self.term)
        
        votes_received = 1  # Vote for itself.
        
        for peer in self.peers:
            try:
                if peer not in self.grpc_channels:
                    channel = grpc.insecure_channel(peer)
                    stub = raft_pb2_grpc.RaftStub(channel)
                    response = stub.RequestVote(vote_request)
                    if response.voteGranted:
                        votes_received += 1
                        print(f"Node {peer} granted vote to Node {self.node_id}")
            except Exception as e:
                print(f"Error requesting vote from Node {peer}: {e}")
        
        majority_votes_required = len(self.peers) // 2 + 1
        
        with self.lock:
            if votes_received >= majority_votes_required and self.state == "candidate":
                print(f"Node {self.node_id} is elected LEADER for term {self.term}")
                self.state = "leader"
                threading.Thread(target=self.start_heartbeat, daemon=True).start()
            else:
                print(f"Node {self.node_id} failed to become leader; reverting to FOLLOWER")
                self.state = "follower"

    def stop(self):
        """Stop the node."""
        print(f"Stopping Node {self.node_id}")
        with self.lock:
            self.stop_event.set()

# gRPC server setup.
def serve(node):
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    raft_pb2_grpc.add_RaftServicer_to_server(node, server)
    server.add_insecure_port(f"[::]:{node.port}")
    server.start()
    
    try:
        while True:
            time.sleep(86400)  # Keep server running.
    except KeyboardInterrupt:
        server.stop(0)

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("-id", type=int, required=True, help="Unique Node ID")
    parser.add_argument("-port", type=int, required=True, help="Port to listen on")
    parser.add_argument("-peers", type=str, required=True, help="Comma-separated list of peer addresses (e.g., node2:5002,node3:5003)")
    
    args = parser.parse_args()
    
    peers_list = args.peers.split(",")
    
    node_instance = RaftNode(node_id=args.id, port=args.port, peers=peers_list)
    
    serve(node_instance)

if __name__ == "__main__":
    main()

##TO DO
## correct logs as per requirement