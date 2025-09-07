#!/usr/bin/env bash
set -euo pipefail

if ! command -v go >/dev/null 2>&1; then
  echo "go command not found" >&2
  exit 1
fi

# Build the synnergy CLI with trimmed debug paths
GO111MODULE=on go build -trimpath -o synnergy ../synnergy
