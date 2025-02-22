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

	pb "Q5new/swim" // Adjust this import path as needed

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

const DC_BASE = 60050

// DisseminationServer implements the SwimService interface for dissemination.
type DisseminationServer struct {
	pb.UnimplementedSwimServiceServer

	nodeID          int32
	membershipMutex sync.RWMutex
	membershipList  []*pb.NodeInfo

	// For streaming membership updates.
	subscribers      []chan *pb.JoinResponse
	subscribersMutex sync.Mutex
}

// NewDisseminationServer constructs a new DisseminationServer.
func NewDisseminationServer(nodeID int32, membershipList []*pb.NodeInfo) *DisseminationServer {
	return &DisseminationServer{
		nodeID:         nodeID,
		membershipList: membershipList,
		subscribers:    make([]chan *pb.JoinResponse, 0),
	}
}

// notifySubscribers sends the updated membership list to all subscriber channels.
func (s *DisseminationServer) notifySubscribers() {
	s.membershipMutex.RLock()
	current := &pb.JoinResponse{MembershipList: s.membershipList}
	s.membershipMutex.RUnlock()

	s.subscribersMutex.Lock()
	defer s.subscribersMutex.Unlock()
	for _, ch := range s.subscribers {
		select {
		case ch <- current:
		default:
			// skip if channel is full
		}
	}
}

// StreamMembership implements the streaming RPC for membership updates.
func (s *DisseminationServer) StreamMembership(_ *emptypb.Empty, stream pb.SwimService_StreamMembershipServer) error {
	updateCh := make(chan *pb.JoinResponse, 1)
	s.subscribersMutex.Lock()
	s.subscribers = append(s.subscribers, updateCh)
	s.subscribersMutex.Unlock()

	// Send the current membership list immediately.
	s.membershipMutex.RLock()
	current := &pb.JoinResponse{MembershipList: s.membershipList}
	s.membershipMutex.RUnlock()
	if err := stream.Send(current); err != nil {
		return err
	}

	// Stream updates.
	for {
		update, ok := <-updateCh
		if !ok {
			break
		}
		if err := stream.Send(update); err != nil {
			return err
		}
	}
	return nil
}

// Join is called by a new node to join the cluster.
func (s *DisseminationServer) Join(ctx context.Context, req *pb.JoinRequest) (*pb.JoinResponse, error) {
	fmt.Printf("Component Dissemination of Node %d runs RPC Join called by Node %d\n", s.nodeID, req.SenderId)
	newNode := &pb.NodeInfo{
		NodeId: req.SenderId,
		Host:   req.Host,
		Port:   req.Port,
	}

	s.membershipMutex.Lock()
	defer s.membershipMutex.Unlock()
	for _, n := range s.membershipList {
		if n.NodeId == req.SenderId {
			return &pb.JoinResponse{MembershipList: s.membershipList}, nil
		}
	}
	s.membershipList = append(s.membershipList, newNode)
	s.notifySubscribers()
	return &pb.JoinResponse{MembershipList: s.membershipList}, nil
}

// BroadcastFailure is invoked when a node detects a failure.
// It updates the membership list and notifies subscribers.
// Only the originating node forwards the failure.
func (s *DisseminationServer) BroadcastFailure(ctx context.Context, req *pb.BroadcastFailureRequest) (*pb.BroadcastFailureResponse, error) {
	fmt.Printf("Component Dissemination of Node %d runs RPC BroadcastFailure called by Node %d\n", s.nodeID, req.SenderId)

	s.membershipMutex.Lock()
	var updated []*pb.NodeInfo
	for _, n := range s.membershipList {
		if n.NodeId != req.FailedNodeId {
			updated = append(updated, n)
		}
	}
	s.membershipList = updated
	s.membershipMutex.Unlock()

	s.notifySubscribers()

	if req.SenderId == s.nodeID {
		go s.BroadcastFailureToAll(req.FailedNodeId)
	}

	return &pb.BroadcastFailureResponse{Success: true}, nil
}

// BroadcastFailureToAll multicasts the failure notification to all other nodes.
func (s *DisseminationServer) BroadcastFailureToAll(failedNodeId int32) {
	s.membershipMutex.RLock()
	nodes := make([]*pb.NodeInfo, len(s.membershipList))
	copy(nodes, s.membershipList)
	s.membershipMutex.RUnlock()

	for _, node := range nodes {
		if node.NodeId == s.nodeID {
			continue
		}
		address := fmt.Sprintf("%s:%d", node.Host, node.Port)
		fmt.Printf("Component Dissemination of Node %d sends RPC BroadcastFailure to Node %d\n", s.nodeID, node.NodeId)
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			log.Printf("Failed to connect to node %d for BroadcastFailure: %v", node.NodeId, err)
			continue
		}
		client := pb.NewSwimServiceClient(conn)
		req := &pb.BroadcastFailureRequest{
			SenderId:     s.nodeID,
			FailedNodeId: failedNodeId,
		}
		_, err = client.BroadcastFailure(context.Background(), req)
		conn.Close()
		if err != nil {
			log.Printf("BroadcastFailure RPC to node %d failed: %v", node.NodeId, err)
		}
	}
}

// JoinCluster allows a new node to join by contacting a bootstrap node.
func (s *DisseminationServer) JoinCluster(bootstrapHost string, bootstrapPort int, myHost string, myPort int, myNodeID int32) {
	address := fmt.Sprintf("%s:%d", bootstrapHost, bootstrapPort)
	fmt.Printf("Component Dissemination of Node %d sends RPC Join to bootstrap node\n", myNodeID)
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
	s.notifySubscribers()
}

func main() {
	nodeIDStr := os.Getenv("NODE_ID")
	var nodeID int32 = 1
	if nodeIDStr != "" {
		if val, err := strconv.Atoi(nodeIDStr); err == nil {
			nodeID = int32(val)
		}
	}
	dcPort := DC_BASE + int(nodeID)

	membershipStr := os.Getenv("MEMBERSHIP") // Format: "1:node1:60051,2:node2:60052,3:node3:60053,4:node4:60054,5:node5:60055"
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
	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", dcPort))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", dcPort, err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSwimServiceServer(grpcServer, srv)
	log.Printf("Dissemination component of Node %d listening on port %d", nodeID, dcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
