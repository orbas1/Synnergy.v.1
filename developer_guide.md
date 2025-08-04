
# Synnergy Network

Synnergy is a research blockchain exploring modular components for
decentralised finance, data storage and AI‑powered compliance.  It is written in
Go and ships a collection of command line utilities alongside the core
libraries.  The project focuses on extensibility rather than production
readiness.  Further background can be found in [`WHITEPAPER.md`](WHITEPAPER.md).

The primary entrypoint is the `synnergy` binary which exposes a large set of
sub‑commands through the
[Cobra](https://github.com/spf13/cobra) framework. Development helpers in
`core/helpers.go` allow the CLI to operate without a full node while modules are
implemented incrementally.

## Key Features

- **Modular architecture** – each blockchain capability lives in its own Go
  package and can be enabled or disabled independently.
- **AI enhanced compliance** – optional machine learning modules analyse
  transactions for risk and assist with on‑chain regulation.
- **Cross‑chain communication** – bridges, rollups and consensus hopping modules
  allow assets to move between networks and adapt to different environments.
- **Advanced governance** – quadratic voting, DAO tooling and reputation systems
  support enterprise decision making.
- **Extensive GUI suite** – wallet, explorers and marketplaces provide rich web
  interfaces backed by REST services.

## Role-based Financial Controls

- Authority nodes register with dedicated wallets and receive unique job keys
  for decrypting randomly assigned governance tasks.
- LoanPool grants distribute 5% of each payout to five selected authority
  wallets. Proposals require concurrent approval from authority and public
  voters to pass.
- ID tokens remain inactive until validated by an authority node, preventing
  double voting across governance systems.
- Only **Central Bank** nodes may deploy the SYN‑10/11/12 token standards.
- Regulated instruments such as ETFs, bonds and real‑estate tokens may only be
  issued by Government, Regulator, Creditor Bank or Central Bank nodes. Bill
  tokens are limited to creditor nodes while benefit tokens are exclusive to the
  government.
- Monetary and fiscal controls on SYN‑10/11/12 tokens are reserved for the
  government role and cannot be applied elsewhere.
- Regulator nodes can initiate security upgrades but final activation requires a
  community vote for decentralised acceptance.

## Directory Layout

The repository is organised as follows:

| Path | Description |
|------|-------------|
| `core/` | Core blockchain modules such as consensus, storage, networking and smart contract logic. See [`core/module_guide.md`](core/module_guide.md) for a file-by-file summary. |
| `cmd/cli/` | CLI command implementations using Cobra. Each file registers a command group. |
| `cmd/config/` | Default YAML configuration files and schema documentation. |
| `cmd/synnergy/` | Entry point for the `synnergy` binary and network setup walkthroughs. |
| `cmd/scripts/` | Shell helpers for common development tasks. See [`cmd/scripts/script_guide.md`](cmd/scripts/script_guide.md). |
| `cmd/smart_contracts/` | Example WebAssembly contracts bundled with the CLI. |
| `scripts/` | Top-level scripts to start local devnets or testnets. |
| `tests/` | Go unit tests covering each module. Run with `go test ./...`. |
| `GUI/` | Web front-ends including `wallet`, `explorer`, `ai-marketplace` and more. |
| `walletserver/` | Go HTTP backend powering the wallet GUI. |
| `internal/` | Shared utilities used across CLI and services. |
| `ai.proto` | gRPC definitions for AI-related services. |
| `smart_contract_guide.md` | Documentation and examples for writing contracts. |
| `third_party/` | Vendored dependencies such as a libp2p fork used during early development. |
| `setup_synn.sh` | Convenience script that installs Go and builds the CLI. |
| `Synnergy.env.sh` | Optional environment bootstrap script that downloads tools and loads variables from `.env`. |
| `WHITEPAPER.md` | Project vision and design goals. |

## Environment Setup

1. Install [Go](https://go.dev/dl/) 1.20 or newer and ensure `$GOPATH/bin` is on
   your `PATH`.
2. Clone the repository and download dependencies:

   ```bash
   git clone https://github.com/synnergy_network/Synnergy.git
   cd Synnergy/synnergy-network
   go mod tidy
   ```
3. Run `../setup_synn.sh` for a minimal toolchain or `../Synnergy.env.sh` for a
   full environment with additional dependencies and variables loaded from
   `.env`.
4. For GUI projects install Node.js 18+ and your preferred package manager.


## Building

Once the environment is prepared, compile the main binary from the repository
root:

```bash
cd synnergy-network
go build ./cmd/synnergy
```

The resulting binary `synnergy` can then be executed with any of the available
sub‑commands. Development stubs in `core/helpers.go` expose `InitLedger` and
`NewTFStubClient` so the CLI can run without a full node during testing.

The provided script `./setup_synn.sh` automates dependency installation and
builds the CLI.  For a fully configured environment with additional tooling and
variables loaded from `.env`, run `./Synnergy.env.sh` after cloning the
repository.

## Command Groups

Each file in `cmd/cli` registers its own group of commands with the root
`cobra.Command`. The main program consolidates them so the CLI surface mirrors
all modules from the core library. Highlights include:

- `ai` – publish models and run inference jobs
- `ai_contract` – manage AI enhanced contracts with risk checks
- `ai-train` – manage on-chain AI model training
- `ai_mgmt` – manage AI model marketplace listings
- `ai_infer` – advanced inference and transaction analysis
- `amm` – swap tokens and manage liquidity pools
- `authority_node` – validator registration and voting; each authority registers
  with a dedicated wallet address and activates after meeting voting thresholds
- `access` – manage role based access permissions
- `authority_apply` – submit and approve authority node applications
- `charity_pool` – query and disburse community funds
- `charity_mgmt` – manage donations and internal payouts
- `identity` – manage verified addresses
- `coin` – mint and transfer the native SYNN coin
- `compliance` – perform KYC/AML checks
- `audit` – manage on-chain audit logs
- `compliance_management` – suspend or whitelist addresses
- `consensus` – control the consensus engine
- `consensus_hop` – dynamically switch consensus mode
- `adaptive` – dynamically adjust consensus weights
- `stake` – adjust validator stake and record penalties
- `contracts` – deploy and invoke smart contracts
- `contractops` – administrative tasks such as pausing and upgrading contracts
- `cross_chain` – configure asset bridges
- `ccsn` – manage cross-consensus scaling networks
- `xcontract` – register cross-chain contract mappings
- `cross_tx` – execute cross-chain transfers
- `cross_chain_connection` – manage chain-to-chain links
- `cross_chain_agnostic_protocols` – register cross-chain protocols
- `cross_chain_bridge` – manage cross-chain transfers
- `data` – inspect and manage raw data storage
- `fork` – inspect and resolve chain forks
- `messages` – queue, process and broadcast network messages
- `partition` – partition data and apply compression
- `distribution` – publish datasets and handle paid access
- `oracle_management` – monitor oracle performance and sync feeds
- `data_ops` – advanced data feed operations
- `anomaly_detection` – analyse transactions for suspicious activity
- `resource` – manage stored data and VM gas allocations
- `immutability` – verify the canonical chain state
- `fault_tolerance` – simulate faults and backups
- `plasma` – deposit tokens and process exits on the plasma bridge
- `resource_allocation` – manage per-contract gas limits
- `failover` – manage ledger snapshots and trigger recovery
- `employment` – manage on-chain employment contracts and salaries
- `governance` – DAO style governance commands
- `token_vote` – cast token weighted governance votes
- `qvote` – quadratic voting on governance proposals
- `polls_management` – create and vote on community polls
- `governance_management` – register and control governance contracts
- `reputation_voting` – weighted voting using SYN-REP tokens
- `timelock` – schedule and execute delayed proposals
- `dao` – create and manage DAOs
- `green_technology` – sustainability features
- `resource_management` – track and charge node resources
- `carbon_credit_system` – manage carbon offset projects and credits
- `energy_efficiency` – track transaction energy usage and efficiency metrics
- `ledger` – low level ledger inspection
- `loanpool` – create and vote on proposals
- `loanmgr` – manage the loan pool (pause, resume, stats)
- `account` – manage basic accounts and balances
- `loanpool_apply` – manage loan applications with on-chain voting
- `network` – libp2p networking helpers
- `bootstrap` – run a dedicated bootstrap node for peers
- `connpool` – manage reusable outbound connections
- `peer` – peer discovery and connection utilities
 - `replication` – snapshot and replicate data
 - `high_availability` – manage standby nodes and promote backups
 - `rollups` – manage rollup batches
- `plasma` – manage plasma deposits and exits
- `replication` – snapshot and replicate data
- `coordination` – coordinate distributed nodes and broadcast ledger state
 - `rollups` – manage rollup batches and aggregator state
- `initrep` – bootstrap a new node by synchronizing the ledger
- `synchronization` – manage blockchain catch-up and progress
- `rollups` – manage rollup batches
- `compression` – save and load compressed ledger snapshots
- `security` – cryptographic utilities
- `firewall` – manage block lists for addresses, tokens and peers
- `biometrics` – manage biometric authentication templates
- `sharding` – shard management
- `sidechain` – launch, manage and interact with sidechains
- `state_channel` – open and settle payment channels
- `grant` – create and release loan pool grants
- `plasma` – manage plasma deposits and block commitments
- `state_channel_mgmt` – pause, resume and force-close channels
- `zero_trust_data_channels` – encrypted data channels with ledger-backed escrows
- `swarm` – orchestrate groups of nodes for high availability
- `storage` – interact with on‑chain storage providers
- `legal` – manage Ricardian contracts and signatures
- `ipfs` – manage IPFS pins and retrieval through the gateway
- `resource` – rent computing resources via the marketplace
- `staking` – lock and release tokens for governance
- `dao_access` – manage DAO membership roles
- `sensor` – manage external sensor inputs and webhooks
- `real_estate` – manage tokenised real estate
- `escrow` – manage multi-party escrow accounts
- `marketplace` – buy and sell items using escrow
- `healthcare` – manage on‑chain healthcare records
- `tangible` – register and transfer tangible asset records
- `warehouse` – manage on‑chain inventory records
- `tokens` – ERC‑20 style token commands
- `defi` – insurance, betting and other DeFi operations
- `event_management` – record and query on-chain events
- `token_management` – high level token creation and administration
- `gaming` – create and join simple on-chain games
- `transactions` – build and sign transactions
- `private_tx` – encrypt and submit private transactions
- `transactionreversal` – request authority-backed reversals
- `transaction_distribution` – split transaction fees between miner and treasury
- `faucet` – dispense test tokens or coins with rate limits
- `utility_functions` – assorted helpers
- `geolocation` – record node location information
- `distribution` – bulk token distribution and airdrops
- `finalization_management` – finalize blocks, rollup batches and channels
- `quorum` – simple quorum tracker management
- `virtual_machine` – run the on‑chain VM service

### Authority Node Policy

Authority nodes are specialised participants that oversee sensitive network
operations. When registering, each node must provide a wallet address for
payments and fees. Approval requires both public and authority votes to reach
configured thresholds. Their multi-signature approval is required for actions
like transaction reversals, loan pool authorisation and regulated token
issuance. Privileged roles – Government, Regulator, Creditor Bank and Central
Bank – restrict deployment of certain financial instruments and SYN‑10/11/12
tokens.
- `sandbox` – manage VM sandboxes
- `workflow` – automate multi-step tasks with triggers and webhooks
- `supply` – manage supply chain assets on chain
- `wallet` – mnemonic generation and signing
- `execution` – orchestrate block creation and transaction execution
- `regulator` – manage approved regulators and enforce rules
- `feedback` – submit and query user feedback
- `system_health` – monitor runtime metrics and write logs
- `idwallet` – register ID-token wallets and verify status
- `offwallet` – offline wallet utilities
- `recovery` – register and invoke account recovery
- `wallet_mgmt` – manage wallets and send SYNN directly via the ledger

Quadratic voting allows token holders to weight their governance votes by the
square root of the staked amount. The `qvote` command submits these weighted
votes and queries results alongside standard governance commands.

More details for each command can be found in `cmd/cli/cli_guide.md`.

## Core Modules

The Go packages under `core/` implement the blockchain runtime. Key modules
include consensus, storage, networking and the virtual machine.  A summary of
every file is maintained in [`core/module_guide.md`](core/module_guide.md). New
contributors should review that document to understand dependencies between
packages.

## Configuration

Runtime settings are defined using YAML files in `cmd/config/`.  The CLI loads
`default.yaml` by default and merges any environment specific file if the
`SYNN_ENV` environment variable is set (for example `SYNN_ENV=prod`).
`bootstrap.yaml` provides a template for running a dedicated bootstrap node. You can launch it via `synnergy bootstrap start` and point other nodes to its address.
The configuration schema is documented in [`cmd/config/config_guide.md`](cmd/config/config_guide.md).

## Running a Local Network

Once the CLI has been built you can initialise a test ledger and start the core
services locally.  A detailed walk‑through is provided in
[`cmd/synnergy/synnergy_set_up.md`](cmd/synnergy/synnergy_set_up.md), but the
basic steps are:

```bash
synnergy ledger init --path ./ledger.db
synnergy network start &
synnergy wallet create --out wallet.json
synnergy coin mint $(jq -r .address wallet.json) 1000
```

By default the node attempts to open the listening port on your router using
UPnP or NAT‑PMP. You can inspect and manage these mappings with the new `nat`
command group.

Additional helper scripts live under `cmd/scripts`.  Running
`start_synnergy_network.sh` will build the CLI, launch networking, consensus and
other daemons, then run a demo security command.

Two top level scripts provide larger network setups:
`scripts/devnet_start.sh` spins up a local multi-node developer network, while
`scripts/testnet_start.sh` starts an ephemeral testnet defined by a YAML
configuration. Both build the CLI automatically and clean up all processes on
`Ctrl+C`.

## Scripts

Reusable shell helpers in `cmd/scripts` automate common tasks such as funding
accounts (`faucet_fund.sh`), deploying contracts (`contracts_deploy.sh`),
starting network services and more. The complete reference is available in
[`cmd/scripts/script_guide.md`](cmd/scripts/script_guide.md).

Top level scripts are provided for multi-node setups:

- `scripts/devnet_start.sh` – start a local developer network.
- `scripts/testnet_start.sh` – launch an ephemeral testnet from YAML.

## GUI Projects

Browser-based interfaces reside under `GUI/`:

- `wallet` – manage accounts using the `walletserver` backend.
- `explorer` – query balances and transactions.
- `ai-marketplace` – browse and purchase AI models.
- `smart-contract-marketplace` – deploy and trade contracts.
- `storage-marketplace` – pay for decentralized storage.
- `nft_marketplace` – create and trade NFTs.
- `dao-explorer` – interact with DAO governance.
- `token-creation-tool` – generate new tokens.
- `dex-screener` – view decentralized exchange listings.
- `authority-node-index` – monitor authority nodes.
- `cross-chain-management` – manage bridges and transfers.

The `walletserver/` directory contains the Go HTTP backend for the wallet GUI.

## Smart Contracts

Example contracts live throughout the repository and are documented in
[`smart_contract_guide.md`](smart_contract_guide.md). They demonstrate opcode
usage for token faucets, storage markets, DAO governance and more. Contracts are
compiled to WebAssembly and deployed via the `synnergy` CLI.

## Docker

A `Dockerfile` at the repository root allows running Synnergy without a local Go installation.
Build the image and start a node with:

```bash
docker build -t synnergy ..
docker run --rm -it synnergy
```

The container launches networking, consensus, replication and VM services automatically.

## Testing

Unit tests are located in the `tests` directory and mirror the structure of the
`core` packages. For any code change run the standard toolchain:

```bash
go fmt ./...
go vet ./...
go build ./...
go test ./...
```

Some modules depend on live services such as the network or security daemon.
Those tests may require additional environment variables or mock
implementations; consult the comments within each test file for setup
instructions.

## User Feedback and UX

Community input guides interface design across the project.

- Report usability and accessibility issues through the [UX Feedback issue template](../.github/ISSUE_TEMPLATE/ux_feedback.md), including environment details, user type and impact.
- Follow the [UX Feedback Guide](UX_FEEDBACK_GUIDE.md) for the feedback lifecycle, accessibility checklist and metrics.
- Issues are triaged weekly and outcomes are summarised in a monthly community call where priorities are set and metrics reviewed.

## Contributing

Development tasks are organised in stages described in
[`AGENTS.md`](../AGENTS.md). When contributing code:

1. Work through the stages sequentially and modify no more than three files per
   commit or pull request.
2. Run `go fmt`, `go vet`, `go build` and `go test` on the packages you touch
   before committing.
3. Update any relevant documentation and cross-reference new features in this
   README or the module guide.
4. Mark completed files in `AGENTS.md` so others know which tasks are in
   progress.

Use the `setup_synn.sh` script when preparing a new environment; it installs
Go, fetches dependencies and builds the CLI.

## License

Synnergy is provided for research and educational purposes.  Third‑party
dependencies located under `third_party/` retain their original licenses.  Refer
to the respective `LICENSE` files in those directories for details.
