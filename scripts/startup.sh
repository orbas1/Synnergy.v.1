#!/usr/bin/env bash
set -e
SCRIPT_DIR="$(dirname "$0")/../cmd/scripts"
"$SCRIPT_DIR/start_synnergy_network.sh" "$@"
