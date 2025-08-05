#!/usr/bin/env bash
set -e
if [ -z "$1" ]; then
  echo "Usage: $0 <file> [provider] [price] [capacity]"
  exit 1
fi
provider=${2:-provider}
price=${3:-0}
capacity=${4:-0}
./synnergy marketplace pin "$1" "$provider" "$price" "$capacity"
