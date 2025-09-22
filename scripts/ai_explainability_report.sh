#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# shellcheck disable=SC1091
. "$SCRIPT_DIR/lib/runtime.sh"

usage() {
    cat <<'USAGE'
Usage: ai_explainability_report.sh --input FILE [options]

Generate a lightweight explainability report for inference payloads.

Options:
  --input FILE              JSON file with model input features (object or array).
  --model HASH              Optional model hash for metadata.
  --output FILE             Persist report to FILE instead of STDOUT.
  --top-k N                 Limit feature breakdown to the top N contributors.
  --label TEXT              Human readable scenario description.
  --dry-run                 Validate inputs without producing a report.
  -h, --help                Display this help text.
USAGE
}

INPUT=""
MODEL=""
OUTPUT=""
TOPK=0
LABEL=""

while [[ $# -gt 0 ]]; do
    case "$1" in
        --input)
            INPUT="$2"
            shift 2
            ;;
        --model)
            MODEL="$2"
            shift 2
            ;;
        --output)
            OUTPUT="$2"
            shift 2
            ;;
        --top-k)
            TOPK="$2"
            shift 2
            ;;
        --label)
            LABEL="$2"
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

if [[ -z "$INPUT" ]]; then
    syn_fail "--input is required"
fi

syn_require_file "$INPUT"

if [[ $SYN_DRY_RUN -eq 1 ]]; then
    syn_info "[dry-run] validated input $INPUT"
    exit 0
fi

report=$(python3 - "$INPUT" "$MODEL" "$TOPK" "$LABEL" <<'PY'
import json
import math
import sys
from pathlib import Path
from statistics import mean

input_path = Path(sys.argv[1])
model_hash = sys.argv[2]
top_k = int(sys.argv[3]) if sys.argv[3].isdigit() else 0
label = sys.argv[4]

def flatten(data):
    if isinstance(data, dict):
        return data
    if isinstance(data, list):
        return {f"feature_{idx}": value for idx, value in enumerate(data)}
    raise ValueError("Unsupported input structure; expected JSON object or array")

payload = json.loads(input_path.read_text(encoding="utf-8"))
features = flatten(payload)

scores = []
for name, value in features.items():
    try:
        numeric = float(value)
    except (TypeError, ValueError):
        numeric = float(len(str(value)))
    contribution = abs(numeric)
    scores.append((name, numeric, contribution))

if not scores:
    raise SystemExit("No features available for explanation")

total = sum(score for _, _, score in scores)
if math.isclose(total, 0.0):
    normalised = [(name, val, 0.0) for name, val, _ in scores]
else:
    normalised = [(name, val, score / total) for name, val, score in scores]

normalised.sort(key=lambda item: item[2], reverse=True)
if top_k > 0:
    normalised = normalised[:top_k]

report = {
    "model_hash": model_hash or None,
    "label": label or None,
    "feature_importance": [
        {
            "feature": name,
            "raw_value": value,
            "relative_importance": round(weight, 6),
        }
        for name, value, weight in normalised
    ],
    "feature_count": len(features),
    "mean_feature_value": mean([abs(v) for _, v, _ in normalised]) if normalised else 0,
}

print(json.dumps(report, indent=2))
PY
)

if [[ -n "$OUTPUT" ]]; then
    OUTPUT=$(syn_expand_path "$OUTPUT")
    syn_info "Writing explainability report to $OUTPUT"
    syn_write_json "$OUTPUT" "$report"
else
    echo "$report"
fi
