# Two-Phase Commit Protocol (2PC) with Docker + gRPC

This project demonstrates a simplified implementation of the Two-Phase Commit protocol using **Python**, **gRPC**, and **Docker**. It features:

1. A **Coordinator** node that sends vote-requests to participants.
2. **Participant** nodes that respond with a vote to commit or abort.
3. Optional **force commit** mode where participants always vote commit.

## Table of Contents

- [Overview of 2PC](#overview-of-2pc)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Building the Docker Images](#building-the-docker-images)
- [Running the Project](#running-the-project)
  - [Random Mode](#random-mode)
  - [Force Commit Mode](#force-commit-mode)
- [How It Works](#how-it-works)
  - [Coordinator](#coordinator)
  - [Participant](#participant)
  - [Script Behavior](#script-behavior)
- [Cleaning Up](#cleaning-up)
- [Extensions](#extensions)
- [License](#license)

---

## Overview of 2PC

Two-Phase Commit (2PC) is a distributed protocol to ensure all nodes in a system either commit a transaction or all abort, achieving atomicity across multiple nodes. It involves:

1. **Voting Phase**:
   - The Coordinator sends a vote-request to all participants.
   - Each Participant votes to **commit** or **abort** based on local conditions (or random choice in this demo).

2. **Decision Phase**:
   - If **any** participant votes **abort**, the Coordinator sends a *global abort* message.  
   - If **all** participants vote **commit**, the Coordinator sends a *global commit* message.

Here, we’re only focusing on the vote phase (and a simplified “commit or abort” outcome).

---

## Project Structure

- **`two_phase_commit.proto`**  
  Defines the gRPC services and messages for 2PC.

- **`participant.py`**  
  A gRPC server representing a participant. By default, votes randomly commit or abort. If `FORCE_COMMIT=true` is set, it always commits.

- **`coordinator.py`**  
  A Python script that connects to participants (via gRPC), requests votes, and logs whether the transaction commits or aborts.

- **`Dockerfile_coordinator`**, **`Dockerfile_participant`**  
  Docker build files for Coordinator and Participant images.

- **`run_twopc.sh`**  
  A convenience script to:
  1. Remove any existing containers,
  2. Create a Docker network (if it doesn’t exist),
  3. Spin up 5 participant containers (all random or all forced commit),
  4. Launch the coordinator container in interactive mode (so logs are visible),
  5. Clean up containers when finished.

---

## Prerequisites

1. **Docker** installed (and Docker daemon running).
2. **Python 3.11** or any version compatible with your environment (only needed if you’re directly testing locally without Docker).
3. **gRPC Tools** (the Dockerfiles already handle installing them internally).

---

## Building the Docker Images

From the root of the project directory, run:

```bash
docker build -t coordinator -f Dockerfile_coordinator .
docker build -t participant -f Dockerfile_participant .
```

This will produce two local Docker images named:
- **`coordinator`**
- **`participant`**

---

## Running the Project

### Random Mode

If you want **all participants** to vote **randomly** (commit or abort), run:

```bash
./run_twopc.sh
```

### Force Commit Mode

If you want **all participants** to **always commit**, run:

```bash
./run_twopc.sh force
```

In either case, you’ll see the coordinator logs in your terminal. Press **Ctrl + C** to stop the coordinator, which will then trigger container cleanup.

---

## How It Works

### Coordinator

- The coordinator loads a list of participant addresses (e.g., `participant1:50051`, etc.).
- It sends a `VoteRequest` to each participant with a transaction ID (like `"txn123"`).
- It collects `VoteResponse` messages. If **any** votes are `ABORT`, the coordinator prints “Aborting transaction.” Otherwise, it prints “All participants voted COMMIT. Proceeding to commit.”

### Participant

- Each participant runs a gRPC server on port `50051`.
- By default, it randomly chooses between **COMMIT** and **ABORT**.
- If it detects `FORCE_COMMIT=true` in the environment, it **always** returns **COMMIT** (0).
- Participants log their decisions to stdout, which you can see via `docker logs` or by attaching to the container.

### Script Behavior

1. **Stop old containers** if they exist (coordinator, participant1..5).  
2. **Create** the Docker network `twopc_net` if missing.  
3. **Spin up** five participant containers:
   - If `force` is passed, they run with `-e FORCE_COMMIT=true`.
   - Otherwise, they run without that env variable (random commits/aborts).
4. **Launch** the coordinator container in interactive mode. You see logs like:
   ```
   Coordinator initiating vote for transaction txn123
   Received vote from participant1:50051: 0
   Received vote from participant2:50051: 1
   ...
   At least one participant voted ABORT. Aborting transaction.
   ```
5. **After** you press `Ctrl + C` to stop the coordinator, the script removes all containers.

---

## Cleaning Up

- The **`run_twopc.sh`** script automatically cleans up containers after you exit the coordinator. 
- If you want to manually remove containers at any point, run:
  ```bash
  docker rm -f coordinator participant1 participant2 participant3 participant4 participant5
  ```
- To remove the Docker network:
  ```bash
  docker network rm twopc_net
  ```
