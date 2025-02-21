package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	pb "q5pranav/swim"

	"google.golang.org/grpc"
)

// DisseminationComponent implements the gRPC server methods for Q5.
type DisseminationComponent struct {
	pb.UnimplementedSwimServiceServer
	nodeID                int32
	membershipList        []int32
	mu                    sync.Mutex
	failureDetectorClient pb.InternalServiceClient
}

// NewDisseminationComponent is the constructor for DisseminationComponent.
func NewDisseminationComponent(nodeID int32, failureDetectorClient pb.InternalServiceClient) *DisseminationComponent {
	return &DisseminationComponent{
		nodeID:                nodeID,
		membershipList:        []int32{nodeID},
		failureDetectorClient: failureDetectorClient,
	}
}

// Multicast handles incoming multicast about a failed node.
func (d *DisseminationComponent) Multicast(ctx context.Context, req *pb.MulticastRequest) (*pb.MulticastResponse, error) {
	fmt.Printf("Component Dissemination of Node %d runs RPC Multicast called by Component Dissemination of Node %d\n", d.nodeID, req.SenderId)

	d.mu.Lock()
	defer d.mu.Unlock()

	// Remove the failed node from the membership list
	for i, nodeID := range d.membershipList {
		if nodeID == req.FailedNodeId {
			d.membershipList = append(d.membershipList[:i], d.membershipList[i+1:]...)
			break
		}
	}

	// Update the Failure Detector Component
	_, err := d.failureDetectorClient.UpdateMembershipList(ctx, &pb.UpdateMembershipListRequest{
		Nodes: []*pb.NodeStatus{{NodeId: req.FailedNodeId, IsAlive: false}},
	})
	if err != nil {
		log.Printf("Failed to update Failure Detector: %v", err)
	}

	return &pb.MulticastResponse{Success: true}, nil
}

// Join handles a new node's request to join the network.
func (d *DisseminationComponent) Join(ctx context.Context, req *pb.JoinRequest) (*pb.JoinResponse, error) {
	fmt.Printf("Component Dissemination of Node %d runs RPC Join called by Component Dissemination of Node %d\n", d.nodeID, req.NewNodeId)

	d.mu.Lock()
	defer d.mu.Unlock()

	// Add the new node to the membership list
	d.membershipList = append(d.membershipList, req.NewNodeId)

	// Update the Failure Detector Component
	_, err := d.failureDetectorClient.UpdateMembershipList(ctx, &pb.UpdateMembershipListRequest{
		Nodes: []*pb.NodeStatus{{NodeId: req.NewNodeId, IsAlive: true}},
	})
	if err != nil {
		log.Printf("Failed to update Failure Detector: %v", err)
	}

	return &pb.JoinResponse{MembershipList: d.membershipList}, nil
}

// multicastFailure illustrates how you would notify other nodes about a failed node.
// (Not strictly called yet unless integrated with your failure detection logic.)
func (d *DisseminationComponent) multicastFailure(failedNodeID int32) {
	d.mu.Lock()
	membershipList := make([]int32, len(d.membershipList))
	copy(membershipList, d.membershipList)
	d.mu.Unlock()

	for _, nodeID := range membershipList {
		if nodeID != d.nodeID && nodeID != failedNodeID {
			go func(targetID int32) {
				// Example: If you wanted to dial each node's DC on port 60050 + targetID
				targetDPort := 60050 + targetID
				conn, err := grpc.Dial(fmt.Sprintf("localhost:%d", targetDPort), grpc.WithInsecure())
				if err != nil {
					log.Printf("Failed to connect to node %d DC: %v", targetID, err)
					return
				}
				defer conn.Close()

				client := pb.NewSwimServiceClient(conn)
				fmt.Printf("Component Dissemination of Node %d sends RPC Multicast to Component Dissemination of Node %d\n", d.nodeID, targetID)
				_, err = client.Multicast(context.Background(), &pb.MulticastRequest{
					SenderId:     d.nodeID,
					FailedNodeId: failedNodeID,
				})
				if err != nil {
					log.Printf("Failed to multicast to node %d: %v", targetID, err)
				}
			}(nodeID)
		}
	}
}

// Run starts the gRPC server for the Dissemination Component.
func (d *DisseminationComponent) Run() {
	// Compute the port for this node's Dissemination Component
	dPort := 60050 + d.nodeID
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", dPort))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", dPort, err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSwimServiceServer(grpcServer, d)

	log.Printf("Dissemination Component of Node %d is running on port %d", d.nodeID, dPort)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Graceful shutdown logic
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	grpcServer.GracefulStop()
	log.Println("Dissemination Component stopped gracefully")
}

func main() {
	// Seed the RNG (only used if NODE_ID is not set).
	rand.Seed(time.Now().UnixNano())

	// 1. Read NODE_ID from environment, fallback to random if not set
	var nodeID int32
	if envID := os.Getenv("NODE_ID"); envID != "" {
		id, err := strconv.Atoi(envID)
		if err != nil {
			log.Fatalf("Invalid NODE_ID: %v", err)
		}
		nodeID = int32(id)
	} else {
		nodeID = int32(rand.Intn(1000))
	}

	// If this node is not the bootstrap node (1), wait before connecting
	// to ensure node 1 has time to start
	if nodeID != 1 {
		log.Printf("Node %d sleeping for 10s to let bootstrap node 1 come up...", nodeID)
		time.Sleep(10 * time.Second)
	}

	// 2. Compute the Failure Detector port and the Dissemination port for this node
	fdPort := 50050 + nodeID
	// The DC port is computed inside Run()

	// 3. Connect to the Failure Detector Component (Python) on localhost:<fdPort>,
	//    but allow more time for the connection.
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	fdAddr := fmt.Sprintf("localhost:%d", fdPort)
	conn, err := grpc.DialContext(ctx, fdAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to Failure Detector (%s): %v", fdAddr, err)
	}
	defer conn.Close()

	failureDetectorClient := pb.NewInternalServiceClient(conn)
	disseminationComponent := NewDisseminationComponent(nodeID, failureDetectorClient)

	// 4. Simulate joining the network.
	//    The bootstrap node is assumed to be node 1, so its DC is on port 60051.
	bootstrapNodeID := int32(1)
	bootstrapDPort := 60050 + bootstrapNodeID

	// Only attempt to join if this node is not the bootstrap node itself
	if nodeID != bootstrapNodeID {
		ctxJoin, cancelJoin := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancelJoin()
		bootstrapAddr := fmt.Sprintf("localhost:%d", bootstrapDPort)

		joinConn, err := grpc.DialContext(ctxJoin, bootstrapAddr, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Printf("Failed to connect to bootstrap node %d at %s: %v", bootstrapNodeID, bootstrapAddr, err)
		} else {
			defer joinConn.Close()
			client := pb.NewSwimServiceClient(joinConn)

			fmt.Printf("Component Dissemination of Node %d sends RPC Join to Component Dissemination of Node %d\n", nodeID, bootstrapNodeID)
			resp, err := client.Join(ctxJoin, &pb.JoinRequest{NewNodeId: nodeID})
			if err != nil {
				log.Printf("Failed to join the network: %v", err)
			} else {
				// Update our membership list with the bootstrap node's response
				disseminationComponent.mu.Lock()
				disseminationComponent.membershipList = resp.MembershipList
				disseminationComponent.mu.Unlock()
			}
		}
	} else {
		// If this node is node 1, it acts as the bootstrap node and has no one to join to.
		// Its membershipList already contains itself.
		log.Printf("Node %d is the bootstrap node; skipping join call.", nodeID)
	}

	// 5. Run the Dissemination Component
	go disseminationComponent.Run()

	// Graceful shutdown in main
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	fmt.Println("Main shutting down gracefully...")
	// Implement any additional cleanup logic here
}
