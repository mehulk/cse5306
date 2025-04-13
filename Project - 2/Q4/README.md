# Raft Log Replicator (Q4)

This project implements a simplified version of Raft’s log replication mechanism using Go for the log replicator and Python for leader election and client handling. The system is containerized and demonstrates how client requests flow through the Python nodes (which handle leader election, heartbeat, and log appending) and then are replicated to corresponding local Go nodes.

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

This project is composed of two main parts:

- **Python Raft Nodes:**  
  These nodes handle leader election using randomized election timeouts (between 1.5 and 3 seconds) and maintain the Raft log state. All nodes begin as followers until one becomes the leader. Clients submit requests to any node; if the node is not the leader, it forwards the request to the elected leader. When processing a client request, the leader appends the new operation to its log and propagates it using heartbeat-based RPCs.

- **Go Log Replicator Nodes:**  
  Each Python node replicates its current log to an associated Go node via a gRPC service. The Go server receives these replicated logs, updates its commit index, and executes pending log entries. This process is part of the "local replication" mechanism within the Raft protocol.

The overall design ensures that:
- Client operations are directed to the current Python leader.
- Log entries are consistently broadcasted to all Python nodes.
- Each Python node continuously replicates its log state to its paired Go node.

---

## Project Structure

```
├── Dockerfile                # Dockerfile for the Go log replicator server
├── Dockerfile.client         # Dockerfile for the Python test client
├── docker-compose.yml        # Orchestration for Python nodes, Go nodes, and the test client
├── raft.proto                # Protobuf definition for gRPC messages and services
├── server.go                 # Go implementation of the log replicator server (Q4)
├── raft-election.py          # Python implementation for Raft leader election and client-forwarding
├── testclient.py             # Python test client to submit operations to the Raft cluster
└── sitecustomise.py          # Polyfill for Python generated stubs
```

**Key Files:**

- **raft.proto:**  
  Defines RPC messages such as `LogEntry`, `AppendEntriesRequest/Response`, `VoteRequest/Response`, and `ClientRequest/Response`. It also declares two services: **Raft** (for leader election and client handling) and **LogReplicator** (for log replication to Go nodes).

- **server.go:**  
  Contains the Go implementation of the log replicator. It listens for incoming RPCs, updates the log, commit index, and prints execution and acknowledgement messages.

- **raft-election.py:**  
  Handles the Raft leader election in Python, manages client requests (including forwarding if the node is not the leader), appends log entries, sends heartbeats, and replicates the current log to a corresponding Go node.

- **docker-compose.yml:**  
  Manages the containerized services for at least 5 Python nodes, 5 Go log replicator nodes, and a test client.

---

## Requirements

- **Docker** and **Docker Compose**
- Familiarity with **Go** and **Python**
- Open ports for mapping:
  - Python nodes: **5001–5005**
  - Go nodes: **6001–6005**

---

## Setup & Build

1. **Clone the Repository:**

   ```bash
   git clone <repository_url>
   cd <repository_directory>
   ```

2. **Build & Start the Containers:**

   From the project’s root directory (where the `docker-compose.yml` file is located), run:

   ```bash
   docker-compose up --build
   ```

   This will:
   - Build the Go server image using the provided `Dockerfile`.
   - Build the Python client image using `Dockerfile.client`.
   - Start 5 containerized Python nodes for Raft election and log management.
   - Start 5 Go log replicator nodes.
   - Start a test client container that periodically submits client requests.

---

## Running the System

When you start the system:
- **Python Nodes:**  
  They participate in leader election. One becomes the leader and handles incoming client requests. Each node uses gRPC to send and receive `AppendEntries` and `RequestVote` RPCs.
  
- **Go Log Replicator Nodes:**  
  Each Python node continuously replicates its log state to a paired Go node using the `ReplicateLog` RPC. The Go nodes update their commit index and execute pending operations.

- **Test Client:**  
  The test client periodically submits operations (e.g., every 10 seconds). These operations follow the full data path—from the client to the Python leader, then to the followers, and finally to the Go nodes for replication.

---

## Detailed Data Flow

The data flow in the project is as follows:

1. **Client Request Submission:**
   - A client sends an RPC (`SubmitClientRequest`) to a Python node.
   - **If the Python node is not the leader**, it forwards the request to the current leader.
   - **Example Log:**
     ```
     Client sends RPC SubmitClientRequest to Node Cluster Leader.
     [CLIENT] Leader node1:5001 accepted → Operation 'operation‑1' committed at index 1
     ```

2. **Python Leader Handling the Request:**
   - The leader appends the client's operation to its log with the format `<operation, term, index>`.
   - It then broadcasts an `AppendEntries` RPC (heartbeats) to all Python followers.
   - **Example Log:**
     ```
     Leader Node 2 sends AppendEntries (Heartbeat) to Node 3
     Node 3 runs RPC AppendEntries called by Node 2
     ```

3. **Python Followers Updating Their State:**
   - Followers update their local log and commit index on receiving an `AppendEntries` RPC.
   - They also update their last heartbeat timestamp.
   - **Forwarding Behavior:** If a follower receives a client request, it forwards the request to the known leader.
 
4. **Replication to Local Go Nodes:**
   - After processing any log update (heartbeat or new log entry), each Python node invokes the helper function (`_replicate_to_go`).
   - This helper sends the current log to the associated Go node via the `ReplicateLog` RPC.
   - **Example Log:**
     ```
     Node 1 sends RPC ReplicateLog to Node q4node1
     q4node1-1  | Node 1 runs RPC ReplicateLog called by Node 1
     q4node1-1  | Node 1: commitIndex=1, total entries=1
     ```

5. **Go Node Log Replication & Execution:**
   - The Go log replicator receives the replicated log, compares its current log with the incoming one, and updates its state.
   - The Go node then commits the operation(s) if the commit index (provided in the RPC) is advanced.
   - **Example Log:**
     ```
     q4node2-1  | Node 2 runs RPC ReplicateLog called by Node 2
     q4node2-1  | Node 2: commitIndex=1, total entries=1
     q4node2-1  | Node 2 executes op=operation‑1 at index=1
     ```

6. **Acknowledgements & Continuous Heartbeats:**
   - The Go nodes return acknowledgements, and the Python nodes log these replication actions.
   - The cyclic nature of the heartbeat RPCs ensures continuous synchronization between Python and Go nodes.
   - Multiple invocations of the replication RPC confirm that the log consistency is maintained over repeated heartbeats.

Thus, the complete data path in the system is:

**Client → Python Node (Leader Election & Request Handling) → Python Followers (Log Update via Heartbeats) → Local Go Nodes (Log Replication, Commit & Execution)**

---

## Technical Details

- **gRPC Communication:**  
  Both Python and Go nodes use gRPC to exchange RPC messages as defined in `raft.proto`.
  
- **Timeout Settings:**  
  - **Heartbeat Timeout:** 1 second  
  - **Election Timeout:** Randomized between 1.5 and 3 seconds per Python node
  
- **Containerization:**  
  Docker Compose is used to deploy the environment. Nodes communicate using the defined networks, and individual ports (5001–5005 for Python and 6001–6005 for Go) are mapped appropriately.

---

## Troubleshooting

- **Network Errors:**  
  Ensure the Docker network is properly created (e.g., `raft-net` or `raft-network`). Use `docker-compose logs` to view detailed output.
  
- **Port Conflicts:**  
  If the mapped ports are in use, modify the `ports` section in `docker-compose.yml`.
  
- **Rebuild After Changes:**  
  If you update any code, rebuild the containers with:
  ```bash
  docker-compose down
  docker-compose up --build
  ```

- **Debug Logs:**  
  Check container logs for detailed messages regarding RPC calls, state transitions, and replication events.

---

## Test Case

To validate the end-to-end functionality, the following test case can be executed:

1. **Scenario: Single Operation Replication Test**

   **Test Steps:**
   - Start the entire system:
     ```bash
     docker-compose up --build
     ```
   - Allow the cluster to stabilize for approximately 10 seconds so that the Python nodes perform leader election.
   - The test client (run inside the container built from `Dockerfile.client`) will start submitting operations (e.g., `operation-1`, `operation-2`, …) every 10 seconds.
   - **Expected Outcome:**
     - The client log shows:
       ```
       Client sends RPC SubmitClientRequest to Node Cluster Leader.
       [CLIENT] Leader nodeX:500X accepted → Operation 'operation‑1' committed at index 1
       ```
     - In the Python node logs:
       - If the request is not received on the leader, the node forwards it to the leader.
       - The leader appends the operation to its log and sends out heartbeats:
         ```
         Leader Node Y sends AppendEntries (Heartbeat) to Node Z
         Node Z runs RPC AppendEntries called by Node Y
         ```
     - In the Go node logs:
       - The corresponding Go node (e.g., `q4nodeY`) prints that it received the replication RPC, has updated its commit index, and executes the operation:
         ```
         q4nodeY-1  | Node Y runs RPC ReplicateLog called by Node Y
         q4nodeY-1  | Node Y: commitIndex=1, total entries=1
         q4nodeY-1  | Node Y executes op=operation‑1 at index=1
         ```

   **Verification:**
   - Check the logs in each container (Python and Go nodes) to ensure:
     - The client request was processed.
     - The leader appropriately forwarded and replicated the operation.
     - The local Go nodes received, updated, and executed the log entry.
   
2. **Additional Checks:**
   - Use `docker-compose logs -f` to continuously view live output.
   - If multiple operations are sent, each should be visible in the respective logs in sequential order.

This test confirms the complete data flow:
**Client → Python (Leader Election, Log Append, Heartbeat) → Python Followers → Go Node Replication & Execution**
