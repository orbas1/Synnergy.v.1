#!/usr/bin/env bash
set -e
if [ $# -lt 2 ]; then
  echo "Usage: $0 <creator> <recipient> [type] [amount] [desc]"
  exit 1
fi
type=${3:-personal}
amount=${4:-0}
desc=${5:-""}
./synnergy loanpool submit "$1" "$2" "$type" "$amount" "$desc"
