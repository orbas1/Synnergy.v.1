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
Build the Synnergy Docker image with repeatable logging, retries and optional
registry pushes.

Usage: docker_build.sh [options]

Options:
  --context DIR          Build context directory (default: project root)
  --file DOCKERFILE      Dockerfile path (default: docker/Dockerfile)
  --tag NAME[:TAG]       Image tag to build (default: synnergy/devnet:latest)
  --platform PLAT        Target platform, e.g. linux/amd64 (default: host)
  --build-arg KEY=VAL    Additional build arguments (repeatable)
  --cache-from IMAGE     Use an existing image as build cache
  --push                 Push the image after a successful build
  --registry HOST/NS     Registry namespace used when pushing
  --load                 Load the image into the local Docker engine (for buildx)
  -h, --help             Show this help text

Common flags:
  --dry-run              Print the docker commands instead of executing them
  --timeout SEC          Override the per-command timeout (default: 120)
  --log-file PATH        Append logs to PATH instead of scripts/logs/
USAGE
}

CONTEXT=${SYN_DOCKER_CONTEXT:-$PROJECT_ROOT}
DOCKERFILE=${SYN_DOCKER_FILE:-$PROJECT_ROOT/docker/Dockerfile}
TAG=${SYN_DOCKER_TAG:-synnergy/devnet:latest}
PLATFORM=${SYN_DOCKER_PLATFORM:-}
PUSH=false
LOAD_IMAGE=false
REGISTRY_NAMESPACE=${SYN_DOCKER_REGISTRY:-}
CACHE_FROM=${SYN_DOCKER_CACHE_FROM:-}
BUILD_ARGS=()

while [[ $# -gt 0 ]]; do
  case "$1" in
    --context)
      CONTEXT="$2"
      shift 2
      ;;
    --file)
      DOCKERFILE="$2"
      shift 2
      ;;
    --tag|-t)
      TAG="$2"
      shift 2
      ;;
    --platform)
      PLATFORM="$2"
      shift 2
      ;;
    --build-arg)
      BUILD_ARGS+=("$2")
      shift 2
      ;;
    --cache-from)
      CACHE_FROM="$2"
      shift 2
      ;;
    --push)
      PUSH=true
      shift
      ;;
    --registry)
      REGISTRY_NAMESPACE="$2"
      shift 2
      ;;
    --load)
      LOAD_IMAGE=true
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

require_command docker

if [[ ! -d "$CONTEXT" ]]; then
  log_error "Build context does not exist: $CONTEXT"
  exit 1
fi
if [[ ! -f "$DOCKERFILE" ]]; then
  log_error "Dockerfile not found: $DOCKERFILE"
  exit 1
fi

log_info "Building Docker image" "context=$CONTEXT" "dockerfile=$DOCKERFILE" "tag=$TAG"

cmd=(docker build --progress=plain -f "$DOCKERFILE" -t "$TAG")
if [[ -n "$PLATFORM" ]]; then
  cmd+=(--platform "$PLATFORM")
fi
for arg in "${BUILD_ARGS[@]}"; do
  cmd+=(--build-arg "$arg")
  log_info "Adding build arg" "$arg"
done
if [[ -n "$CACHE_FROM" ]]; then
  cmd+=(--cache-from "$CACHE_FROM")
fi
if [[ "$LOAD_IMAGE" == true ]]; then
  cmd+=(--load)
fi
cmd+=("$CONTEXT")

if [[ "$DRY_RUN" == true ]]; then
  log_info "[dry-run] ${cmd[*]}"
else
  with_timeout "docker build" "${cmd[@]}"
fi

maybe_tag_for_registry() {
  local src="$1"
  local registry="$2"
  if [[ -z "$registry" ]]; then
    printf '%s\n' "$src"
    return 0
  fi
  local name="${src%%:*}"
  local ref="${src#*:}"
  if [[ "$src" == "$name" ]]; then
    ref="latest"
  fi
  registry="${registry%/}"
  local qualified="$registry/$name:$ref"
  if [[ "$DRY_RUN" == true ]]; then
    log_info "[dry-run] docker tag $src $qualified"
  else
    with_timeout "docker tag" docker tag "$src" "$qualified"
  fi
  printf '%s\n' "$qualified"
}

if [[ "$PUSH" == true ]]; then
  target_tag="$TAG"
  target_tag="$(maybe_tag_for_registry "$target_tag" "$REGISTRY_NAMESPACE")"
  if [[ "$DRY_RUN" == true ]]; then
    log_info "[dry-run] docker push $target_tag"
  else
    with_timeout "docker push" docker push "$target_tag"
  fi
fi

log_info "Docker build completed" "tag=$TAG" "pushed=$PUSH"
