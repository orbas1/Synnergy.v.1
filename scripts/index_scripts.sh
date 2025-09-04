#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

for script in "$SCRIPT_DIR"/*.sh; do
  if [ "$(basename "$script")" != "index_scripts.sh" ]; then
    bash "$script"
  fi
done
