#!/bin/bash
#
# Usage:
#   ./run_twopc.sh          # => all participants vote random commit/abort
#   ./run_twopc.sh force    # => all participants FORCE_COMMIT=true (always commit)
#
# After coordinator exits, all containers are removed automatically.

# 1. Remove old containers (ignore errors if they don't exist)
docker rm -f coordinator participant1 participant2 participant3 participant4 participant5 2>/dev/null

# 2. Create Docker network if it doesn't already exist (ignore error if it does)
docker network create twopc_net 2>/dev/null

# 3. Check if "force" mode was requested
if [ "$1" == "force" ]; then
  echo "Starting ALL participants with FORCE_COMMIT=true (always commit)."
  for i in {1..5}; do
    docker run -d -e FORCE_COMMIT=true --name participant$i --network twopc_net participant
  done
else
  echo "Starting ALL participants with random commit/abort."
  for i in {1..5}; do
    docker run -d --name participant$i --network twopc_net participant
  done
fi

# 4. Run coordinator in interactive mode so we can see logs
echo "Running coordinator container (press Ctrl+C to stop)..."
docker run -it --name coordinator --network twopc_net coordinator

# 5. Once coordinator is stopped (Ctrl+C), clean up again
echo "Coordinator exited. Cleaning up containers..."
docker rm -f coordinator participant1 participant2 participant3 participant4 participant5 2>/dev/null
echo "All done!"
