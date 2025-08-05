#!/usr/bin/env bash
set -e
title=${1:-Proposal}
body=${2:-proposal.md}
./synnergy governance propose "$title" "$body"
