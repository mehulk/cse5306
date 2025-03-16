import os
import grpc
from concurrent import futures
import random
import time
from two_phase_commit_pb2 import VoteRequest, VoteResponse, COMMIT, ABORT
from two_phase_commit_pb2_grpc import TwoPhaseCommitServicer, add_TwoPhaseCommitServicer_to_server

class Participant(TwoPhaseCommitServicer):
    def RequestVote(self, request, context):
        # Check environment variable
        force_val = os.getenv("FORCE_COMMIT", "false").lower()

        if force_val == "true":
            # Always commit if FORCE_COMMIT=true
            decision = COMMIT
        else:
            # Otherwise random choice
            decision = random.choice([COMMIT, ABORT])

        transaction_id = request.transaction_id
        print(f"Participant received vote request for {transaction_id}. "
              f"FORCE_COMMIT={force_val}, Decision={decision}")

        return VoteResponse(transaction_id=transaction_id, vote=decision)

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    add_TwoPhaseCommitServicer_to_server(Participant(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("Participant node running on port 50051.")
    try:
        while True:
            time.sleep(86400)
    except KeyboardInterrupt:
        server.stop(0)

if __name__ == "__main__":
    serve()
