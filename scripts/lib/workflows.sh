#!/usr/bin/env bash
# Workflow orchestration helpers for Synnergy operational scripts.

if [[ "${BASH_SOURCE[0]}" == "$0" ]]; then
  echo "workflows.sh must be sourced" >&2
  exit 1
fi

if [[ -n "${__SYN_WORKFLOWS_SH:-}" ]]; then
  return
fi
__SYN_WORKFLOWS_SH=1

: "${SCRIPTS_ROOT:?SCRIPTS_ROOT must be defined before sourcing workflows.sh}"

STATE_DIR="${SYN_WORKFLOW_STATE_DIR:-$SCRIPTS_ROOT/state}"
mkdir -p "$STATE_DIR"

# Internal state shared between the launcher and the builders.
WF_ARG_DELIM=$'\x1f'
WF_STEP_TYPE=()
WF_STEP_DESC=()
WF_STEP_DATA1=()
WF_STEP_DATA2=()
WF_STEP_STATUS=()
WF_STEP_MESSAGE=()

WF_CATEGORY=""
WF_ACTION=""
WF_SUBJECT=""
WF_SUBJECT_TITLE=""
WF_DESCRIPTION=""
WF_PARAM_HELP=""

WORKFLOW_PLAN_MODE=${WORKFLOW_PLAN_MODE:-false}
WORKFLOW_OUTPUT="${WORKFLOW_OUTPUT:-}"  # May be set by the launcher.
WORKFLOW_NOTES=(${WORKFLOW_NOTES[@]:-})

# The launcher provides WORKFLOW_PARAMS as a global associative array.
if ! declare -p WORKFLOW_PARAMS &>/dev/null; then
  declare -gA WORKFLOW_PARAMS=()
fi

reset_workflow_steps() {
  WF_STEP_TYPE=()
  WF_STEP_DESC=()
  WF_STEP_DATA1=()
  WF_STEP_DATA2=()
  WF_STEP_STATUS=()
  WF_STEP_MESSAGE=()
}

add_step() {
  local type="$1"
  local desc="$2"
  local data1="${3:-}"
  local data2="${4:-}"
  WF_STEP_TYPE+=("$type")
  WF_STEP_DESC+=("$desc")
  WF_STEP_DATA1+=("$data1")
  WF_STEP_DATA2+=("$data2")
  WF_STEP_STATUS+=("pending")
  WF_STEP_MESSAGE+=("")
}

pack_args() {
  local result=""
  local arg
  for arg in "$@"; do
    result+="$arg$WF_ARG_DELIM"
  done
  printf '%s' "$result"
}

unpack_args() {
  local data="$1"
  local trimmed="${data%$WF_ARG_DELIM}"
  local IFS="$WF_ARG_DELIM"
  # shellcheck disable=SC2206
  WF_UNPACKED_ARGS=($trimmed)
}

slugify() {
  local input="$1"
  input="${input// /-}"
  input="${input//_/ -}"
  input="${input//-/-}"
  printf '%s' "${input// /-}"
}

title_case() {
  local input="$1"
  local IFS='_ '
  read -r -a parts <<<"$input"
  local out=()
  local part
  for part in "${parts[@]}"; do
    [[ -z "$part" ]] && continue
    out+=("${part^}")
  done
  printf '%s' "${out[*]}"
}

params_get() {
  local key="$1"
  local default="$2"
  if [[ -n "${WORKFLOW_PARAMS[$key]:-}" ]]; then
    printf '%s' "${WORKFLOW_PARAMS[$key]}"
  else
    printf '%s' "$default"
  fi
}

is_plan_mode() {
  if [[ "$WORKFLOW_PLAN_MODE" == true ]]; then
    return 0
  fi
  if [[ "${DRY_RUN:-false}" == true ]]; then
    return 0
  fi
  return 1
}

maybe_plan_status() {
  if is_plan_mode; then
    printf 'planned'
  else
    printf 'ok'
  fi
}

append_note() {
  local message="$1"
  WORKFLOW_NOTES+=("$message")
}

workflow_subject_dir() {
  local slug
  if [[ -n "$WF_SUBJECT" ]]; then
    slug="${WF_SUBJECT// /-}"
  else
    slug="$WF_CATEGORY"
  fi
  printf '%s/%s' "$STATE_DIR" "$slug"
}

workflow_subject_title() {
  if [[ -n "$WF_SUBJECT_TITLE" ]]; then
    printf '%s' "$WF_SUBJECT_TITLE"
  elif [[ -n "$WF_SUBJECT" ]]; then
    title_case "$WF_SUBJECT"
  else
    title_case "$WF_CATEGORY"
  fi
}

syn_cli_executable() {
  if command -v synnergy >/dev/null 2>&1; then
    return 0
  fi
  local compiled="$PROJECT_ROOT/bin/synnergy"
  if [[ -x "$compiled" ]]; then
    return 0
  fi
  return 1
}

record_step_status() {
  local idx="$1"
  local status="$2"
  local message="$3"
  WF_STEP_STATUS[$idx]="$status"
  WF_STEP_MESSAGE[$idx]="$message"
}

subject_slug() {
  if [[ -z "$WF_SUBJECT" ]]; then
    printf '%s' "$WF_CATEGORY"
  else
    local input="${WF_SUBJECT//_/ }"
    input="${input// /-}"
    printf '%s' "$input"
  fi
}

subject_safe_name() {
  local slug
  slug="$(subject_slug)"
  printf '%s' "${slug//[^a-zA-Z0-9_-]/-}"
}

ensure_parent_dir() {
  local file="$1"
  local dir
  dir="$(dirname "$file")"
  mkdir -p "$dir"
}

determine_workflow_context() {
  local base="$1"
  WF_SUBJECT=""
  WF_SUBJECT_TITLE=""

  case "$base" in
    *_node_setup)
      WF_CATEGORY="node"
      WF_ACTION="setup"
      WF_SUBJECT="${base%_node_setup}"
      ;;
    cross_chain_*)
      WF_CATEGORY="cross_chain"
      WF_ACTION="${base#cross_chain_}"
      ;;
    cross_consensus_network)
      WF_CATEGORY="consensus"
      WF_ACTION="cross_network"
      ;;
    consensus_*)
      WF_CATEGORY="consensus"
      WF_ACTION="${base#consensus_}"
      ;;
    dynamic_consensus_hopping)
      WF_CATEGORY="consensus"
      WF_ACTION="dynamic_hopping"
      ;;
    bridge_fallback_recovery)
      WF_CATEGORY="recovery"
      WF_ACTION="bridge_fallback"
      ;;
    bridge_verification)
      WF_CATEGORY="cross_chain"
      WF_ACTION="bridge_verification"
      ;;
    wallet_*)
      WF_CATEGORY="wallet"
      WF_ACTION="${base#wallet_}"
      ;;
    metrics_*)
      WF_CATEGORY="monitoring"
      WF_ACTION="${base#metrics_}"
      ;;
    system_health_logging|logs_collect|alerting_setup|anomaly_detection)
      WF_CATEGORY="monitoring"
      WF_ACTION="$base"
      ;;
    multi_factor_setup|firewall_setup|secure_node_hardening|secure_store_setup|key_backup|key_rotation_schedule|zero_trust_data_channels|tamper_alert|pki_setup|certificate_issue)
      WF_CATEGORY="security"
      WF_ACTION="$base"
      ;;
    aml_kyc_process|identity_verification|credential_revocation|idwallet_register|biometric_enroll|biometric_verify|private_transactions|regulatory_report)
      WF_CATEGORY="compliance"
      WF_ACTION="$base"
      ;;
    dao_* )
      WF_CATEGORY="dao"
      WF_ACTION="${base#dao_}"
      ;;
    proposal_lifecycle|stake_penalty)
      WF_CATEGORY="dao"
      WF_ACTION="$base"
      ;;
    treasury_* )
      WF_CATEGORY="treasury"
      WF_ACTION="${base#treasury_}"
      ;;
    grant_distribution|grant_reporting|financial_prediction)
      WF_CATEGORY="treasury"
      WF_ACTION="$base"
      ;;
    data_distribution|generate_mock_data|merkle_proof_generator)
      WF_CATEGORY="data"
      WF_ACTION="$base"
      ;;
    holographic_storage)
      WF_CATEGORY="storage"
      WF_ACTION="$base"
      ;;
    disaster_recovery_backup|restore_disaster_recovery|restore_ledger)
      WF_CATEGORY="recovery"
      WF_ACTION="$base"
      ;;
    block_integrity_check|immutable_audit_log_export|immutable_audit_verify|immutable_log_snapshot|immutability_verifier|forensic_data_export)
      WF_CATEGORY="audit"
      WF_ACTION="$base"
      ;;
    ansible_deploy|helm_deploy|k8s_deploy|terraform_apply|cd_deploy)
      WF_CATEGORY="deployment"
      WF_ACTION="$base"
      ;;
    configure_environment|format_code|install_dependencies|update_dependencies|cli_tooling|artifact_checksum|script_completion_setup|script_launcher)
      WF_CATEGORY="devops"
      WF_ACTION="$base"
      ;;
    contract_coverage_report|contract_language_compatibility_test|contract_static_analysis|contract_test_suite|integration_test_suite|fuzz_testing|performance_regression|stress_test_network|benchmarks|gui_wallet_test|release_sign_verify|scripts_test)
      WF_CATEGORY="testing"
      WF_ACTION="$base"
      ;;
    network_diagnostics|network_harness|network_migration|network_partition_test|shutdown_network)
      WF_CATEGORY="network"
      WF_ACTION="$base"
      ;;
    virtual_machine|vm_sandbox_management)
      WF_CATEGORY="virtualization"
      WF_ACTION="$base"
      ;;
    ai_training)
      WF_CATEGORY="ai"
      WF_ACTION="$base"
      ;;
    tutorial_scripts)
      WF_CATEGORY="documentation"
      WF_ACTION="$base"
      ;;
    faq_autoresolve)
      WF_CATEGORY="support"
      WF_ACTION="$base"
      ;;
    multi_node_cluster_setup)
      WF_CATEGORY="cluster"
      WF_ACTION="setup"
      WF_SUBJECT="multi_node_cluster"
      ;;
    high_availability_setup|ha_failover_test|ha_immutable_verification)
      WF_CATEGORY="high_availability"
      WF_ACTION="$base"
      ;;
    mint_nft|upgrade_contract)
      WF_CATEGORY="contract"
      WF_ACTION="$base"
      ;;
    mining_node_setup|mobile_mining_node_setup)
      WF_CATEGORY="node"
      WF_ACTION="setup"
      WF_SUBJECT="${base%_node_setup}"
      ;;
    optimization_node_setup)
      WF_CATEGORY="node"
      WF_ACTION="setup"
      WF_SUBJECT="optimization"
      ;;
    content_node_setup)
      WF_CATEGORY="node"
      WF_ACTION="setup"
      WF_SUBJECT="content"
      ;;
    *)
      WF_CATEGORY="misc"
      WF_ACTION="$base"
      ;;
  esac

  if [[ -z "$WF_SUBJECT" && $WF_CATEGORY == node ]]; then
    WF_SUBJECT="${base%_node_setup}"
  fi

  if [[ -z "$WF_CATEGORY" ]]; then
    WF_CATEGORY="misc"
    WF_ACTION="$base"
  fi

  WF_SUBJECT_TITLE="$(workflow_subject_title)"
}

describe_workflow() {
  local base="$1"
  local subject_title="$(workflow_subject_title)"
  case "$WF_CATEGORY" in
    node)
      WF_DESCRIPTION="Provision and activate a ${subject_title} node using declarative manifests and the Synnergy CLI."
      WF_PARAM_HELP=$'  --set network=ID        Target network identifier (default: synnergy-devnet)\n  --set region=NAME       Deployment region hint (default: primary)\n  --set cpu=CORES        Requested CPU allocation (default: 2)\n  --set memory=SIZE      Requested memory allocation (default: 4Gi)\n  --set replicas=COUNT   Desired replica count (default: 1)'
      ;;
    consensus)
      case "$WF_ACTION" in
        start)
          WF_DESCRIPTION="Start the consensus engine, enrol validators and verify finality health."
          WF_PARAM_HELP=$'  --set profile=NAME     Tuning profile such as balanced or aggressive\n  --set quorum=VALUE    Minimum quorum threshold (default: 0.67)'
          ;;
        finality_check)
          WF_DESCRIPTION="Execute a finality audit against the active consensus ledger."
          WF_PARAM_HELP=$'  --set window=BLOCKS    Number of recent blocks to sample (default: 128)'
          ;;
        adaptive_manage)
          WF_DESCRIPTION="Adjust consensus demand and stake weighting according to adaptive inputs."
          WF_PARAM_HELP=$'  --set demand=VALUE     Target demand ratio (default: 0.55)\n  --set stake=VALUE      Target stake ratio (default: 0.45)\n  --set window=SIZE      Moving window used for smoothing (default: 20)'
          ;;
        difficulty_adjust)
          WF_DESCRIPTION="Tune block production difficulty based on the supplied target and variance."
          WF_PARAM_HELP=$'  --set delta=VALUE      Difficulty delta to apply (default: 0.02)\n  --set target=VALUE     Target block time variance (default: 0.65)'
          ;;
        recovery)
          WF_DESCRIPTION="Recover consensus state from a trusted snapshot and restart validators."
          WF_PARAM_HELP=$'  --set snapshot=PATH    Snapshot archive to restore (default: state/snapshots/latest.tar.gz)\n  --set integrity=HASH   Expected checksum of the snapshot (optional)'
          ;;
        specific_node)
          WF_DESCRIPTION="Inspect and manage a specific consensus participant."
          WF_PARAM_HELP=$'  --set node=ID          Validator identifier to operate on (default: validator-1)\n  --set action=NAME      Node action such as promote or drain (default: inspect)'
          ;;
        validator_manage)
          WF_DESCRIPTION="Rotate or scale the consensus validator set following policy constraints."
          WF_PARAM_HELP=$'  --set mode=ACTION      Validator management mode: rotate, expand, contract (default: rotate)\n  --set count=NUMBER     Number of validators to affect (default: 2)'
          ;;
        dynamic_hopping)
          WF_DESCRIPTION="Coordinate dynamic consensus hopping between shards for resiliency testing."
          WF_PARAM_HELP=$'  --set interval=SECONDS Migration interval between shards (default: 45)\n  --set shards=COUNT     Number of shards to cycle through (default: 3)'
          ;;
        cross_network)
          WF_DESCRIPTION="Synchronise consensus checkpoints across heterogeneous networks."
          WF_PARAM_HELP=$'  --set source=NETWORK   Source network identifier (default: devnet)\n  --set target=NETWORK   Target network identifier (default: auditnet)'
          ;;
        *)
          WF_DESCRIPTION="Execute a consensus workflow (${WF_ACTION})."
          WF_PARAM_HELP=""
          ;;
      esac
      ;;
    cross_chain)
      case "$WF_ACTION" in
        agnostic_protocols)
          WF_DESCRIPTION="Apply cross-chain agnostic protocol manifests and verify compatibility."
          WF_PARAM_HELP=$'  --set manifest=PATH    Protocol manifest file (default: configs/cross-chain/agnostic.yaml)'
          ;;
        contracts_deploy)
          WF_DESCRIPTION="Deploy cross-chain smart contracts and register the resulting endpoints."
          WF_PARAM_HELP=$'  --set manifest=PATH    Deployment manifest (default: configs/cross-chain/contracts.yaml)'
          ;;
        transactions)
          WF_DESCRIPTION="Simulate cross-chain transaction relays and collect latency metrics."
          WF_PARAM_HELP=$'  --set batch=COUNT      Number of transactions per batch (default: 10)\n  --set channel=NAME     Logical channel identifier (default: cross-link)'
          ;;
        setup)
          WF_DESCRIPTION="Initialise cross-chain connectivity, trust roots and channel registries."
          WF_PARAM_HELP=$'  --set topology=NAME    Connectivity topology profile (default: mesh)'
          ;;
        connection)
          WF_DESCRIPTION="Validate a cross-chain connection handshake and capture diagnostics."
          WF_PARAM_HELP=$'  --set endpoint=URL     Remote endpoint to validate (default: https://example)'
          ;;
        bridge)
          WF_DESCRIPTION="Provision the Synnergy reference bridge and publish its manifest."
          WF_PARAM_HELP=$'  --set name=NAME        Bridge name (default: synnergy-bridge)'
          ;;
        bridge_verification)
          WF_DESCRIPTION="Perform attestation checks against active bridge channels."
          WF_PARAM_HELP=$'  --set channel=NAME     Bridge channel to verify (default: channel-0)'
          ;;
        *)
          WF_DESCRIPTION="Run cross-chain workflow ${WF_ACTION}."
          WF_PARAM_HELP=""
          ;;
      esac
      ;;
    wallet)
      case "$WF_ACTION" in
        hardware_integration)
          WF_DESCRIPTION="Provision hardware wallet integration with deterministic metadata exports."
          WF_PARAM_HELP=$'  --set device=NAME      Hardware device label (default: ledger-nano)\n  --set operator=ID      Operator identifier recorded in the manifest (default: ops-team)'
          ;;
        init)
          WF_DESCRIPTION="Create a baseline wallet, encrypt keys and record metadata."
          WF_PARAM_HELP=$'  --set password=VALUE   Wallet encryption password (default: synnergy-dev)\n  --set dir=PATH         Directory for wallet artifacts (default: var/wallets)'
          ;;
        key_rotation)
          WF_DESCRIPTION="Prepare and document a wallet key rotation event."
          WF_PARAM_HELP=$'  --set next-key=HEX     Predetermined replacement key (optional)\n  --set effective=DATE   ISO-8601 activation date (default: +24h)'
          ;;
        multisig_setup)
          WF_DESCRIPTION="Generate multisignature configuration files and coordinator plans."
          WF_PARAM_HELP=$'  --set participants=N    Number of multisig participants (default: 3)\n  --set threshold=M      Required approvals to finalise (default: 2)'
          ;;
        offline_sign)
          WF_DESCRIPTION="Produce an offline signing manifest for an unsigned transaction payload."
          WF_PARAM_HELP=$'  --set tx=PATH          Unsigned transaction file (default: pending.tx.json)\n  --set output=PATH      Destination for the signature (default: artifacts/offline.sig)'
          ;;
        server_setup)
          WF_DESCRIPTION="Configure the wallet server with TLS credentials and archival policies."
          WF_PARAM_HELP=$'  --set host=HOST        Listen address (default: 127.0.0.1)\n  --set port=PORT        Listen port (default: 8546)'
          ;;
        transfer)
          WF_DESCRIPTION="Prepare a wallet transfer request and capture the resulting manifest."
          WF_PARAM_HELP=$'  --set from=ADDR        Sender wallet address (default: wallet-1)\n  --set to=ADDR          Recipient wallet address (default: wallet-2)\n  --set amount=VALUE     Amount to transfer (default: 10)'
          ;;
        *)
          WF_DESCRIPTION="Execute wallet workflow ${WF_ACTION}."
          WF_PARAM_HELP=""
          ;;
      esac
      ;;
    monitoring)
      case "$WF_ACTION" in
        alerting_setup)
          WF_DESCRIPTION="Bootstrap alerting destinations and severity policies."
          WF_PARAM_HELP=$'  --set webhook=URL      Alert webhook endpoint (default: https://hooks.example)\n  --set severity=LEVEL   Minimum severity to forward (default: warning)'
          ;;
        logs_collect)
          WF_DESCRIPTION="Collect node logs into a structured bundle for incident review."
          WF_PARAM_HELP=$'  --set duration=MIN     Minutes of history to gather (default: 30)\n  --set target=DIR       Output directory (default: logs/bundles)'
          ;;
        metrics_alert_dispatch)
          WF_DESCRIPTION="Dispatch synthetic metrics alerts to validate routing and escalation."
          WF_PARAM_HELP=$'  --set rule=NAME        Rule identifier (default: high-cpu)\n  --set repeat=COUNT     Number of dispatch iterations (default: 1)'
          ;;
        metrics_export)
          WF_DESCRIPTION="Generate a metrics snapshot and expose it as JSON for dashboards."
          WF_PARAM_HELP=$'  --set sink=PATH        Output path for metrics JSON (default: metrics/export.json)'
          ;;
        system_health_logging)
          WF_DESCRIPTION="Capture system health indicators and persist them alongside metadata."
          WF_PARAM_HELP=$'  --set interval=SEC     Polling interval (default: 60)\n  --set samples=COUNT    Number of samples to collect (default: 5)'
          ;;
        anomaly_detection)
          WF_DESCRIPTION="Simulate the anomaly detection pipeline and persist evaluation metrics."
          WF_PARAM_HELP=$'  --set window=SIZE      Sliding window size (default: 50)\n  --set sensitivity=VAL  Detection sensitivity (default: 0.85)'
          ;;
        *)
          WF_DESCRIPTION="Execute monitoring workflow ${WF_ACTION}."
          WF_PARAM_HELP=""
          ;;
      esac
      ;;
    security)
      case "$WF_ACTION" in
        multi_factor_setup)
          WF_DESCRIPTION="Generate seed secrets and enrol operators for multi-factor authentication."
          WF_PARAM_HELP=$'  --set operator=ID      Operator identifier (default: ops-team)\n  --set issuer=NAME      Issuer label for OTP seeds (default: Synnergy)'
          ;;
        firewall_setup)
          WF_DESCRIPTION="Render firewall policies and record the resulting rule set."
          WF_PARAM_HELP=$'  --set profile=NAME     Policy profile such as default or hardened (default: hardened)'
          ;;
        secure_node_hardening)
          WF_DESCRIPTION="Apply hardening recommendations to validator nodes and track compliance."
          WF_PARAM_HELP=$'  --set checklist=PATH   Custom checklist path (optional)'
          ;;
        secure_store_setup)
          WF_DESCRIPTION="Initialise an encrypted secure store for custodial secrets."
          WF_PARAM_HELP=$'  --set store=PATH       Secure store path (default: var/secure-store)\n  --set password=VALUE   Encryption password (default: changeit)'
          ;;
        key_backup)
          WF_DESCRIPTION="Create encrypted key backups with audit-friendly metadata."
          WF_PARAM_HELP=$'  --set destination=PATH Backup archive destination (default: backups/keys.tar.gz)'
          ;;
        key_rotation_schedule)
          WF_DESCRIPTION="Document a key rotation schedule and export it as JSON."
          WF_PARAM_HELP=$'  --set cadence=DAYS     Rotation cadence in days (default: 30)'
          ;;
        zero_trust_data_channels)
          WF_DESCRIPTION="Configure zero-trust data channels with policy manifests."
          WF_PARAM_HELP=$'  --set policy=PATH      Policy manifest path (default: configs/security/zt-data.yaml)'
          ;;
        tamper_alert)
          WF_DESCRIPTION="Define tamper alert hooks and notification endpoints."
          WF_PARAM_HELP=$'  --set channel=NAME     Alert channel (default: security-ops)'
          ;;
        pki_setup)
          WF_DESCRIPTION="Bootstrap a certificate authority and issue initial node certificates."
          WF_PARAM_HELP=$'  --set ca-name=NAME     Certificate authority common name (default: Synnergy Root CA)'
          ;;
        certificate_issue)
          WF_DESCRIPTION="Issue a certificate from the configured PKI and record artefacts."
          WF_PARAM_HELP=$'  --set subject=NAME     Certificate subject common name (default: syn-node)\n  --set san=DNS[,DNS]   Subject alternative names (optional)'
          ;;
        *)
          WF_DESCRIPTION="Execute security workflow ${WF_ACTION}."
          WF_PARAM_HELP=""
          ;;
      esac
      ;;
    compliance)
      case "$WF_ACTION" in
        aml_kyc_process)
          WF_DESCRIPTION="Simulate AML/KYC checks and record their audit trace."
          WF_PARAM_HELP=$'  --set applicant=ID     Applicant identifier (default: applicant-001)'
          ;;
        identity_verification)
          WF_DESCRIPTION="Perform identity verification and archive signed attestations."
          WF_PARAM_HELP=$'  --set subject=ID       Identity subject identifier (default: user-001)'
          ;;
        credential_revocation)
          WF_DESCRIPTION="Revoke issued credentials and publish revocation lists."
          WF_PARAM_HELP=$'  --set credential=ID    Credential identifier (default: cred-001)'
          ;;
        idwallet_register)
          WF_DESCRIPTION="Register a new ID wallet with initial trust anchors."
          WF_PARAM_HELP=$'  --set holder=ID        Holder identifier (default: holder-001)'
          ;;
        biometric_enroll)
          WF_DESCRIPTION="Enroll biometric templates and produce integrity hashes."
          WF_PARAM_HELP=$'  --set template=PATH    Biometric template source (default: data/biometric.bin)'
          ;;
        biometric_verify)
          WF_DESCRIPTION="Verify biometric templates against stored references."
          WF_PARAM_HELP=$'  --set template=PATH    Verification template path (default: data/verify.bin)'
          ;;
        private_transactions)
          WF_DESCRIPTION="Run privacy-preserving transaction workflows and capture transcripts."
          WF_PARAM_HELP=$'  --set circuit=NAME     Zero-knowledge circuit name (default: zk-transfer)'
          ;;
        regulatory_report)
          WF_DESCRIPTION="Compile a regulatory compliance report with traceable metadata."
          WF_PARAM_HELP=$'  --set period=RANGE     Reporting period descriptor (default: Q1)'
          ;;
        *)
          WF_DESCRIPTION="Execute compliance workflow ${WF_ACTION}."
          WF_PARAM_HELP=""
          ;;
      esac
      ;;
    dao)
      case "$WF_ACTION" in
        init)
          WF_DESCRIPTION="Bootstrap DAO governance parameters and treasury configuration."
          WF_PARAM_HELP=$'  --set name=NAME        DAO name (default: Synnergy DAO)\n  --set quorum=VALUE     Minimum quorum threshold (default: 0.6)'
          ;;
        proposal_submit)
          WF_DESCRIPTION="Submit a DAO proposal manifest and archive submission receipts."
          WF_PARAM_HELP=$'  --set file=PATH        Proposal manifest (default: proposals/example.yaml)'
          ;;
        proposal_lifecycle)
          WF_DESCRIPTION="Document the lifecycle of a DAO proposal from draft to execution."
          WF_PARAM_HELP=$'  --set proposal=ID      Proposal identifier (default: prop-001)'
          ;;
        token_manage)
          WF_DESCRIPTION="Manage DAO token supply and governance weights."
          WF_PARAM_HELP=$'  --set mint=AMOUNT      Tokens to mint (optional)\n  --set burn=AMOUNT      Tokens to burn (optional)'
          ;;
        vote)
          WF_DESCRIPTION="Record DAO votes and calculate tallies using the CLI."
          WF_PARAM_HELP=$'  --set proposal=ID      Proposal identifier (default: prop-001)\n  --set weight=VALUE     Voting weight (default: 1)'
          ;;
        offchain_vote_tally)
          WF_DESCRIPTION="Aggregate off-chain votes and reconcile with on-chain state."
          WF_PARAM_HELP=$'  --set source=PATH      Off-chain vote export (default: data/offchain_votes.json)'
          ;;
        stake_penalty)
          WF_DESCRIPTION="Apply stake penalties for underperforming validators per DAO policy."
          WF_PARAM_HELP=$'  --set validator=ID     Validator identifier (default: validator-1)\n  --set penalty=VALUE    Penalty ratio (default: 0.05)'
          ;;
        *)
          WF_DESCRIPTION="Execute DAO workflow ${WF_ACTION}."
          WF_PARAM_HELP=""
          ;;
      esac
      ;;
    treasury)
      case "$WF_ACTION" in
        manage)
          WF_DESCRIPTION="Reconcile treasury balances and produce an allocation manifest."
          WF_PARAM_HELP=$'  --set target=AMOUNT    Target reserve balance (default: 100000)\n  --set currency=CODE    Treasury currency code (default: SYN)'
          ;;
        investment_sh)
          WF_DESCRIPTION="Simulate treasury investment placements with scenario modelling."
          WF_PARAM_HELP=$'  --set strategy=NAME    Strategy name (default: conservative)\n  --set allocation=VALUE Allocation percentage (default: 0.2)'
          ;;
        grant_distribution)
          WF_DESCRIPTION="Authorise and document grant distributions from the treasury."
          WF_PARAM_HELP=$'  --set programme=NAME   Grant programme identifier (default: ecosystem)'
          ;;
        grant_reporting)
          WF_DESCRIPTION="Compile grant performance reports with milestone tracking."
          WF_PARAM_HELP=$'  --set period=RANGE     Reporting period (default: Q1)'
          ;;
        financial_prediction)
          WF_DESCRIPTION="Run financial prediction models and store the resulting projections."
          WF_PARAM_HELP=$'  --set horizon=MONTHS   Forecast horizon in months (default: 6)\n  --set model=NAME       Prediction model (default: arima)'
          ;;
        *)
          WF_DESCRIPTION="Execute treasury workflow ${WF_ACTION}."
          WF_PARAM_HELP=""
          ;;
      esac
      ;;
    data)
      case "$WF_ACTION" in
        data_distribution)
          WF_DESCRIPTION="Distribute dataset manifests and update catalogues."
          WF_PARAM_HELP=$'  --set manifest=PATH    Data manifest (default: data/distribution.yaml)'
          ;;
        generate_mock_data)
          WF_DESCRIPTION="Generate reproducible mock data sets for testing pipelines."
          WF_PARAM_HELP=$'  --set rows=COUNT       Number of synthetic rows (default: 500)\n  --set schema=NAME      Schema template (default: standard)'
          ;;
        merkle_proof_generator)
          WF_DESCRIPTION="Produce Merkle proofs for supplied dataset fingerprints."
          WF_PARAM_HELP=$'  --set input=PATH       File containing hashes (default: data/hashes.txt)'
          ;;
        *)
          WF_DESCRIPTION="Execute data workflow ${WF_ACTION}."
          WF_PARAM_HELP=""
          ;;
      esac
      ;;
    storage)
      WF_DESCRIPTION="Manage holographic storage replicas and audit references."
      WF_PARAM_HELP=$'  --set replicas=COUNT    Number of replicas (default: 3)\n  --set checksum=HASH    Expected checksum (optional)'
      ;;
    recovery)
      case "$WF_ACTION" in
        disaster_recovery_backup)
          WF_DESCRIPTION="Create a disaster recovery backup package with integrity metadata."
          WF_PARAM_HELP=$'  --set destination=PATH Backup destination (default: backups/disaster.tar.gz)'
          ;;
        restore_disaster_recovery)
          WF_DESCRIPTION="Restore services using a disaster recovery backup manifest."
          WF_PARAM_HELP=$'  --set source=PATH      Backup source archive (default: backups/disaster.tar.gz)'
          ;;
        restore_ledger)
          WF_DESCRIPTION="Restore the ledger state from archival checkpoints."
          WF_PARAM_HELP=$'  --set checkpoint=PATH  Ledger checkpoint archive (default: backups/ledger.tar.gz)'
          ;;
        bridge_fallback)
          WF_DESCRIPTION="Execute a fallback plan for cross-chain bridges after disruption."
          WF_PARAM_HELP=$'  --set channel=NAME     Bridge channel identifier (default: channel-0)'
          ;;
        *)
          WF_DESCRIPTION="Execute recovery workflow ${WF_ACTION}."
          WF_PARAM_HELP=""
          ;;
      esac
      ;;
    audit)
      case "$WF_ACTION" in
        block_integrity_check)
          WF_DESCRIPTION="Scan recent blocks and record integrity results."
          WF_PARAM_HELP=$'  --set window=BLOCKS    Number of blocks to audit (default: 256)'
          ;;
        immutable_audit_log_export)
          WF_DESCRIPTION="Export immutable audit logs for archival storage."
          WF_PARAM_HELP=$'  --set output=PATH      Export destination (default: audits/immutable.log)' 
          ;;
        immutable_audit_verify)
          WF_DESCRIPTION="Verify audit log signatures against recorded checkpoints."
          WF_PARAM_HELP=$'  --set source=PATH      Audit log source (default: audits/immutable.log)'
          ;;
        immutable_log_snapshot)
          WF_DESCRIPTION="Capture an immutable log snapshot and reference metadata."
          WF_PARAM_HELP=$'  --set output=PATH      Snapshot destination (default: audits/snapshot.json)'
          ;;
        immutability_verifier)
          WF_DESCRIPTION="Run the immutability verifier toolchain and record results."
          WF_PARAM_HELP=$'  --set report=PATH      Report destination (default: audits/immutability.json)'
          ;;
        forensic_data_export)
          WF_DESCRIPTION="Export forensic data bundles with hash manifests."
          WF_PARAM_HELP=$'  --set destination=PATH Output bundle (default: forensic/export.tar.gz)'
          ;;
        *)
          WF_DESCRIPTION="Execute audit workflow ${WF_ACTION}."
          WF_PARAM_HELP=""
          ;;
      esac
      ;;
    deployment)
      WF_DESCRIPTION="Orchestrate deployment automation (${WF_ACTION//_/ })."
      WF_PARAM_HELP=$'  --set manifest=PATH    Deployment manifest (default: configs/deploy/default.yaml)\n  --set env=NAME        Target environment (default: staging)'
      ;;
    devops)
      case "$WF_ACTION" in
        configure_environment)
          WF_DESCRIPTION="Configure shell environment variables and tooling defaults."
          WF_PARAM_HELP=$'  --set profile=NAME     Environment profile (default: developer)'
          ;;
        format_code)
          WF_DESCRIPTION="Run formatting checks and record affected files."
          WF_PARAM_HELP=$'  --set gofmt=BOOL       Whether to run gofmt (default: true)'
          ;;
        install_dependencies)
          WF_DESCRIPTION="Install project dependencies and capture install manifests."
          WF_PARAM_HELP=$'  --set section=NAME     Dependency section to install (default: all)'
          ;;
        update_dependencies)
          WF_DESCRIPTION="Update go module dependencies and record the diff."
          WF_PARAM_HELP=$'  --set module=NAME      Specific module to update (optional)'
          ;;
        cli_tooling)
          WF_DESCRIPTION="Bootstrap helper CLI tooling used by operators."
          WF_PARAM_HELP=$'  --set tools=LIST       Comma separated list of tools (default: synctl,vmctl)'
          ;;
        artifact_checksum)
          WF_DESCRIPTION="Compute checksums for release artifacts and store them alongside metadata."
          WF_PARAM_HELP=$'  --set path=PATH        Artifact directory (default: dist)'
          ;;
        script_completion_setup)
          WF_DESCRIPTION="Prepare shell completion scripts for the Synnergy CLI."
          WF_PARAM_HELP=$'  --set shell=NAME       Target shell (default: bash)'
          ;;
        script_launcher)
          WF_DESCRIPTION="Dispatch Synnergy workflow scripts through a common launcher."
          WF_PARAM_HELP=$'  --target SCRIPT       Explicitly run the named workflow script'
          ;;
        *)
          WF_DESCRIPTION="Execute devops workflow ${WF_ACTION}."
          WF_PARAM_HELP=""
          ;;
      esac
      ;;
    testing)
      case "$WF_ACTION" in
        contract_static_analysis)
          WF_DESCRIPTION="Run contract static analysis and capture the findings."
          WF_PARAM_HELP=$'  --set suite=NAME       Analysis suite (default: default)'
          ;;
        contract_test_suite)
          WF_DESCRIPTION="Execute the contract integration test suite."
          WF_PARAM_HELP=$'  --set package=PATH     Optional package filter (default: ./...)'
          ;;
        contract_coverage_report)
          WF_DESCRIPTION="Generate a smart contract coverage report."
          WF_PARAM_HELP=$'  --set output=PATH      Coverage output path (default: reports/coverage.out)'
          ;;
        contract_language_compatibility_test)
          WF_DESCRIPTION="Validate smart contract language compatibility across toolchains."
          WF_PARAM_HELP=$'  --set languages=LIST   Languages to test (default: solidity,vyper)'
          ;;
        integration_test_suite)
          WF_DESCRIPTION="Run the integration test suite for the platform."
          WF_PARAM_HELP=$'  --set focus=PATTERN    go test -run pattern (optional)'
          ;;
        fuzz_testing)
          WF_DESCRIPTION="Execute fuzz testing workflows and collate the results."
          WF_PARAM_HELP=$'  --set package=PATH     Package to fuzz (default: ./...)\n  --set duration=SEC     Duration per target (default: 30)'
          ;;
        performance_regression)
          WF_DESCRIPTION="Simulate a performance regression run and record key metrics."
          WF_PARAM_HELP=$'  --set benchmark=NAME   Benchmark identifier (default: default)'
          ;;
        stress_test_network)
          WF_DESCRIPTION="Coordinate the stress test harness across the network."
          WF_PARAM_HELP=$'  --set duration=SEC     Stress test duration (default: 120)'
          ;;
        benchmarks)
          WF_DESCRIPTION="Execute benchmark suites and archive performance summaries."
          WF_PARAM_HELP=$'  --set suite=NAME       Benchmark suite (default: core)'
          ;;
        gui_wallet_test)
          WF_DESCRIPTION="Run GUI wallet smoke tests and export screenshots."
          WF_PARAM_HELP=$'  --set profile=NAME     Test profile (default: smoke)'
          ;;
        scripts_test)
          WF_DESCRIPTION="Execute go test for the scripts package with structured reporting."
          WF_PARAM_HELP=$'  --set run=PATTERN      go test -run pattern (optional)'
          ;;
        release_sign_verify)
          WF_DESCRIPTION="Sign release artifacts and verify the resulting signatures."
          WF_PARAM_HELP=$'  --set artifact=PATH    Artifact to sign (default: dist/release.tar.gz)'
          ;;
        *)
          WF_DESCRIPTION="Execute testing workflow ${WF_ACTION}."
          WF_PARAM_HELP=""
          ;;
      esac
      ;;
    network)
      WF_DESCRIPTION="Operate network tooling for diagnostics and orchestration (${WF_ACTION//_/ })."
      WF_PARAM_HELP=$'  --set profile=NAME     Network profile (default: devnet)'
      ;;
    virtualization)
      WF_DESCRIPTION="Manage virtual machine profiles and sandbox policies (${WF_ACTION//_/ })."
      WF_PARAM_HELP=$'  --set profile=NAME     VM profile (default: heavy)'
      ;;
    ai)
      WF_DESCRIPTION="Train or evaluate AI models used by the Synnergy stack."
      WF_PARAM_HELP=$'  --set dataset=PATH     Training dataset (default: data/training.csv)\n  --set epochs=COUNT     Number of epochs (default: 5)'
      ;;
    documentation)
      WF_DESCRIPTION="Generate tutorial scaffolding and publish script usage guides."
      WF_PARAM_HELP=$'  --set format=NAME      Output format (default: markdown)'
      ;;
    support)
      WF_DESCRIPTION="Simulate FAQ auto-resolution flows and capture analytics."
      WF_PARAM_HELP=$'  --set topic=NAME       FAQ topic (default: onboarding)'
      ;;
    cluster)
      WF_DESCRIPTION="Provision multi-node clusters with shared configuration baselines."
      WF_PARAM_HELP=$'  --set size=COUNT       Number of nodes (default: 4)\n  --set profile=NAME     Cluster profile (default: balanced)'
      ;;
    high_availability)
      case "$WF_ACTION" in
        high_availability_setup)
          WF_DESCRIPTION="Configure active-active high availability policies and replication targets."
          WF_PARAM_HELP=$'  --set regions=LIST     Comma separated regions (default: primary,secondary)'
          ;;
        ha_failover_test)
          WF_DESCRIPTION="Execute a high availability failover exercise and record the findings."
          WF_PARAM_HELP=$'  --set scenario=NAME    Failover scenario (default: validator-outage)'
          ;;
        ha_immutable_verification)
          WF_DESCRIPTION="Validate immutability guarantees under failover scenarios."
          WF_PARAM_HELP=$'  --set checkpoint=PATH  Checkpoint reference (default: audits/snapshot.json)'
          ;;
        *)
          WF_DESCRIPTION="Execute high availability workflow ${WF_ACTION}."
          WF_PARAM_HELP=""
          ;;
      esac
      ;;
    contract)
      case "$WF_ACTION" in
        mint_nft)
          WF_DESCRIPTION="Mint an example NFT using the Synnergy CLI and record metadata."
          WF_PARAM_HELP=$'  --set name=NAME        NFT name (default: Synnergy NFT)\n  --set recipient=ADDR  Recipient address (default: wallet-1)'
          ;;
        upgrade_contract)
          WF_DESCRIPTION="Apply a smart contract upgrade and produce audit artefacts."
          WF_PARAM_HELP=$'  --set address=ADDR     Contract address (default: 0xdeadbeef)\n  --set version=TAG     Upgrade version tag (default: v2.0)'
          ;;
        *)
          WF_DESCRIPTION="Execute contract workflow ${WF_ACTION}."
          WF_PARAM_HELP=""
          ;;
      esac
      ;;
    storage|documentation|support|cluster|high_availability|contract)
      # Already handled above in specific cases.
      ;;
    misc)
      WF_DESCRIPTION="Run workflow ${base} (generic handler)."
      WF_PARAM_HELP=""
      ;;
  esac
}

build_node_steps() {
  local node_slug
  node_slug="$(subject_safe_name)"
  local node_dir="$STATE_DIR/nodes/$node_slug"
  local manifest="$node_dir/${node_slug}_manifest.json"
  local network="$(params_get network synnergy-devnet)"
  local region="$(params_get region primary)"
  local cpu="$(params_get cpu 2)"
  local memory="$(params_get memory 4Gi)"
  local replicas="$(params_get replicas 1)"

  add_step ensure_directory "Ensure ${WF_SUBJECT_TITLE} node workspace" "$node_dir"
  add_step write_file "Render ${WF_SUBJECT_TITLE} manifest" "$manifest" "{\"node_type\":\"$WF_SUBJECT\",\"network\":\"$network\",\"region\":\"$region\",\"resources\":{\"cpu\":\"$cpu\",\"memory\":\"$memory\"},\"replicas\":\"$replicas\"}"
  add_step syn_cli "Register ${WF_SUBJECT_TITLE} node" "nodes" "$(pack_args register --type "$WF_SUBJECT" --manifest "$manifest")"
  add_step syn_cli "Activate ${WF_SUBJECT_TITLE} node" "nodes" "$(pack_args activate "$WF_SUBJECT")"
  add_step note "${WF_SUBJECT_TITLE} manifest recorded" "$manifest"
}

build_consensus_steps() {
  local status_file="$STATE_DIR/consensus/status.json"
  case "$WF_ACTION" in
    start)
      local profile="$(params_get profile balanced)"
      local quorum="$(params_get quorum 0.67)"
      add_step ensure_directory "Ensure consensus state directory" "$STATE_DIR/consensus"
      add_step syn_cli "Query consensus status" "consensus" "$(pack_args status --format json)"
      add_step syn_cli "Start consensus engine" "consensus" "$(pack_args start --profile "$profile")"
      add_step syn_cli "Adjust quorum" "consensus" "$(pack_args quorum set "$quorum")"
      add_step write_file "Record consensus profile" "$status_file" "{\"profile\":\"$profile\",\"quorum\":\"$quorum\"}"
      add_step note "Consensus startup summary" "$status_file"
      ;;
    finality_check)
      local window="$(params_get window 128)"
      add_step syn_cli "Fetch finality metrics" "consensus" "$(pack_args finality --window "$window" --json)"
      add_step write_file "Persist finality window" "$status_file" "{\"finality_window\":$window}"
      ;;
    adaptive_manage)
      local demand="$(params_get demand 0.55)"
      local stake="$(params_get stake 0.45)"
      local window="$(params_get window 20)"
      add_step syn_cli "Apply adaptive demand" "consensus" "$(pack_args adjust demand "$demand" --window "$window")"
      add_step syn_cli "Apply adaptive stake" "consensus" "$(pack_args adjust stake "$stake" --window "$window")"
      add_step write_file "Document adaptive policy" "$status_file" "{\"demand\":$demand,\"stake\":$stake,\"window\":$window}"
      ;;
    difficulty_adjust)
      local delta="$(params_get delta 0.02)"
      local target="$(params_get target 0.65)"
      add_step syn_cli "Inspect current difficulty" "consensus" "$(pack_args difficulty show --json)"
      add_step syn_cli "Adjust difficulty" "consensus" "$(pack_args difficulty adjust "$delta" --target "$target")"
      add_step write_file "Difficulty tuning record" "$status_file" "{\"delta\":$delta,\"target\":$target}"
      ;;
    recovery|consensus_recovery)
      local snapshot="$(params_get snapshot state/snapshots/latest.tar.gz)"
      local integrity="$(params_get integrity none)"
      add_step ensure_directory "Ensure recovery directory" "$STATE_DIR/recovery"
      add_step note "Verify snapshot checksum" "$snapshot"
      add_step syn_cli "Stop consensus" "consensus" "$(pack_args stop)"
      add_step syn_cli "Restore consensus snapshot" "consensus" "$(pack_args restore "$snapshot")"
      add_step syn_cli "Restart consensus" "consensus" "$(pack_args start --resume)"
      add_step write_file "Recovery metadata" "$STATE_DIR/recovery/consensus.json" "{\"snapshot\":\"$snapshot\",\"integrity\":\"$integrity\"}"
      ;;
    specific_node)
      local node_id="$(params_get node validator-1)"
      local action="$(params_get action inspect)"
      add_step syn_cli "Inspect node status" "consensus" "$(pack_args node status "$node_id" --json)"
      add_step syn_cli "Apply node action" "consensus" "$(pack_args node "$action" "$node_id")"
      add_step note "Node action executed" "$node_id:$action"
      ;;
    validator_manage)
      local mode="$(params_get mode rotate)"
      local count="$(params_get count 2)"
      add_step syn_cli "Fetch validator set" "consensus" "$(pack_args validators list --json)"
      add_step syn_cli "Manage validators" "consensus" "$(pack_args validators "$mode" --count "$count")"
      add_step write_file "Validator management record" "$status_file" "{\"mode\":\"$mode\",\"count\":$count}"
      ;;
    dynamic_hopping)
      local interval="$(params_get interval 45)"
      local shards="$(params_get shards 3)"
      add_step syn_cli "Plan shard hopping" "consensus" "$(pack_args shards plan --count "$shards")"
      add_step syn_cli "Execute shard hop" "consensus" "$(pack_args shards hop --interval "$interval" --count "$shards")"
      add_step write_file "Shard hopping plan" "$status_file" "{\"interval\":$interval,\"shards\":$shards}"
      ;;
    cross_network)
      local source="$(params_get source devnet)"
      local target="$(params_get target auditnet)"
      add_step syn_cli "Export checkpoints" "consensus" "$(pack_args checkpoint export --network "$source" --out "$STATE_DIR/consensus/${source}_chk.json")"
      add_step syn_cli "Import checkpoints" "consensus" "$(pack_args checkpoint import --network "$target" --file "$STATE_DIR/consensus/${source}_chk.json")"
      add_step note "Cross-network sync completed" "$source->$target"
      ;;
    *)
      add_step note "No explicit consensus steps for $WF_ACTION" ""
      ;;
  esac
}

build_cross_chain_steps() {
  local base_dir="$STATE_DIR/cross-chain"
  local manifest="$(params_get manifest configs/cross-chain/default.yaml)"
  add_step ensure_directory "Ensure cross-chain workspace" "$base_dir"
  case "$WF_ACTION" in
    agnostic_protocols)
      add_step syn_cli "Apply agnostic protocols" "cross-chain" "$(pack_args protocols apply --manifest "$manifest")"
      add_step write_file "Protocol manifest copy" "$base_dir/agnostic_manifest.json" "{\"source\":\"$manifest\"}"
      ;;
    contracts_deploy)
      add_step syn_cli "Deploy cross-chain contracts" "cross-chain" "$(pack_args contracts deploy --manifest "$manifest")"
      add_step write_file "Deployment receipt" "$base_dir/deployment.json" "{\"manifest\":\"$manifest\"}"
      ;;
    transactions)
      local batch="$(params_get batch 10)"
      local channel="$(params_get channel cross-link)"
      add_step syn_cli "Relay transactions" "cross-chain" "$(pack_args tx relay --channel "$channel" --batch "$batch")"
      add_step write_file "Transaction plan" "$base_dir/transactions.json" "{\"channel\":\"$channel\",\"batch\":$batch}"
      ;;
    setup)
      local topology="$(params_get topology mesh)"
      add_step syn_cli "Initialise topology" "cross-chain" "$(pack_args setup --topology "$topology")"
      add_step write_file "Topology summary" "$base_dir/topology.json" "{\"topology\":\"$topology\"}"
      ;;
    connection)
      local endpoint="$(params_get endpoint https://example)"
      add_step syn_cli "Validate connection" "cross-chain" "$(pack_args connection check --endpoint "$endpoint")"
      add_step note "Connection check complete" "$endpoint"
      ;;
    bridge)
      local name="$(params_get name synnergy-bridge)"
      add_step syn_cli "Provision bridge" "cross-chain" "$(pack_args bridge create "$name")"
      add_step write_file "Bridge record" "$base_dir/${name}.json" "{\"bridge\":\"$name\"}"
      ;;
    bridge_verification)
      local channel="$(params_get channel channel-0)"
      add_step syn_cli "Verify bridge channel" "cross-chain" "$(pack_args bridge verify --channel "$channel")"
      add_step note "Bridge verification complete" "$channel"
      ;;
    bridge_fallback)
      local channel="$(params_get channel channel-0)"
      add_step syn_cli "Trigger bridge fallback" "cross-chain" "$(pack_args bridge fallback --channel "$channel")"
      add_step note "Bridge fallback simulated" "$channel"
      ;;
    *)
      add_step note "No explicit cross-chain steps for $WF_ACTION" ""
      ;;
  esac
}

build_wallet_steps() {
  local wallet_dir="$(params_get dir var/wallets)"
  local summary="$STATE_DIR/wallet/${WF_ACTION}.json"
  ensure_directory "$wallet_dir"
  case "$WF_ACTION" in
    hardware_integration)
      local device="$(params_get device ledger-nano)"
      local operator="$(params_get operator ops-team)"
      add_step ensure_directory "Ensure wallet integration dir" "$STATE_DIR/wallet"
      add_step write_file "Hardware integration manifest" "$summary" "{\"device\":\"$device\",\"operator\":\"$operator\"}"
      add_step note "Hardware integration recorded" "$device"
      ;;
    init)
      local password="$(params_get password synnergy-dev)"
      add_step ensure_directory "Ensure wallet directory" "$wallet_dir"
      add_step syn_cli "Create wallet" "wallet" "$(pack_args new --out "$wallet_dir/init.wallet" --password "$password")"
      add_step note "Wallet initialised" "$wallet_dir/init.wallet"
      ;;
    key_rotation)
      local next_key="$(params_get next-key pending)"
      local effective="$(params_get effective +24h)"
      add_step write_file "Key rotation plan" "$summary" "{\"next_key\":\"$next_key\",\"effective\":\"$effective\"}"
      add_step note "Key rotation scheduled" "$effective"
      ;;
    multisig_setup)
      local participants="$(params_get participants 3)"
      local threshold="$(params_get threshold 2)"
      add_step write_file "Multisig configuration" "$summary" "{\"participants\":$participants,\"threshold\":$threshold}"
      add_step note "Multisig plan ready" "$participants/$threshold"
      ;;
    offline_sign)
      local tx="$(params_get tx pending.tx.json)"
      local output="$(params_get output artifacts/offline.sig)"
      add_step syn_cli "Prepare offline signing" "wallet" "$(pack_args offline sign --tx "$tx" --out "$output")"
      add_step note "Offline signing prepared" "$output"
      ;;
    server_setup)
      local host="$(params_get host 127.0.0.1)"
      local port="$(params_get port 8546)"
      add_step ensure_directory "Ensure wallet server dir" "$STATE_DIR/wallet"
      add_step write_file "Wallet server config" "$summary" "{\"host\":\"$host\",\"port\":$port}"
      add_step note "Wallet server configuration ready" "$host:$port"
      ;;
    transfer)
      local from="$(params_get from wallet-1)"
      local to="$(params_get to wallet-2)"
      local amount="$(params_get amount 10)"
      add_step syn_cli "Prepare transfer" "wallet" "$(pack_args transfer --from "$from" --to "$to" --amount "$amount")"
      add_step write_file "Transfer manifest" "$summary" "{\"from\":\"$from\",\"to\":\"$to\",\"amount\":$amount}"
      ;;
    key_rotation_schedule)
      local cadence="$(params_get cadence 30)"
      add_step write_file "Key rotation schedule" "$summary" "{\"cadence\":$cadence}"
      ;;
    *)
      add_step note "No explicit wallet steps for $WF_ACTION" ""
      ;;
  esac
}

build_monitoring_steps() {
  local summary="$STATE_DIR/monitoring/${WF_ACTION}.json"
  add_step ensure_directory "Ensure monitoring state" "$STATE_DIR/monitoring"
  case "$WF_ACTION" in
    alerting_setup)
      local webhook="$(params_get webhook https://hooks.example)"
      local severity="$(params_get severity warning)"
      add_step write_file "Alerting configuration" "$summary" "{\"webhook\":\"$webhook\",\"severity\":\"$severity\"}"
      ;;
    logs_collect)
      local duration="$(params_get duration 30)"
      local target="$(params_get target logs/bundles)"
      add_step ensure_directory "Ensure log bundle directory" "$target"
      add_step note "Collect logs" "$duration minutes"
      add_step write_file "Log collection summary" "$summary" "{\"duration\":$duration,\"target\":\"$target\"}"
      ;;
    metrics_alert_dispatch)
      local rule="$(params_get rule high-cpu)"
      local repeat="$(params_get repeat 1)"
      add_step write_file "Metrics alert dispatch" "$summary" "{\"rule\":\"$rule\",\"repeat\":$repeat}"
      ;;
    metrics_export)
      local sink="$(params_get sink metrics/export.json)"
      add_step write_file "Metrics export" "$summary" "{\"sink\":\"$sink\"}"
      ;;
    system_health_logging)
      local interval="$(params_get interval 60)"
      local samples="$(params_get samples 5)"
      add_step write_file "System health plan" "$summary" "{\"interval\":$interval,\"samples\":$samples}"
      ;;
    anomaly_detection)
      local window="$(params_get window 50)"
      local sensitivity="$(params_get sensitivity 0.85)"
      add_step write_file "Anomaly detection run" "$summary" "{\"window\":$window,\"sensitivity\":$sensitivity}"
      ;;
    *)
      add_step note "No monitoring steps for $WF_ACTION" ""
      ;;
  esac
}

build_security_steps() {
  local summary="$STATE_DIR/security/${WF_ACTION}.json"
  add_step ensure_directory "Ensure security state" "$STATE_DIR/security"
  case "$WF_ACTION" in
    multi_factor_setup)
      local operator="$(params_get operator ops-team)"
      local issuer="$(params_get issuer Synnergy)"
      add_step write_file "MFA enrolment" "$summary" "{\"operator\":\"$operator\",\"issuer\":\"$issuer\"}"
      ;;
    firewall_setup)
      local profile="$(params_get profile hardened)"
      add_step write_file "Firewall profile" "$summary" "{\"profile\":\"$profile\"}"
      ;;
    secure_node_hardening)
      local checklist="$(params_get checklist none)"
      add_step write_file "Hardening checklist" "$summary" "{\"checklist\":\"$checklist\"}"
      ;;
    secure_store_setup)
      local store="$(params_get store var/secure-store)"
      local password="$(params_get password changeit)"
      add_step ensure_directory "Ensure secure store" "$store"
      add_step write_file "Secure store setup" "$summary" "{\"store\":\"$store\",\"password\":\"$password\"}"
      ;;
    key_backup)
      local destination="$(params_get destination backups/keys.tar.gz)"
      add_step write_file "Key backup plan" "$summary" "{\"destination\":\"$destination\"}"
      ;;
    key_rotation_schedule)
      local cadence="$(params_get cadence 30)"
      add_step write_file "Key rotation schedule" "$summary" "{\"cadence\":$cadence}"
      ;;
    zero_trust_data_channels)
      local policy="$(params_get policy configs/security/zt-data.yaml)"
      add_step write_file "Zero trust policy" "$summary" "{\"policy\":\"$policy\"}"
      ;;
    tamper_alert)
      local channel="$(params_get channel security-ops)"
      add_step write_file "Tamper alert config" "$summary" "{\"channel\":\"$channel\"}"
      ;;
    pki_setup)
      local ca_name="$(params_get ca-name "Synnergy Root CA")"
      add_step write_file "PKI configuration" "$summary" "{\"ca_name\":\"$ca_name\"}"
      ;;
    certificate_issue)
      local subject="$(params_get subject syn-node)"
      local san="$(params_get san none)"
      add_step write_file "Certificate issuance" "$summary" "{\"subject\":\"$subject\",\"san\":\"$san\"}"
      ;;
    *)
      add_step note "No security steps for $WF_ACTION" ""
      ;;
  esac
}

build_compliance_steps() {
  local summary="$STATE_DIR/compliance/${WF_ACTION}.json"
  add_step ensure_directory "Ensure compliance state" "$STATE_DIR/compliance"
  case "$WF_ACTION" in
    aml_kyc_process)
      local applicant="$(params_get applicant applicant-001)"
      add_step write_file "AML/KYC record" "$summary" "{\"applicant\":\"$applicant\"}"
      ;;
    identity_verification)
      local subject="$(params_get subject user-001)"
      add_step write_file "Identity verification" "$summary" "{\"subject\":\"$subject\"}"
      ;;
    credential_revocation)
      local credential="$(params_get credential cred-001)"
      add_step write_file "Credential revocation" "$summary" "{\"credential\":\"$credential\"}"
      ;;
    idwallet_register)
      local holder="$(params_get holder holder-001)"
      add_step write_file "ID wallet registration" "$summary" "{\"holder\":\"$holder\"}"
      ;;
    biometric_enroll)
      local template="$(params_get template data/biometric.bin)"
      add_step write_file "Biometric enrolment" "$summary" "{\"template\":\"$template\"}"
      ;;
    biometric_verify)
      local template="$(params_get template data/verify.bin)"
      add_step write_file "Biometric verification" "$summary" "{\"template\":\"$template\"}"
      ;;
    private_transactions)
      local circuit="$(params_get circuit zk-transfer)"
      add_step write_file "Private transaction run" "$summary" "{\"circuit\":\"$circuit\"}"
      ;;
    regulatory_report)
      local period="$(params_get period Q1)"
      add_step write_file "Regulatory report" "$summary" "{\"period\":\"$period\"}"
      ;;
    *)
      add_step note "No compliance steps for $WF_ACTION" ""
      ;;
  esac
}

build_dao_steps() {
  local summary="$STATE_DIR/dao/${WF_ACTION}.json"
  add_step ensure_directory "Ensure DAO state" "$STATE_DIR/dao"
  case "$WF_ACTION" in
    init)
      local name="$(params_get name "Synnergy DAO")"
      local quorum="$(params_get quorum 0.6)"
      add_step write_file "DAO bootstrap" "$summary" "{\"name\":\"$name\",\"quorum\":$quorum}"
      ;;
    proposal_submit)
      local file="$(params_get file proposals/example.yaml)"
      add_step write_file "Proposal submission" "$summary" "{\"proposal\":\"$file\"}"
      ;;
    proposal_lifecycle)
      local proposal="$(params_get proposal prop-001)"
      add_step write_file "Proposal lifecycle" "$summary" "{\"proposal\":\"$proposal\"}"
      ;;
    token_manage)
      local mint="$(params_get mint 0)"
      local burn="$(params_get burn 0)"
      add_step write_file "Token management" "$summary" "{\"mint\":$mint,\"burn\":$burn}"
      ;;
    vote)
      local proposal="$(params_get proposal prop-001)"
      local weight="$(params_get weight 1)"
      add_step write_file "Vote record" "$summary" "{\"proposal\":\"$proposal\",\"weight\":$weight}"
      ;;
    offchain_vote_tally)
      local source="$(params_get source data/offchain_votes.json)"
      add_step write_file "Off-chain tally" "$summary" "{\"source\":\"$source\"}"
      ;;
    stake_penalty)
      local validator="$(params_get validator validator-1)"
      local penalty="$(params_get penalty 0.05)"
      add_step write_file "Stake penalty" "$summary" "{\"validator\":\"$validator\",\"penalty\":$penalty}"
      ;;
    *)
      add_step note "No DAO steps for $WF_ACTION" ""
      ;;
  esac
}

build_treasury_steps() {
  local summary="$STATE_DIR/treasury/${WF_ACTION}.json"
  add_step ensure_directory "Ensure treasury state" "$STATE_DIR/treasury"
  case "$WF_ACTION" in
    manage)
      local target="$(params_get target 100000)"
      local currency="$(params_get currency SYN)"
      add_step write_file "Treasury management" "$summary" "{\"target\":$target,\"currency\":\"$currency\"}"
      ;;
    investment_sh)
      local strategy="$(params_get strategy conservative)"
      local allocation="$(params_get allocation 0.2)"
      add_step write_file "Investment scenario" "$summary" "{\"strategy\":\"$strategy\",\"allocation\":$allocation}"
      ;;
    grant_distribution)
      local programme="$(params_get programme ecosystem)"
      add_step write_file "Grant distribution" "$summary" "{\"programme\":\"$programme\"}"
      ;;
    grant_reporting)
      local period="$(params_get period Q1)"
      add_step write_file "Grant report" "$summary" "{\"period\":\"$period\"}"
      ;;
    financial_prediction)
      local horizon="$(params_get horizon 6)"
      local model="$(params_get model arima)"
      add_step write_file "Financial prediction" "$summary" "{\"horizon\":$horizon,\"model\":\"$model\"}"
      ;;
    *)
      add_step note "No treasury steps for $WF_ACTION" ""
      ;;
  esac
}

build_data_steps() {
  local summary="$STATE_DIR/data/${WF_ACTION}.json"
  add_step ensure_directory "Ensure data state" "$STATE_DIR/data"
  case "$WF_ACTION" in
    data_distribution)
      local manifest="$(params_get manifest data/distribution.yaml)"
      add_step write_file "Data distribution" "$summary" "{\"manifest\":\"$manifest\"}"
      ;;
    generate_mock_data)
      local rows="$(params_get rows 500)"
      local schema="$(params_get schema standard)"
      add_step write_file "Mock data configuration" "$summary" "{\"rows\":$rows,\"schema\":\"$schema\"}"
      ;;
    merkle_proof_generator)
      local input="$(params_get input data/hashes.txt)"
      add_step write_file "Merkle proof request" "$summary" "{\"input\":\"$input\"}"
      ;;
    *)
      add_step note "No data steps for $WF_ACTION" ""
      ;;
  esac
}

build_storage_steps() {
  local replicas="$(params_get replicas 3)"
  local checksum="$(params_get checksum none)"
  local summary="$STATE_DIR/storage/holographic.json"
  add_step ensure_directory "Ensure storage state" "$STATE_DIR/storage"
  add_step write_file "Holographic storage plan" "$summary" "{\"replicas\":$replicas,\"checksum\":\"$checksum\"}"
}

build_recovery_steps() {
  local summary="$STATE_DIR/recovery/${WF_ACTION}.json"
  add_step ensure_directory "Ensure recovery state" "$STATE_DIR/recovery"
  case "$WF_ACTION" in
    disaster_recovery_backup)
      local destination="$(params_get destination backups/disaster.tar.gz)"
      add_step write_file "Disaster backup" "$summary" "{\"destination\":\"$destination\"}"
      ;;
    restore_disaster_recovery)
      local source="$(params_get source backups/disaster.tar.gz)"
      add_step write_file "Disaster restore" "$summary" "{\"source\":\"$source\"}"
      ;;
    restore_ledger)
      local checkpoint="$(params_get checkpoint backups/ledger.tar.gz)"
      add_step write_file "Ledger restore" "$summary" "{\"checkpoint\":\"$checkpoint\"}"
      ;;
    bridge_fallback)
      local channel="$(params_get channel channel-0)"
      add_step write_file "Bridge fallback" "$summary" "{\"channel\":\"$channel\"}"
      ;;
    *)
      add_step note "No recovery steps for $WF_ACTION" ""
      ;;
  esac
}

build_audit_steps() {
  local summary="$STATE_DIR/audit/${WF_ACTION}.json"
  add_step ensure_directory "Ensure audit state" "$STATE_DIR/audit"
  case "$WF_ACTION" in
    block_integrity_check)
      local window="$(params_get window 256)"
      add_step write_file "Block integrity" "$summary" "{\"window\":$window}"
      ;;
    immutable_audit_log_export)
      local output="$(params_get output audits/immutable.log)"
      add_step write_file "Audit log export" "$summary" "{\"output\":\"$output\"}"
      ;;
    immutable_audit_verify)
      local source="$(params_get source audits/immutable.log)"
      add_step write_file "Audit verify" "$summary" "{\"source\":\"$source\"}"
      ;;
    immutable_log_snapshot)
      local output="$(params_get output audits/snapshot.json)"
      add_step write_file "Audit snapshot" "$summary" "{\"output\":\"$output\"}"
      ;;
    immutability_verifier)
      local report="$(params_get report audits/immutability.json)"
      add_step write_file "Immutability verifier" "$summary" "{\"report\":\"$report\"}"
      ;;
    forensic_data_export)
      local destination="$(params_get destination forensic/export.tar.gz)"
      add_step write_file "Forensic export" "$summary" "{\"destination\":\"$destination\"}"
      ;;
    *)
      add_step note "No audit steps for $WF_ACTION" ""
      ;;
  esac
}

build_deployment_steps() {
  local summary="$STATE_DIR/deployment/${WF_ACTION}.json"
  local manifest="$(params_get manifest configs/deploy/default.yaml)"
  local env="$(params_get env staging)"
  add_step ensure_directory "Ensure deployment state" "$STATE_DIR/deployment"
  add_step write_file "Deployment configuration" "$summary" "{\"manifest\":\"$manifest\",\"environment\":\"$env\"}"
}

build_devops_steps() {
  local summary="$STATE_DIR/devops/${WF_ACTION}.json"
  add_step ensure_directory "Ensure devops state" "$STATE_DIR/devops"
  case "$WF_ACTION" in
    configure_environment)
      local profile="$(params_get profile developer)"
      add_step write_file "Environment profile" "$summary" "{\"profile\":\"$profile\"}"
      ;;
    format_code)
      local gofmt="$(params_get gofmt true)"
      add_step write_file "Format plan" "$summary" "{\"gofmt\":\"$gofmt\"}"
      add_step note "Format operations prepared" "gofmt=$gofmt"
      ;;
    install_dependencies)
      local section="$(params_get section all)"
      add_step write_file "Dependency install" "$summary" "{\"section\":\"$section\"}"
      add_step note "Install dependencies" "$section"
      ;;
    update_dependencies)
      local module="$(params_get module all)"
      add_step write_file "Dependency update" "$summary" "{\"module\":\"$module\"}"
      ;;
    cli_tooling)
      local tools="$(params_get tools synctl,vmctl)"
      add_step write_file "CLI tooling" "$summary" "{\"tools\":\"$tools\"}"
      ;;
    artifact_checksum)
      local path="$(params_get path dist)"
      add_step write_file "Checksum plan" "$summary" "{\"path\":\"$path\"}"
      ;;
    script_completion_setup)
      local shell="$(params_get shell bash)"
      add_step write_file "Completion setup" "$summary" "{\"shell\":\"$shell\"}"
      ;;
    script_launcher)
      add_step note "Invoke script launcher" "Use --target to run specific workflows"
      ;;
    *)
      add_step note "No devops steps for $WF_ACTION" ""
      ;;
  esac
}

build_testing_steps() {
  local summary="$STATE_DIR/testing/${WF_ACTION}.json"
  add_step ensure_directory "Ensure testing state" "$STATE_DIR/testing"
  case "$WF_ACTION" in
    contract_static_analysis)
      local suite="$(params_get suite default)"
      add_step write_file "Static analysis" "$summary" "{\"suite\":\"$suite\"}"
      ;;
    contract_test_suite)
      local package="$(params_get package ./...)"
      add_step write_file "Contract tests" "$summary" "{\"package\":\"$package\"}"
      add_step note "Contract test command prepared" "go test $package"
      ;;
    contract_coverage_report)
      local output="$(params_get output reports/coverage.out)"
      add_step write_file "Coverage report" "$summary" "{\"output\":\"$output\"}"
      ;;
    contract_language_compatibility_test)
      local languages="$(params_get languages "solidity,vyper")"
      add_step write_file "Language compatibility" "$summary" "{\"languages\":\"$languages\"}"
      ;;
    integration_test_suite)
      local focus="$(params_get focus all)"
      add_step write_file "Integration tests" "$summary" "{\"focus\":\"$focus\"}"
      add_step note "Integration test command prepared" "go test ./... -run $focus"
      ;;
    fuzz_testing)
      local package="$(params_get package ./...)"
      local duration="$(params_get duration 30)"
      add_step write_file "Fuzz testing" "$summary" "{\"package\":\"$package\",\"duration\":$duration}"
      ;;
    performance_regression)
      local benchmark="$(params_get benchmark default)"
      add_step write_file "Performance regression" "$summary" "{\"benchmark\":\"$benchmark\"}"
      ;;
    stress_test_network)
      local duration="$(params_get duration 120)"
      add_step write_file "Stress test" "$summary" "{\"duration\":$duration}"
      ;;
    benchmarks)
      local suite="$(params_get suite core)"
      add_step write_file "Benchmarks" "$summary" "{\"suite\":\"$suite\"}"
      ;;
    gui_wallet_test)
      local profile="$(params_get profile smoke)"
      add_step write_file "GUI wallet test" "$summary" "{\"profile\":\"$profile\"}"
      ;;
    scripts_test)
      local run_pattern="$(params_get run .)"
      add_step write_file "Scripts go test" "$summary" "{\"run\":\"$run_pattern\"}"
      add_step note "Scripts go test command prepared" "go test ./scripts -run $run_pattern"
      ;;
    release_sign_verify)
      local artifact="$(params_get artifact dist/release.tar.gz)"
      add_step write_file "Release sign" "$summary" "{\"artifact\":\"$artifact\"}"
      ;;
    *)
      add_step note "No testing steps for $WF_ACTION" ""
      ;;
  esac
}

build_network_steps() {
  local summary="$STATE_DIR/network/${WF_ACTION}.json"
  add_step ensure_directory "Ensure network state" "$STATE_DIR/network"
  add_step write_file "Network workflow" "$summary" "{\"profile\":\"$(params_get profile devnet)\"}"
}

build_virtualization_steps() {
  local summary="$STATE_DIR/virtualization/${WF_ACTION}.json"
  local profile="$(params_get profile heavy)"
  add_step ensure_directory "Ensure virtualization state" "$STATE_DIR/virtualization"
  add_step write_file "Virtualization workflow" "$summary" "{\"profile\":\"$profile\"}"
}

build_ai_steps() {
  local dataset="$(params_get dataset data/training.csv)"
  local epochs="$(params_get epochs 5)"
  local summary="$STATE_DIR/ai/${WF_ACTION}.json"
  add_step ensure_directory "Ensure AI state" "$STATE_DIR/ai"
  add_step write_file "AI training" "$summary" "{\"dataset\":\"$dataset\",\"epochs\":$epochs}"
}

build_documentation_steps() {
  local format="$(params_get format markdown)"
  local summary="$STATE_DIR/documentation/${WF_ACTION}.json"
  add_step ensure_directory "Ensure documentation state" "$STATE_DIR/documentation"
  add_step write_file "Documentation workflow" "$summary" "{\"format\":\"$format\"}"
}

build_support_steps() {
  local topic="$(params_get topic onboarding)"
  local summary="$STATE_DIR/support/${WF_ACTION}.json"
  add_step ensure_directory "Ensure support state" "$STATE_DIR/support"
  add_step write_file "Support workflow" "$summary" "{\"topic\":\"$topic\"}"
}

build_cluster_steps() {
  local size="$(params_get size 4)"
  local profile="$(params_get profile balanced)"
  local summary="$STATE_DIR/cluster/${WF_ACTION}.json"
  add_step ensure_directory "Ensure cluster state" "$STATE_DIR/cluster"
  add_step write_file "Cluster workflow" "$summary" "{\"size\":$size,\"profile\":\"$profile\"}"
}

build_high_availability_steps() {
  local summary="$STATE_DIR/high-availability/${WF_ACTION}.json"
  add_step ensure_directory "Ensure HA state" "$STATE_DIR/high-availability"
  case "$WF_ACTION" in
    high_availability_setup)
      local regions="$(params_get regions primary,secondary)"
      add_step write_file "HA setup" "$summary" "{\"regions\":\"$regions\"}"
      ;;
    ha_failover_test)
      local scenario="$(params_get scenario validator-outage)"
      add_step write_file "Failover test" "$summary" "{\"scenario\":\"$scenario\"}"
      ;;
    ha_immutable_verification)
      local checkpoint="$(params_get checkpoint audits/snapshot.json)"
      add_step write_file "Immutable verification" "$summary" "{\"checkpoint\":\"$checkpoint\"}"
      ;;
    *)
      add_step note "No HA steps for $WF_ACTION" ""
      ;;
  esac
}

build_contract_steps() {
  local summary="$STATE_DIR/contract/${WF_ACTION}.json"
  add_step ensure_directory "Ensure contract state" "$STATE_DIR/contract"
  case "$WF_ACTION" in
    mint_nft)
      local name="$(params_get name "Synnergy NFT")"
      local recipient="$(params_get recipient wallet-1)"
      add_step write_file "Mint NFT" "$summary" "{\"name\":\"$name\",\"recipient\":\"$recipient\"}"
      ;;
    upgrade_contract)
      local address="$(params_get address 0xdeadbeef)"
      local version="$(params_get version v2.0)"
      add_step write_file "Upgrade contract" "$summary" "{\"address\":\"$address\",\"version\":\"$version\"}"
      ;;
    *)
      add_step note "No contract steps for $WF_ACTION" ""
      ;;
  esac
}

build_misc_steps() {
  local summary="$STATE_DIR/misc/${WF_ACTION}.json"
  add_step ensure_directory "Ensure misc state" "$STATE_DIR/misc"
  add_step write_file "Misc workflow" "$summary" "{\"action\":\"$WF_ACTION\"}"
}

build_workflow_steps() {
  reset_workflow_steps
  case "$WF_CATEGORY" in
    node) build_node_steps ;;
    consensus) build_consensus_steps ;;
    cross_chain) build_cross_chain_steps ;;
    wallet) build_wallet_steps ;;
    monitoring) build_monitoring_steps ;;
    security) build_security_steps ;;
    compliance) build_compliance_steps ;;
    dao) build_dao_steps ;;
    treasury) build_treasury_steps ;;
    data) build_data_steps ;;
    storage) build_storage_steps ;;
    recovery) build_recovery_steps ;;
    audit) build_audit_steps ;;
    deployment) build_deployment_steps ;;
    devops) build_devops_steps ;;
    testing) build_testing_steps ;;
    network) build_network_steps ;;
    virtualization) build_virtualization_steps ;;
    ai) build_ai_steps ;;
    documentation) build_documentation_steps ;;
    support) build_support_steps ;;
    cluster) build_cluster_steps ;;
    high_availability) build_high_availability_steps ;;
    contract) build_contract_steps ;;
    misc) build_misc_steps ;;
    *) build_misc_steps ;;
  esac
}

workflow_plan_mode() {
  if is_plan_mode; then
    printf 'true'
  else
    printf 'false'
  fi
}

execute_step() {
  local idx="$1"
  local type="${WF_STEP_TYPE[$idx]}"
  local desc="${WF_STEP_DESC[$idx]}"
  local data1="${WF_STEP_DATA1[$idx]}"
  local data2="${WF_STEP_DATA2[$idx]}"
  local status="planned"
  local message=""

  case "$type" in
    ensure_directory)
      if is_plan_mode; then
        message="would ensure directory $data1"
      else
        ensure_directory "$data1"
        status="ok"
        message="ensured $data1"
      fi
      ;;
    write_file)
      if is_plan_mode; then
        message="would write file $data1"
      else
        ensure_parent_dir "$data1"
        printf '%s\n' "$data2" >"$data1"
        status="ok"
        message="wrote $data1"
      fi
      ;;
    syn_cli)
      unpack_args "$data2"
      local -a args=("${WF_UNPACKED_ARGS[@]}")
      local cli_cmd="synnergy $data1"
      if ((${#args[@]})); then
        local IFS=' '
        cli_cmd+=" ${args[*]}"
      fi
      if is_plan_mode || [[ "${SYN_WORKFLOW_EXECUTE_CLI:-0}" -ne 1 ]]; then
        status="planned"
        message="$cli_cmd"
      else
        if syn_cli_executable; then
          if synnergy_cli "$data1" "${args[@]}"; then
            status="ok"
            message="$cli_cmd"
          else
            status="failed"
            message="command failed"
          fi
        else
          status="skipped"
          message="cli unavailable: $cli_cmd"
        fi
      fi
      ;;
    note)
      status="info"
      message="$data1"
      ;;
    *)
      status="info"
      message="$type: $data1"
      ;;
  esac

  case "$status" in
    ok)
      log_info "$desc" "$message"
      ;;
    planned)
      log_info "[plan] $desc" "$message"
      ;;
    info)
      log_info "$desc" "$message"
      ;;
    skipped)
      log_warn "$desc skipped: $message"
      ;;
    failed)
      log_error "$desc failed: $message"
      ;;
    *)
      log_info "$desc" "$message"
      ;;
  esac

  record_step_status "$idx" "$status" "$message"
  if [[ "$status" == "failed" ]]; then
    return 1
  fi
  return 0
}

write_workflow_manifest() {
  local base="$1"
  local output="$2"
  local overall_status="$3"
  local plan_mode="$4"

  if [[ -z "$output" ]]; then
    output="$STATE_DIR/${base}.json"
  fi

  ensure_parent_dir "$output"
  local tmp_steps
  tmp_steps="$(mktemp)"
  local idx
  for idx in "${!WF_STEP_TYPE[@]}"; do
    printf '%s\t%s\t%s\t%s\n' "${WF_STEP_DESC[$idx]}" "${WF_STEP_TYPE[$idx]}" "${WF_STEP_STATUS[$idx]}" "${WF_STEP_MESSAGE[$idx]}" >>"$tmp_steps"
  done

  local notes_payload=""
  if ((${#WORKFLOW_NOTES[@]})); then
    local note
    for note in "${WORKFLOW_NOTES[@]}"; do
      notes_payload+="$note$WF_ARG_DELIM"
    done
  fi

  local params_payload=""
  if ((${#WORKFLOW_PARAMS[@]})); then
    local key
    for key in "${!WORKFLOW_PARAMS[@]}"; do
      params_payload+="$key=${WORKFLOW_PARAMS[$key]}$WF_ARG_DELIM"
    done
  fi

  python3 - "$output" "$tmp_steps" <<'PY'
import json
import os
import sys
from pathlib import Path

output = Path(sys.argv[1])
steps_file = Path(sys.argv[2])
notes_payload = os.environ.get('NOTES_PAYLOAD', '')
params_payload = os.environ.get('PARAMS_PAYLOAD', '')
category = os.environ.get('WF_CATEGORY', '')
action = os.environ.get('WF_ACTION', '')
subject = os.environ.get('WF_SUBJECT', '')
description = os.environ.get('WF_DESCRIPTION', '')
param_help = os.environ.get('WF_PARAM_HELP', '')
plan_mode = os.environ.get('WF_PLAN_MODE', 'false')
overall_status = os.environ.get('WF_OVERALL_STATUS', 'ok')
log_file = os.environ.get('WF_LOG_FILE', '')

steps = []
with steps_file.open('r', encoding='utf-8') as fh:
    for line in fh:
        desc, step_type, status, message = line.rstrip('\n').split('\t')
        steps.append({
            'description': desc,
            'type': step_type,
            'status': status,
            'message': message,
        })

notes = [n for n in notes_payload.split('\x1f') if n]
params = {}
for entry in params_payload.split('\x1f'):
    if not entry:
        continue
    if '=' in entry:
        key, value = entry.split('=', 1)
        params[key] = value

payload = {
    'category': category,
    'action': action,
    'subject': subject,
    'description': description,
    'parameter_help': param_help,
    'plan': plan_mode == 'true',
    'status': overall_status,
    'log_file': log_file,
    'parameters': params,
    'notes': notes,
    'steps': steps,
}

output.write_text(json.dumps(payload, indent=2) + '\n', encoding='utf-8')
PY
  local status=$?
  rm -f "$tmp_steps"
  return $status
}

run_workflow_engine() {
  local base="$1"
  determine_workflow_context "$base"
  describe_workflow "$base"
  build_workflow_steps

  local overall_status="ok"
  local idx
  for idx in "${!WF_STEP_TYPE[@]}"; do
    if ! execute_step "$idx"; then
      overall_status="failed"
      if ! is_plan_mode; then
        log_warn "Continuing after failure in step $idx"
      fi
    fi
  done

  local plan="$(workflow_plan_mode)"
  export NOTES_PAYLOAD
  export PARAMS_PAYLOAD
  export WF_CATEGORY
  export WF_ACTION
  export WF_SUBJECT
  export WF_DESCRIPTION
  export WF_PARAM_HELP
  export WF_PLAN_MODE="$plan"
  export WF_OVERALL_STATUS="$overall_status"
  export WF_LOG_FILE="${LOG_FILE:-}"

  local notes_payload=""
  if ((${#WORKFLOW_NOTES[@]})); then
    local note
    for note in "${WORKFLOW_NOTES[@]}"; do
      notes_payload+="$note$WF_ARG_DELIM"
    done
  fi
  NOTES_PAYLOAD="$notes_payload"

  local params_payload=""
  if ((${#WORKFLOW_PARAMS[@]})); then
    local key
    for key in "${!WORKFLOW_PARAMS[@]}"; do
      params_payload+="$key=${WORKFLOW_PARAMS[$key]}$WF_ARG_DELIM"
    done
  fi
  PARAMS_PAYLOAD="$params_payload"

  write_workflow_manifest "$base" "$WORKFLOW_OUTPUT" "$overall_status" "$plan"
  return $?
}

