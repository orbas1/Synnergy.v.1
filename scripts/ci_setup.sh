#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Usage: $(basename "$0")
Prepares the repository for continuous integration by tidying modules,
verifying dependencies, linting, testing and building binaries.
USAGE
}

if [[ "${1:-}" == "--help" || "${1:-}" == "-h" ]]; then
  usage
  exit 0
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$SCRIPT_DIR/.."
cd "$ROOT_DIR"

echo "Tidying modules..."
go mod tidy

echo "Verifying modules..."
go mod verify

echo "Linting source..."
"$SCRIPT_DIR/lint.sh"

echo "Running tests..."
"$SCRIPT_DIR/run_tests.sh"

echo "Building binaries..."
"$SCRIPT_DIR/build_all.sh"

echo "CI setup complete"

