#!/usr/bin/env bash
# Common helper library for Synnergy operational scripts.
# shellcheck disable=SC2317

if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    echo "runtime.sh must be sourced, not executed" >&2
    exit 1
fi

if [[ -n "${__SYN_RUNTIME_LIB_SOURCED:-}" ]]; then
    return
fi
__SYN_RUNTIME_LIB_SOURCED=1

SYN_DRY_RUN=${SYN_DRY_RUN:-0}
SYN_CLI=${SYN_CLI:-synnergy}
ALLOW_CLI_FALLBACK=${ALLOW_CLI_FALLBACK:-1}
SYN_CLI_AVAILABLE=0

_syn_timestamp() {
    date -u +"%Y-%m-%dT%H:%M:%SZ"
}

syn_log() {
    local level=$1
    shift
    printf '%s [%s] %s\n' "$(_syn_timestamp)" "$level" "$*" >&2
}

syn_info() {
    syn_log INFO "$*"
}

syn_warn() {
    syn_log WARN "$*"
}

syn_error() {
    syn_log ERROR "$*"
}

syn_fail() {
    syn_error "$*"
    exit 1
}

syn_require_tool() {
    local tool=$1
    if ! command -v "$tool" >/dev/null 2>&1; then
        syn_fail "required tool '$tool' not found in PATH"
    fi
}

syn_require_file() {
    local file=$1
    if [[ ! -f "$file" ]]; then
        syn_fail "required file '$file' does not exist"
    fi
}

syn_ensure_cli() {
    if command -v "$SYN_CLI" >/dev/null 2>&1; then
        SYN_CLI_AVAILABLE=1
        return
    fi

    SYN_CLI_AVAILABLE=0
    if [[ "$ALLOW_CLI_FALLBACK" -eq 1 ]]; then
        syn_warn "CLI '$SYN_CLI' not found; switching to dry-run mode"
        SYN_DRY_RUN=1
    else
        syn_fail "CLI '$SYN_CLI' not available"
    fi
}

syn_run() {
    local desc=$1
    shift
    if [[ "$SYN_DRY_RUN" -eq 1 ]]; then
        syn_info "[dry-run] $desc :: $*"
        return 0
    fi
    syn_info "$desc"
    "$@"
}

syn_capture() {
    local desc=$1
    shift
    if [[ "$SYN_DRY_RUN" -eq 1 ]]; then
        syn_info "[dry-run] $desc :: $*"
        return 0
    fi
    syn_info "$desc"
    "$@"
}

syn_write_json() {
    local file=$1
    local payload=$2
    printf '%s\n' "$payload" >"$file"
}

syn_expand_path() {
    local path=$1
    if [[ -z "$path" ]]; then
        printf '\n'
        return
    fi
    case $path in
        ~*) eval printf '%s\n' "${path}" ;;
        *) printf '%s\n' "$path" ;;
    esac
}
