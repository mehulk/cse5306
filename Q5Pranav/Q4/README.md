# SWIM Protocol Implementation

This project implements the Failure Detector Component of the SWIM (Scalable Weakly-consistent Infection-style Process Group Membership) protocol using Python and gRPC.

## Project Structure

```
swim_protocol/
│
├── proto/
│   └── swim.proto
│
├── src/
│   ├── node.py
│   ├── failure_detector.py
│   └── main.py
│
├── Dockerfile
├── docker-compose.yml
├── requirements.txt
└── README.md
```

## Prerequisites

- Docker
- Docker Compose

## Setup and Running

1. Build and start the containers:
   ```
   docker compose build
   docker compose up -d
   ```

2. View logs:
   - All nodes:
     ```
     docker compose logs -f
     ```
   - Specific node (e.g., node1):
     ```
     docker compose logs -f node1
     ```

3. Stop and remove containers:
   ```
   docker compose down
   ```

## Testing Failure Detection

1. To simulate a node failure, stop one of the containers:
   ```
   docker compose stop node3
   ```

2. Observe the logs to see how other nodes detect the failure:
   ```
   docker compose logs -f
   ```

3. Restart the stopped node:
   ```
   docker compose start node3
   ```

## Implementation Details

- The project uses gRPC for communication between nodes.
- Each node runs a failure detector that periodically pings other nodes.
- If a node doesn't respond to a ping, an indirect ping is initiated through other nodes.
- Nodes that fail to respond to both direct and indirect pings are marked as failed.

## Files Description

- `proto/swim.proto`: Defines the gRPC service and message types.
- `src/node.py`: Implements the Node class with gRPC service methods.
- `src/failure_detector.py`: Implements the FailureDetector class with ping logic.
- `src/main.py`: Sets up and runs the nodes.
- `Dockerfile`: Defines the Docker image for the nodes.
- `docker-compose.yml`: Defines the multi-container Docker application.
- `requirements.txt`: Lists the Python dependencies.

## Customization

- Adjust the ping interval and number of indirect ping nodes in `src/failure_detector.py`.
- Modify the number of nodes by editing `docker-compose.yml`.

## Output Format

The implementation prints messages for each RPC call in the following formats:

1. On the client side (sender):
   ```
   Component <component_name> of Node <node_id> sends RPC <rpc_name> to Component <component_name> of Node <node_id>
   ```

2. On the server side (receiver):
   ```
   Component <component_name> of Node <node_id> runs RPC <rpc_name> called by Component <component_name> of Node <node_id>
   ```

These messages help track the communication between nodes and the execution of RPC methods.