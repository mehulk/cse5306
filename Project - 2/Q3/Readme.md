# Raft Leader Election

This project implements a simplified version of the Raft consensus algorithm for leader election using Python and gRPC. The implementation includes features such as randomized election timeouts, leader heartbeats, and handling split votes.

---

## Features

- **Leader Election**: Nodes elect a leader when the current leader fails or when the cluster initializes.
- **Randomized Election Timeout**: Election timeouts are randomized to avoid split votes.
- **Heartbeat Mechanism**: The leader sends periodic heartbeats to maintain its leadership.
- **gRPC Communication**: Nodes communicate using gRPC for efficiency and scalability.
- **Containerized Deployment**: Nodes are containerized using Docker, enabling distributed deployment.

---

## Project Structure

### Files

1. `raft-election.py`: Main implementation of the Raft protocol.
2. `.proto/raft.proto`: Protocol buffer file defining gRPC services and messages.
3. `Dockerfile`: Dockerfile to containerize the application.
4. `docker-compose.yml`: Docker Compose file to run multiple nodes in containers.
5. `README.md`: Documentation for the project.

---

## Setup Instructions

### Prerequisites

- Python 3.9 or later
- gRPC and gRPC Tools (`grpcio`, `grpcio-tools`)
- Docker and Docker Compose

### Steps

1. **Clone the Repository**
git clone https://github.com/mehulk/cse5306.git
cd cse5306

2. **Install Dependencies**
Install Python dependencies:
pip install grpcio grpcio-tools

3. **Generate gRPC Code**
Compile the `.proto` file to generate Python gRPC code:
python -m grpc_tools.protoc -I./proto --python_out=. --grpc_python_out=. ./proto/raft.proto

4. **Build Docker Images**
Build the Docker image for the application:
docker-compose build

5. **Run the Application**
Start the cluster with 5 nodes using Docker Compose:
docker-compose up

6. **Observe Logs**
Logs will show nodes transitioning between follower, candidate, and leader states:
docker-compose logs -f

---

## Usage

### Running a Single Node Locally

To run a single node locally, use:
python raft-election.py -id=1 -port=5001 -peers=node2:5002,node3:5003,node4:5004,node5:5005

### Running Multiple Nodes with Docker Compose

The `docker-compose.yml` file defines 5 nodes (`node1`, `node2`, ..., `node5`). Each node communicates with its peers over a shared Docker network.

---

## Key Concepts

### 1. States of a Node:

- **Follower**: Listens for heartbeats from a leader or waits for an election timeout.
- **Candidate**: Starts an election after an election timeout and requests votes from peers.
- **Leader**: Sends heartbeats to followers and handles client requests.

### 2. Election Timeout:

- Randomized between 1.5s and 3s to reduce split votes.

### 3. Heartbeat Interval:

- Fixed at 1 second to maintain leadership.

### 4. RPC Methods:

- `AppendEntries`: Used by leaders to send heartbeats and replicate logs.
- `RequestVote`: Used by candidates to request votes during elections.
- `SubmitClientRequest`: Used by clients to submit operations (simplified in this implementation).

---

## Logging Format

### Client-Side Logging:

For each RPC method being called, logs are printed in the format:
Node <node_id> sends RPC <rpc_name> to Node <node_id>.

text

### Server-Side Logging:

For each RPC method being received, logs are printed in the format:
Node <node_id> runs RPC <rpc_name> called by Node <node_id>.

Example:
Node 1 sends RPC RequestVote to Node 2.
Node 2 runs RPC RequestVote called by Node 1.

---

## Example Output

After starting the cluster, you will see logs like this:
Node 1 starting as FOLLOWER with election timeout 2.04s
Node 2 starting as FOLLOWER with election timeout 2.55s
Node 3 starting as FOLLOWER with election timeout 1.70s

Node 3 election timeout expired; starting election
Node 3 becomes CANDIDATE for term 1
Node 3 sends RPC RequestVote to Node 1
Node 1 runs RPC RequestVote called by Node 3
Node 1 grants vote to Node 3

Node 3 is elected LEADER for term 1
Leader Node 3 sends HEARTBEAT to all nodes

---

## Testing

You can test leader failure by stopping a node (e.g., `node3`) and observing how other nodes elect a new leader.

To stop a node:
docker stop node3

To restart it:
docker start node3