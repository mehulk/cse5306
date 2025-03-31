package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	"google.golang.org/grpc"

	pb "raftserver/proto"
)

type LogReplicatorServer struct {
	pb.UnimplementedLogReplicatorServer

	nodeID      int32
	peers       []string
	commitIndex int32
	logs        []*pb.LogEntry
	lock        sync.Mutex

	grpcChannels map[string]pb.LogReplicatorClient
}

func NewLogReplicatorServer(nodeID int32, peers []string) *LogReplicatorServer {
	return &LogReplicatorServer{
		nodeID:       nodeID,
		peers:        peers,
		commitIndex:  0,
		logs:         []*pb.LogEntry{},
		grpcChannels: make(map[string]pb.LogReplicatorClient),
	}
}

func (s *LogReplicatorServer) ReplicateLog(ctx context.Context, req *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	// Server-side print: Node <local> runs RPC ReplicateLog called by Node <req.LeaderId>
	fmt.Printf("Node %d runs RPC ReplicateLog called by Node %d\n", s.nodeID, req.LeaderId)

	s.lock.Lock()
	s.logs = req.Log
	s.commitIndex = req.CommitIndex
	s.lock.Unlock()

	fmt.Printf("Node %d: logs updated. commitIndex=%d. Total log entries=%d\n",
		s.nodeID, req.CommitIndex, len(req.Log))

	// "Apply" all committed entries
	for i := 0; i < int(req.CommitIndex); i++ {
		e := req.Log[i]
		fmt.Printf("Node %d: executing committed operation at index=%d op=%s\n",
			s.nodeID, e.Index, e.Operation)
	}

	return &pb.AppendEntriesResponse{
		Success: true,
		NodeId:  s.nodeID,
	}, nil
}

func (s *LogReplicatorServer) AckReplication(ctx context.Context, resp *pb.AppendEntriesResponse) (*pb.ClientResponse, error) {
	fmt.Printf("Node %d runs RPC AckReplication called by Node %d\n", s.nodeID, resp.NodeId)
	return &pb.ClientResponse{Result: "ACK received"}, nil
}

func (s *LogReplicatorServer) getOrCreateChannel(peer string) pb.LogReplicatorClient {
	s.lock.Lock()
	defer s.lock.Unlock()

	client, ok := s.grpcChannels[peer]
	if !ok {
		conn, err := grpc.Dial(peer, grpc.WithInsecure())
		if err != nil {
			fmt.Printf("Node %d error connecting to peer %s: %v\n", s.nodeID, peer, err)
			return nil
		}
		client = pb.NewLogReplicatorClient(conn)
		s.grpcChannels[peer] = client
	}
	return client
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: go run server.go <nodeId> <port> <peer1:port,peer2:port,...>")
		return
	}

	nodeIDStr := os.Args[1]
	portStr := os.Args[2]
	peerStr := os.Args[3]

	var nodeID int32
	fmt.Sscanf(nodeIDStr, "%d", &nodeID)

	peers := strings.Split(peerStr, ",")

	server := grpc.NewServer()
	replicator := NewLogReplicatorServer(nodeID, peers)
	pb.RegisterLogReplicatorServer(server, replicator)

	lis, err := net.Listen("tcp", ":"+portStr)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", portStr, err)
	}
	log.Printf("Node %d LogReplicator server listening on port %s", nodeID, portStr)

	// optional: dial peers
	for _, p := range peers {
		replicator.getOrCreateChannel(p)
	}

	if err := server.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
