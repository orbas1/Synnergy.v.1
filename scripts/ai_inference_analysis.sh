#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# shellcheck disable=SC1091
. "$SCRIPT_DIR/lib/runtime.sh"

usage() {
    cat <<'USAGE'
Usage: ai_inference_analysis.sh [options]

Aggregate inference telemetry and produce performance statistics.

Options:
  --metrics FILE            JSON or CSV file of latency values (milliseconds).
  --latency VALUE           Inline latency sample (may be repeated).
  --window SIZE             Rolling window size for percentile computation (default: all).
  --save PATH               Persist aggregated report to PATH.
  --tag NAME                Identifier for the batch being analysed.
  --dry-run                 Validate inputs without writing output.
  -h, --help                Display this help text.
USAGE
}

METRICS_FILE=""
WINDOW=0
OUTPUT=""
TAG=""
declare -a LATENCIES

while [[ $# -gt 0 ]]; do
    case "$1" in
        --metrics)
            METRICS_FILE="$2"
            shift 2
            ;;
        --latency)
            LATENCIES+=("$2")
            shift 2
            ;;
        --window)
            WINDOW="$2"
            shift 2
            ;;
        --save)
            OUTPUT="$2"
            shift 2
            ;;
        --tag)
            TAG="$2"
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

samples=()

if [[ -n "$METRICS_FILE" ]]; then
    syn_require_file "$METRICS_FILE"
    mapfile -t file_values < <(python3 - "$METRICS_FILE" <<'PY'
import csv
import json
import sys
from pathlib import Path

path = Path(sys.argv[1])
content = path.read_text(encoding='utf-8').strip()
if not content:
    sys.exit(0)

values = []
try:
    data = json.loads(content)
    if isinstance(data, dict):
        for value in data.values():
            try:
                values.append(float(value))
            except (TypeError, ValueError):
                continue
    else:
        for value in data:
            try:
                values.append(float(value))
            except (TypeError, ValueError):
                continue
except json.JSONDecodeError:
    reader = csv.reader(content.splitlines())
    for row in reader:
        for cell in row:
            try:
                values.append(float(cell))
            except ValueError:
                continue

for v in values:
    print(v)
PY
)
    samples+=(${file_values[@]})
fi

for value in "${LATENCIES[@]}"; do
    samples+=("$value")
done

if [[ ${#samples[@]} -eq 0 ]]; then
    syn_fail "no metrics provided"
fi

report=$(python3 - "$WINDOW" "$TAG" <<'PY'
import json
import math
import statistics
import sys

window = int(sys.argv[1])
tag = sys.argv[2]
values = []
for line in sys.stdin:
    line = line.strip()
    if not line:
        continue
    try:
        values.append(float(line))
    except ValueError:
        continue

if not values:
    print(json.dumps({"count": 0}))
    sys.exit(0)

subset = values[-window:] if window > 0 and window <= len(values) else values
subset.sort()

percentile = lambda p: subset[int(max(0, min(len(subset) - 1, round((p / 100) * (len(subset) - 1)))))]

report = {
    "tag": tag or None,
    "count": len(values),
    "window": len(subset),
    "mean_ms": statistics.fmean(subset),
    "p50_ms": percentile(50),
    "p95_ms": percentile(95),
    "p99_ms": percentile(99),
    "min_ms": min(subset),
    "max_ms": max(subset),
}

print(json.dumps(report, indent=2))
PY
<<<"${samples[*]}"
)

if [[ $SYN_DRY_RUN -eq 1 ]]; then
    syn_info "[dry-run] computed metrics"
    echo "$report"
    exit 0
fi

if [[ -n "$OUTPUT" ]]; then
    OUTPUT=$(syn_expand_path "$OUTPUT")
    syn_write_json "$OUTPUT" "$report"
else
    echo "$report"
fi
