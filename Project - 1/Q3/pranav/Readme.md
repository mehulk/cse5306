gRPC Calculator Project
This project demonstrates a simple gRPC-based calculator service with a Go server and a Python client. The service provides two functionalities:

Unary RPC: Calculates the average of two numbers.
Server Streaming RPC: Computes the running average of a stream of numbers.
Both the server and client are containerized using Docker for easy deployment and portability.

Project Structure
 
 
Q3Pranav/
├── go-grpc-server/
│   ├── proto/
│   │   ├── average.proto          # Protocol Buffers definition file
│   │   ├── average.pb.go          # Generated Go code for messages
│   │   ├── average_grpc.pb.go     # Generated Go code for gRPC service
│   ├── Dockerfile.server          # Dockerfile for the Go server
│   ├── go.mod                     # Go module dependencies
│   ├── go.sum                     # Go module checksum file
│   ├── server.go                  # Go server implementation
└── python client/
    ├── __pycache__/               # Compiled Python files (auto-generated)
    ├── average_pb2.py             # Generated Python code for messages
    ├── average_pb2_grpc.py        # Generated Python code for gRPC service
    ├── calculator_client.py       # Python client implementation
    ├── Dockerfile.client          # Dockerfile for the Python client
    ├── requirements.txt           # Python dependencies
How to Run the Project
Step 1: Generate Code from the .proto File
Ensure that the necessary gRPC code is generated before building the application.

For Go (Server):
 
 
protoc --go_out=. --go-grpc_out=. proto/average.proto
For Python (Client):
 
 
python -m grpc_tools.protoc -I=proto --python_out="python client" --grpc_python_out="python client" proto/average.proto
Step 2: Build Docker Images
Navigate to the project’s root directory and build the Docker images for both the server and client.

Build Server Image:
From the go-grpc-server directory (or from the root if the Dockerfile path is adjusted):

 
 
docker build -t grpc-calculator-server -f Dockerfile.server .
Build Client Image:
From the root directory:

 
 
docker build -t grpc-calculator-client -f "python client/Dockerfile.client" "python client/"
Note:
The Python client’s Dockerfile uses an ENTRYPOINT so that any arguments passed to the container will be forwarded to the client script.

Step 3: Create a Docker Network
Create a custom Docker network to allow the server and client containers to communicate.

 
 
docker network create grpc-net
Step 4: Run the Server Container
Run the Go server container on port 50051 and attach it to the grpc-net network.

 
 
docker run --rm --name grpc-server --network grpc-net -p 50051:50051 grpc-calculator-server
Step 5: Run the Client Container
Since both containers are on the same network, update the client code to connect to grpc-server:50051 (instead of localhost:50051). After updating and rebuilding the client image, run the container with the desired arguments:

 
 
docker run --rm --name grpc-client --network grpc-net grpc-calculator-client 10 20 30 40 50
The first two numbers (10 and 20) are used for the unary RPC (calculating the average).
The remaining numbers (30, 40, 50) are used for the streaming RPC (running average).
How It Works
1. Unary RPC (Average)
The Average method calculates the average of two numbers provided by the client.

Example Input: 10, 20
Example Output: 15.0
2. Server Streaming RPC (RunningAverage)
The RunningAverage method computes a running average from a stream of numbers sent by the client.

Example Input: 30, 40, 50
Example Output:
After receiving 30: Running Average = 30.0
After receiving 40: Running Average = 35.0
After receiving 50: Running Average = 40.0
The final result is sent back after all numbers are processed.
Key Features
gRPC Communication:

Uses HTTP/2 for efficient communication.
Supports both unary and streaming RPCs.
Protocol Buffers:

Defines structured data exchange between server and client.
Enables cross-language compatibility (Go server, Python client).
Dockerized Deployment:

Ensures consistent environments across systems.
Simplifies deployment using containerized applications.
Scalability:

Can be scaled horizontally by running multiple instances of the server behind a load balancer.
Testing Scenarios
1. Test Unary RPC with different pairs of numbers:
 
 
docker run --rm --network grpc-net grpc-calculator-client 10 20
2. Test Streaming RPC with multiple numbers:
 
 
docker run --rm --network grpc-net grpc-calculator-client 10 20 30 40 50
3. Observe Logs to Understand the Request/Response Flow:
 
 
docker logs grpc-server
docker logs grpc-client
Clean Up
After testing, remove all containers and the custom network:

 
 
docker rm -f grpc-server grpc-client
docker network rm grpc-net