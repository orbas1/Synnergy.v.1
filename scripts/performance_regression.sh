#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
# shellcheck source=lib/common.sh
source "$SCRIPT_DIR/lib/common.sh"

usage() {
  cat <<'USAGE'
Run performance regression checks across core packages.

Usage: performance_regression.sh [--bench BENCH] [--output FILE]
                                 [--packages PKG,...]

Options:
  --bench BENCH     go test -bench filter (default: .).
  --output FILE     File to write benchmark results.
  --packages LIST   Comma separated package list (default: ./core,...).
  -h, --help        Show this help message.
USAGE
}

BENCH="."
OUTPUT=""
PACKAGES="./core,..."

while [[ $# -gt 0 ]]; do
  case "$1" in
    --bench)
      BENCH=$2
      shift 2
      ;;
    --output)
      OUTPUT=$2
      shift 2
      ;;
    --packages)
      PACKAGES=$2
      shift 2
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

require_commands go tee

IFS=',' read -r -a pkg_array <<<"$PACKAGES"

log_info "Running benchmarks with filter $BENCH"

{
  for pkg in "${pkg_array[@]}"; do
    log_info "Benchmarking $pkg"
    (cd "$ROOT_DIR" && go test -bench="$BENCH" -run=^$ "$pkg")
  done
} | {
  if [[ -n $OUTPUT ]]; then
    tee "$OUTPUT"
  else
    cat
  fi
}

log_info "Performance regression run complete"
