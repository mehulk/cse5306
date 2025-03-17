package main

import (
	"bufio"
	"context"
	"fmt"
	"grpc-average/generated"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
)

func main() {
	// Connect to the server
	serverAddr := os.Getenv("SERVER_ADDRESS")
	if serverAddr == "" {
		serverAddr = "localhost:50051" // Fallback to localhost if not set
	}

	conn, err := grpc.Dial(serverAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %v", err)
	}
	defer conn.Close()

	// Create the gRPC client
	client := generated.NewCalculatorClient(conn)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nChoose functionality:")
		fmt.Println("1. Average of two numbers")
		fmt.Println("2. Running average (streaming)")
		fmt.Println("3. Exit")
		fmt.Print("Enter your choice (1, 2, or 3): ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			calculateAverage(client, reader)
		case "2":
			calculateRunningAverage(client, reader)
		case "3":
			fmt.Println("Exiting the program.")
			return
		default:
			fmt.Println("Invalid choice. Please select 1, 2, or 3.")
		}
	}
}

func calculateAverage(client generated.CalculatorClient, reader *bufio.Reader) {
	var num1, num2 int32
	fmt.Print("Enter first number: ")
	fmt.Scanf("%d", &num1)
	fmt.Print("Enter second number: ")
	fmt.Scanf("%d", &num2)

	response, err := client.Average(context.Background(), &generated.AverageRequest{
		Num1: num1,
		Num2: num2,
	})
	if err != nil {
		log.Printf("Could not calculate average: %v", err)
		return
	}
	fmt.Printf("The average of %d and %d is: %.2f\n", num1, num2, response.Result)
}

func calculateRunningAverage(client generated.CalculatorClient, reader *bufio.Reader) {
	stream, err := client.RunningAverage(context.Background())
	if err != nil {
		log.Printf("Could not start stream: %v", err)
		return
	}

	fmt.Println("Enter numbers one by one (type 'done' to finish):")

	for {
		fmt.Print("Enter a number: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "done" {
			break
		}

		var number int32
		_, err := fmt.Sscanf(input, "%d", &number)
		if err != nil {
			fmt.Println("Invalid input. Please enter a number.")
			continue
		}

		if err := stream.Send(&generated.NumberRequest{Number: number}); err != nil {
			log.Printf("Failed to send number: %v", err)
			return
		}
	}

	response, err := stream.CloseAndRecv()
	if err != nil {
		log.Printf("Failed to receive response: %v", err)
		return
	}

	fmt.Printf("Running average is: %.2f\n", response.Result)
}
