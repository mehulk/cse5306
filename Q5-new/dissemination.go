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

	pb "github.com/pranavsdamu/damu-swim-disseem/generated/swim"

	"google.golang.org/grpc"
)

// DisseminationServer implements the SwimService gRPC interface for dissemination.
type DisseminationServer struct {
	pb.UnimplementedSwimServiceServer

	nodeID          int32
	membershipMutex sync.RWMutex
	membershipList  []*pb.NodeInfo
}

// NewDisseminationServer constructs a new DisseminationServer.
func NewDisseminationServer(nodeID int32, membershipList []*pb.NodeInfo) *DisseminationServer {
	return &DisseminationServer{
		nodeID:         nodeID,
		membershipList: membershipList,
	}
}

// Join is called by a new node that wants to join the cluster.
func (s *DisseminationServer) Join(ctx context.Context, req *pb.JoinRequest) (*pb.JoinResponse, error) {
	fmt.Printf(
		"Component Dissemination of Node %d runs RPC Join called by Component Dissemination of Node %d\n",
		s.nodeID, req.SenderId,
	)

	// Add the new node if not already present.
	newNode := &pb.NodeInfo{
		NodeId: req.SenderId,
		Host:   req.Host,
		Port:   req.Port,
	}

	s.membershipMutex.Lock()
	defer s.membershipMutex.Unlock()

	for _, n := range s.membershipList {
		if n.NodeId == req.SenderId {
			return &pb.JoinResponse{
				MembershipList: s.membershipList,
			}, nil
		}
	}

	s.membershipList = append(s.membershipList, newNode)
	return &pb.JoinResponse{
		MembershipList: s.membershipList,
	}, nil
}

// BroadcastFailure is invoked when a node detects a failure.
// It updates the local membership list and, if this node detected the failure,
// multicasts the failure info to all other nodes.
func (s *DisseminationServer) BroadcastFailure(ctx context.Context, req *pb.BroadcastFailureRequest) (*pb.BroadcastFailureResponse, error) {
	fmt.Printf(
		"Component Dissemination of Node %d runs RPC BroadcastFailure called by Component Dissemination of Node %d\n",
		s.nodeID, req.SenderId,
	)

	// Update local membership: remove the failed node.
	s.membershipMutex.Lock()
	var updated []*pb.NodeInfo
	for _, n := range s.membershipList {
		if n.NodeId != req.FailedNodeId {
			updated = append(updated, n)
		}
	}
	s.membershipList = updated
	s.membershipMutex.Unlock()

	// If this node is the one that detected the failure (i.e. sender equals self),
	// then multicast the failure info to all other nodes.
	if req.SenderId == s.nodeID {
		go s.BroadcastFailureToAll(req.FailedNodeId)
	}

	return &pb.BroadcastFailureResponse{Success: true}, nil
}

// BroadcastFailureToAll sends the failure information to all other nodes.
func (s *DisseminationServer) BroadcastFailureToAll(failedNodeId int32) {
	s.membershipMutex.RLock()
	defer s.membershipMutex.RUnlock()

	for _, node := range s.membershipList {
		if node.NodeId == s.nodeID {
			continue
		}

		address := fmt.Sprintf("%s:%d", node.Host, node.Port)
		fmt.Printf(
			"Component Dissemination of Node %d sends RPC BroadcastFailure to Component Dissemination of Node %d\n",
			s.nodeID, node.NodeId,
		)

		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Printf("Failed to connect to node %d for BroadcastFailure: %v", node.NodeId, err)
			continue
		}

		client := pb.NewSwimServiceClient(conn)
		_, err = client.BroadcastFailure(context.Background(), &pb.BroadcastFailureRequest{
			SenderId:     s.nodeID,
			FailedNodeId: failedNodeId,
		})
		conn.Close()

		if err != nil {
			log.Printf("BroadcastFailure RPC to node %d failed: %v", node.NodeId, err)
		}
	}
}

// JoinCluster allows a new node to join the cluster by contacting a bootstrap node.
func (s *DisseminationServer) JoinCluster(bootstrapHost string, bootstrapPort int, myHost string, myPort int, myNodeID int32) {
	address := fmt.Sprintf("%s:%d", bootstrapHost, bootstrapPort)
	fmt.Printf(
		"Component Dissemination of Node %d sends RPC Join to Component Dissemination of Node (bootstrap)\n",
		myNodeID,
	)

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Printf("Unable to dial bootstrap node: %v", err)
		return
	}
	defer conn.Close()

	client := pb.NewSwimServiceClient(conn)
	res, err := client.Join(context.Background(), &pb.JoinRequest{
		SenderId: myNodeID,
		Host:     myHost,
		Port:     int32(myPort),
	})
	if err != nil {
		log.Printf("Failed to call Join on bootstrap node: %v", err)
		return
	}

	s.membershipMutex.Lock()
	s.membershipList = res.MembershipList
	s.membershipMutex.Unlock()
}

// main starts the gRPC server for the Dissemination component.
func main() {
	nodeIDStr := os.Getenv("NODE_ID")
	portStr := os.Getenv("DISSEM_PORT")
	membershipStr := os.Getenv("MEMBERSHIP") // e.g., "1:node1:50051,2:node2:50052,3:node3:50053"

	var nodeID int32 = 1
	if nodeIDStr != "" {
		if val, err := strconv.Atoi(nodeIDStr); err == nil {
			nodeID = int32(val)
		}
	}

	dissemPort := 6000
	if portStr != "" {
		if val, err := strconv.Atoi(portStr); err == nil {
			dissemPort = val
		}
	}

	var membershipList []*pb.NodeInfo
	if membershipStr != "" {
		nodes := strings.Split(membershipStr, ",")
		for _, n := range nodes {
			parts := strings.Split(n, ":")
			if len(parts) == 3 {
				if id, errID := strconv.Atoi(parts[0]); errID == nil {
					if prt, errPort := strconv.Atoi(parts[2]); errPort == nil {
						membershipList = append(membershipList, &pb.NodeInfo{
							NodeId: int32(id),
							Host:   parts[1],
							Port:   int32(prt),
						})
					}
				}
			}
		}
	}

	srv := NewDisseminationServer(nodeID, membershipList)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", dissemPort))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", dissemPort, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSwimServiceServer(grpcServer, srv)

	log.Printf("Dissemination component of Node %d listening on port %d", nodeID, dissemPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
