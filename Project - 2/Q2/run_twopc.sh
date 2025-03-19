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

# Start participant containers for voting (50051-50055) and global decision (6001-6005)
echo "Starting participant1 (voting: 50051, decision: 6001)..."
docker run -d --name participant1 --network twopc_net -p 50051:50051 -p 6001:6001 -e LISTEN_PORT=6001 $FORCE_ENV twopc_participant_image ./global_decision -role=participant -port=6001

echo "Starting participant2 (voting: 50052, decision: 6002)..."
docker run -d --name participant2 --network twopc_net -p 50052:50052 -p 6002:6002 -e LISTEN_PORT=6002 $FORCE_ENV twopc_participant_image ./global_decision -role=participant -port=6002

echo "Starting participant3 (voting: 50053, decision: 6003)..."
docker run -d --name participant3 --network twopc_net -p 50053:50053 -p 6003:6003 -e LISTEN_PORT=6003 $FORCE_ENV twopc_participant_image ./global_decision -role=participant -port=6003

echo "Starting participant4 (voting: 50054, decision: 6004)..."
docker run -d --name participant4 --network twopc_net -p 50054:50054 -p 6004:6004 -e LISTEN_PORT=6004 $FORCE_ENV twopc_participant_image ./global_decision -role=participant -port=6004

echo "Starting participant5 (voting: 50055, decision: 6005)..."
docker run -d --name participant5 --network twopc_net -p 50055:50055 -p 6005:6005 -e LISTEN_PORT=6005 $FORCE_ENV twopc_participant_image ./global_decision -role=participant -port=6005

# Wait a bit to ensure participants have started listening
echo "Waiting 5 seconds for participants to start..."
sleep 5

PARTICIPANT_ADDRS="participant1:6001,participant2:6002,participant3:6003,participant4:6004,participant5:6005"

# Run the coordinator container (it will connect to all participants)
echo "Starting coordinator..."
docker run -it --name coordinator \
    --network twopc_net \
    -p 50050:50050 \
    -p 6000:6000 \
    -e PARTICIPANT_ADDRS="$PARTICIPANT_ADDRS" \
    twopc_coordinator_image ./global_decision -role=coordinator -port=6000

# After the coordinator exits, clean up all containers
echo "Cleaning up containers..."
docker rm -f coordinator participant1 participant2 participant3 participant4 participant5 2>/dev/null
echo "All done!"
