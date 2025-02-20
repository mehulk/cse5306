package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"sync"
	"time"
	"os"
	"syscall"
	"os/signal"

	"google.golang.org/grpc"
	pb "q5pranav/swim"
)

type DisseminationComponent struct {
    pb.UnimplementedSwimServiceServer
    nodeID         int32
    membershipList []int32
    mu             sync.Mutex
    failureDetectorClient pb.InternalServiceClient
}

func NewDisseminationComponent(nodeID int32, failureDetectorClient pb.InternalServiceClient) *DisseminationComponent {
    return &DisseminationComponent{
        nodeID:         nodeID,
        membershipList: []int32{nodeID},
        failureDetectorClient: failureDetectorClient,
    }
}

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

func (d *DisseminationComponent) multicastFailure(failedNodeID int32) {
    d.mu.Lock()
    membershipList := make([]int32, len(d.membershipList))
    copy(membershipList, d.membershipList)
    d.mu.Unlock()

    for _, nodeID := range membershipList {
        if nodeID != d.nodeID && nodeID != failedNodeID {
            go func(targetID int32) {
                conn, err := grpc.Dial(fmt.Sprintf("node%d:50051", targetID), grpc.WithInsecure())
                if err != nil {
                    log.Printf("Failed to connect to node %d: %v", targetID, err)
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

func (d *DisseminationComponent) Run() {
    lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50051))
    if err != nil {
        log.Fatalf("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterSwimServiceServer(grpcServer, d)

    log.Printf("Dissemination Component of Node %d is running on port 50051", d.nodeID)
    
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
    log.Println("Server stopped gracefully")
}


func main() {
    rand.Seed(time.Now().UnixNano())
    nodeID := int32(rand.Intn(1000))
    
    // Connect to the Failure Detector Component
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    conn, err := grpc.DialContext(ctx, "localhost:50052", grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        log.Fatalf("Failed to connect to Failure Detector: %v", err)
    }
    defer conn.Close()
    failureDetectorClient := pb.NewInternalServiceClient(conn)

    disseminationComponent := NewDisseminationComponent(nodeID, failureDetectorClient)

    // Simulate joining the network
    bootstrapNodeID := int32(1) // Assume node 1 is the bootstrap node
    ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    conn, err = grpc.DialContext(ctx, fmt.Sprintf("node%d:50051", bootstrapNodeID), grpc.WithInsecure(), grpc.WithBlock())
    if err != nil {
        log.Fatalf("Failed to connect to bootstrap node: %v", err)
    }
    defer conn.Close()

    client := pb.NewSwimServiceClient(conn)
    fmt.Printf("Component Dissemination of Node %d sends RPC Join to Component Dissemination of Node %d\n", nodeID, bootstrapNodeID)
    resp, err := client.Join(ctx, &pb.JoinRequest{NewNodeId: nodeID})
    if err != nil {
        log.Fatalf("Failed to join the network: %v", err)
    }

    disseminationComponent.membershipList = resp.MembershipList

    // Run the Dissemination Component
    go disseminationComponent.Run()

    // Graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    <-sigChan

    fmt.Println("Shutting down gracefully...")
    // Implement any cleanup logic here
}
