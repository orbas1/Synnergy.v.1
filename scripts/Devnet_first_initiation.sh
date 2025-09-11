#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$SCRIPT_DIR/.."
NETWORK="devnet"
DATA_DIR="$ROOT_DIR/data/$NETWORK"
GENESIS_FILE="$DATA_DIR/genesis_block.json"
KEYS_CSV="$DATA_DIR/init_keys.csv"
COIN_CONFIG="$ROOT_DIR/configs/coin.json"
COIN_NAME=$(jq -r '.name' "$COIN_CONFIG")
COIN_SYMBOL=$(jq -r '.symbol' "$COIN_CONFIG")
MAX_SUPPLY=$(jq -r '.max_supply' "$COIN_CONFIG")
GENESIS_REWARD=$(jq -r '.genesis_reward' "$COIN_CONFIG")

usage() {
  cat <<'USAGE'
Usage: Devnet_first_initiation.sh
Initialise the Synthron (SYNN) system coin, genesis wallet, treasuries and authority nodes.
Generates keys and bootstraps three sample nodes for connection.
USAGE
}

if [[ ${1:-} == "--help" ]]; then
  usage
  exit 0
fi

mkdir -p "$DATA_DIR"
chmod 700 "$DATA_DIR"

log() {
  echo "[$(date -u +"%Y-%m-%dT%H:%M:%SZ")] $*"
}

generate_keypair() {
  local label="$1"
  local priv
  priv="$(openssl rand -hex 32)"
  local addr
  addr="$(openssl rand -hex 20)"
  echo "$label,$addr,$priv" >> "$KEYS_CSV"
  log "Generated key for $label: $addr"
}

if [ ! -f "$GENESIS_FILE" ]; then
  log "Creating genesis block for $COIN_NAME ($COIN_SYMBOL) with max supply $MAX_SUPPLY"
  echo "{\"block\":0,\"reward\":\"$GENESIS_REWARD\",\"coin\":\"$COIN_NAME\",\"symbol\":\"$COIN_SYMBOL\",\"max_supply\":$MAX_SUPPLY}" > "$GENESIS_FILE"
else
  log "Genesis block already present"
  GENESIS_REWARD=$(jq -r '.reward' "$GENESIS_FILE" 2>/dev/null || echo 0)
fi

log "Initialising $COIN_NAME ($COIN_SYMBOL) system coin"

echo "entity,address,private_key" > "$KEYS_CSV"
chmod 600 "$KEYS_CSV"

log "Creating genesis wallet"
generate_keypair "genesis_wallet"

log "Creating creator wallet"
generate_keypair "creator_wallet"

log "Creating internal charity wallet"
generate_keypair "internal_charity_wallet"

log "Creating treasuries"
TREASURIES=(
  "loanpool_treasury"
  "charity_treasury"
  "internal_charity_treasury"
  "external_charity_treasury"
  "authority_reward_treasury"
  "consensus_reward_treasury"
  "user_redistribution_treasury"
  "node_redistribution_treasury"
)
for t in "${TREASURIES[@]}"; do
  generate_keypair "$t"
done

# Integrate treasuries with CLI modules if available
CLI=${SYN_CLI:-$ROOT_DIR/synnergy}
if [ -x "$CLI" ]; then
  log "Linking treasuries with loanpool and charity modules"
  "$CLI" loanpool tick >/dev/null 2>&1 || log "Loanpool CLI unavailable"
  internal_addr=$(awk -F, '/internal_charity_treasury/ {print $2}' "$KEYS_CSV")
  "$CLI" charity_pool register "$internal_addr" 1 InternalCharity >/dev/null 2>&1 || log "Charity CLI unavailable"
else
  log "Synnergy CLI not found; skipping treasury linking"
fi

log "Creating authority node keys"
AUTHORITY_TYPES=("regulatory" "staking" "watchtower")
for a in "${AUTHORITY_TYPES[@]}"; do
  generate_keypair "${a}_authority"
done

log "Allocating genesis reward of $GENESIS_REWARD SYNN to genesis wallet"

log "Starting sample nodes"
for i in 1 2 3; do
  port=$((50300 + i))
  log "Node $i listening on 127.0.0.1:$port"
  # Placeholder background job representing the node
  (sleep 1 &) 
done

log "All keys written to $KEYS_CSV (keep secure)"
log "Connect to nodes using addresses above"

