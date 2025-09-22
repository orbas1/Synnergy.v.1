#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# shellcheck disable=SC1091
. "$SCRIPT_DIR/lib/runtime.sh"

usage() {
    cat <<'USAGE'
Usage: access_control_setup.sh [options]

Bootstrap Synnergy role-based access control policies through the CLI.

Options:
  --assign ROLE:ADDRESS       Grant ROLE to ADDRESS (may be repeated).
  --revoke ROLE:ADDRESS       Revoke ROLE from ADDRESS (may be repeated).
  --from-file PATH            Read ROLE,ADDRESS pairs from PATH (CSV).
  --verify                    Validate role assignments after execution.
  --cli PATH                  Path to the synnergy CLI binary.
  --dry-run                   Print intended actions without executing them.
  -h, --help                  Display this help message.

The CSV format for --from-file expects two columns without headers:
  role,address
Blank lines and lines starting with '#" are ignored.
USAGE
}

assignments=()
revocations=()
verify=0
policy_file=""

while [[ $# -gt 0 ]]; do
    case "$1" in
        --assign)
            assignments+=("$2")
            shift 2
            ;;
        --revoke)
            revocations+=("$2")
            shift 2
            ;;
        --from-file)
            policy_file="$2"
            shift 2
            ;;
        --verify)
            verify=1
            shift
            ;;
        --cli)
            SYN_CLI="$2"
            shift 2
            ;;
        --dry-run)
            SYN_DRY_RUN=1
            shift
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            syn_fail "unknown option: $1"
            ;;
    esac
done

if [[ -n "$policy_file" ]]; then
    syn_require_file "$policy_file"
    while IFS=',' read -r role addr; do
        [[ -z "$role" ]] && continue
        [[ "$role" =~ ^# ]] && continue
        role=$(echo "$role" | xargs)
        addr=$(echo "${addr:-}" | xargs)
        if [[ -z "$role" || -z "$addr" ]]; then
            continue
        fi
        assignments+=("$role:$addr")
    done < "$policy_file"
fi

if [[ ${#assignments[@]} -eq 0 && ${#revocations[@]} -eq 0 ]]; then
    syn_fail "no assignments or revocations provided"
fi

syn_ensure_cli()

granted=0
revoked=0

for entry in "${assignments[@]}"; do
    role=${entry%%:*}
    addr=${entry#*:}
    [[ -z "$role" || -z "$addr" ]] && continue
    syn_run "Granting role $role to $addr" "$SYN_CLI" access grant "$role" "$addr"
    ((granted++))
done

for entry in "${revocations[@]}"; do
    role=${entry%%:*}
    addr=${entry#*:}
    [[ -z "$role" || -z "$addr" ]] && continue
    syn_run "Revoking role $role from $addr" "$SYN_CLI" access revoke "$role" "$addr"
    ((revoked++))
done

if [[ $verify -eq 1 && $SYN_DRY_RUN -eq 0 ]]; then
    for entry in "${assignments[@]}"; do
        role=${entry%%:*}
        addr=${entry#*:}
        mapfile -t roles < <(syn_capture "Listing roles for $addr" "$SYN_CLI" access list "$addr")
        found=0
        for r in "${roles[@]}"; do
            [[ "$r" == "$role" ]] && found=1
        done
        if [[ $found -eq 0 ]]; then
            syn_warn "role $role not detected for $addr after grant"
        fi
    done
fi

summary=$(cat <<JSON
{
  "granted": $granted,
  "revoked": $revoked,
  "dry_run": $SYN_DRY_RUN
}
JSON
)

echo "$summary"
