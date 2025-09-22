#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
# shellcheck source=lib/common.sh
source "$SCRIPT_DIR/lib/common.sh"

usage() {
  cat <<'USAGE'
Automate a DAO proposal lifecycle including creation, voting and execution.

Usage: proposal_lifecycle.sh [--dao DAO_NAME] [--description TEXT]
                             [--creator ADDRESS] [--voter ADDRESS]

Options:
  --dao NAME        DAO name (default: Stage103DAO).
  --description TX  Proposal description (default: Enable stage 103 automation).
  --creator ADDR    Existing address to use as creator (optional).
  --voter ADDR      Existing address to use as voter (optional).
  -h, --help        Show help message.
USAGE
}

DAO_NAME="Stage103DAO"
DESCRIPTION="Enable stage 103 automation"
CREATOR=""
VOTER=""

while [[ $# -gt 0 ]]; do
  case "$1" in
    --dao)
      DAO_NAME=$2
      shift 2
      ;;
    --description)
      DESCRIPTION=$2
      shift 2
      ;;
    --creator)
      CREATOR=$2
      shift 2
      ;;
    --voter)
      VOTER=$2
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      usage
      fail "unknown flag $1"
      ;;
  esac
done

require_commands jq go

sign_message() {
  local message=$1
  local json
  json=$(GO111MODULE=on go run <<'GO' "$message"
package main
import (
  "crypto/ecdsa"
  "crypto/elliptic"
  "crypto/rand"
  "crypto/sha256"
  "encoding/hex"
  "encoding/json"
  "fmt"
  "math/big"
  "os"
)
func pad(b []byte) []byte {
  if len(b) >= 32 {
    return b
  }
  padded := make([]byte, 32)
  copy(padded[32-len(b):], b)
  return padded
}
func main() {
  msg := []byte(os.Args[1])
  priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
  if err != nil { panic(err) }
  digest := sha256.Sum256(msg)
  r, s, err := ecdsa.Sign(rand.Reader, priv, digest[:])
  if err != nil { panic(err) }
  pub := append([]byte{4}, pad(priv.PublicKey.X.Bytes())...)
  pub = append(pub, pad(priv.PublicKey.Y.Bytes())...)
  sig := append(pad(r.Bytes()), pad(s.Bytes())...)
  out := map[string]string{
    "pub": hex.EncodeToString(pub),
    "msg": hex.EncodeToString(msg),
    "sig": hex.EncodeToString(sig),
  }
  enc, _ := json.Marshal(out)
  fmt.Println(string(enc))
}
GO
)
  echo "$json"
}

if [[ -z $CREATOR ]]; then
  creator_wallet=$(run_cli_json wallet new)
  CREATOR=$(json_extract '.address' "$creator_wallet")
fi

if [[ -z $VOTER ]]; then
  voter_wallet=$(run_cli_json wallet new)
  VOTER=$(json_extract '.address' "$voter_wallet")
fi

log_info "Creating DAO $DAO_NAME with creator $CREATOR"
dao_resp=$(run_cli_json dao create "$DAO_NAME" "$CREATOR")
dao_id=$(json_extract '.id' "$dao_resp")
run_cli_json dao join "$dao_id" "$VOTER" >/dev/null

sig_create=$(sign_message "create-$dao_id")
create_pub=$(echo "$sig_create" | jq -r '.pub')
create_msg=$(echo "$sig_create" | jq -r '.msg')
create_sig=$(echo "$sig_create" | jq -r '.sig')

log_info "Submitting proposal"
proposal_out=$(run_cli dao-proposal create "$dao_id" "$CREATOR" "$DESCRIPTION" --pub "$create_pub" --msg "$create_msg" --sig "$create_sig" --json)
proposal_id=$(echo "$proposal_out" | jq -r '.id')

sig_vote=$(sign_message "vote-$proposal_id")
vote_pub=$(echo "$sig_vote" | jq -r '.pub')
vote_msg=$(echo "$sig_vote" | jq -r '.msg')
vote_sig=$(echo "$sig_vote" | jq -r '.sig')

log_info "Casting vote"
run_cli dao-proposal vote "$proposal_id" "$VOTER" "10" yes --pub "$vote_pub" --msg "$vote_msg" --sig "$vote_sig" --json

results=$(run_cli dao-proposal results "$proposal_id" --json)
log_info "Results: $results"

sig_exec=$(sign_message "exec-$proposal_id")
exec_pub=$(echo "$sig_exec" | jq -r '.pub')
exec_msg=$(echo "$sig_exec" | jq -r '.msg')
exec_sig=$(echo "$sig_exec" | jq -r '.sig')

log_info "Executing proposal"
run_cli dao-proposal execute "$proposal_id" "$CREATOR" --pub "$exec_pub" --msg "$exec_msg" --sig "$exec_sig" --json

log_info "Proposal lifecycle complete"
