#!/usr/bin/env bash
set -e
if [ $# -lt 3 ]; then
  echo "Usage: $0 <token> <from> <to> [amount]"
  exit 1
fi
token=$1
from=$2
to=$3
amt=${4:-1}
./synnergy "$token" transfer "$from" "$to" "$amt"
