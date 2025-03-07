package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "go-grpc-server/proto" // Import the generated proto package
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCalculatorServer
}

// Unary RPC: Computes the average of two numbers
func (s *server) Average(ctx context.Context, req *pb.AverageRequest) (*pb.AverageResponse, error) {
	log.Printf("Received Average request: num1=%d, num2=%d", req.Num1, req.Num2)
	result := float64(req.Num1+req.Num2) / 2.0
	log.Printf("Sending Average response: result=%.2f", result)
	return &pb.AverageResponse{Result: result}, nil
}

// Server Streaming RPC: Computes the running average of a stream of numbers
func (s *server) RunningAverage(stream pb.Calculator_RunningAverageServer) error {
	var sum int32
	var count int32

	log.Println("Started RunningAverage stream")
	for {
		req, err := stream.Recv()
		if err != nil {
			result := float64(sum) / float64(count)
			log.Printf("Closing RunningAverage stream: sum=%d, count=%d, result=%.2f", sum, count, result)
			return stream.SendAndClose(&pb.AverageResponse{Result: result})
		}
		sum += req.Number
		count++
		log.Printf("Received number: %d, current sum=%d, count=%d", req.Number, sum, count)
	}
}

func main() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterCalculatorServer(grpcServer, &server{})

	fmt.Println("gRPC server is running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
