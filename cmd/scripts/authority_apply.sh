#!/usr/bin/env bash
set -e
if [ -z "$1" ]; then
  echo "Usage: $0 <address> [role]"
  exit 1
fi
role=${2:-validator}
./synnergy authority_apply submit "$1" "$role" "auto-generated"
