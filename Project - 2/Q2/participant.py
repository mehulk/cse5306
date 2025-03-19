import os
import grpc
from concurrent import futures
import random
import time
from two_phase_commit_pb2 import VoteRequest, VoteResponse, COMMIT, ABORT
from two_phase_commit_pb2_grpc import TwoPhaseCommitServicer, add_TwoPhaseCommitServicer_to_server

PHASE_NAME = "VOTING"

def vote_enum_to_str(vote_enum):
    return "COMMIT" if vote_enum == COMMIT else "ABORT"

class Participant(TwoPhaseCommitServicer):
    def RequestVote(self, request, context):
        node_id = os.getenv("HOSTNAME", "PARTICIPANT")
        listen_port = os.getenv("LISTEN_PORT", "50051")

        # Before deciding, let's log the arrival
        print(f"Phase {PHASE_NAME} of Node {node_id}:{listen_port} "
              f"receives RPC RequestVote from Phase {PHASE_NAME} of Node COORDINATOR "
              f"for transaction {request.transaction_id}")

        # Decide
        force_val = os.getenv("FORCE_COMMIT", "false").lower()
        if force_val == "true":
            decision = COMMIT
        else:
            decision = random.choice([COMMIT, ABORT])
        decision_str = vote_enum_to_str(decision)

        # Log the response
        print(f"Phase {PHASE_NAME} of Node {node_id}:{listen_port} "
              f"sends RPC RequestVoteResponse to Phase {PHASE_NAME} of Node COORDINATOR "
              f"with vote={decision_str} for transaction {request.transaction_id}")

        return VoteResponse(transaction_id=request.transaction_id, vote=decision)

def serve():
    port = os.getenv("VOTING_PORT", "50051")  # <-- Ensure this is VOTING_PORT
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    add_TwoPhaseCommitServicer_to_server(Participant(), server)
    server.add_insecure_port(f'[::]:{port}')
    server.start()
    print(f"Participant voting node running on port {port}.")
    server.wait_for_termination()


if __name__ == "__main__":
    serve()
