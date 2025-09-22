#!/usr/bin/env bash

set -euo pipefail
IFS=$'\n\t'

SCRIPT_NAME=$(basename "$0")

usage() {
  cat <<'USAGE'
Usage: certificate_renew.sh --cert <path> --key <path> [options]

Enterprise-grade certificate renewal helper for Synnergy Network nodes.

Required arguments:
  --cert PATH           Existing PEM encoded certificate to inspect for expiry.
  --key PATH            Private key used to sign the renewal CSR. Use --rotate-key
                        to generate a fresh key instead.

Optional arguments:
  --output-dir PATH     Directory to write CSR, rotated keys, and logs. Defaults
                        to <cert directory>/renewals.
  --csr-out PATH        Explicit path for generated certificate signing request.
  --renew-within DAYS   Renew when certificate expires within DAYS (default: 30).
  --force               Force renewal regardless of expiry window.
  --dry-run             Generate CSR and metrics without invoking the CLI.
  --cli PATH            Explicit Synnergy CLI binary to invoke for renewal.
                        Defaults to $SYN_CLI, then synnergy, then synnergy-cli.
  --endpoint URL        Override API endpoint when invoking the CLI.
  --profile NAME        CLI profile (default: production).
  --rotate-key          Generate a brand-new private key for the CSR.
  --key-type TYPE       Key type for rotation: rsa (default) or ec.
  --key-bits BITS       RSA key size when rotating (default: 4096).
  --ec-curve CURVE      EC curve name when rotating (default: prime256v1).
  --subject SUBJECT     CSR subject. Defaults to the current certificate subject.
  --hook PATH           Executable notified post-renewal with CSR, cert, key paths.
  --lock-file PATH      Use flock-based locking to avoid concurrent renewals.
  --metrics-file PATH   Emit JSON metrics describing the renewal activity.
  --help, -h            Display this help message.

Examples:
  # Dry-run renewal and emit metrics.
  certificate_renew.sh --cert node.pem --key node.key --dry-run \\
    --metrics-file /var/log/synnergy/cert_renewal.json

  # Force renewal, rotate key, and submit CSR via custom CLI binary.
  certificate_renew.sh --cert node.pem --rotate-key --output-dir /etc/synnergy/certs \\
    --cli /usr/local/bin/synnergy --force
USAGE
}

log() {
  local level=$1
  local message=$2
  printf '%s [%s] %s\n' "$(date -u +'%Y-%m-%dT%H:%M:%SZ')" "$level" "$message" >&2
}

fatal() {
  local message=$1
  local code=${2:-1}
  log "ERROR" "$message"
  exit "$code"
}

cleanup() {
  local status=$?
  if [[ -n "${TMP_WORKDIR:-}" && -d "${TMP_WORKDIR:-}" ]]; then
    rm -rf "$TMP_WORKDIR"
  fi
  if [[ -n "${LOCK_FD:-}" ]]; then
    eval "exec ${LOCK_FD}>&-"
  fi
  if [[ $status -ne 0 ]]; then
    log "ERROR" "${SCRIPT_NAME} terminated with status $status"
  fi
}
trap cleanup EXIT
trap 'fatal "${SCRIPT_NAME} interrupted" 2' INT TERM

require_command() {
  local cmd=$1
  local description=${2:-$1}
  if ! command -v "$cmd" >/dev/null 2>&1; then
    fatal "$description is required but not installed"
  fi
}

resolve_cli() {
  local candidate=$1
  if [[ -z "$candidate" ]]; then
    if [[ -n "${SYN_CLI:-}" ]]; then
      candidate=$SYN_CLI
    elif command -v synnergy >/dev/null 2>&1; then
      candidate=$(command -v synnergy)
    elif command -v synnergy-cli >/dev/null 2>&1; then
      candidate=$(command -v synnergy-cli)
    else
      fatal "CLI binary not found. Provide --cli or set SYN_CLI"
    fi
  fi

  if [[ -x "$candidate" ]]; then
    echo "$candidate"
    return 0
  fi

  if command -v "$candidate" >/dev/null 2>&1; then
    command -v "$candidate"
    return 0
  fi

  fatal "CLI binary '$candidate' not found"
}

generate_rsa_key() {
  local path=$1
  local bits=$2
  log "INFO" "Generating RSA key (${bits} bits) at $path"
  openssl genpkey -quiet -algorithm RSA -pkeyopt "rsa_keygen_bits:${bits}" -out "$path"
}

generate_ec_key() {
  local path=$1
  local curve=$2
  log "INFO" "Generating EC key (${curve}) at $path"
  openssl ecparam -name "$curve" -genkey -noout -out "$path"
}

emit_metrics() {
  local metrics_path=$1
  local csr=$2
  local key_path=$3
  local cert_path=$4
  local days_remaining=$5
  local dry_run=$6
  local cli_invoked=$7
  local profile=$8
  local endpoint=$9

  cat >"$metrics_path" <<METRICS
{
  "certificate": "${cert_path}",
  "csr": "${csr}",
  "key": "${key_path}",
  "days_remaining": ${days_remaining},
  "dry_run": ${dry_run},
  "cli_invoked": ${cli_invoked},
  "profile": "${profile}",
  "endpoint": "${endpoint}",
  "renewed_at": "$(date -u +'%Y-%m-%dT%H:%M:%SZ')"
}
METRICS
}

normalize_subject() {
  local raw=$1
  if [[ -z "$raw" ]]; then
    echo ""
    return
  fi

  raw=${raw#subject=}
  # Normalize spacing around separators without touching escaped spaces inside values.
  raw=$(echo "$raw" | sed -E 's/, */,/g; s/ *= */=/g')

  if [[ "$raw" == /* ]]; then
    echo "$raw"
    return
  fi

  if [[ "$raw" == *"="* ]]; then
    echo "/${raw//,/\/}"
    return
  fi

  echo "$raw"
}

if [[ $# -eq 0 ]]; then
  usage
  exit 1
fi

CERT_PATH=""
KEY_PATH=""
OUTPUT_DIR=""
CSR_OUT=""
RENEW_WITHIN_DAYS=30
FORCE_RENEW=0
DRY_RUN=0
CLI_PATH=${SYN_CLI:-}
ENDPOINT=""
PROFILE="production"
ROTATE_KEY=0
KEY_TYPE="rsa"
KEY_BITS=4096
EC_CURVE="prime256v1"
SUBJECT_OVERRIDE=""
HOOK_SCRIPT=""
LOCK_FILE=""
LOCK_FD=""
METRICS_FILE=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --help|-h)
      usage
      exit 0
      ;;
    --cert)
      CERT_PATH=$2
      shift 2
      ;;
    --key)
      KEY_PATH=$2
      shift 2
      ;;
    --output-dir)
      OUTPUT_DIR=$2
      shift 2
      ;;
    --csr-out)
      CSR_OUT=$2
      shift 2
      ;;
    --renew-within)
      RENEW_WITHIN_DAYS=$2
      shift 2
      ;;
    --force)
      FORCE_RENEW=1
      shift
      ;;
    --dry-run)
      DRY_RUN=1
      shift
      ;;
    --cli)
      CLI_PATH=$2
      shift 2
      ;;
    --endpoint)
      ENDPOINT=$2
      shift 2
      ;;
    --profile)
      PROFILE=$2
      shift 2
      ;;
    --rotate-key)
      ROTATE_KEY=1
      shift
      ;;
    --key-type)
      KEY_TYPE=$2
      shift 2
      ;;
    --key-bits)
      KEY_BITS=$2
      shift 2
      ;;
    --ec-curve)
      EC_CURVE=$2
      shift 2
      ;;
    --subject)
      SUBJECT_OVERRIDE=$2
      shift 2
      ;;
    --hook)
      HOOK_SCRIPT=$2
      shift 2
      ;;
    --lock-file)
      LOCK_FILE=$2
      shift 2
      ;;
    --metrics-file)
      METRICS_FILE=$2
      shift 2
      ;;
    *)
      printf 'Unknown option: %s\n' "$1" >&2
      usage
      exit 1
      ;;
  esac
  done

require_command openssl "OpenSSL"
require_command sed "sed"

if [[ -z "$CERT_PATH" ]]; then
  fatal "--cert is required"
fi

if [[ -z "$OUTPUT_DIR" ]]; then
  OUTPUT_DIR="$(dirname "$CERT_PATH")/renewals"
fi

if [[ -n "$LOCK_FILE" ]]; then
  require_command flock "flock"
  LOCK_FD=200
  eval "exec ${LOCK_FD}>\"${LOCK_FILE}\""
  if ! flock -n "$LOCK_FD"; then
    fatal "Unable to acquire lock at ${LOCK_FILE}; another renewal may be running"
  fi
fi

mkdir -p "$OUTPUT_DIR"
umask 077

if [[ ! -f "$CERT_PATH" ]]; then
  fatal "certificate '$CERT_PATH' not found"
fi

TMP_WORKDIR=$(mktemp -d "${OUTPUT_DIR%/}/.renew-XXXXXX")

if [[ $ROTATE_KEY -eq 1 ]]; then
  KEY_PATH="${TMP_WORKDIR}/rotated.key"
  case "$KEY_TYPE" in
    rsa)
      generate_rsa_key "$KEY_PATH" "$KEY_BITS"
      ;;
    ec)
      generate_ec_key "$KEY_PATH" "$EC_CURVE"
      ;;
    *)
      fatal "Unsupported key type '$KEY_TYPE'. Use rsa or ec."
      ;;
  esac
elif [[ -z "$KEY_PATH" ]]; then
  fatal "--key is required unless --rotate-key is supplied"
fi

if [[ ! -f "$KEY_PATH" ]]; then
  fatal "private key '$KEY_PATH' not found"
fi

if [[ -z "$CSR_OUT" ]]; then
  CSR_OUT="${OUTPUT_DIR%/}/$(basename "${CERT_PATH}").$(date -u +'%Y%m%dT%H%M%SZ').csr"
fi

expiry_raw=$(openssl x509 -enddate -noout -in "$CERT_PATH" | cut -d= -f2 || true)
if [[ -z "$expiry_raw" ]]; then
  fatal "unable to read expiry from certificate $CERT_PATH"
fi

if ! expiry_epoch=$(date -d "$expiry_raw" +%s 2>/dev/null); then
  fatal "unable to parse certificate expiry '${expiry_raw}'"
fi

current_epoch=$(date +%s)
seconds_remaining=$((expiry_epoch - current_epoch))
days_remaining=$((seconds_remaining / 86400))

if [[ $FORCE_RENEW -ne 1 && $days_remaining -gt $RENEW_WITHIN_DAYS ]]; then
  log "INFO" "Certificate has ${days_remaining} day(s) left; renewal not required."
  exit 0
fi

if [[ -z "$SUBJECT_OVERRIDE" ]]; then
  SUBJECT_OVERRIDE=$(openssl x509 -subject -nameopt RFC2253 -noout -in "$CERT_PATH")
fi

SUBJECT_OVERRIDE=$(normalize_subject "$SUBJECT_OVERRIDE")

if [[ -z "$SUBJECT_OVERRIDE" ]]; then
  fatal "unable to determine CSR subject"
fi

CSR_TEMP="${TMP_WORKDIR}/request.csr"
openssl req -new -key "$KEY_PATH" -out "$CSR_TEMP" -subj "$SUBJECT_OVERRIDE" -sha256
mv "$CSR_TEMP" "$CSR_OUT"

log "INFO" "Generated CSR at ${CSR_OUT}"

CLI_INVOKED=0
CLI_EXEC=""

if [[ $DRY_RUN -eq 0 ]]; then
  CLI_EXEC=$(resolve_cli "$CLI_PATH")
  CLI_ARGS=(certificate renew --csr "$CSR_OUT" --profile "$PROFILE")
  if [[ -n "$ENDPOINT" ]]; then
    CLI_ARGS+=(--endpoint "$ENDPOINT")
  fi

  log "INFO" "Invoking CLI: ${CLI_EXEC} ${CLI_ARGS[*]}"
  if ! "$CLI_EXEC" "${CLI_ARGS[@]}"; then
    fatal "CLI renewal command failed"
  fi
  CLI_INVOKED=1
fi

if [[ -n "$HOOK_SCRIPT" ]]; then
  if [[ ! -x "$HOOK_SCRIPT" ]]; then
    fatal "hook '$HOOK_SCRIPT' is not executable"
  fi
  "$HOOK_SCRIPT" "$CSR_OUT" "$CERT_PATH" "$KEY_PATH"
fi

if [[ -n "$METRICS_FILE" ]]; then
  mkdir -p "$(dirname "$METRICS_FILE")"
  emit_metrics "$METRICS_FILE" "$CSR_OUT" "$KEY_PATH" "$CERT_PATH" "$days_remaining" "$DRY_RUN" "$CLI_INVOKED" "$PROFILE" "$ENDPOINT"
  log "INFO" "Metrics written to ${METRICS_FILE}"
fi

TARGET_KEY_PATH="${OUTPUT_DIR%/}/$(basename "$KEY_PATH")"
if [[ "$KEY_PATH" != "$TARGET_KEY_PATH" ]]; then
  cp "$KEY_PATH" "$TARGET_KEY_PATH"
else
  log "INFO" "Key already resides at ${TARGET_KEY_PATH}; skipping copy"
fi

cat <<SUMMARY
CSR_PATH=${CSR_OUT}
KEY_PATH=${TARGET_KEY_PATH}
CERT_PATH=${CERT_PATH}
DAYS_REMAINING=${days_remaining}
DRY_RUN=${DRY_RUN}
CLI_INVOKED=${CLI_INVOKED}
SUMMARY

log "INFO" "${SCRIPT_NAME} completed successfully"
