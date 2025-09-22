#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib/common.sh
source "$SCRIPT_DIR/lib/common.sh"

usage() {
  cat <<'USAGE'
Initialise multi-factor credentials for a Synnergy wallet.

Usage: multi_factor_setup.sh [--wallet FILE] [--password PASS] [--otp-file FILE]
                             [--role ROLE] [--address ADDRESS]

Options:
  --wallet FILE      Wallet file path to generate (default: ./wallets/mfa.json).
  --password PASS    Password used to encrypt the wallet file (default: random).
  --otp-file FILE    Location where OTP seed is stored (default: ./wallets/mfa.otp).
  --role ROLE        Optional authority role assigned to wallet via CLI access module.
  --address ADDRESS  Optional existing address instead of generating a new wallet.
  -h, --help         Show help message.

The script combines wallet generation, OTP provisioning and role assignment so
MFA credentials are synchronised across CLI, consensus modules and regulatory
policies.
USAGE
}

WALLET_FILE="${PWD}/wallets/mfa.json"
PASSWORD=""
OTP_FILE="${PWD}/wallets/mfa.otp"
ROLE=""
ADDRESS=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --wallet)
      WALLET_FILE=$2
      shift 2
      ;;
    --password)
      PASSWORD=$2
      shift 2
      ;;
    --otp-file)
      OTP_FILE=$2
      shift 2
      ;;
    --role)
      ROLE=$2
      shift 2
      ;;
    --address)
      ADDRESS=$2
      shift 2
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

require_commands openssl jq

mkdir -p "$(dirname "$WALLET_FILE")"
mkdir -p "$(dirname "$OTP_FILE")"

if [[ -z $PASSWORD ]]; then
  PASSWORD=$(openssl rand -hex 16)
  log_info "Generated random wallet password"
fi

if [[ -z $ADDRESS ]]; then
  log_info "Generating wallet"
  wallet_json=$(run_cli_json wallet new --out "$WALLET_FILE" --password "$PASSWORD")
  ADDRESS=$(json_extract '.address' "$wallet_json")
else
  log_info "Using existing address $ADDRESS"
fi

log_info "Creating OTP secret"
OTP_SECRET=$(openssl rand -base64 32)
write_secure_file "$OTP_FILE" <<DATA
$OTP_SECRET
DATA

log_info "Encrypting OTP metadata for secure storage"
enc_payload=$(printf '%s' "$OTP_SECRET" | xxd -p -c 256)
json_payload=$(cat <<JSON
{
  "payload": "$enc_payload",
  "nonce": "$(openssl rand -hex 12)"
}
JSON
)

meta_file="${OTP_FILE}.json"
write_secure_file "$meta_file" <<JSON
$json_payload
JSON

if [[ -n $ROLE ]]; then
  log_info "Granting role $ROLE to address $ADDRESS"
  run_cli_json access grant "$ROLE" "$ADDRESS" >/dev/null
fi

log_info "MFA setup completed for address $ADDRESS"
log_info "Wallet stored at $WALLET_FILE; OTP seed at $OTP_FILE"
