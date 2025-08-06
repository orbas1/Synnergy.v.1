#!/usr/bin/env bash
set -e
# Build the synnergy CLI
GO111MODULE=on go build -trimpath -o synnergy ./cmd/synnergy
# Use network configuration and genesis file
export SYN_CONFIG=../configs/network.yaml
# Show configured genesis wallets for operators
./synnergy genesis show
# Stake initial validator, authority and host nodes
./synnergy node stake d22b7fffcb6a72c94713f5c2e2f2142565b8402299842eb4c1039cea3c293ff0 1000000
./synnergy node stake 5baf968c0972eab52046bbca74763d61c923d6e81e549f0f7f0db66f8cfddad9 500000
./synnergy node stake 20d2698d356868e08411ce2ea8ac824d07011aebe788a6458d33e244172b8c34 500000
# Launch network and consensus services
./synnergy network start &
NET_PID=$!
./synnergy consensus-service start 3000 &
CONS_PID=$!
trap 'kill $NET_PID $CONS_PID' INT TERM
wait
