# gRPC Average Calculator

This project implements a gRPC-based client-server application to calculate averages. The server is written in Python, and the client is implemented in Go.

## Project Structure
- `server.py` - Python gRPC server
- `client.go` - Go gRPC client
- `average.proto` - Protocol Buffer definition
- `generated/` - Contains generated gRPC and Protobuf files
- `Dockerfile.server` - Dockerfile for the server
- `Dockerfile.client` - Dockerfile for the client
- `docker-compose.yml` - Defines multi-container setup

## Requirements
- Python 3
- Go
- gRPC and Protobuf tools
- Docker & Docker Compose

## Setup and Usage

### Generate gRPC Code
Before running the application, generate gRPC code:

#### For Python Server:
python -m grpc_tools.protoc -I. --python_out=generated --grpc_python_out=generated average.proto


#### For Go Client:
protoc --go_out=generated --go-grpc_out=generated average.proto


### Running the Server
python server.py


### Running the Client
go run client.go

## Dockerization

### Build Docker Images
docker build -t grpc-server -f Dockerfile.server .
docker build -t grpc-client -f Dockerfile.client .

### Running Docker Containers
#### For Server (needs to be run first)
docker run -d --name grpc-server -p 50051:50051 grpc-server

#### For Client
docker run -it --name grpc-client grpc-client     

## gRPC Services

### 1. Unary RPC: `Average`
- Computes the average of two numbers.
- Request: `AverageRequest { num1, num2 }`
- Response: `AverageResponse { result }`

### 2. Client Streaming RPC: `RunningAverage`
- Computes the running average of a stream of numbers.
- Request: `NumberRequest { number }`
- Response: `AverageResponse { result }`


