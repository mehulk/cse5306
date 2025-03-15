# Distributed SWIM Protocol Implementation

## Project Overview
This project implements the SWIM protocol in two parts using different programming languages:

### Q4 – Failure Detector Component (FD)
- Implemented in Python.
- Periodically sends direct and indirect pings to detect failed nodes.
- Logs RPC calls, e.g., “Component FailureDetector of Node X sends RPC Ping to Component FailureDetector of Node Y.”

### Q5 – Dissemination Component (DC)
- Implemented in Go.
- Maintains the membership list and disseminates updates.
- Supports:
  - **Join:** New nodes contact a bootstrap node to receive the membership list.
  - **BroadcastFailure:** Propagates failure updates.
  - **StreamMembership:** Streams membership updates to FDs.
- Both components use gRPC and are containerized to support multiple nodes (at least five).

## Authors
- **Pranav S Damu (1002158458)**
- **Mehul Kanotra (1002166093)**

## Project Structure

### **swim.proto**
Defines gRPC messages and service methods for both FD and DC:
- **Failure Detector RPCs:** Ping, IndirectPing
- **Dissemination RPCs:** Join, BroadcastFailure, StreamMembership

### **dissemination.go**
Go implementation of DC:
- Maintains membership list.
- Implements Join, BroadcastFailure, and StreamMembership RPCs.

### **failure_detector.py**
Python implementation of FD:
- Periodically pings random nodes.
- Uses direct and indirect pings for health checks.
- Reports failures to DC.
- Subscribes to membership updates via StreamMembership.

### **node.py**
Implements FD server to handle Ping and IndirectPing RPCs.

### **main.py**
- Starts the FD server and launches the FD client, subscribing to membership updates.

### **Dockerfile & docker-compose.yml**
- Containerizes the application with FD and DC communicating over localhost within the same container.

## Environment Setup & Build Instructions

### 1. Install Prerequisites
- Go (version 1.24 or later)
- Python 3.9
- Docker
- protoc (Protocol Buffers Compiler)

### 2. Install Go Plugins
```sh
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
export PATH="$PATH:$(go env GOPATH)/bin"
```

### 3. Initialize and Download Go Modules
```sh
go mod init github.com/pranavsdamu/damu-swim-disseem
go get google.golang.org/grpc
go get google.golang.org/protobuf
go mod download
go mod tidy
```

### 4. Generate gRPC Stubs
#### Go Stubs:
```sh
protoc \
  --go_out=./swim --go_opt=paths=source_relative \
  --go-grpc_out=./swim --go-grpc_opt=paths=source_relative \
  swim.proto
```
#### Python Stubs:
```sh
python3 -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. swim.proto
```

### 5. Build Docker Images
```sh
docker compose build
```

### 6. Running the Cluster
#### Start a 5-node cluster:
```sh
docker compose up
```

#### Add a new node (e.g., node6):
```sh
docker build -t swim-node .  # Build your image first

docker run -it --rm --name node6 --hostname node6 \
  --network q5_swim_network \
  -e NODE_ID=6 \
  -e BOOTSTRAP_NEEDED=true \
  -e BOOTSTRAP_ADDRESS=node1:60051 \
  swim-node

```

### 7. Managing the Cluster
#### List Docker Networks:
```sh
docker network ls
```

#### Stop a Specific Node:
```sh
docker compose stop nodeX  # Replace nodeX with the node name
```

#### Tear Down the Cluster:
```sh
docker compose down
```

## Project Details

### Q4 – Failure Detector Component (Python)
#### Functionality:
- Each node's FD pings a random node every T seconds.
- If direct ping fails, requests indirect pings from k other nodes.
- If no responses within T seconds, marks the node as failed.
- Logs all RPC calls.

#### Containerization:
- Runs in the same container as DC, communicating via gRPC on separate ports.

### Q5 – Dissemination Component (Go)
#### Functionality:
- Maintains membership list.
- Supports:
  - **Join:** New nodes contact a bootstrap node to receive the list.
  - **BroadcastFailure:** Disseminates failure information.
  - **StreamMembership:** Streams updates to FDs.

#### Containerization:
- Runs alongside FD in the same container.
- Uses Docker Compose to launch multiple nodes in a shared network.

## Useful Commands
### Generating Stubs
#### Go:
```sh
protoc \
  --go_out=./swim --go_opt=paths=source_relative \
  --go-grpc_out=./swim --go-grpc_opt=paths=source_relative \
  swim.proto
```

#### Python:
```sh
python3 -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. swim.proto
```

### Docker Commands
#### Build Images:
```sh
docker compose build
```

#### Start Cluster:
```sh
docker compose up -d
```

#### Stop Node:
```sh
docker compose stop nodeX  # Replace nodeX with the node name
```

#### Tear Down:
```sh
docker compose down
```

### Go Environment Setup
```sh
go mod init github.com/pranavsdamu/damu-swim-disseem
go mod download
go mod tidy

go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Python Environment Setup
Ensure Python 3.9 is installed and install dependencies:
```sh
pip install --no-cache-dir grpcio grpcio-tools
```

## Conclusion
This project implements a distributed SWIM protocol with failure detection (Python) and membership dissemination (Go). It leverages gRPC for communication and Docker for containerized deployment, enabling multiple nodes to run and interact efficiently.

