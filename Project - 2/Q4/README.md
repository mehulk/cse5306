# Raft Log Replicator (Q4)

This project implements a simplified version of Raft’s log replication mechanism using two languages and containerization. The Python components handle leader election, heartbeat-based log appending, and client request forwarding; whereas the Go components provide the log replicator that receives logs from Python nodes and commits the operations.

---

## Table of Contents

1. [Overview](#overview)
2. [Project Structure](#project-structure)
3. [Requirements](#requirements)
4. [Setup & Build](#setup--build)
5. [Running the System](#running-the-system)
6. [Detailed Data Flow](#detailed-data-flow)
7. [Technical Details](#technical-details)
8. [Troubleshooting](#troubleshooting)
9. [Test Case](#test-case)
10. [License](#license)

---

## Overview

The system is divided into two main parts:

- **Python Raft Nodes:**  
  - **Role:** Handle leader election using randomized election timeouts (between 1.5 and 3 seconds), maintain a Raft log, process client requests, and forward requests to the current leader if needed.
  - **Heartbeats:** The leader sends periodic `AppendEntries` RPCs (every 1 second) as heartbeats and to propagate log entries to followers.

- **Go Log Replicator Nodes:**  
  - **Role:** Each Python node replicates its current log to its paired Go node using the `ReplicateLog` gRPC call. The Go node maintains the log, updates the commit index, and executes pending operations.
  - **Commitment:** Upon receiving the replicated log with an advanced commit index, the Go node prints execution messages to confirm that the operation has been applied.

This architecture ensures that:
- Client operations are correctly processed and committed.
- Log entries are distributed from the elected Python leader to all followers.
- Every Python node synchronizes its log with its corresponding Go node for local replication and operation execution.

---

## Project Structure

Below is the directory layout of the project:

```
Q4/
├── .DS_Store
├── README.md
├── go/
│   ├── Dockerfile             # Dockerfile for the Go log replicator server
│   ├── docker-compose.yml     # Docker Compose for the Go nodes
│   ├── go.mod
│   ├── go.sum
│   ├── raft.proto             # Protobuf definition (can be used by Go)
│   ├── server.go              # Go implementation of the log replication server
│   └── proto/                 # Generated proto files for Go
│       ├── raft.pb.go
│       └── raft_grpc.pb.go
└── python/
    ├── .python-version
    ├── .venv/                 # Virtual environment directory
    ├── Dockerfile             # Dockerfile for the Python nodes
    ├── Dockerfile.client      # Dockerfile for the test client image
    ├── docker-compose.yml     # Docker Compose for the Python nodes and client
    ├── proto/                 # Protobuf definitions for Python (e.g., raft.proto)
    │   └── raft.proto
    ├── raft/                  # (Optional) additional Raft-related source files
    ├── raft-election.py       # Python implementation for Raft election and log replication forwarding
    ├── raft_pb2.py            # Generated proto module for Python
    ├── raft_pb2_grpc.py       # Generated gRPC module for Python
    ├── requirements.txt
    ├── sitecustomize.py       # Polyfill for generated stubs that expect protobuf‑5.x runtime
    └── testclient.py          # Python test client to submit operations to the Raft cluster
```

---

## Requirements

- **Docker** and **Docker Compose**
- **Go** (for the log replicator)
- **Python 3.9+**
- Open ports (ensure external ports are available):
  - Python nodes: **5001–5005**
  - Go nodes: **6001–6005**

---

## Setup & Build

1. **Clone the Repository:**

   ```bash
   git clone <repository_url>
   cd Q4
   ```

2. **Build & Start the Containers:**

   Use Docker Compose in both the Go and Python directories as needed. At the project root (if a top-level docker-compose exists) or within each subdirectory, run:

   ```bash
   docker-compose up --build
   ```

   This will:
   - Build the Go server image using `Dockerfile` in the **go/** directory.
   - Build the Python client and server images using `Dockerfile` and `Dockerfile.client` in the **python/** directory.
   - Launch containerized services for at least 5 Python nodes (Raft election and request handling), 5 Go log replicator nodes, and the test client.

---

## Running the System

When the system starts:

- **Python Nodes:**
  - Start as followers and, after the randomized election timeout, one becomes the leader.
  - Handle client operations via `SubmitClientRequest` RPC. Non-leader nodes forward requests to the leader.
  - Broadcast heartbeats (using `AppendEntries`) to synchronize log entries.

- **Go Log Replicator Nodes:**
  - Each Python node invokes a helper function to call `ReplicateLog` on its paired Go node.
  - The Go nodes update their commit index and execute log entries, printing confirmations.

- **Test Client:**
  - The test client container submits operations (e.g., `operation-1`, `operation-2`, …) every 10 seconds.
  - This submission triggers the complete data flow through the cluster.

---

## Detailed Data Flow

The complete data flow in the system is as follows:

1. **Client Request Submission:**
   - The client sends an RPC (`SubmitClientRequest`) to a Python node.
   - **If the node is not the leader,** it forwards the request to the current leader.
   - **Log Example:**
     ```
     Client sends RPC SubmitClientRequest to Node Cluster Leader.
     [CLIENT] Leader node1:5001 accepted → Operation 'operation‑1' committed at index 1
     ```

2. **Handling by the Python Leader:**
   - The leader appends the client operation to its log, using the format `<operation, term, index>`.
   - It broadcasts an `AppendEntries` RPC (as a heartbeat) to all followers.
   - **Log Example:**
     ```
     Leader Node 2 sends AppendEntries (Heartbeat) to Node 3
     Node 3 runs RPC AppendEntries called by Node 2
     ```

3. **Update by Python Followers:**
   - On receiving the `AppendEntries` RPC, followers update their log, commit index, and reset their election timer.
   - They may also forward any client request received to the leader.

4. **Replication to Local Go Nodes:**
   - After processing log updates, each Python node calls its `_replicate_to_go` helper.
   - This helper issues a `ReplicateLog` RPC to the corresponding Go node (e.g., `q4node1` for Python node 1).
   - **Log Example:**
     ```
     Node 1 sends RPC ReplicateLog to Node q4node1
     q4node1-1  | Node 1 runs RPC ReplicateLog called by Node 1
     q4node1-1  | Node 1: commitIndex=1, total entries=1
     ```

5. **Go Node Log Replication & Execution:**
   - The Go log replicator receives the RPC, compares and updates its log and commit index.
   - When the commit index advances, the Go node executes the corresponding log entry.
   - **Log Example:**
     ```
     q4node2-1  | Node 2 runs RPC ReplicateLog called by Node 2
     q4node2-1  | Node 2: commitIndex=1, total entries=1
     q4node2-1  | Node 2 executes op=operation‑1 at index=1
     ```

6. **Continuous Heartbeats & Acknowledgements:**
   - The leader and followers continue to exchange heartbeats to maintain synchronization.
   - The continuous replication RPCs help ensure that all Go nodes eventually reach the same state and commit the operations.
   - This redundancy in replication helps maintain consistency across the system.

Thus, the complete data path is:

**Client → Python Node (Leader Election, Log Append & Heartbeats) → Python Followers (State Update) → Local Go Nodes (Log Replication, Commit & Execution)**

---

## Technical Details

- **gRPC Communication:**  
  Both Python and Go components use gRPC as defined in `raft.proto` to exchange messages such as `AppendEntries`, `RequestVote`, and `ReplicateLog`.

- **Timeouts:**  
  - **Heartbeat Timeout:** 1 second (all nodes)  
  - **Election Timeout:** Randomized between 1.5 and 3 seconds for each Python node

- **Containerization:**  
  Docker Compose is used to run:
  - 5 Python nodes for leader election and request processing.
  - 5 Go nodes for log replication.
  - A client container for test submissions.

---

## Troubleshooting

- **Network Issues:**  
  Verify the creation of the Docker network (e.g., `raft-net` or `raft-network`) and that nodes can resolve container names.

- **Port Conflicts:**  
  If ports 5001–5005 (Python) or 6001–6005 (Go) are in use, adjust the port mappings in the respective `docker-compose.yml`.

- **Rebuild After Changes:**  
  To rebuild images after code modifications, run:
  ```bash
  docker-compose down
  docker-compose up --build
  ```

- **Viewing Logs:**  
  Use:
  ```bash
  docker-compose logs -f
  ```
  to monitor real-time log output from all containers for debugging.

---

## Test Case

**Scenario: Single Operation End-to-End Replication**

**Test Steps:**

1. **Start the Cluster:**

   ```bash
   docker-compose up --build
   ```
   Allow 10 seconds for the Python nodes to complete leader election and stabilize.

2. **Submit a Client Operation:**

   The test client (container built from `Dockerfile.client`) automatically submits an operation (e.g., `operation-1`) every 10 seconds. You should see in the logs:
   ```
   Client sends RPC SubmitClientRequest to Node Cluster Leader.
   [CLIENT] Leader nodeX:500X accepted → Operation 'operation‑1' committed at index 1
   ```

3. **Monitor Python Node Logs:**

   - If a node other than the leader receives the request, it forwards it to the leader:
     ```
     Node 1 runs RPC SubmitClientRequest called by a client.
     Node 1 forwards client req to Leader 2
     ```
   - The leader appends the operation to its log and sends heartbeats:
     ```
     Leader Node 2 sends AppendEntries (Heartbeat) to Node 3
     Node 3 runs RPC AppendEntries called by Node 2
     ```

4. **Verify Go Node Replication:**

   - Check the logs of the corresponding Go node (e.g., `q4node2` if node 2 is the leader):
     ```
     q4node2-1  | Node 2 runs RPC ReplicateLog called by Node 2
     q4node2-1  | Node 2: commitIndex=1, total entries=1
     q4node2-1  | Node 2 executes op=operation‑1 at index=1
     ```

**Expected Outcome:**

- The client operation is processed and the leader commits the operation.
- The operation is replicated to the associated Go node which updates its commit index and executes the operation.
- The logs at every stage confirm the complete flow:
  
  **Client → Python Leader (and any forwarding) → Python Followers (via AppendEntries) → Local Go Node (via ReplicateLog)**
