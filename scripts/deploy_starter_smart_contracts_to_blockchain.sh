#!/usr/bin/env bash
set -euo pipefail

CONTRACT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../smart-contracts/solidity" && pwd)"

for contract in "$CONTRACT_DIR"/*.sol; do
  name=$(basename "$contract")
  echo "Deploying $name to blockchain..."
  # Placeholder: replace with actual deployment command
  echo "npx hardhat run --network localhost \"$contract\""
done
