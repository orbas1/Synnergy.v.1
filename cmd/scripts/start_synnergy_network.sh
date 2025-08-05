#!/usr/bin/env bash
set -e
# Compile the synnergy CLI
GO111MODULE=on go build -trimpath -o synnergy ../synnergy
# Start core services
./synnergy network start &
NET_PID=$!
./synnergy consensus-service start 1000 &
CONS_PID=$!
./synnergy replication start &
REP_PID=$!
./synnergy simplevm start &
VM_PID=$!
# Run a sample security command
./synnergy security merkle "deadbeef,baadf00d"
# Wait for Ctrl+C
trap 'kill $NET_PID $CONS_PID $REP_PID $VM_PID' INT TERM
wait
