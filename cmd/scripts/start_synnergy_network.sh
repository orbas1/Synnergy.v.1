#!/usr/bin/env bash
set -euo pipefail

CLI=${SYN_CLI:-./synnergy}
if ! command -v go >/dev/null 2>&1; then
  echo "go command not found" >&2
  exit 1
fi

# Compile the synnergy CLI
GO111MODULE=on go build -trimpath -o "$CLI" ../synnergy

# Start core services
"$CLI" network start &
NET_PID=$!
"$CLI" consensus-service start 1000 &
CONS_PID=$!
"$CLI" replication start &
REP_PID=$!
"$CLI" simplevm start &
VM_PID=$!

# Run a sample security command
"$CLI" security merkle "deadbeef,baadf00d"

cleanup() {
  kill "$NET_PID" "$CONS_PID" "$REP_PID" "$VM_PID" 2>/dev/null || true
}
trap cleanup INT TERM

wait
