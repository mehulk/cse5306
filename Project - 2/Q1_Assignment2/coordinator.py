import grpc
import time
from two_phase_commit_pb2 import VoteRequest, VoteResponse, COMMIT, ABORT
from two_phase_commit_pb2_grpc import TwoPhaseCommitStub

PHASE_NAME = "VOTING"
NODE_ID = "COORDINATOR"

PARTICIPANTS = [
    "participant1:50051",
    "participant2:50052",
    "participant3:50053",
    "participant4:50054",
    "participant5:50055"
]

def vote_enum_to_str(vote_enum):
    """Helper to convert numeric enum to 'COMMIT' or 'ABORT' string."""
    return "COMMIT" if vote_enum == COMMIT else "ABORT"

def request_votes(transaction_id):
    votes = []
    for address in PARTICIPANTS:
        # Print before sending
        print(f"Phase {PHASE_NAME} of Node {NODE_ID} sends RPC RequestVote "
              f"to Phase {PHASE_NAME} of Node {address} for transaction {transaction_id}")

        with grpc.insecure_channel(address) as channel:
            stub = TwoPhaseCommitStub(channel)
            request = VoteRequest(transaction_id=transaction_id)
            response = stub.RequestVote(request)

            # Convert the numeric vote (0 or 1) to a string
            vote_str = vote_enum_to_str(response.vote)

            # Print after receiving response
            print(f"Phase {PHASE_NAME} of Node {NODE_ID} receives RPC RequestVoteResponse "
                  f"from Phase {PHASE_NAME} of Node {address} with vote={vote_str}")

            votes.append(response.vote)
    return votes

def main():
    transaction_id = "txn123"
    print("Waiting for participants to start up...")
    time.sleep(5)
    print(f"Coordinator initiating vote for transaction {transaction_id}")
    
    votes = request_votes(transaction_id)

    # Convert numeric votes to strings for final log
    final_votes_str = [vote_enum_to_str(v) for v in votes]
    if all(v == COMMIT for v in votes):
        print(f"All participants voted {', '.join(final_votes_str)}. Proceeding to COMMIT.")
    else:
        print(f"At least one participant voted ABORT among {', '.join(final_votes_str)}. Aborting transaction.")

if __name__ == "__main__":
    main()
