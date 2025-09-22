#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# shellcheck disable=SC1091
. "$SCRIPT_DIR/lib/runtime.sh"

usage() {
    cat <<'USAGE'
Usage: ai_privacy_preservation.sh --dataset FILE [options]

Sanitise training datasets for compliant AI workflows.

Options:
  --dataset FILE            Source CSV dataset (required).
  --columns COL1,COL2       Columns to anonymise (default: all columns).
  --salt VALUE              Salt applied to hashed values.
  --policy FILE             JSON policy describing forbidden columns.
  --output FILE             Write sanitised dataset to FILE.
  --report FILE             Persist compliance report to FILE.
  --dry-run                 Validate without writing output files.
  -h, --help                Display this help message.
USAGE
}

DATASET=""
COLUMNS=""
SALT=""
POLICY=""
OUTPUT=""
REPORT=""

while [[ $# -gt 0 ]]; do
    case "$1" in
        --dataset)
            DATASET="$2"
            shift 2
            ;;
        --columns)
            COLUMNS="$2"
            shift 2
            ;;
        --salt)
            SALT="$2"
            shift 2
            ;;
        --policy)
            POLICY="$2"
            shift 2
            ;;
        --output)
            OUTPUT="$2"
            shift 2
            ;;
        --report)
            REPORT="$2"
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

[[ -z "$DATASET" ]] && syn_fail "--dataset is required"

syn_require_file "$DATASET"
if [[ -n "$POLICY" ]]; then
    syn_require_file "$POLICY"
fi

if [[ $SYN_DRY_RUN -eq 1 ]]; then
    syn_info "[dry-run] validated dataset $DATASET"
fi

summary=$(python3 - "$DATASET" "$COLUMNS" "$SALT" "$POLICY" "$OUTPUT" "$REPORT" "$SYN_DRY_RUN" <<'PY'
import csv
import hashlib
import json
import sys
from pathlib import Path

source = Path(sys.argv[1])
columns_arg = sys.argv[2]
salt = sys.argv[3]
policy_path = sys.argv[4]
output_path = sys.argv[5]
report_path = sys.argv[6]
dry_run = bool(int(sys.argv[7]))

policy = {}
if policy_path:
    policy = json.loads(Path(policy_path).read_text(encoding='utf-8'))

forbidden = set(policy.get('forbidden_columns', []))
required = set(policy.get('required_columns', []))

with source.open('r', encoding='utf-8', newline='') as fh:
    reader = csv.DictReader(fh)
    headers = reader.fieldnames or []
    missing_required = sorted(required - set(headers))
    violations = sorted(forbidden & set(headers))
    if missing_required:
        print(json.dumps({
            'status': 'error',
            'reason': 'missing required columns',
            'details': missing_required,
        }))
        sys.exit(1)
    if violations:
        print(json.dumps({
            'status': 'error',
            'reason': 'forbidden columns present',
            'details': violations,
        }))
        sys.exit(1)

    if columns_arg:
        target_columns = [c.strip() for c in columns_arg.split(',') if c.strip()]
    else:
        target_columns = headers

    hashed_rows = []
    digest_preview = []
    for row in reader:
        hashed_row = {}
        for key, value in row.items():
            value = value or ''
            if key in target_columns:
                digest = hashlib.sha256((salt + value).encode('utf-8')).hexdigest()
                hashed_row[key] = digest
                digest_preview.append(digest)
            else:
                hashed_row[key] = value
        hashed_rows.append(hashed_row)

if output_path and not dry_run:
    out_file = Path(output_path).expanduser()
    out_file.parent.mkdir(parents=True, exist_ok=True)
    with out_file.open('w', encoding='utf-8', newline='') as fh:
        writer = csv.DictWriter(fh, fieldnames=headers)
        writer.writeheader()
        writer.writerows(hashed_rows)

report = {
    'status': 'ok',
    'source_rows': len(hashed_rows),
    'hashed_columns': target_columns,
    'salted': bool(salt),
    'dataset': str(source),
    'preview_hashes': digest_preview[:5],
}

report_payload = json.dumps(report, indent=2)

if report_path and not dry_run:
    out = Path(report_path).expanduser()
    out.parent.mkdir(parents=True, exist_ok=True)
    out.write_text(report_payload, encoding='utf-8')

print(report_payload)
PY
)

if [[ $SYN_DRY_RUN -eq 1 ]]; then
    syn_info "[dry-run] generated compliance summary"
fi

echo "$summary"
