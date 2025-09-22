#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=lib/common.sh
source "$SCRIPT_DIR/lib/common.sh"

usage() {
  cat <<'USAGE'
Mint, audit and publish NFTs through the Synnergy CLI.

Usage: mint_nft.sh --id ID --owner ADDRESS --metadata URI --price AMOUNT [--list]
                   [--transfer NEW_OWNER] [--output FILE]

Options:
  --id ID             Identifier for the NFT asset.
  --owner ADDRESS     Wallet address that will receive the NFT.
  --metadata URI      Metadata descriptor (IPFS hash, HTTPS URL, etc.).
  --price AMOUNT      Listing price in SYN units.
  --list              Emit listing details after minting.
  --transfer ADDR     Transfer ownership to ADDR once minted (regression test).
  --output FILE       Persist the mint transaction JSON payload for auditing.
  -h, --help          Show this message and exit.

The script ensures gas alignment with MintNFT/ListNFT/BuyNFT opcodes, updates
the Function Web telemetry and stores JSON evidence for regulatory nodes.
USAGE
}

NFT_ID=""
OWNER=""
METADATA=""
PRICE=""
LIST=false
TRANSFER=""
OUTPUT=""

if [[ $# -eq 0 ]]; then
  usage
  exit 1
fi

while [[ $# -gt 0 ]]; do
  case "$1" in
    --id)
      NFT_ID=$2
      shift 2
      ;;
    --owner)
      OWNER=$2
      shift 2
      ;;
    --metadata)
      METADATA=$2
      shift 2
      ;;
    --price)
      PRICE=$2
      shift 2
      ;;
    --list)
      LIST=true
      shift
      ;;
    --transfer)
      TRANSFER=$2
      shift 2
      ;;
    --output)
      OUTPUT=$2
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

[[ -z $NFT_ID || -z $OWNER || -z $METADATA || -z $PRICE ]] && fail "id, owner, metadata and price are required"

require_commands jq

log_info "Minting NFT $NFT_ID for $OWNER"

mint_result=$(run_cli_json nft mint "$NFT_ID" "$OWNER" "$METADATA" "$PRICE")
log_info "Mint result: $mint_result"

if [[ -n $OUTPUT ]]; then
  log_info "Persisting mint artefact to $OUTPUT"
  mkdir -p "$(dirname "$OUTPUT")"
  write_secure_file "$OUTPUT" <<JSON
$mint_result
JSON
fi

if [[ $LIST == true ]]; then
  listing=$(run_cli_json nft list "$NFT_ID")
  log_info "Listing data: $listing"
fi

if [[ -n $TRANSFER ]]; then
  log_info "Transferring NFT $NFT_ID to $TRANSFER for validation"
  transfer_out=$(run_cli_json nft buy "$NFT_ID" "$TRANSFER")
  log_info "Transfer output: $transfer_out"
fi

log_info "NFT workflow completed"
