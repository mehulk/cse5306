import grpc
import time
import os
from two_phase_commit_pb2 import VoteRequest, VotesReport, COMMIT, ABORT
from two_phase_commit_pb2_grpc import TwoPhaseCommitStub

PHASE_NAME = "VOTING"
NODE_ID = "COORDINATOR"

def vote_enum_to_str(vote_enum):
    """Helper to convert numeric enum to 'COMMIT' or 'ABORT' string."""
    return "COMMIT" if vote_enum == COMMIT else "ABORT"

def get_participant_addresses():
    """Dynamically calculate participant addresses based on NODE_ID."""
    participant_count = int(os.getenv("PARTICIPANT_COUNT", "5"))  # Default to 5 participants
    participants = []
    for i in range(1, participant_count + 1):
        voting_port = f"5005{i}"
        participants.append(f"participant{i}:{voting_port}")
    return participants

def request_votes(transaction_id):
    participants = get_participant_addresses()
    votes = []
    for address in participants:
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

def send_votes_to_go(transaction_id, votes):
    """Send collected votes to Go decision service."""
    print("Sending votes from Python coordinator to Go decision service...")
    decision_port = "60050"  # Fixed decision port for coordinator
    with grpc.insecure_channel(f'localhost:{decision_port}') as channel:
        stub = TwoPhaseCommitStub(channel)
        stub.SendVotes(VotesReport(transaction_id=transaction_id, votes=votes))
    print("Votes sent successfully.")

def main():
    transaction_id = "txn123"
    print("Waiting for participants to start up...")
    time.sleep(5)
    print(f"Coordinator initiating vote for transaction {transaction_id}")
    
    votes = request_votes(transaction_id)

    # Send collected votes to Go decision service
    send_votes_to_go(transaction_id, votes)

if __name__ == "__main__":
    main()
