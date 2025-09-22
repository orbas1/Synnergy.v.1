#!/usr/bin/env bash

set -euo pipefail
IFS=$'\n\t'

SCRIPT_NAME=$(basename "$0")

usage() {
  cat <<'USAGE'
Usage: cli_help_generator.sh --output <path> [options]

Generate consolidated CLI documentation for the Synnergy toolchain. The
generator can source definitions from recorded fixtures or query live binaries
and produces Markdown or JSON for integration with docs, web portals, and
regulatory submissions.

Required arguments:
  --output PATH             Destination file.

Optional arguments:
  --cli PATH                Explicit CLI binary (default: resolve synnergy/synnergy-cli).
  --commands LIST           Comma-separated commands/subcommands to document.
  --command-file PATH       File containing one command per line.
  --definition-file PATH    JSON file with pre-rendered help output.
  --format FORMAT           Output format: markdown or json (default: markdown).
  --title TEXT              Title for generated documentation (default: Synnergy CLI Reference).
  --metrics-file PATH       Write JSON metrics about generated commands.
  --lock-file PATH          Flock-based mutex to prevent concurrent generation.
  --dry-run                 Skip CLI invocations and generation (requires definition file).
  --help, -h                Display this help message.

Examples:
  cli_help_generator.sh --output docs/cli.md --commands "status,ledger sync"
  cli_help_generator.sh --output docs/cli.json --format json --definition-file fixtures/help.json
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
  if [[ -n "${CMD_FILE:-}" && -f "$CMD_FILE" ]]; then
    rm -f "$CMD_FILE"
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

resolve_cli() {
  local candidate=$1
  if [[ -n "$candidate" ]]; then
    if [[ -x "$candidate" ]]; then
      echo "$candidate"
      return 0
    fi
    if command -v "$candidate" >/dev/null 2>&1; then
      command -v "$candidate"
      return 0
    fi
    fatal "CLI binary '$candidate' not found"
  fi
  if command -v synnergy >/dev/null 2>&1; then
    command -v synnergy
    return 0
  fi
  if command -v synnergy-cli >/dev/null 2>&1; then
    command -v synnergy-cli
    return 0
  fi
  fatal "CLI binary not found. Provide --cli or set SYN_CLI"
}

acquire_lock() {
  local path=$1
  if [[ -z "$path" ]]; then
    return
  fi
  mkdir -p "$(dirname "$path")"
  LOCK_FD=216
  eval "exec ${LOCK_FD}>\"$path\""
  if ! flock -n "$LOCK_FD"; then
    fatal "another cli_help_generator.sh invocation holds the lock at $path"
  fi
}

OUTPUT=""
CLI=""
COMMANDS=()
COMMAND_FILE=""
DEFINITION_FILE=""
FORMAT="markdown"
TITLE="Synnergy CLI Reference"
METRICS_FILE=""
LOCK_FILE=""
DRY_RUN=0

while [[ $# -gt 0 ]]; do
  case "$1" in
    --output)
      OUTPUT=${2:-}
      shift 2
      ;;
    --cli)
      CLI=${2:-}
      shift 2
      ;;
    --commands)
      IFS=',' read -r -a tmp <<<"${2:-}"
      for entry in "${tmp[@]}"; do
        entry=$(echo "$entry" | xargs)
        [[ -n "$entry" ]] && COMMANDS+=("$entry")
      done
      shift 2
      ;;
    --command-file)
      COMMAND_FILE=${2:-}
      shift 2
      ;;
    --definition-file)
      DEFINITION_FILE=${2:-}
      shift 2
      ;;
    --format)
      FORMAT=${2:-}
      shift 2
      ;;
    --title)
      TITLE=${2:-}
      shift 2
      ;;
    --metrics-file)
      METRICS_FILE=${2:-}
      shift 2
      ;;
    --lock-file)
      LOCK_FILE=${2:-}
      shift 2
      ;;
    --dry-run)
      DRY_RUN=1
      shift
      ;;
    --help|-h)
      usage
      exit 0
      ;;
    *)
      log "ERROR" "unknown argument: $1"
      usage
      exit 1
      ;;
  esac
done

if [[ -z "$OUTPUT" ]]; then
  usage
  fatal "--output is required"
fi

if [[ $FORMAT != "markdown" && $FORMAT != "json" ]]; then
  fatal "--format must be markdown or json"
fi

if [[ -n "$COMMAND_FILE" ]]; then
  if [[ ! -f "$COMMAND_FILE" ]]; then
    fatal "command file $COMMAND_FILE does not exist"
  fi
  while IFS= read -r line; do
    line=$(echo "$line" | xargs)
    [[ -z "$line" || $line == \#* ]] && continue
    COMMANDS+=("$line")
  done <"$COMMAND_FILE"
fi

if [[ ${#COMMANDS[@]} -eq 0 && -z "$DEFINITION_FILE" ]]; then
  fatal "no commands specified; provide --commands, --command-file, or --definition-file"
fi

if [[ $DRY_RUN -eq 1 && -z "$DEFINITION_FILE" ]]; then
  fatal "--dry-run requires --definition-file"
fi

acquire_lock "$LOCK_FILE"

if [[ -n "$METRICS_FILE" ]]; then
  mkdir -p "$(dirname "$METRICS_FILE")"
fi

if [[ $DRY_RUN -eq 0 || -z "$DEFINITION_FILE" ]]; then
  CLI=$(resolve_cli "$CLI")
fi

CMD_FILE=$(mktemp)
if [[ ${#COMMANDS[@]} -gt 0 ]]; then
  printf '%s\n' "${COMMANDS[@]}" >"$CMD_FILE"
fi

python3 - "$OUTPUT" "$FORMAT" "$TITLE" "$CLI" "$DEFINITION_FILE" "$CMD_FILE" "$METRICS_FILE" "$DRY_RUN" <<'PY'
import json
import subprocess
import sys
import time
from pathlib import Path

output_path = Path(sys.argv[1])
output_format = sys.argv[2]
title = sys.argv[3]
cli = sys.argv[4]
definition_file = sys.argv[5]
command_file = Path(sys.argv[6])
metrics_file = sys.argv[7]
dry_run = sys.argv[8] == "1"

commands = []
if command_file.exists():
    commands = [line.strip() for line in command_file.read_text(encoding="utf-8").splitlines() if line.strip()]
help_entries = []

def load_definitions(path: str):
    data = Path(path).read_text(encoding="utf-8")
    parsed = json.loads(data)
    entries = []
    if isinstance(parsed, list):
        for item in parsed:
            if isinstance(item, dict) and "command" in item and "help" in item:
                entries.append({"command": item["command"], "help": item["help"], "source": "definition"})
    elif isinstance(parsed, dict) and "commands" in parsed:
        for item in parsed["commands"]:
            if isinstance(item, dict) and "command" in item and "help" in item:
                entries.append({"command": item["command"], "help": item["help"], "source": "definition"})
    return entries

if definition_file:
    if not Path(definition_file).is_file():
        print(f"definition file {definition_file} does not exist", file=sys.stderr)
        sys.exit(3)
    help_entries.extend(load_definitions(definition_file))
    if not commands:
        commands = [entry["command"] for entry in help_entries]

captured_commands = {entry["command"] for entry in help_entries}

if not dry_run:
    for command in commands:
        if command in captured_commands:
            continue
        tokens = command.split()
        try:
            result = subprocess.run([cli, *tokens, "--help"], check=True, capture_output=True, text=True)
            help_entries.append({"command": command, "help": result.stdout.strip(), "source": "cli"})
        except subprocess.CalledProcessError as exc:
            print(f"failed to capture help for {command}: {exc.stderr}", file=sys.stderr)
            sys.exit(exc.returncode or 4)
else:
    if not help_entries:
        print("dry-run with no definitions has nothing to do", file=sys.stderr)
        sys.exit(0)

def emit_markdown(entries):
    lines = [f"# {title}", "", f"_Generated at {time.strftime('%Y-%m-%dT%H:%M:%SZ', time.gmtime())}_", ""]
    for entry in sorted(entries, key=lambda e: e["command"]):
        lines.append(f"## {entry['command']}")
        lines.append("")
        lines.append("```text")
        lines.append(entry["help"])
        lines.append("```")
        lines.append("")
    output_path.parent.mkdir(parents=True, exist_ok=True)
    output_path.write_text("\n".join(lines).rstrip() + "\n", encoding="utf-8")


def emit_json(entries):
    payload = {
        "title": title,
        "generated_at": time.strftime('%Y-%m-%dT%H:%M:%SZ', time.gmtime()),
        "commands": sorted(entries, key=lambda e: e["command"]),
    }
    output_path.parent.mkdir(parents=True, exist_ok=True)
    output_path.write_text(json.dumps(payload, indent=2) + "\n", encoding="utf-8")

if output_format == "markdown":
    emit_markdown(help_entries)
else:
    emit_json(help_entries)

if metrics_file:
    summary = {
        "output": str(output_path),
        "format": output_format,
        "title": title,
        "command_count": len(help_entries),
        "sources": sorted({entry["source"] for entry in help_entries}),
    }
    Path(metrics_file).parent.mkdir(parents=True, exist_ok=True)
    Path(metrics_file).write_text(json.dumps(summary, indent=2) + "\n", encoding="utf-8")
PY

status=$?
if [[ $status -ne 0 ]]; then
  exit $status
fi

log "INFO" "cli_help_generator.sh wrote documentation to $OUTPUT"
exit 0
