#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)
# shellcheck disable=SC1091
. "$SCRIPT_DIR/lib/runtime.sh"

usage() {
    cat <<'USAGE'
Usage: ai_secure_storage.sh <command> [options]

Commands:
  encrypt   Encrypt an artifact with AES-256-GCM.
  decrypt   Decrypt an artifact previously encrypted by this script.
  sign      Sign a file using an RSA or EC private key.
  verify    Verify a signature with the corresponding public key.

Set SYN_DRY_RUN=1 to preview commands without executing them.
USAGE
}

syn_require_tool openssl
syn_require_tool sha256sum

command=${1:-}
if [[ -z "$command" ]]; then
    usage
    exit 1
fi
shift

get_passphrase() {
    local pass="$1"
    local file="$2"
    if [[ -n "$pass" ]]; then
        echo "$pass"
    elif [[ -n "$file" ]]; then
        syn_require_file "$file"
        cat "$file"
    else
        syn_fail "passphrase is required"
    fi
}

case "$command" in
    encrypt)
        encrypt_usage() {
            cat <<'USAGE'
Usage: ai_secure_storage.sh encrypt --input FILE --output FILE [options]

Options:
  --input FILE             Source file to encrypt.
  --output FILE            Destination encrypted file.
  --passphrase VALUE       Passphrase used for key derivation.
  --passphrase-file PATH   File containing passphrase.
  --metadata FILE          Additional JSON metadata to embed alongside ciphertext.
USAGE
        }

        input=""
        output=""
        passphrase=""
        passphrase_file=""
        metadata=""

        while [[ $# -gt 0 ]]; do
            case "$1" in
                --input)
                    input="$2"
                    shift 2
                    ;;
                --output)
                    output="$2"
                    shift 2
                    ;;
                --passphrase)
                    passphrase="$2"
                    shift 2
                    ;;
                --passphrase-file)
                    passphrase_file="$2"
                    shift 2
                    ;;
                --metadata)
                    metadata="$2"
                    shift 2
                    ;;
                -h|--help)
                    encrypt_usage
                    exit 0
                    ;;
                *)
                    syn_fail "unknown option: $1"
                    ;;
            esac
        done

        [[ -z "$input" || -z "$output" ]] && syn_fail "--input and --output are required"
        syn_require_file "$input"
        [[ -n "$metadata" ]] && syn_require_file "$metadata"

        pass=$(get_passphrase "$passphrase" "$passphrase_file")

        if [[ $SYN_DRY_RUN -eq 1 ]]; then
            syn_info "[dry-run] would encrypt $input to $output"
            exit 0
        fi

        tmp_output=$(syn_expand_path "$output")
        mkdir -p "$(dirname "$tmp_output")"
        pass_file=$(mktemp)
        printf '%s' "$pass" >"$pass_file"
        cleanup() {
            rm -f "$pass_file"
        }
        trap cleanup EXIT
        syn_run "Encrypting artifact" openssl enc -aes-256-gcm -pbkdf2 -iter 200000 -salt -in "$input" -out "$tmp_output" -pass file:"$pass_file"
        trap - EXIT
        cleanup

        checksum=$(sha256sum "$tmp_output" | awk '{print $1}')
        meta_payload=$(python3 - "$input" "$checksum" "$metadata" <<'PY'
import json
import sys
from pathlib import Path

source = Path(sys.argv[1])
checksum = sys.argv[2]
meta_path = sys.argv[3]
meta = {}
if meta_path:
    meta = json.loads(Path(meta_path).read_text(encoding='utf-8'))

print(json.dumps({
    "source": str(source),
    "checksum": checksum,
    "metadata": meta,
}))
PY
)
        echo "$meta_payload" >"$tmp_output.meta.json"
        echo "$meta_payload"
        ;;

    decrypt)
        decrypt_usage() {
            cat <<'USAGE'
Usage: ai_secure_storage.sh decrypt --input FILE --output FILE [options]

Options:
  --input FILE             Encrypted input file.
  --output FILE            Destination plaintext file.
  --passphrase VALUE       Passphrase used for key derivation.
  --passphrase-file PATH   File containing passphrase.
USAGE
        }

        input=""
        output=""
        passphrase=""
        passphrase_file=""

        while [[ $# -gt 0 ]]; do
            case "$1" in
                --input)
                    input="$2"
                    shift 2
                    ;;
                --output)
                    output="$2"
                    shift 2
                    ;;
                --passphrase)
                    passphrase="$2"
                    shift 2
                    ;;
                --passphrase-file)
                    passphrase_file="$2"
                    shift 2
                    ;;
                -h|--help)
                    decrypt_usage
                    exit 0
                    ;;
                *)
                    syn_fail "unknown option: $1"
                    ;;
            esac
        done

        [[ -z "$input" || -z "$output" ]] && syn_fail "--input and --output are required"
        syn_require_file "$input"

        pass=$(get_passphrase "$passphrase" "$passphrase_file")

        if [[ $SYN_DRY_RUN -eq 1 ]]; then
            syn_info "[dry-run] would decrypt $input to $output"
            exit 0
        fi

        out_path=$(syn_expand_path "$output")
        mkdir -p "$(dirname "$out_path")"
        pass_file=$(mktemp)
        printf '%s' "$pass" >"$pass_file"
        cleanup() {
            rm -f "$pass_file"
        }
        trap cleanup EXIT
        syn_run "Decrypting artifact" openssl enc -d -aes-256-gcm -pbkdf2 -iter 200000 -salt -in "$input" -out "$out_path" -pass file:"$pass_file"
        trap - EXIT
        cleanup
        echo "{\"status\":\"decrypted\"}"
        ;;

    sign)
        sign_usage() {
            cat <<'USAGE'
Usage: ai_secure_storage.sh sign --input FILE --key PRIVATE_KEY --signature FILE
USAGE
        }

        input=""
        key=""
        signature=""

        while [[ $# -gt 0 ]]; do
            case "$1" in
                --input)
                    input="$2"
                    shift 2
                    ;;
                --key)
                    key="$2"
                    shift 2
                    ;;
                --signature)
                    signature="$2"
                    shift 2
                    ;;
                -h|--help)
                    sign_usage
                    exit 0
                    ;;
                *)
                    syn_fail "unknown option: $1"
                    ;;
            esac
        done

        [[ -z "$input" || -z "$key" || -z "$signature" ]] && syn_fail "--input, --key and --signature are required"
        syn_require_file "$input"
        syn_require_file "$key"

        if [[ $SYN_DRY_RUN -eq 1 ]]; then
            syn_info "[dry-run] would sign $input"
            exit 0
        fi

        sig_path=$(syn_expand_path "$signature")
        mkdir -p "$(dirname "$sig_path")"
        syn_run "Signing artifact" openssl dgst -sha256 -sign "$key" -out "$sig_path" "$input"
        echo "{\"status\":\"signed\",\"signature\":\"$sig_path\"}"
        ;;

    verify)
        verify_usage() {
            cat <<'USAGE'
Usage: ai_secure_storage.sh verify --input FILE --signature FILE --pubkey FILE
USAGE
        }

        input=""
        signature=""
        pubkey=""

        while [[ $# -gt 0 ]]; do
            case "$1" in
                --input)
                    input="$2"
                    shift 2
                    ;;
                --signature)
                    signature="$2"
                    shift 2
                    ;;
                --pubkey)
                    pubkey="$2"
                    shift 2
                    ;;
                -h|--help)
                    verify_usage
                    exit 0
                    ;;
                *)
                    syn_fail "unknown option: $1"
                    ;;
            esac
        done

        [[ -z "$input" || -z "$signature" || -z "$pubkey" ]] && syn_fail "--input, --signature and --pubkey are required"
        syn_require_file "$input"
        syn_require_file "$signature"
        syn_require_file "$pubkey"

        if [[ $SYN_DRY_RUN -eq 1 ]]; then
            syn_info "[dry-run] would verify signature"
            exit 0
        fi

        if openssl dgst -sha256 -verify "$pubkey" -signature "$signature" "$input" >/dev/null 2>&1; then
            echo '{"status":"verified"}'
        else
            echo '{"status":"invalid"}'
            exit 1
        fi
        ;;

    *)
        usage
        exit 1
        ;;
esac
