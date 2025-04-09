package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	pb "q4/proto" // Import your generated package.
	// If your go_package = "q4/proto/raft;raft",
	// you might do: import pb "q4/proto/raft"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type logEntry = pb.LogEntry

type Node struct {
	pb.UnimplementedLogReplicatorServer

	id    int32
	port  string
	peers []string

	stateMu     sync.Mutex
	log         []*logEntry
	commitIndex int32
	kv          map[string]string

	grpcClients map[string]pb.LogReplicatorClient
}

func NewNode(id int32, port string, peers []string) *Node {
	return &Node{
		id:          id,
		port:        port,
		peers:       peers,
		kv:          make(map[string]string),
		grpcClients: make(map[string]pb.LogReplicatorClient),
	}
}

func (n *Node) ReplicateLog(ctx context.Context, req *pb.AppendEntriesRequest) (*pb.AppendEntriesResponse, error) {
	fmt.Printf("Node %d runs RPC ReplicateLog called by Node %d\n", n.id, req.LeaderId)

	n.stateMu.Lock()
	defer n.stateMu.Unlock()

	// 1) Replace local log with incoming log
	n.log = make([]*logEntry, len(req.Log))
	copy(n.log, req.Log)

	// 2) Apply newly committed entries
	for i := n.commitIndex; i < req.CommitIndex; i++ {
		entry := n.log[i]
		fmt.Printf("Node %d applying: %s (term=%d, idx=%d)\n", n.id, entry.Operation, entry.Term, entry.Index)

		opParts := strings.SplitN(entry.Operation, "=", 2)
		if len(opParts) == 2 {
			key := opParts[0]
			val := opParts[1]
			n.kv[key] = val
			fmt.Printf("Node %d  •  KV store update: %q => %q\n", n.id, key, val)
		}
	}

	n.commitIndex = req.CommitIndex

	return &pb.AppendEntriesResponse{
		Success: true,
		NodeId:  n.id,
	}, nil
}

func (n *Node) getClient(addr string) pb.LogReplicatorClient {
	n.stateMu.Lock()
	defer n.stateMu.Unlock()

	if cl, ok := n.grpcClients[addr]; ok {
		return cl
	}
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("Node %d failed to dial %s: %v\n", n.id, addr, err)
		return nil
	}
	client := pb.NewLogReplicatorClient(conn)
	n.grpcClients[addr] = client
	return client
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("usage: q4server <nodeId> <port> <peer1,peer2,...>")
		os.Exit(1)
	}
	id64, _ := strconv.ParseInt(os.Args[1], 10, 32)
	id := int32(id64)
	port := os.Args[2]
	peerStr := os.Args[3]
	peers := strings.Split(peerStr, ",")

	node := NewNode(id, port, peers)

	srv := grpc.NewServer()
	pb.RegisterLogReplicatorServer(srv, node)

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Node %d listen error on port %s: %v", id, port, err)
	}
	log.Printf("Node %d (Go) listening on %s", id, port)
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("Node %d Serve failed: %v", id, err)
	}
}
