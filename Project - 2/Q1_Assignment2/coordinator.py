import grpc
from two_phase_commit_pb2 import VoteRequest, VoteResponse, COMMIT, ABORT
from two_phase_commit_pb2_grpc import TwoPhaseCommitStub

# Participant addresses
PARTICIPANTS = [
    'participant1:50051',
    'participant2:50051',
    'participant3:50051',
    'participant4:50051',
    'participant5:50051'
]

def request_votes(transaction_id):
    votes = []
    for address in PARTICIPANTS:
        with grpc.insecure_channel(address) as channel:
            stub = TwoPhaseCommitStub(channel)
            request = VoteRequest(transaction_id=transaction_id)
            response = stub.RequestVote(request)
            print(f"Received vote from {address}: {response.vote}")
            votes.append(response.vote)
    return votes

def main():
    transaction_id = "txn123"
    print(f"Coordinator initiating vote for transaction {transaction_id}")
    votes = request_votes(transaction_id)
    if all(vote == COMMIT for vote in votes):
        print("All participants voted COMMIT. Proceeding to commit.")
    else:
        print("At least one participant voted ABORT. Aborting transaction.")

if __name__ == "__main__":
    main()
