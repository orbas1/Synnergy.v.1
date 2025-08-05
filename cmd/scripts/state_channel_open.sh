#!/usr/bin/env bash
set -e
if [ $# -lt 2 ]; then
  echo "Usage: $0 <from> <to> [amount]"
  exit 1
fi
amt=${3:-1}
./synnergy state-channel open "$1" "$2" "$amt"
