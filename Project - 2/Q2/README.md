# Two-Phase Commit (2PC) with Docker + gRPC

This project implements the Two-Phase Commit (2PC) protocol using Docker and gRPC. It includes a **Coordinator** node and multiple **Participant** nodes that communicate via gRPC to decide whether to commit or abort a transaction.

---

## **Table of Contents**
- [Overview](#overview)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Generating Proto Stubs](#generating-proto-stubs)
- [Building Docker Images](#building-docker-images)
- [Running the Project](#running-the-project)
  - [Random Mode](#random-mode)
  - [Force Commit Mode](#force-commit-mode)
- [Stopping and Cleaning Up](#stopping-and-cleaning-up)
- [File Structure](#file-structure)
- [Extensions](#extensions)

---

## **Overview**

The Two-Phase Commit protocol consists of two phases:
1. **Voting Phase**:  
   The coordinator sends a `RequestVote` RPC to each participant. Participants respond with either "COMMIT" or "ABORT."
2. **Decision Phase**:  
   If all participants vote "COMMIT," the coordinator sends a global commit decision. Otherwise, it sends a global abort decision.

Participants can be configured to always commit by setting the `FORCE_COMMIT=true` environment variable.

---

## **Project Structure**

```
.
├── coordinator.py          # Coordinator logic for initiating transactions
├── participant.py          # Participant logic for handling votes
├── global_decision.go      # Go-based decision phase logic
├── two_phase_commit.proto  # Protobuf definitions for gRPC services
├── Dockerfile.coordinator  # Dockerfile for the coordinator service
├── Dockerfile.participant  # Dockerfile for the participant service
├── docker-compose.yml      # Docker Compose file to manage services
└── README.md               # Project documentation
```

---

## **Prerequisites**

1. **Docker** installed and running.
2. **Python 3.x** installed (optional, if you want to run scripts locally).
3. **Go** installed (optional, if you want to modify and compile `global_decision.go` locally).

---

## **Generating Proto Stubs**

To generate Python and Go gRPC stubs from the `two_phase_commit.proto` file:

1. Install the required tools:
   ```
   pip install grpcio grpcio-tools
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```

2. Generate Python stubs:
   ```
   python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. two_phase_commit.proto
   ```

3. Generate Go stubs:
   ```
   protoc --go_out=./twopc_go_pb --go_opt=paths=source_relative \
          --go-grpc_out=./twopc_go_pb --go-grpc_opt=paths=source_relative \
          two_phase_commit.proto
   ```

---

## **Building Docker Images**

Build the coordinator and participant images using the following commands:

```
docker build -t twopc_coordinator_image -f Dockerfile.coordinator .
docker build -t twopc_participant_image -f Dockerfile.participant .
```

---

## **Running the Project**

### **Random Mode**

To start the project in random mode (default behavior where participants randomly vote COMMIT or ABORT):

```
docker-compose up
```

### **Force Commit Mode**

To force all participants to vote COMMIT during both phases, use:

```
FORCE_COMMIT=true docker-compose up
```

---

## **Stopping and Cleaning Up**

To stop all containers and clean up resources:

```
docker-compose down --remove-orphans
```

If you want to manually remove containers and networks:

```
docker rm -f coordinator participant1 participant2 participant3 participant4 participant5
docker network rm twopc_net
```

---

## **File Structure**

### Coordinator (`coordinator.py`)
- Sends `RequestVote` RPCs to participants during the voting phase.
- Collects votes and sends them to the Go decision service via `SendVotes`.
- Connects to participants on ports `50051–50055`.

### Participant (`participant.py`)
- Handles voting requests (`RequestVote`) from the coordinator.
- Can be configured with `FORCE_COMMIT=true` to always vote COMMIT.

### Decision Service (`global_decision.go`)
- Implements the decision phase logic.
- Listens on ports `60051–60055` for global decisions.

### Protobuf Definitions (`two_phase_commit.proto`)
Defines gRPC services:
1. `RequestVote`
2. `GlobalDecision`
3. `SendVotes`

---