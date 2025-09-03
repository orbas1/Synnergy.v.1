#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Usage: $(basename "$0") <ledger_dir> <backup_dir>
Creates a timestamped tar.gz backup of the ledger directory.

Arguments:
  ledger_dir  Path to the ledger data directory.
  backup_dir  Destination directory for backups.
USAGE
}

if [[ "${1:-}" == "--help" || "${1:-}" == "-h" ]]; then
  usage
  exit 0
fi

if [[ $# -lt 2 ]]; then
  echo "Missing arguments" >&2
  usage
  exit 1
fi

LEDGER_DIR="$1"
DEST_DIR="$2"

if [[ ! -d "$LEDGER_DIR" ]]; then
  echo "Ledger directory not found: $LEDGER_DIR" >&2
  exit 1
fi

mkdir -p "$DEST_DIR"
TIMESTAMP="$(date +%Y%m%d%H%M%S)"
ARCHIVE="$DEST_DIR/ledger_backup_$TIMESTAMP.tar.gz"

tar -C "$LEDGER_DIR" -czf "$ARCHIVE" .

echo "Ledger backup created at $ARCHIVE"

