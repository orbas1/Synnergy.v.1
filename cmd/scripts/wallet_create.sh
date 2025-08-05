#!/usr/bin/env bash
set -e
out=${1:-wallet.json}
pass=${2:-password}
./synnergy wallet new > "$out"
echo "password: $pass" >> "$out"
