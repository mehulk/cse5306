#!/usr/bin/env bash
set -e

##############################################
# 1) Start Q3’s Python cluster for leader election
##############################################
echo "=== Starting Q3 Python cluster (leader election) ==="
pushd Q3 >/dev/null

# Build & start Q3 containers in background
docker compose up -d

popd >/dev/null

echo "Q3 cluster started in the background."
echo "Sleeping 10 seconds to allow leader election to finish..."
sleep 10

##############################################
# 2) Ask user which node was elected leader
##############################################
echo "Please check logs of Q3 to see which node became leader. Then enter that node ID (1..5)."
echo -n "Enter the Elected Leader ID (e.g., 2 or 3): "
read LEADER_ID

if [ -z "$LEADER_ID" ]; then
  echo "No LEADER_ID entered, defaulting to 1."
  LEADER_ID=1
fi

##############################################
# 3) Start Q4’s Go cluster, passing LEADER_ID
##############################################
echo "=== Starting Q4 Go cluster (log replication) with LEADER_ID=$LEADER_ID ==="
pushd Q4 >/dev/null

# Build & start Q4 containers in foreground or background
# Here, let's do background so the script continues
LEADER_ID=$LEADER_ID docker compose up

popd >/dev/null

echo "Q4 cluster started in the background."

##############################################
# 4) Instructions for the user to run client
##############################################
echo "=========================================="
echo "Both Q3 (Python) and Q4 (Go) clusters are running."
echo "Q3 handles leader election, but Q4 also requires the LEADER_ID environment (you just entered '$LEADER_ID')."
echo ""
echo "To see Q3 logs, do:  docker compose logs -f  (inside Q3/ folder)"
echo "To see Q4 logs, do:  docker compose logs -f  (inside Q4/ folder)"
echo ""
echo "Finally, open another terminal and run the Q4 client (which calls SubmitClientRequest), e.g.:"
echo "  cd Q4 && ./client"
echo "or"
echo "  cd Q4 && go run client.go"
echo "=========================================="
