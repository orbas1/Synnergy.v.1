#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# shellcheck disable=SC1091
. "$SCRIPT_DIR/lib/runtime.sh"

usage() {
    cat <<'USAGE'
Usage: active_active_sync.sh --primary ID [options]

Configure an active-active failover cluster through the Synnergy CLI.

Options:
  --primary ID              Primary node identifier (required).
  --timeout SECONDS         Failover heartbeat timeout (default: 30).
  --backup ID               Register backup node ID (may be repeated).
  --heartbeat ID            Emit a heartbeat for the specified node.
  --cli PATH                Path to the synnergy CLI binary.
  --dry-run                 Print intended operations without executing them.
  --json                    Emit a JSON summary to STDOUT.
  -h, --help                Display this help message.
USAGE
}

primary=""
timeout=30
json_output=0
declare -a backups
declare -a heartbeats

while [[ $# -gt 0 ]]; do
    case "$1" in
        --primary)
            primary="$2"
            shift 2
            ;;
        --timeout)
            timeout="$2"
            shift 2
            ;;
        --backup)
            backups+=("$2")
            shift 2
            ;;
        --heartbeat)
            heartbeats+=("$2")
            shift 2
            ;;
        --cli)
            SYN_CLI="$2"
            shift 2
            ;;
        --dry-run)
            SYN_DRY_RUN=1
            shift
            ;;
        --json)
            json_output=1
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            syn_fail "unknown option: $1"
            ;;
    esac
done

if [[ -z "$primary" ]]; then
    syn_fail "--primary is required"
fi

syn_ensure_cli()

syn_run "Initialising failover manager" "$SYN_CLI" highavailability init "$primary" "$timeout"
for node in "${backups[@]}"; do
    syn_run "Registering backup $node" "$SYN_CLI" highavailability add "$node"
done

for hb in "${heartbeats[@]}"; do
    syn_run "Recording heartbeat from $hb" "$SYN_CLI" highavailability heartbeat "$hb"
done

active="$primary"
if [[ $SYN_DRY_RUN -eq 0 ]]; then
    if out=$(syn_capture "Querying active node" "$SYN_CLI" highavailability active); then
        active=$(echo "$out" | tr -d '\r')
    fi
else
    syn_info "Dry run: assuming $primary remains active"
fi

if [[ $json_output -eq 1 ]]; then
    python3 - "$primary" "$active" "$SYN_DRY_RUN" "${backups[@]:-}" <<'PY'
import json
import sys

primary = sys.argv[1]
active = sys.argv[2]
dry_run = bool(int(sys.argv[3]))
backups = list(dict.fromkeys(sys.argv[4:])) if len(sys.argv) > 4 else []

print(json.dumps({
    "primary": primary,
    "active": active,
    "backups": backups,
    "dry_run": dry_run,
}))
PY
else
    echo "Active node: $active"
fi
