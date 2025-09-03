#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$SCRIPT_DIR/.."
cd "$ROOT_DIR"

mkdir -p bin

go mod tidy

go mod verify

for dir in cmd/*; do
  if [[ -d "$dir" ]]; then
    go build -o "bin/$(basename "$dir")" "./$dir"
  fi

done

echo "Binaries built in $ROOT_DIR/bin"
