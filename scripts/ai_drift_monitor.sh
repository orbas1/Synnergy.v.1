#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# shellcheck disable=SC1091
. "$SCRIPT_DIR/lib/runtime.sh"

usage() {
    cat <<'USAGE'
Usage: ai_drift_monitor.sh [options]

Maintain AI drift baselines and evaluate model health.

Options:
  --model HASH              Target model hash (required for update/check).
  --metric VALUE            Latest metric value used for drift detection.
  --threshold VALUE         Acceptable deviation from baseline (default: 0.05).
  --baseline VALUE          Override baseline value when updating.
  --update-baseline         Persist the provided --baseline or --metric as baseline.
  --check                   Evaluate drift using the provided --metric.
  --list                    Display known baselines.
  --state PATH              Custom baseline store (default: ~/.synnergy/drift_baselines.json).
  --dry-run                 Log operations without mutating state.
  -h, --help                Show this help text.
USAGE
}

STATE_FILE=${DRIFT_STATE_FILE:-"$HOME/.synnergy/drift_baselines.json"}
MODEL=""
METRIC=""
BASELINE=""
THRESHOLD=0.05
DO_UPDATE=0
DO_CHECK=0
DO_LIST=0

while [[ $# -gt 0 ]]; do
    case "$1" in
        --model)
            MODEL="$2"
            shift 2
            ;;
        --metric)
            METRIC="$2"
            shift 2
            ;;
        --threshold)
            THRESHOLD="$2"
            shift 2
            ;;
        --baseline)
            BASELINE="$2"
            shift 2
            ;;
        --update-baseline)
            DO_UPDATE=1
            shift
            ;;
        --check)
            DO_CHECK=1
            shift
            ;;
        --list)
            DO_LIST=1
            shift
            ;;
        --state)
            STATE_FILE="$2"
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

STATE_FILE=$(syn_expand_path "$STATE_FILE")

ensure_state_file() {
    local file=$1
    if [[ ! -f "$file" ]]; then
        syn_info "Creating baseline state file at $file"
        if [[ $SYN_DRY_RUN -eq 0 ]]; then
            mkdir -p "$(dirname "$file")"
            echo '{}' >"$file"
        fi
    fi
}

ensure_state_file "$STATE_FILE"

if [[ $DO_LIST -eq 1 ]]; then
    python3 - "$STATE_FILE" <<'PY'
import json
import sys

path = sys.argv[1]
try:
    with open(path, 'r', encoding='utf-8') as fh:
        data = json.load(fh)
except FileNotFoundError:
    data = {}
print(json.dumps(data, indent=2))
PY
    exit 0
fi

if [[ $DO_UPDATE -eq 1 ]]; then
    if [[ -z "$MODEL" ]]; then
        syn_fail "--model is required when updating baselines"
    fi
    value="$BASELINE"
    if [[ -z "$value" ]]; then
        value="$METRIC"
    fi
    if [[ -z "$value" ]]; then
        syn_fail "provide --baseline or --metric for update"
    fi
    if [[ $SYN_DRY_RUN -eq 1 ]]; then
        syn_info "[dry-run] would set baseline for $MODEL to $value"
    else
        python3 - "$STATE_FILE" "$MODEL" "$value" <<'PY'
import json
import sys
from pathlib import Path

path = Path(sys.argv[1])
model = sys.argv[2]
value = float(sys.argv[3])

if path.exists():
    data = json.loads(path.read_text(encoding='utf-8'))
else:
    data = {}

data[model] = value
path.write_text(json.dumps(data, indent=2, sort_keys=True), encoding='utf-8')
print(json.dumps({"model": model, "baseline": value}))
PY
    fi
fi

if [[ $DO_CHECK -eq 1 ]]; then
    if [[ -z "$MODEL" ]]; then
        syn_fail "--model is required for drift checks"
    fi
    if [[ -z "$METRIC" ]]; then
        syn_fail "--metric is required for drift checks"
    fi
    python3 - "$STATE_FILE" "$MODEL" "$METRIC" "$THRESHOLD" <<'PY'
import json
import math
import sys
from pathlib import Path

path = Path(sys.argv[1])
model = sys.argv[2]
metric = float(sys.argv[3])
threshold = float(sys.argv[4])

if not path.exists():
    print(json.dumps({"model": model, "baseline": None, "drift": False, "reason": "no baseline"}))
    sys.exit(0)

data = json.loads(path.read_text(encoding='utf-8'))
base = data.get(model)
if base is None:
    print(json.dumps({"model": model, "baseline": None, "drift": False, "reason": "baseline missing"}))
    sys.exit(0)

diff = abs(metric - base)
print(json.dumps({
    "model": model,
    "baseline": base,
    "metric": metric,
    "threshold": threshold,
    "drift": diff > threshold,
    "delta": diff,
}))
PY
fi
