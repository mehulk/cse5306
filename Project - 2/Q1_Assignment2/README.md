# Two-Phase Commit (2PC) with Docker + gRPC

This repository provides a simplified implementation of the Two-Phase Commit protocol:

1. A **Coordinator** node that gathers votes from participants.  
2. Multiple **Participant** nodes that each vote to commit or abort.  
3. Optional **force commit** mode, where all participants always vote COMMIT.

## Table of Contents

- [Overview](#overview)
- [Project Structure](#project-structure)
- [Prerequisites](#prerequisites)
- [Building the Docker Images](#building-the-docker-images)
- [Running the Project Manually (Without Script)](#running-the-project-manually-without-script)
  - [Random Mode](#random-mode-manual)
  - [Force Commit Mode](#force-commit-mode-manual)
- [Using the Shell Script](#using-the-shell-script)
  - [Random Mode (Default)](#random-mode-default)
  - [Force Commit Mode](#force-commit-mode)
- [Stopping and Cleaning Up](#stopping-and-cleaning-up)
- [License](#license)

---

## Overview

### Protocol

1. **Voting Phase**:  
   - The coordinator sends a *RequestVote* RPC to each participant.  
   - Each participant returns either “COMMIT” or “ABORT.” (Random or forced commit if `FORCE_COMMIT=true`.)

2. **Decision**:  
   - If all participants say “COMMIT,” coordinator logs “Proceeding to commit.”  
   - If any says “ABORT,” coordinator logs “Aborting transaction.”

### Docker + gRPC

- Each node (1 coordinator + multiple participants) is containerized with Docker.  
- They communicate via gRPC, with a `.proto` definition generating Python stubs.

---

## Project Structure

- **`two_phase_commit.proto`**  
  Protobuf definitions for the RPC messages (`VoteRequest`, `VoteResponse`) and the `TwoPhaseCommit` service.

- **`coordinator.py`**  
  The coordinator script that starts a transaction, requests votes from participants, and decides commit or abort.

- **`participant.py`**  
  The participant script that starts a gRPC server to receive “RequestVote” calls and respond with commit/abort.

- **`Dockerfile_coordinator`** and **`Dockerfile_participant`**  
  Docker images for each node type. They install Python, gRPC tools, copy the code, and build the protobuf stubs.

- **`run_twopc.sh`**  
  A convenience shell script to:
  1. Build / remove old containers,
  2. Create a Docker network (if needed),
  3. Spin up 5 participant containers on ports 50051–50055,
  4. Start the coordinator,
  5. Clean up after exit.  
  It supports an optional “force” argument to make participants always commit.

---

## Prerequisites

1. **Docker** installed and running.  
2. (Optionally) **Python 3.x** if you want to test or modify code outside of containers, but not required to run everything in Docker.  

---

## Building the Docker Images

From your project’s root directory (where the `Dockerfile_*` files are located), run:

```bash
docker build -t coordinator -f Dockerfile_coordinator .
docker build -t participant -f Dockerfile_participant .
```

This produces two images locally:

- An image named **`coordinator`**  
- An image named **`participant`**

---

## Running the Project **Manually (Without Script)**

If you prefer to **manually** start containers (rather than using `run_twopc.sh`), follow these steps:

1. **Create** the Docker network (if you haven’t already):

   ```bash
   docker network create twopc_net
   ```

2. **Run** five participant containers, each with a unique port:

   ```bash
   docker run -d --name participant1 --network twopc_net \
       -e LISTEN_PORT=50051 participant

   docker run -d --name participant2 --network twopc_net \
       -e LISTEN_PORT=50052 participant

   docker run -d --name participant3 --network twopc_net \
       -e LISTEN_PORT=50053 participant

   docker run -d --name participant4 --network twopc_net \
       -e LISTEN_PORT=50054 participant

   docker run -d --name participant5 --network twopc_net \
       -e LISTEN_PORT=50055 participant
   ```

3. **Wait** a few seconds to let them start.  

4. **Run** the coordinator container in interactive mode:

   ```bash
   docker run -it --name coordinator --network twopc_net coordinator
   ```
   - This container calls `coordinator.py`, which tries to contact:
     - `participant1:50051`
     - `participant2:50052`
     - `participant3:50053`
     - `participant4:50054`
     - `participant5:50055`

5. **Observe** the logs in your terminal. You’ll see whether each participant voted commit or abort, and then a final commit/abort decision.

### Random Mode (Manual)

If you start participants with only the `LISTEN_PORT=####` environment variable (no `FORCE_COMMIT`), they randomly choose COMMIT or ABORT.

### Force Commit Mode (Manual)

If you want to **force commit** from all participants (so the coordinator definitely commits), run each participant container with `FORCE_COMMIT=true` as well:

```bash
docker run -d --name participant1 --network twopc_net \
    -e LISTEN_PORT=50051 \
    -e FORCE_COMMIT=true \
    participant
```

*(Repeat for participant2..5.)* Now each participant always votes COMMIT.

Then run the coordinator as before:
```bash
docker run -it --name coordinator --network twopc_net coordinator
```
All participants should vote COMMIT, causing a global commit.

---

## Using the Shell Script

We also provide a **`run_twopc.sh`** script that automates everything:

### Random Mode (Default)

Just run:

```bash
./run_twopc.sh
```

It will:

1. Remove old containers (coordinator, participant1..5),  
2. Create the `twopc_net` network if needed,  
3. Start participants (random commit/abort),  
4. Start the coordinator,  
5. After you press Ctrl+C to stop the coordinator, it cleans up containers.

### Force Commit Mode

To force commit from all participants, run:

```bash
./run_twopc.sh force
```

All participants will start with `FORCE_COMMIT=true`, returning COMMIT votes.

---

## Stopping and Cleaning Up

- If you’re running **manually**:  
  - Press **Ctrl+C** in the coordinator’s terminal to stop it.  
  - Then remove containers:
    ```bash
    docker rm -f coordinator participant1 participant2 participant3 participant4 participant5
    docker network rm twopc_net
    ```
- If you’re running **via `run_twopc.sh`**:
  - Press **Ctrl+C** to stop the coordinator logs.
  - The script **automatically** does `docker rm -f ...` for you.