#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$SCRIPT_DIR/.."
cd "$ROOT_DIR"

# Format Go files
find . -name '*.go' -not -path './vendor/*' -print0 | xargs -0 gofmt -w

# Run basic linters
if command -v golangci-lint >/dev/null 2>&1; then
  golangci-lint run
else
  go vet ./...
fi

echo "Linting completed"
