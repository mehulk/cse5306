package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	pb "twopc_project/twopc_go_pb"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedTwoPhaseCommitServer
	nodeID        string
	transactionID string
	votes         []pb.Vote
	mu            sync.Mutex
}

// Coordinator logic for deciding global commit or abort
func coordinatorDecision(votes []pb.Vote) pb.Decision {
	for _, vote := range votes {
		if vote == pb.Vote_ABORT {
			return pb.Decision_GLOBAL_ABORT
		}
	}
	return pb.Decision_GLOBAL_COMMIT
}

// Coordinator sends global decision to participants
func sendGlobalDecision(transactionID string, decision pb.Decision, participants []string, nodeID string) {
	var wg sync.WaitGroup
	for _, addr := range participants {
		wg.Add(1)
		go func(address string) {
			defer wg.Done()
			conn, err := grpc.Dial(address, grpc.WithInsecure())
			if err != nil {
				log.Printf("Failed to connect to %s: %v", address, err)
				return
			}
			defer conn.Close()

			client := pb.NewTwoPhaseCommitClient(conn)

			fmt.Printf("Phase DECISION of Node %s sends RPC GlobalDecision to Phase DECISION of Node %s\n", nodeID, address)

			resp, err := client.GlobalDecision(context.Background(), &pb.GlobalDecisionRequest{
				TransactionId: transactionID,
				Decision:      decision,
			})
			if err != nil {
				log.Printf("Error sending global decision to %s: %v", address, err)
				return
			}

			fmt.Printf("Participant %s acknowledged global decision for transaction %s with ACK=%v\n", address, transactionID, resp.Ack)
		}(addr)
	}
	wg.Wait()
	fmt.Println("All participants have received the global decision.")
}

// Coordinator receives votes from Python coordinator via gRPC
func (s *server) SendVotes(ctx context.Context, req *pb.VotesReport) (*pb.AckResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	fmt.Printf("Coordinator received votes for transaction %s: %v\n", req.TransactionId, req.Votes)
	s.transactionID = req.TransactionId
	s.votes = req.Votes

	return &pb.AckResponse{Received: true}, nil
}

// Participant receives global decision from coordinator via gRPC
func (s *server) GlobalDecision(ctx context.Context, req *pb.GlobalDecisionRequest) (*pb.GlobalDecisionResponse, error) {
	fmt.Printf("Phase DECISION of Node %s receives RPC GlobalDecision from Phase DECISION of Node COORDINATOR\n", s.nodeID)

	if req.Decision == pb.Decision_GLOBAL_COMMIT {
		fmt.Printf("Participant %s locally commits transaction %s\n", s.nodeID, req.TransactionId)
	} else {
		fmt.Printf("Participant %s locally aborts transaction %s\n", s.nodeID, req.TransactionId)
	}

	return &pb.GlobalDecisionResponse{
		TransactionId: req.TransactionId,
		Ack:           pb.Ack_SUCCESS,
	}, nil
}

func main() {
	nodeID := os.Getenv("NODE_ID")
	if nodeID == "" {
		log.Fatalf("NODE_ID environment variable is not set!")
	}

	var port string
	if nodeID == "COORDINATOR" {
		port = "60050" // Fixed port for the coordinator
	} else {
		port = fmt.Sprintf("6005%s", nodeID) // Dynamic port for participants
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	s := &server{nodeID: nodeID}
	gRPCServer := grpc.NewServer()
	pb.RegisterTwoPhaseCommitServer(gRPCServer, s)

	log.Printf("Participant gRPC server running on port %s...", port)

	go func() {
		if nodeID == "COORDINATOR" { // Coordinator-specific logic
			fmt.Println("Coordinator waiting for votes...")

			waitForVotes(s)

			participantsEnv := os.Getenv("PARTICIPANT_ADDRS") // e.g., "participant1:60051,..."
			if participantsEnv == "" {
				log.Fatalf("PARTICIPANT_ADDRS environment variable is not set!")
			}
			participants := strings.Split(participantsEnv, ",")

			decision := coordinatorDecision(s.votes)
			sendGlobalDecision(s.transactionID, decision, participants, nodeID)

			log.Println("Coordinator completed decision phase.")
			os.Exit(0)
		}
	}()

	if err := gRPCServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}

func waitForVotes(s *server) {
	for len(s.votes) == 0 { // Wait until votes are received via gRPC
		time.Sleep(1 * time.Second)
	}
}
