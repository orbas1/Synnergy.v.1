#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'USAGE'
Usage: $(basename "$0") [version]
Builds the synnergy binary and packages it with checksums for release.

Arguments:
  version  Optional version tag to embed in archive name.
USAGE
}

if [[ "${1:-}" == "--help" || "${1:-}" == "-h" ]]; then
  usage
  exit 0
fi

VERSION="${1:-latest}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$SCRIPT_DIR/.."
DIST_DIR="$ROOT_DIR/dist/$VERSION"
BIN_NAME="synnergy"

mkdir -p "$DIST_DIR"
cd "$ROOT_DIR"

echo "Building $BIN_NAME..."
go build -o "$DIST_DIR/$BIN_NAME" ./cmd/synnergy

echo "Packaging artifacts..."
tar -C "$DIST_DIR" -czf "$DIST_DIR/${BIN_NAME}_${VERSION}.tar.gz" "$BIN_NAME"

echo "Generating checksums..."
sha256sum "$DIST_DIR/${BIN_NAME}_${VERSION}.tar.gz" > "$DIST_DIR/${BIN_NAME}_${VERSION}.tar.gz.sha256"

echo "Release package created at $DIST_DIR"
