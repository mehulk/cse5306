package main

import (
	context "context"
	fmt "fmt"
	log "log"
	net "net"
	os "os"
	strings "strings"
	"sync"

	grpc "google.golang.org/grpc"

	pb "raftserver/proto"
)

type server struct {
	pb.UnimplementedLogReplicatorServer

	nodeID      int32
	peers       []string
	commitIndex int32
	logs        []*pb.LogEntry
	lock        sync.Mutex
}

func (s *server) ReplicateLog(ctx context.Context, req *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	fmt.Printf("Node %d runs RPC ReplicateLog called by Node %d\n", s.nodeID, req.LeaderId)

	s.lock.Lock()
	defer s.lock.Unlock()

	// keep longest log
	if len(req.Log) > len(s.logs) {
		s.logs = make([]*pb.LogEntry, len(req.Log))
		copy(s.logs, req.Log)
	}
	// commit index only moves forward
	if req.CommitIndex > s.commitIndex {
		s.commitIndex = req.CommitIndex
	}

	fmt.Printf("Node %d: commitIndex=%d, total entries=%d\n", s.nodeID, s.commitIndex, len(s.logs))
	for i := 0; i < int(s.commitIndex); i++ {
		entry := s.logs[i]
		fmt.Printf("Node %d executes op=%s at index=%d\n", s.nodeID, entry.Operation, entry.Index)
	}
	return &pb.AppendEntriesResponse{Success: true, NodeId: s.nodeID}, nil
}

func (s *server) AckReplication(ctx context.Context, resp *pb.AppendEntriesResponse) (*pb.ClientResponse, error) {
	fmt.Printf("Node %d runs RPC AckReplication called by Node %d\n", s.nodeID, resp.NodeId)
	return &pb.ClientResponse{Result: "ACK"}, nil
}

func main() {
	if len(os.Args) < 4 {
		log.Fatalf("usage: q4server <id> <port> <peer1,peer2,...>")
	}
	var id int32
	fmt.Sscanf(os.Args[1], "%d", &id)
	port := os.Args[2]
	peers := strings.Split(os.Args[3], ",")

	s := &server{nodeID: id, peers: peers, logs: []*pb.LogEntry{}}

	grpcServer := grpc.NewServer()
	pb.RegisterLogReplicatorServer(grpcServer, s)
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	log.Printf("Node %d LogReplicator listening on %s", id, port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
