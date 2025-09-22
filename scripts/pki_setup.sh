#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib/common.sh
source "$SCRIPT_DIR/lib/common.sh"

usage() {
  cat <<'USAGE'
Provision a local PKI hierarchy for authority and validator nodes.

Usage: pki_setup.sh [--ca-dir DIR] [--org NAME] [--nodes N]

Options:
  --ca-dir DIR   Directory to store certificates (default: ./pki).
  --org NAME     Organisation name for certificates (default: Synnergy).
  --nodes N      Number of node certificates to generate (default: 3).
  -h, --help     Show help message.
USAGE
}

CA_DIR="${PWD}/pki"
ORG="Synnergy"
NODE_COUNT=3

while [[ $# -gt 0 ]]; do
  case "$1" in
    --ca-dir)
      CA_DIR=$2
      shift 2
      ;;
    --org)
      ORG=$2
      shift 2
      ;;
    --nodes)
      NODE_COUNT=$2
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

require_commands openssl

mkdir -p "$CA_DIR"
CA_KEY="$CA_DIR/ca.key.pem"
CA_CERT="$CA_DIR/ca.cert.pem"

if [[ ! -f $CA_KEY ]]; then
  log_info "Generating root CA"
  openssl req -x509 -nodes -newkey rsa:4096 -keyout "$CA_KEY" -out "$CA_CERT" \
    -subj "/CN=$ORG Root CA" -days 3650
fi

for i in $(seq 1 "$NODE_COUNT"); do
  node_key="$CA_DIR/node-$i.key.pem"
  node_csr="$CA_DIR/node-$i.csr.pem"
  node_cert="$CA_DIR/node-$i.cert.pem"
  log_info "Generating certificate for node $i"
  openssl req -new -nodes -newkey rsa:2048 -keyout "$node_key" -out "$node_csr" \
    -subj "/CN=$ORG Node $i"
  openssl x509 -req -in "$node_csr" -CA "$CA_CERT" -CAkey "$CA_KEY" \
    -CAcreateserial -out "$node_cert" -days 825 -sha256
  rm -f "$node_csr"
done

log_info "PKI setup complete in $CA_DIR"
