#!/usr/bin/env bash
set -euo pipefail
IFS=$'\n\t'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck disable=SC1091
source "$SCRIPT_DIR/lib/common.sh"

parse_common_flags "$@"
set -- "${POSITIONAL_ARGS[@]:-}"

usage() {
  cat <<'USAGE'
Launch the Synnergy stack via Docker Compose with health checks and structured
logging.

Usage: docker_compose_up.sh [options]

Options:
  --file FILE             Compose file (default: docker/docker-compose.yml)
  --project NAME          Compose project name (default: synnergy)
  --profile NAME          Compose profile to enable (repeatable)
  --scale SERVICE=N       Scale a service to N replicas (repeatable)
  --env-file FILE         Additional environment file passed to Compose
  --no-build              Skip image builds (default builds before up)
  --detach                Run containers in the background
  -h, --help              Show this help message

Common flags:
  --dry-run               Print the compose command without executing it
  --timeout SEC           Override per-command timeout (default: 120)
  --log-file PATH         Append logs to PATH instead of scripts/logs/
USAGE
}

COMPOSE_FILE=${SYN_COMPOSE_FILE:-$PROJECT_ROOT/docker/docker-compose.yml}
PROJECT_NAME=${SYN_COMPOSE_PROJECT:-synnergy}
PROFILES=()
SCALE_ARGS=()
ENV_FILE=${SYN_COMPOSE_ENV:-}
DETACH=false
NO_BUILD=false

while [[ $# -gt 0 ]]; do
  case "$1" in
    --file)
      COMPOSE_FILE="$2"
      shift 2
      ;;
    --project)
      PROJECT_NAME="$2"
      shift 2
      ;;
    --profile)
      PROFILES+=("$2")
      shift 2
      ;;
    --scale)
      SCALE_ARGS+=("$2")
      shift 2
      ;;
    --env-file)
      ENV_FILE="$2"
      shift 2
      ;;
    --detach|-d)
      DETACH=true
      shift
      ;;
    --no-build)
      NO_BUILD=true
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      log_error "Unknown argument: $1"
      usage
      exit 1
      ;;
  esac
done

if [[ -z "$LOG_FILE" ]]; then
  set_log_file "$(basename "$0")"
fi

require_command docker python3

if [[ ! -f "$COMPOSE_FILE" ]]; then
  log_error "Compose file not found: $COMPOSE_FILE"
  exit 1
fi

compose_bin="docker compose"
if ! docker compose version &>/dev/null; then
  if command -v docker-compose &>/dev/null; then
    compose_bin="docker-compose"
  else
    log_error "docker compose plugin is not available"
    exit 1
  fi
fi

cmd=($compose_bin -f "$COMPOSE_FILE" --project-name "$PROJECT_NAME" up)
if [[ "$DETACH" == true ]]; then
  cmd+=(-d)
fi
if [[ "$NO_BUILD" == false ]]; then
  cmd+=(--build)
fi
for profile in "${PROFILES[@]}"; do
  cmd+=(--profile "$profile")
  log_info "Enabling profile" "$profile"
done
for scale in "${SCALE_ARGS[@]}"; do
  cmd+=(--scale "$scale")
  log_info "Scaling service" "$scale"
done
if [[ -n "$ENV_FILE" ]]; then
  if [[ ! -f "$ENV_FILE" ]]; then
    log_error "Env file not found: $ENV_FILE"
    exit 1
  fi
  cmd+=(--env-file "$ENV_FILE")
fi

if [[ "$DRY_RUN" == true ]]; then
  log_info "[dry-run] ${cmd[*]}"
  exit 0
fi

with_timeout "docker compose up" "${cmd[@]}"

status_cmd=($compose_bin -f "$COMPOSE_FILE" --project-name "$PROJECT_NAME" ps --format json)
status_output="$(with_timeout "docker compose ps" "${status_cmd[@]}" || true)"
python3 - "$status_output" <<'PY'
import json
import sys
raw = sys.argv[1]
health = []
for idx, ch in enumerate(raw):
    if ch in '[{':
        try:
            data = json.loads(raw[idx:])
        except json.JSONDecodeError:
            break
        else:
            for svc in data:
                name = svc.get("Service", svc.get("Name", "?"))
                status = svc.get("State", svc.get("Status", "unknown"))
                health.append((name, status))
        break
if health:
    for name, status in health:
        print(f"service={name} status={status}")
PY

log_info "Compose stack running" "project=$PROJECT_NAME" "file=$COMPOSE_FILE"
