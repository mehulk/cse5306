package main

import (
	"context"
	"flag"
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

// Participant receives global decision from coordinator
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
	role := flag.String("role", "", "Role of this node: coordinator or participant")
	port := flag.String("port", "6000", "Port to run the Go gRPC server")
	flag.Parse()

	nodeID := os.Getenv("HOSTNAME")
	if nodeID == "" {
		nodeID = strings.ToUpper(*role)
	}

	lis, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", *port, err)
	}

	s := &server{nodeID: nodeID}
	gRPCServer := grpc.NewServer()
	pb.RegisterTwoPhaseCommitServer(gRPCServer, s)

	go func() {
		log.Printf("%s gRPC server running on port %s...", strings.Title(*role), *port)
		if err := gRPCServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC server: %v", err)
		}
	}()

	if *role == "coordinator" {
		fmt.Println("Coordinator waiting for votes...")

		// Wait for votes to be sent via gRPC (handled by `SendVotes`)
		waitForVotes(s)

		// Process votes and send global decision
		participantsEnv := os.Getenv("PARTICIPANT_ADDRS") // e.g., "participant1:6001,..."
		participants := strings.Split(participantsEnv, ",")

		decision := coordinatorDecision(s.votes)
		sendGlobalDecision(s.transactionID, decision, participants, nodeID)

		log.Println("Coordinator completed decision phase.")
		os.Exit(0)
	}

	select {}
}

func waitForVotes(s *server) {
	for len(s.votes) == 0 { // Wait until votes are received via gRPC
		time.Sleep(1 * time.Second)
	}
}
