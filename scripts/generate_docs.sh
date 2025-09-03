#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Usage: $(basename "$0")
Builds the project documentation using MkDocs.
USAGE
}

if [[ "${1:-}" == "--help" || "${1:-}" == "-h" ]]; then
  usage
  exit 0
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$SCRIPT_DIR/.."
cd "$ROOT_DIR"

if command -v mkdocs >/dev/null 2>&1; then
  mkdocs build
  echo "Documentation generated at site/"
else
  echo "mkdocs command not found" >&2
  exit 1
fi

