# How to Use the CLI

The Synnergy command-line interface (CLI) is the primary gateway for interacting with the Synnergy blockchain. Developed and maintained by **Neto Solaris**, the `synnergy` binary enables administrators, validators, institutions, and developers to operate nodes, manage on-chain assets, and integrate with external systems. This guide provides an end‑to‑end overview of installation, command structure, and the most common operational workflows.

## 1. Installation
### 1.1 Build from Source
```bash
go build ./cmd/synnergy   # or: make build
```
The resulting `synnergy` executable is placed in the repository root. Build tags such as `-tags prod` can be supplied for environment‑specific defaults. Verify the binary and configuration with:
```bash
./synnergy --version
echo $SYN_CONFIG  # confirms active config file
```

### 1.2 Configuration
Synnergy reads a YAML configuration file whose path is resolved from the `SYN_CONFIG` environment variable. Development nodes default to `configs/dev.yaml` while production builds point to `configs/prod.yaml`.

## 2. Command Structure and Conventions
Synnergy uses the [Cobra](https://github.com/spf13/cobra) framework. Commands follow the pattern:
```bash
synnergy <module> <subcommand> [flags] [arguments]
```
Global flags include `--json` for machine‑readable output and `--help` for in‑line documentation. All modules are registered beneath the root command defined in `cli/root.go`.

## 3. Core Modules
### 3.1 Network Operations
Control the peer‑to‑peer layer, publish messages, and monitor connectivity:
```bash
synnergy network start
synnergy network peers
synnergy network broadcast topic "payload"
```
Supported subcommands: `start`, `stop`, `peers`, `broadcast`, `subscribe`.

### 3.2 Wallet Management
Generate encrypted wallets for signing transactions:
```bash
synnergy wallet new --out wallet.json --password secret
```
The command prints the address and optionally persists an encrypted keystore to disk.

### 3.3 Authority Governance
Register authority nodes, cast votes, and inspect the registry:
```bash
# Application lifecycle
synnergy authority_apply submit node1 validator "candidate node"
synnergy authority_apply vote voterA <appID> true
synnergy authority_apply finalize <appID>
synnergy authority_apply get <appID> --json
synnergy authority_apply list --json
synnergy authority_apply tick

# Registry management
synnergy authority register node1 validator
synnergy authority vote voterA node1
synnergy authority info node1 --json
synnergy authority elect 5
synnergy authority is node1
synnergy authority list --json
synnergy authority deregister node1
```
`authority_apply` manages candidacy, voting, and finalisation, while `authority` maintains the active registry.

### Stage 79 Enterprise Bootstrap
Use `synnergy orchestrator bootstrap --node-id node79 --consensus Synnergy-PBFT --governance SYN-Gov --authority treasury=governor` to automate node provisioning. The command invokes `core.EnterpriseOrchestrator.BootstrapNetwork`, emitting a signed bootstrap signature, registering authority roles, enabling ledger replication and returning diagnostics so subsequent modules start from a trusted baseline.【F:cli/orchestrator.go†L58-L117】【F:core/enterprise_orchestrator.go†L71-L209】【F:core/enterprise_orchestrator_test.go†L73-L178】 Stage 79 gas costs are synchronised during startup and the web control panel mirrors the same workflow, keeping CLI automation, browser operators and documentation aligned on pricing and readiness signals.【F:cmd/synnergy/main.go†L63-L106】【F:web/pages/index.js†L1-L214】【F:web/pages/api/bootstrap.js†L1-L45】

### 3.4 Banking Integration
Institutions can register participation and query supported node types:
```bash
synnergy bankinst register MyBank
synnergy bankinst list --json
synnergy bankinst is MyBank
synnergy banknodes types
```
These modules expose the institutional registry and enumerate allowed banking node roles.

### 3.5 Liquidity Pools
Create pools, add liquidity, execute swaps, and inspect state:
```bash
synnergy liquidity_pools create TOKENA TOKENB 25
synnergy liquidity_pools add TOKENA-TOKENB provider1 1000 1000
synnergy liquidity_pools swap TOKENA-TOKENB TOKENA 50 45
synnergy liquidity_pools remove TOKENA-TOKENB provider1 500
synnergy liquidity_pools info TOKENA-TOKENB
synnergy liquidity_views list
```
`liquidity_pools` performs state‑changing operations; `liquidity_views` provides read‑only inspection.

### 3.6 Cross-Chain Bridge
Lock assets for transfer to other chains and release them using proofs:
```bash
synnergy cross_chain_bridge deposit bridge1 alice bob 100 TOKENX
synnergy cross_chain_bridge list --json
synnergy cross_chain_bridge get <transferID> --json
synnergy cross_chain_bridge claim <transferID> <proof>
```
Transfers can be listed or fetched individually, each reporting deterministic gas costs.

### 3.7 Coin Economics
Query monetary parameters and perform economic calculations:
```bash
synnergy coin info
synnergy coin reward 100000
synnergy coin supply 100000
synnergy coin price C R M V T E
synnergy coin alpha 0.3 0.8 0.5 0.1
synnergy coin minstake 500000 10 1000000 0.2
```
Each subcommand validates numeric inputs and emits either plain text or JSON depending on `--json`.

### 3.8 Consensus Coordination
Tune the adaptive consensus engine or measure transition thresholds:
```bash
synnergy consensus mine 5
synnergy consensus weights
synnergy consensus adjust 0.7 0.3
synnergy consensus threshold 0.5 0.6
synnergy consensus difficulty 1.0 0.9 1.0
synnergy consensus availability true true false
synnergy consensus powrewards true
```
These commands mine test blocks, display and adjust PoW/PoS/PoH weights, and toggle reward or availability parameters.

### 3.9 Smart Contract Lifecycle
Compile, deploy, invoke, and manage WebAssembly contracts:
```bash
synnergy contracts compile contract.wat
synnergy contracts deploy --wasm build/contract.wasm --gas 200000 --owner addr1
synnergy contracts invoke <address> --method run --args "payload" --gas 5000
synnergy contracts list
synnergy contracts deploy-template --name token_faucet --owner addr1
synnergy contracts list-templates

# Administrative management
synnergy contract-mgr pause <address>
synnergy contract-mgr resume <address>
synnergy contract-mgr upgrade <address> <wasmHex> 250000
synnergy contract-mgr info <address>
```
`contracts` handles compilation and deployment, while `contract-mgr` transfers ownership, pauses execution, and upgrades bytecode.

### 3.10 Testnet Faucet
Bootstrap local testing by minting tokens from the faucet:
```bash
synnergy faucet init --balance 1000 --amount 10 --cooldown 1m
synnergy faucet request addr1
synnergy faucet balance
synnergy faucet config --amount 5 --cooldown 30s
```
The faucet must be initialised before requests are honoured.

### 3.11 Loan Pool Governance
Submit and manage community lending proposals:
```bash
synnergy loanpool submit creatorA recipientB microloan 500 "startup capital"
synnergy loanpool vote voter1 1
synnergy loanpool disburse 1
synnergy loanpool get 1
synnergy loanpool list
synnergy loanpool cancel creatorA 1
synnergy loanpool extend creatorA 1 24
synnergy loanpool tick
```
Loans accrue votes until disbursement or cancellation; `tick` processes expirations.

### 3.12 Charity Pool Distribution
Register charities, collect votes, and audit allocations:
```bash
synnergy charity_pool register addr1 1 "WaterAid"
synnergy charity_pool vote voter1 addr1
synnergy charity_pool winners
synnergy charity_mgmt donate donor1 100
synnergy charity_mgmt withdraw addr1 50
```
Registrations and winners can be queried per cycle to verify transparent fund dispersal.

### 3.13 Identity Verification
Maintain identity profiles and verification logs:
```bash
synnergy identity register addr1 "Alice" 1990-01-01 UK
synnergy identity verify addr1 passport
synnergy identity info addr1
synnergy identity logs addr1
```
Profiles record name, date of birth, nationality, and the methods used to verify an address.

### 3.14 Biometric Security Node
Offload authentication to a dedicated biometric node:
```bash
synnergy bsn enroll addr1 "iris-scan" <pubKeyHex>
synnergy bsn auth addr1 "iris-scan" <sigHex>
synnergy bsn addtx addr1 "iris-scan" <sigHex> fromAddr toAddr 10 1 1
```
The node stores biometric templates, authenticates signatures, and can queue signed transactions for propagation.

### 3.15 Regulatory Management
Govern jurisdiction‑specific transaction policies:
```bash
synnergy regulator add reg1 EU "GDPR limit" 1000
synnergy regulator list
synnergy regulator evaluate 500
synnergy regulator remove reg1
```
Rules can be added or removed dynamically, and every transaction amount can be evaluated against active limits.

### 3.16 Node Operations
Inspect and maintain a running node:
```bash
synnergy node info
synnergy node stake addr1 1000
synnergy node slash addr1 double
synnergy node rehab addr1
synnergy node addtx fromAddr toAddr 10 1 1
synnergy node mempool
synnergy node mine
```
Nodes track validator stakes, record penalties, accept transactions into the mempool, and mine blocks on demand.

### 3.17 Validator Management
Operate the consensus validator set:
```bash
synnergy validator add addr1 1000
synnergy validator remove addr1
synnergy validator slash addr1
synnergy validator eligible
synnergy validator stake addr1
synnergy validator set-min 5000
```
Commands register validators, adjust their stake, and enforce minimum staking thresholds.

### 3.18 Fee and Gas Utilities
Estimate fees and inspect opcode gas tables:
```bash
synnergy fees estimate --type transfer --units 2 --tip 1 --load 0.3
synnergy fees feedback --message "estimates high"
synnergy fees share 100 60 40
synnergy gas list
synnergy gas set 5 50
synnergy gas snapshot --json
```
`fees` calculates cost breakdowns while `gas` snapshots or adjusts opcode pricing for deterministic budgeting.

### 3.19 Compliance Operations
Validate KYC data, track fraud, and run anomaly detection:
```bash
synnergy compliance validate kyc.json
synnergy compliance erase addr1
synnergy compliance fraud addr1 4
synnergy compliance risk addr1
synnergy compliance audit addr1
synnergy compliance monitor tx.json 0.8
synnergy compliance verifyzkp blob.bin 0xcommit 0xproof
```
All commands accept `--json` for structured auditing.

### 3.20 Audit Logs
Record and review event trails:
```bash
synnergy audit log addr1 login ip=1.2.3.4
synnergy audit list addr1 --json
```
Each entry captures timestamped metadata for compliance and forensic analysis.

### 3.21 Private Transactions
Encrypt payloads and submit confidential transactions:
```bash
synnergy private-tx encrypt key "secret"
synnergy private-tx decrypt key <hexdata>
synnergy private-tx send encrypted_tx.json
```
Payloads are symmetrically encrypted and relayed through a privacy-preserving manager.

### 3.22 DAO Governance
Create member‑driven organisations:
```bash
synnergy dao create DevDAO creator1
synnergy dao join daoID addr1
synnergy dao leave daoID addr1
synnergy dao info daoID
synnergy dao list
```
The manager issues unique IDs and enforces membership operations.

### 3.23 Historical Ledger
Archive block summaries and retrieve historical data:
```bash
synnergy historical archive 100 abcd1234
synnergy historical height 100
synnergy historical hash abcd1234
synnergy historical total
```
Archived summaries capture height, hash, and timestamp for audit‑grade lineage.

### 3.24 High Availability Failover
Configure backup nodes for rapid failover:
```bash
synnergy highavailability init primary1 30
synnergy highavailability add backup1
synnergy highavailability heartbeat backup1
synnergy highavailability active
```
The manager tracks heartbeat signals and automatically promotes a standby when the primary is unresponsive.

### 3.25 Watchtower Monitoring
Run a watchtower node to supervise network health:
```bash
synnergy watchtower start
synnergy watchtower fork 500 badhash
synnergy watchtower metrics
synnergy watchtower stop
```
Metrics report CPU, memory, peer count, and block height with timestamped snapshots.

### 3.26 Rollups
Aggregate transactions into batches with fraud proofs:
```bash
synnergy rollups submit tx1 tx2
synnergy rollups challenge batchID 0 proof
synnergy rollups finalize batchID true
synnergy rollups info batchID
synnergy rollups list
synnergy rollups txs batchID
synnergy rollups pause
synnergy rollups resume
synnergy rollups status
```
Rollup batches can be paused, resumed, challenged, or finalised to scale throughput.

### 3.27 Sharding
Manage shard leaders and cross‑shard traffic:
```bash
synnergy sharding leader get 1
synnergy sharding leader set 1 addr1
synnergy sharding map
synnergy sharding submit 0 1 txhash
synnergy sharding pull 1
synnergy sharding reshard 4
synnergy sharding rebalance 75
```
Sharding commands rebalance load and route cross‑shard transaction receipts.

### 3.28 Sidechains
Register and administer auxiliary chains:
```bash
synnergy sidechain register chain1 "meta" val1 val2
synnergy sidechain header chain1 header
synnergy sidechain get-header chain1
synnergy sidechain meta chain1
synnergy sidechain list
synnergy sidechain pause chain1
synnergy sidechain resume chain1
synnergy sidechain update-validators chain1 val3 val4
synnergy sidechain remove chain1
synnergy sidechain deposit chain1 addr1 100
synnergy sidechain withdraw chain1 addr1 50 proof
```
Each side‑chain maintains independent headers, validator sets, and escrow operations for deposits and withdrawals.

### 3.29 Plasma Bridge
Control the Plasma asset bridge:
```bash
synnergy plasma-mgmt pause
synnergy plasma-mgmt resume
synnergy plasma-mgmt status --json
```
Pausing safeguards bridge operations during upgrades or incident response.

## 4. JSON Output and Automation
Most operational commands support a `--json` flag to aid scripting and GUI integration. For example, `synnergy authority_apply list --json` returns a serialized array of applications, and `synnergy coin --json info` outputs key/value pairs describing the coin parameters. Structured output simplifies downstream tooling, log aggregation, and API responses.

## 5. Infrastructure Automation
Infrastructure‑as‑code tooling underpins reproducible deployments:
```bash
cd deploy/terraform
terraform init
terraform apply -var 'ami_id=ami-123456'

ansible-playbook -i inventory deploy/ansible/playbook.yml
```
Terraform provisions cloud resources while Ansible configures nodes with hardened defaults.

## 6. Testing and Continuous Integration
Run the full suite to verify CLI behaviour and runtime contracts:
```bash
go test ./...
```
Integration tests execute commands against in‑memory nodes and wallet services to ensure end‑to‑end flow integrity across releases.

## 7. Operational Guidance
For enterprise deployments, pair CLI workflows with observability and policy controls:
- **Structured logging** – every command emits deterministic strings that can be forwarded to SIEM pipelines.
- **Role separation** – limit privileged subcommands (e.g., `contract-mgr` and `regulator`) to secured service accounts.
- **Environment parity** – maintain dedicated config files for development, staging, and production to ensure reproducible builds.
- **Disaster recovery** – export wallet files and configuration snapshots regularly to support node restoration.

---
Neto Solaris maintains the Synnergy CLI as an enterprise‑grade interface for blockchain administration. Continual testing and automation ensure the tool remains reliable for both manual operators and scripted workflows.
