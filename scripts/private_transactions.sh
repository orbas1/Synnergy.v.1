#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib/common.sh
source "$SCRIPT_DIR/lib/common.sh"

usage() {
  cat <<'USAGE'
Encrypt and submit private transactions through the CLI.

Usage: private_transactions.sh --key KEY --payload DATA [--out FILE]
                               [--nonce HEX] [--submit]

Options:
  --key KEY       Symmetric key for encryption (hex or passphrase).
  --payload DATA  Plaintext payload to encrypt.
  --out FILE      Persist encrypted payload JSON to FILE.
  --nonce HEX     Optional nonce override (hex encoded).
  --submit        Submit the JSON payload via private-tx send.
  -h, --help      Display this message.
USAGE
}

KEY=""
PAYLOAD=""
OUT_FILE=""
NONCE=""
SUBMIT=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    --key)
      KEY=$2
      shift 2
      ;;
    --payload)
      PAYLOAD=$2
      shift 2
      ;;
    --out)
      OUT_FILE=$2
      shift 2
      ;;
    --nonce)
      NONCE=$2
      shift 2
      ;;
    --submit)
      SUBMIT=true
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      usage
      fail "unknown flag $1"
      ;;
  esac
done

[[ -z $KEY || -z $PAYLOAD ]] && fail "key and payload are required"

require_commands jq xxd openssl

log_info "Encrypting payload"
cipher_hex=$(run_cli private-tx encrypt "$KEY" "$PAYLOAD")
nonce_hex=${NONCE:-$(openssl rand -hex 12)}

json_body=$(cat <<JSON
{
  "payload": "$cipher_hex",
  "nonce": "$nonce_hex"
}
JSON
)

if [[ -n $OUT_FILE ]]; then
  write_secure_file "$OUT_FILE" <<JSON
$json_body
JSON
fi

if [[ $SUBMIT == true ]]; then
  tmp_file=$(mktemp)
  write_secure_file "$tmp_file" <<JSON
$json_body
JSON
  log_info "Submitting encrypted transaction"
  run_cli private-tx send "$tmp_file"
  rm -f "$tmp_file"
fi

log_info "Private transaction workflow completed"
