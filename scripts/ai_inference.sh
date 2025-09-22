#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# shellcheck disable=SC1091
. "$SCRIPT_DIR/lib/runtime.sh"

usage() {
    cat <<'USAGE'
Usage: ai_inference.sh --contract ADDRESS [options]

Invoke an AI-enhanced smart contract via the Synnergy CLI.

Options:
  --contract ADDRESS        Contract address (required).
  --payload FILE            Binary or JSON payload file.
  --payload-hex HEX         Pre-encoded payload in hexadecimal.
  --gas LIMIT               Gas limit to supply (default: 500000).
  --output FILE             Write decoded response to FILE (JSON).
  --decode text|base64|hex  Output format for contract response (default: hex).
  --cli PATH                Path to synnergy CLI binary.
  --dry-run                 Print intended command without execution.
  -h, --help                Show this message.
USAGE
}

CONTRACT=""
PAYLOAD_FILE=""
PAYLOAD_HEX=""
GAS_LIMIT=500000
OUTPUT=""
DECODE="hex"

while [[ $# -gt 0 ]]; do
    case "$1" in
        --contract)
            CONTRACT="$2"
            shift 2
            ;;
        --payload)
            PAYLOAD_FILE="$2"
            shift 2
            ;;
        --payload-hex)
            PAYLOAD_HEX="$2"
            shift 2
            ;;
        --gas)
            GAS_LIMIT="$2"
            shift 2
            ;;
        --output)
            OUTPUT="$2"
            shift 2
            ;;
        --decode)
            DECODE="$2"
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
            usage
            exit 0
            ;;
        *)
            syn_fail "unknown option: $1"
            ;;
    esac
done

if [[ -z "$CONTRACT" ]]; then
    syn_fail "--contract is required"
fi

if [[ -n "$PAYLOAD_FILE" && -n "$PAYLOAD_HEX" ]]; then
    syn_fail "specify either --payload or --payload-hex"
fi

if [[ -z "$PAYLOAD_FILE" && -z "$PAYLOAD_HEX" ]]; then
    syn_fail "provide --payload or --payload-hex"
fi

if [[ -n "$PAYLOAD_FILE" ]]; then
    syn_require_file "$PAYLOAD_FILE"
    PAYLOAD_HEX=$(python3 - "$PAYLOAD_FILE" <<'PY'
import binascii
import sys
from pathlib import Path

path = Path(sys.argv[1])
raw = path.read_bytes()
print(binascii.hexlify(raw).decode())
PY
)
fi

syn_ensure_cli()

if [[ $SYN_DRY_RUN -eq 1 ]]; then
    syn_info "[dry-run] would invoke $CONTRACT with gas $GAS_LIMIT"
    echo '{"dry_run": true}'
    exit 0
fi

result=$(syn_capture "Invoking AI contract" "$SYN_CLI" ai_contract invoke "$CONTRACT" "$PAYLOAD_HEX" "$GAS_LIMIT")
if [[ -z "$result" ]]; then
    syn_fail "CLI returned no output"
fi

result=${result//$'\r'/}
output_hex=$(echo "$result" | sed -n 's/.*output=\([0-9a-fA-F]*\).*/\1/p')
gas_used=$(echo "$result" | sed -n 's/.*gas=\([0-9]\+\).*/\1/p')
if [[ -z "$output_hex" ]]; then
    syn_warn "unable to parse output payload from CLI response"
    output_hex="$result"
fi
if [[ -z "$gas_used" ]]; then
    gas_used=null
fi

decode_payload() {
    local mode=$1
    local data=$2
    python3 - "$mode" "$data" <<'PY'
import base64
import binascii
import sys

mode = sys.argv[1]
data = sys.argv[2]

if mode == 'hex':
    print(data)
elif mode == 'base64':
    try:
        raw = binascii.unhexlify(data)
    except binascii.Error:
        raw = data.encode()
    print(base64.b64encode(raw).decode())
elif mode == 'text':
    try:
        raw = bytes.fromhex(data)
        print(raw.decode('utf-8', errors='replace'))
    except ValueError:
        print(data)
else:
    print(data)
PY
}

decoded=$(decode_payload "$DECODE" "$output_hex")

summary=$(python3 - "$DECODE" "$decoded" "$output_hex" "$gas_used" <<'PY'
import json
import sys

decode = sys.argv[1]
content = sys.argv[2]
hex_payload = sys.argv[3]
used = sys.argv[4]

summary = {
    "decode": decode,
    "payload": content,
    "payload_hex": hex_payload,
    "gas_used": None if used == 'null' else int(used),
}

print(json.dumps(summary, indent=2))
PY
)

if [[ -n "$OUTPUT" ]]; then
    OUTPUT=$(syn_expand_path "$OUTPUT")
    syn_write_json "$OUTPUT" "$summary"
else
    echo "$summary"
fi
