#!/usr/bin/env bash
set -e
# Build the synnergy CLI with trimmed debug paths
GO111MODULE=on go build -trimpath -o synnergy ../synnergy
