#!/usr/bin/env bash
set -e
shard=${1:-0}
./synnergy sharding leader get "$shard"
