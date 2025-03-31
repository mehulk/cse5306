#!/bin/bash
#
# Usage:
#   ./run_twopc.sh            # => Random votes for all participants
#   ./run_twopc.sh force      # => All participants forced to commit
#
# After stopping the coordinator (Ctrl+C), it cleans up all containers.

# If "force" is provided, set the env variable for participants
if [ "$1" == "force" ]; then
  echo "Running in FORCE COMMIT mode: All participants will vote COMMIT."
  FORCE_ENV="-e FORCE_COMMIT=true"
else
  echo "Running in RANDOM mode: Participants will randomly vote COMMIT or ABORT."
  FORCE_ENV=""
fi

# Remove old containers (ignore errors if they don't exist)
docker rm -f coordinator participant1 participant2 participant3 participant4 participant5 2>/dev/null

# Create Docker network (ignore error if it exists)
docker network create twopc_net 2>/dev/null

# Start participant containers on ports 50051..50055
echo "Starting participant1 on port 50051..."
docker run -d --name participant1 --network twopc_net -e LISTEN_PORT=50051 $FORCE_ENV participant

echo "Starting participant2 on port 50052..."
docker run -d --name participant2 --network twopc_net -e LISTEN_PORT=50052 $FORCE_ENV participant

echo "Starting participant3 on port 50053..."
docker run -d --name participant3 --network twopc_net -e LISTEN_PORT=50053 $FORCE_ENV participant

echo "Starting participant4 on port 50054..."
docker run -d --name participant4 --network twopc_net -e LISTEN_PORT=50054 $FORCE_ENV participant

echo "Starting participant5 on port 50055..."
docker run -d --name participant5 --network twopc_net -e LISTEN_PORT=50055 $FORCE_ENV participant

# Wait a bit to ensure participants have started listening
echo "Waiting 5 seconds for participants to start..."
sleep 5

# Run the coordinator container (it will connect to participant1..5 on the specified ports)
echo "Starting coordinator..."
docker run -it --name coordinator --network twopc_net coordinator

# After the coordinator exits, clean up all containers
# echo "Cleaning up containers..."
# docker rm -f coordinator participant1 participant2 participant3 participant4 participant5 2>/dev/null
# echo "All done!"
