#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
# shellcheck source=lib/common.sh
source "$SCRIPT_DIR/lib/common.sh"

usage() {
  cat <<'USAGE'
Sign or verify Synnergy release archives.

Usage: release_sign_verify.sh --archive FILE [--key KEYFILE] [--verify]
                              [--signature SIGFILE]

Options:
  --archive FILE     Release archive to sign/verify.
  --key KEYFILE      Private key for signing (PEM). Required unless --verify.
  --verify           Verify instead of sign.
  --signature FILE   Signature file path (default: FILE.sha512).
  -h, --help         Show help message.
USAGE
}

ARCHIVE=""
KEY=""
VERIFY=false
SIGNATURE=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --archive)
      ARCHIVE=$2
      shift 2
      ;;
    --key)
      KEY=$2
      shift 2
      ;;
    --verify)
      VERIFY=true
      shift
      ;;
    --signature)
      SIGNATURE=$2
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

[[ -z $ARCHIVE ]] && fail "--archive is required"
require_commands openssl

if [[ -z $SIGNATURE ]]; then
  SIGNATURE="$ARCHIVE.sha512"
fi

if [[ $VERIFY == true ]]; then
  [[ -n $KEY ]] || fail "--key certificate required for verification"
  log_info "Verifying $ARCHIVE"
  [[ -f $SIGNATURE ]] || fail "signature file $SIGNATURE not found"
  openssl dgst -sha512 -verify <(openssl x509 -in "$KEY" -pubkey -noout) -signature "$SIGNATURE" "$ARCHIVE"
  log_info "Verification complete"
else
  [[ -n $KEY ]] || fail "--key required for signing"
  log_info "Signing $ARCHIVE"
  openssl dgst -sha512 -sign "$KEY" -out "$SIGNATURE" "$ARCHIVE"
  log_info "Signature written to $SIGNATURE"
fi
