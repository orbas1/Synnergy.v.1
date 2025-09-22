#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# shellcheck disable=SC1091
. "$SCRIPT_DIR/lib/runtime.sh"

REGISTRY_FILE=${AI_MODEL_REGISTRY:-"$HOME/.synnergy/models_registry.json"}

usage() {
    cat <<'USAGE'
Usage: ai_model_management.sh <command> [options]

Commands:
  register   Register a new AI-enhanced contract deployment.
  list       List known AI contracts with metadata.
  show       Show metadata for a single contract.
  archive    Mark a contract as archived in the local registry.

Use "ai_model_management.sh <command> --help" to view command-specific options.
USAGE
}

ensure_registry() {
    local file=$1
    file=$(syn_expand_path "$file")
    mkdir -p "$(dirname "$file")"
    if [[ ! -f "$file" ]]; then
        echo '[]' >"$file"
    fi
    echo "$file"
}

command=${1:-}
if [[ -z "$command" ]]; then
    usage
    exit 1
fi
shift

registry_path=$(ensure_registry "$REGISTRY_FILE")

case "$command" in
    register)
        register_usage() {
            cat <<'USAGE'
Usage: ai_model_management.sh register --wasm FILE --model HASH --manifest FILE|STRING --owner ADDRESS [options]

Options:
  --wasm FILE             Compiled WASM artifact to deploy.
  --model HASH            Associated model hash identifier.
  --manifest FILE|STRING  Contract manifest (file path or inline JSON/string).
  --owner ADDRESS         Owner wallet or account identifier.
  --gas LIMIT             Gas limit (default: 750000).
  --note TEXT             Free-form description stored locally.
  --cli PATH              Path to synnergy CLI binary.
  --dry-run               Record metadata without invoking CLI.
  -h, --help              Show this help text.
USAGE
        }

        wasm=""
        model=""
        manifest=""
        owner=""
        gas=750000
        note=""

        while [[ $# -gt 0 ]]; do
            case "$1" in
                --wasm)
                    wasm="$2"
                    shift 2
                    ;;
                --model)
                    model="$2"
                    shift 2
                    ;;
                --manifest)
                    manifest="$2"
                    shift 2
                    ;;
                --owner)
                    owner="$2"
                    shift 2
                    ;;
                --gas)
                    gas="$2"
                    shift 2
                    ;;
                --note)
                    note="$2"
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
                -h|--help)
                    register_usage
                    exit 0
                    ;;
                *)
                    syn_fail "unknown option: $1"
                    ;;
            esac
        done

        [[ -z "$wasm" ]] && syn_fail "--wasm is required"
        [[ -z "$model" ]] && syn_fail "--model is required"
        [[ -z "$manifest" ]] && syn_fail "--manifest is required"
        [[ -z "$owner" ]] && syn_fail "--owner is required"

        syn_require_file "$wasm"
        manifest_payload="$manifest"
        manifest_encoding="inline"
        if [[ -f "$manifest" ]]; then
            manifest_meta=$(python3 - "$manifest" <<'PY'
import base64
import json
import sys
from pathlib import Path

path = Path(sys.argv[1])
raw = path.read_bytes()
try:
    data = json.loads(raw.decode())
    payload = json.dumps(data, separators=(",", ":"))
    encoding = "json"
except Exception:
    payload = base64.b64encode(raw).decode()
    encoding = "base64"

print(json.dumps({"payload": payload, "encoding": encoding}))
PY
)
            manifest_payload=$(python3 -c 'import json,sys;print(json.loads(sys.argv[1])["payload"])' "$manifest_meta")
            manifest_encoding=$(python3 -c 'import json,sys;print(json.loads(sys.argv[1])["encoding"])' "$manifest_meta")
        fi

        syn_ensure_cli()
        address=""
        if [[ $SYN_DRY_RUN -eq 1 ]]; then
            syn_info "[dry-run] would deploy model $model"
            address="dry-run-$model"
        else
            response=$(syn_capture "Deploying AI contract" "$SYN_CLI" ai_contract deploy "$wasm" "$model" "$manifest_payload" "$gas" "$owner")
            address=$(echo "$response" | awk '/contract:/ {print $2}')
            if [[ -z "$address" ]]; then
                syn_fail "unable to parse contract address from CLI response"
            fi
        fi

        python3 - "$registry_path" "$address" "$model" "$owner" "$note" "$manifest" "$manifest_encoding" "$gas" "$SYN_DRY_RUN" <<'PY'
import json
import sys
from datetime import datetime
from pathlib import Path

path = Path(sys.argv[1])
addr, model, owner, note, manifest, encoding, gas, dry = sys.argv[2:9]

dry = bool(int(dry))
if path.exists():
    try:
        registry = json.loads(path.read_text(encoding="utf-8"))
    except json.JSONDecodeError:
        registry = []
else:
    registry = []

record = {
    "contract": addr,
    "model_hash": model,
    "owner": owner,
    "note": note or None,
    "manifest_source": manifest,
    "manifest_encoding": encoding,
    "gas_limit": int(gas),
    "registered_at": datetime.utcnow().isoformat() + "Z",
    "dry_run": dry,
}

registry = [entry for entry in registry if entry.get("contract") != addr]
registry.append(record)
registry.sort(key=lambda item: item.get("registered_at", ""))
path.write_text(json.dumps(registry, indent=2), encoding="utf-8")

print(json.dumps(record, indent=2))
PY
        ;;

    list)
        list_usage() {
            cat <<'USAGE'
Usage: ai_model_management.sh list [--cli PATH]

List AI contracts tracked locally, merging on-chain metadata when available.
USAGE
        }

        while [[ $# -gt 0 ]]; do
            case "$1" in
                --cli)
                    SYN_CLI="$2"
                    shift 2
                    ;;
                -h|--help)
                    list_usage
                    exit 0
                    ;;
                *)
                    syn_fail "unknown option: $1"
                    ;;
            esac
        done

        syn_ensure_cli()

        local_data=$(python3 - "$registry_path" <<'PY'
import json
import sys
from pathlib import Path
path = Path(sys.argv[1])
if path.exists():
    try:
        data = json.loads(path.read_text(encoding='utf-8'))
    except json.JSONDecodeError:
        data = []
else:
    data = []
print(json.dumps(data))
PY
)

        if [[ $SYN_DRY_RUN -eq 1 ]]; then
            echo "$local_data"
            exit 0
        fi

        onchain=$(syn_capture "Listing AI contracts" "$SYN_CLI" ai_contract list --json)
        python3 - <<'PY'
import json
import sys

local_data = json.loads(sys.argv[1])
try:
    onchain = json.loads(sys.argv[2])
except json.JSONDecodeError:
    onchain = []

index = {item.get("address"): item for item in onchain}
for record in local_data:
    addr = record.get("contract")
    chain = index.get(addr, {})
    record["chain"] = chain
print(json.dumps(local_data, indent=2))
PY
"$local_data" "$onchain"
        ;;

    show)
        show_usage() {
            cat <<'USAGE'
Usage: ai_model_management.sh show --contract ADDRESS [--cli PATH]
USAGE
        }

        contract=""
        while [[ $# -gt 0 ]]; do
            case "$1" in
                --contract)
                    contract="$2"
                    shift 2
                    ;;
                --cli)
                    SYN_CLI="$2"
                    shift 2
                    ;;
                -h|--help)
                    show_usage
                    exit 0
                    ;;
                *)
                    syn_fail "unknown option: $1"
                    ;;
            esac
        done
        [[ -z "$contract" ]] && syn_fail "--contract is required"
        syn_ensure_cli()

        local_entry=$(python3 - "$registry_path" "$contract" <<'PY'
import json
import sys
from pathlib import Path

path = Path(sys.argv[1])
contract = sys.argv[2]
if path.exists():
    try:
        data = json.loads(path.read_text(encoding='utf-8'))
    except json.JSONDecodeError:
        data = []
else:
    data = []

for entry in data:
    if entry.get('contract') == contract:
        print(json.dumps(entry))
        break
PY
)

        chain_info="{}"
        if [[ $SYN_DRY_RUN -eq 0 ]]; then
            chain_info=$(syn_capture "Fetching model hash" "$SYN_CLI" ai_contract model "$contract" || true)
        fi

        python3 - "$contract" "$local_entry" "$chain_info" <<'PY'
import json
import sys

contract = sys.argv[1]
local_entry = sys.argv[2]
chain_info = sys.argv[3]

try:
    local = json.loads(local_entry) if local_entry else {}
except json.JSONDecodeError:
    local = {}

try:
    chain = {"model_hash": chain_info.strip()} if chain_info.strip() else {}
except Exception:
    chain = {}

summary = {
    "contract": contract,
    "local": local or None,
    "chain": chain or None,
}

print(json.dumps(summary, indent=2))
PY
        ;;

    archive)
        archive_usage() {
            cat <<'USAGE'
Usage: ai_model_management.sh archive --contract ADDRESS
USAGE
        }
        contract=""
        while [[ $# -gt 0 ]]; do
            case "$1" in
                --contract)
                    contract="$2"
                    shift 2
                    ;;
                -h|--help)
                    archive_usage
                    exit 0
                    ;;
                *)
                    syn_fail "unknown option: $1"
                    ;;
            esac
        done
        [[ -z "$contract" ]] && syn_fail "--contract is required"
        python3 - "$registry_path" "$contract" <<'PY'
import json
import sys
from datetime import datetime
from pathlib import Path

path = Path(sys.argv[1])
contract = sys.argv[2]
if path.exists():
    try:
        data = json.loads(path.read_text(encoding='utf-8'))
    except json.JSONDecodeError:
        data = []
else:
    data = []

archived = None
for entry in data:
    if entry.get('contract') == contract:
        entry['archived_at'] = datetime.utcnow().isoformat() + 'Z'
        archived = entry
        break

if archived is None:
    archived = {
        'contract': contract,
        'archived_at': datetime.utcnow().isoformat() + 'Z',
    }
    data.append(archived)

path.write_text(json.dumps(data, indent=2), encoding='utf-8')
print(json.dumps(archived, indent=2))
PY
        ;;

    *)
        usage
        exit 1
        ;;
esac
