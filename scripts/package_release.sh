#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
# shellcheck source=lib/common.sh
source "$SCRIPT_DIR/lib/common.sh"

usage() {
  cat <<'USAGE'
Build and package the Synnergy CLI for release.

Usage: package_release.sh [--version VERSION] [--output DIR] [--skip-tests]
                          [--sign]

Options:
  --version VERSION  Semantic version tag applied to artifacts (default: 0.0.0-dev).
  --output DIR       Directory to store build artifacts (default: ./dist).
  --skip-tests       Skip running go test ./... prior to packaging.
  --sign             Generate SHA512 digest for release_sign_verify.sh.
  -h, --help         Show this help message.
USAGE
}

VERSION="0.0.0-dev"
OUTPUT_DIR="$ROOT_DIR/dist"
SKIP_TESTS=false
SIGN=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    --version)
      VERSION=$2
      shift 2
      ;;
    --output)
      OUTPUT_DIR=$2
      shift 2
      ;;
    --skip-tests)
      SKIP_TESTS=true
      shift
      ;;
    --sign)
      SIGN=true
      shift
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

require_commands go tar
mkdir -p "$OUTPUT_DIR"

if [[ $SKIP_TESTS == false ]]; then
  log_info "Running go test ./..."
  (
    cd "$ROOT_DIR"
    go test ./...
  )
fi

log_info "Building synnergy CLI"
(
  cd "$ROOT_DIR"
  go build -o "$OUTPUT_DIR/synnergy" ./cmd/synnergy
)

ARCHIVE="$OUTPUT_DIR/synnergy-$VERSION.tar.gz"
log_info "Creating archive $ARCHIVE"
(
  cd "$OUTPUT_DIR"
  tar -czf "$(basename "$ARCHIVE")" synnergy
)

if [[ $SIGN == true ]]; then
  require_commands openssl
  log_info "Generating SHA512 digest"
  (
    cd "$OUTPUT_DIR"
    openssl dgst -sha512 "$(basename "$ARCHIVE")" >"$(basename "$ARCHIVE").sha512"
  )
fi

log_info "Package build complete"
