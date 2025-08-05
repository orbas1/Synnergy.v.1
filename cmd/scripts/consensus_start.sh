#!/usr/bin/env bash
set -e
interval=${1:-1000}
./synnergy consensus-service start "$interval"
