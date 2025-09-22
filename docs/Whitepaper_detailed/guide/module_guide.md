# Synnergy Core Module Guide

The `core` directory contains the fundamental building blocks of the Synnergy
network.  Each Go source file defines an individual **module** that can usually
be compiled and tested in isolation.  Modules favour small, explicit
dependencies and generally only import `common_structs.go` for shared types.

This guide provides an expanded orientation for contributors who want to
understand or extend the runtime.  It groups related modules and highlights
their responsibilities, common interactions and supporting documentation.

## Stage 82 Module Coordination

Stage 82 tightens the integration between modules that previously operated
independently. `bootstrapRuntime` now initialises the authority registry,
consensus manager, VM sandbox, secrets manager, wallet services and token suite
in a single routine. `registerEnterpriseGasMetadata` guarantees that contract,
storage, DAO, wallet, node and orchestrator modules share the same gas metadata,
while `configureLogging` validates logging expectations before any module emits
telemetry. The new `EnterpriseOrchestrator` helpers expose module state through
`synnergy orchestrator bootstrap`, and execution hooks surface VM failures to the
log pipeline so maintainers can trace issues back to the originating module
without enabling verbose debug builds.

## Directory Orientation

```
synnergy-network/
├─ core/                 # runtime modules (this guide)
├─ tests/                # accompanying unit and integration tests
└─ docs/                 # project wide documentation
```

The `tests` folder mirrors the structure of `core` and contains regression
tests for most modules.  When adding new functionality, create a matching test
file in `synnergy-network/tests`.

## Module Categories

### Consensus and Networking

These modules drive block production and peer-to-peer communication.

- **consensus.go** – Hybrid Proof‑of‑History / Proof‑of‑Stake engine.  Works in
  tandem with `ledger.go`, `network.go` and `security.go` to finalise blocks.
- **consensus_adaptive_management.go** and
  **consensus_difficulty.go** – Runtime tuning of validator sets and mining
  difficulty.
- **network.go** – libp2p networking stack offering pub‑sub gossip, peer
  discovery and stream multiplexer helpers.
- **kademlia.go**, **peer_management.go**, **nat_traversal.go** and
  **quorum_tracker.go** – Maintain peer tables, map ports behind NAT and track
  quorum participation.
- **fault_tolerance.go**, **bft_simulation.go** and
  **high_availability.go** – Byzantine fault detection, recovery and simulation
  utilities.
- **replication.go** – Propagates blocks and snapshots to new or recovering
  nodes.

### Ledger and State Management

Persistent chain state and transaction execution live in these modules.

- **ledger.go** – Authoritative blockchain state with a write‑ahead log and
  periodic snapshots.
- **transactions.go** and **tx_types.go** – Core transaction structures,
  signature validation and fee calculations.
- **storage.go** – Backend‑agnostic key/value adapters used by the ledger and
  virtual machine.
- **sharding.go**, **state_channel.go** and **plasma.go** – Partition or move
  state off‑chain for scalability and fast settlement.
- **virtual_machine.go**, **opcode_dispatcher.go** and **gas_table.go** –
  Execute WebAssembly contracts and map opcodes to gas costs.  The
  `opcode_and_gas_guide.md` file provides an in‑depth reference.
- **resource_management.go**, **resource_allocator.go** and
  **resource_allocation_management.go** – Monitor and enforce per‑contract
  resource limits.

### Tokens and Financial Primitives

Modules for coins, tokens and financial operations.

- **coin.go**, **tokens.go** and numerous `tokens_syn*.go` files – Define the
  native coin and token factories used across the platform.
- **token_management.go** and related helpers – Minting, burning and metadata
  management for issued assets.
- **amm.go** and **liquidity_pools.go** – Constant‑product automated market
  maker with pricing helpers.
- **loanpool.go** plus `loanpool_*` files – Collateralised lending, proposal and
  disbursement logic.
- **charity_pool.go**, **staking_node.go** and **stake_penalty.go** – Examples
  of contract controlled pools that redistribute funds or stake.

### Governance and Compliance

Protocol upgrades, on‑chain voting and regulatory hooks.

- **governance.go**, **governance_execution.go**,
  **governance_token_voting.go** and **governance_timelock.go** – Proposal
  creation, voting strategies and enactment of approved changes.
- **dao.go**, **dao_proposal.go**, **dao_quadratic_voting.go**,
  **dao_staking.go** and **dao_token.go** – Building blocks for DAOs that use
  Synnergy as their execution layer. `dao_staking.go` restricts staking to
  registered members for secure governance.
- **compliance.go**, **audit_management.go**, **audit_node.go** and
  **regulatory_management.go** – Optional KYC/AML checks, audit logging and
  jurisdiction specific rules.

### Cross‑Chain and Scaling

Interoperability modules that connect Synnergy to other chains or scaling
environments.

- **cross_chain.go**, **cross_chain_bridge.go** and
  **cross_chain_transactions.go** – Lock/mint and burn/release flows when
  moving assets between chains.
- **cross_consensus_scaling_networks.go** – Coordinate assets across disparate
  consensus systems.
- **rollups.go**, **sidechains.go** and **plasma.go** – Batch transactions or
  execute application specific chains that settle onto the main ledger.
- **gateway_node.go** and **rpc_webrtc.go** – Provide external interfaces for
  other networks and light‑clients.

### Node Roles and Infrastructure

Specialised node behaviours and orchestration utilities.

- **full_node.go**, **validator_node.go**, **master_node.go** and
  **super_node.go** – Canonical node roles participating in consensus or
  providing enhanced services.
- **authority_nodes.go**, **government_authority_node.go** and
  **regulatory_node.go** – Nodes operated by trusted parties or regulators.
- **gateway_node.go**, **mobile_node.go**, **lightning_node.go** and
  **watchtower_node.go** – Edge nodes that provide bridging, mobile access or
  off‑chain security.
- Numerous domain‑specific nodes such as **ai_enhanced_node.go**,
  **biometric_security_node.go** or **storage.go** demonstrate how services can
  be embedded directly into the network.

### AI, Monitoring and Automation

Synnergy includes experimental modules that integrate machine learning and
system monitoring.

- **ai.go**, **ai_model_management.go**, **ai_training.go** and
  **ai_secure_storage.go** – Interface with external ML services and manage model
  lifecycles on‑chain.
- **ai_inference_analysis.go**, **ai_drift_monitor.go** and
  **anomaly_detection.go** – Evaluate transactions or node metrics for abnormal
  behaviour.
- **system_health_logging.go**, **event_management.go** and
  **monitoring/incident_response_runbook.md** (in the `monitoring` directory) –
  Operational tooling for production deployments.

### Security and Access Control

Cryptography, authentication and defensive mechanisms.

- **security.go**, **firewall.go** and **zero_trust_data_channels.go** – Core
  cryptographic helpers, firewall enforcement and secure channel creation.
- **access_control.go**, **biometrics_auth.go**,
  **biometric_security_node.go** and **identity_verification.go** – User and
  node level authentication mechanisms.
- **failover_recovery.go**, **disaster_recovery_node.go** and
  **high_availability.go** – Resilience tooling for handling outages.

### Utilities and Shared Infrastructure

Miscellaneous helpers used across many modules.

- **common_structs.go**, **helpers.go** and **utility_functions.go** – Shared
  types and small utility routines.  These are the preferred dependencies for
  new modules to avoid cycles.
- **module_plugin.go** – Enables dynamic loading of optional modules.
- **data.go**, **data_operations.go** and **data_distribution.go** – Low level
  key/value helpers and content distribution utilities.
## Comprehensive Module Index

Every file under `core/` is listed below with a short description derived from its leading comments.

<!-- generated via scripts -->

- **Nodes/Nodes_Type_manual.md** – Documentation for Nodes Type manual.
- **Nodes/authority_nodes/Authority_node_typr_manual.md** – Documentation for Authority node typr manual.
- **Nodes/authority_nodes/index.go** – AuthorityNodeInterface extends NodeInterface with authority-specific actions.
- **Nodes/bank_nodes/index.go** – BankInstitutionalNode interface extends NodeInterface with signature-verified institution management (register, remove, list).
- **Nodes/consensus_specific.go** – ConsensusNodeInterface defines behaviour for nodes specialised for a specific consensus algorithm.
- **Nodes/elected_authority_node.go** – ElectedAuthorityNode provides enhanced authority capabilities subject to
- **Nodes/experimental_node.go** – ExperimentalNodeInterface extends NodeInterface with methods
- **Nodes/forensic_node.go** – TransactionLite represents the minimal transaction data required for
- **Nodes/geospatial.go** – GeospatialNodeInterface extends NodeInterface with geospatial capabilities.
- **Nodes/historical_node.go** – HistoricalNodeInterface extends NodeInterface with archival functionality.
- **Nodes/holographic_node.go** – HolographicNode provides holographic data distribution and redundancy.
- **Nodes/index.go** – NodeInterface defines minimal node behaviour independent from core types.
- **Nodes/light_node.go** – BlockHeader contains minimal block information used by light nodes.
- **Nodes/military_nodes/index.go** – WarfareNodeInterface extends the base node interface with military specific operations.
- **Nodes/molecular_node.go** – MolecularNodeInterface defines behaviour for nano-scale nodes interfacing with
- **Nodes/optimization_nodes/index.go** – Transaction represents the minimal transaction information required for
- **Nodes/optimization_nodes/optimization.go** – OptimizationNode enhances network performance through transaction ordering
- **Nodes/staking_node_interface.go** – Package Nodes provides shared node interfaces across the network.
- **Nodes/super_node.go** – SuperNodeInterface extends NodeInterface with additional services for
- **Nodes/syn845_node.go** – DebtRecord describes a debt instrument managed by SYN845 nodes.
- **Nodes/types.go** – Address mirrors the core address type without creating a dependency.
- **Nodes/watchtower/index.go** – WatchtowerNodeInterface extends NodeInterface with monitoring capabilities.
- **Nodes/witness/archival_witness_node.go** – witnessRecord is stored in the ledger for each notarised item.
- **SYN1967.go** – Implements SYN1967 functionality.
- **SYN2369.go** – SYN2369Token represents virtual world items and properties.
- **Tokens/SYN1000.go** – go:build tokens
- **Tokens/SYN3000.go** – RentalTokenMetadata defines the on-chain information for a SYN3000 token.
- **Tokens/Tokens_manual.md** – Documentation for Tokens manual.
- **Tokens/balance_table_test.go** – TestBalanceTableAddGet ensures Add and Get work as expected.
- **Tokens/base.go** – TokenID uniquely identifies a token instance within the registry.
- **Tokens/base_test.go** – helper to create distinct addresses
- **Tokens/index.go** – index.go – collection of lightweight interfaces and structs that expose token
- **Tokens/syn10.go** – SYN10Token represents CBDC with exchange rate and issuer info.
- **Tokens/syn1000_index.go** – Stablecoin defines reserve operations for SYN1000 tokens.
- **Tokens/syn1100.go** – HealthcareRecord represents encrypted healthcare data tied to an owner.
- **Tokens/syn12.go** – SYN12Metadata holds CBD Treasury Bill metadata.
- **Tokens/syn200.go** – CarbonCreditMetadata captures details about a carbon credit.
- **Tokens/syn2200.go** – go:build tokens
- **Tokens/syn2600.go** – InvestorTokenMeta defines metadata for SYN2600 investor tokens.
- **Tokens/syn2800.go** – LifePolicy defines the metadata for a SYN2800 life insurance token.
- **Tokens/syn2900.go** – InsurancePolicy mirrors the structure in the core package but avoids
- **Tokens/syn3400.go** – ForexMetadata defines pair specific information for SYN3400 tokens.
- **Tokens/syn70.go** – SYN70Asset represents a single in-game asset tracked by the SYN70 token
- **Tokens/syn845.go** – DebtMetadata stores comprehensive data about a debt instrument.
- **access_control.go** – AccessController manages role based access permissions using the ledger
- **access_control_test.go** – Implements access control test functionality.
- **account_and_balance_operations.go** – AccountManager provides helper operations for creating accounts and
- **account_and_balance_operations_test.go** – Implements account and balance operations test functionality.
- **address_from_common.go** – go:build !tokens
- **address_from_common_tokens.go** – go:build tokens
- **address_zero.go** – AddressZero represents the zero-value address (all 20 bytes set to zero).
- **ai.go** – AI module – on‑chain ML hooks for fraud detection, fee optimisation, and model
- **ai_drift_monitor.go** – ai_drift_monitor.go - simple sliding window drift detector for AI models.
- **ai_enhanced_contract.go** – AIEnhancedContract ties a smart contract to an AI model hash.
- **ai_enhanced_node.go** – AIEnhancedConfig aggregates configuration required to start an AI enhanced node.
- **ai_inference_analysis.go** – ai_inference_analysis.go - advanced AI inference and transaction analysis
- **ai_model_management.go** – GetModelListing fetches a marketplace listing by ID.
- **ai_secure_storage.go** – ai_secure_storage.go - helpers for encrypted model parameter and dataset storage.
- **ai_training.go** – TrainingJob represents a long running training process for an AI model.
- **amm.go** – amm.go – high‑level router and pricing utilities that sit on top of the
- **anomaly_detection.go** – AnomalyService provides anomaly detection helpers that integrate
- **api_node.go** – APINode exposes a HTTP API gateway backed by a network node and
- **audit_management.go** – AuditManager coordinates persistent audit logs stored on the ledger.
- **audit_node.go** – AuditNode ties together a BootstrapNode with the AuditManager.
- **audit_trail_test.go** – Implements audit trail test functionality.
- **authority_apply.go** – AuthorityApplier manages proposals for new authority nodes.
- **authority_nodes.go** – Authority Nodes governance sub‑system.
- **authority_penalty_test.go** – go:build unit
- **autonomous_agent_node.go** – AutonomousRule defines a trigger and action pair executed by the node.
- **bank_institutional_node.go** – BankInstitutionalConfig groups dependencies required by the node.
- **base_node.go** – BaseNode wraps a NodeInterface and exposes common networking behaviour.
- **bft_simulation.go** – bft_simulation.go - quick Monte Carlo estimation of consensus safety under
- **bft_simulation_test.go** – TestSimulateBFTWith verifies deterministic scenarios of the Monte Carlo
- **binary_tree_operations.go** – BinaryTree provides a simple in-memory binary search tree that persists
- **biometric_security_node.go** – BiometricSecurityNode couples a Node with biometric authentication.
- **biometrics_auth.go** – BiometricsAuth manages hashed biometric templates for addresses. It
- **blockchain_compression.go** – CompressLedger returns the gzip-compressed JSON encoding of the provided ledger.
- **blockchain_synchronization.go** – SyncManager coordinates block download and verification to keep a node's
- **bootstrap_node.go** – BootstrapNode bundles networking with optional replication to help new
- **carbon_credit_system.go** – go:build tokens
- **central_banking_node.go** – CentralBankingNode embeds a networking node and exposes hooks for monetary
- **chain_fork_manager.go** – ForkInfo summarizes a fork branch.
- **charity_pool.go** – CharityPool – 5% cut from every gas fee routed to on-chain philanthropy.
- **coin.go** – MaxSupply is the maximum number of Synthron coins that may ever exist.
- **coin_test.go** – TestBlockRewardAt verifies the halving schedule for block rewards.
- **common_structs.go** – common_structs.go – centralised struct definitions referenced across modules.
- **compliance.go** – compliance.go – Regulatory & Data‑Privacy utilities for Synnergy Network.
- **compliance_management.go** – ComplianceManager coordinates account suspensions and whitelists.
- **connection_pool.go** – manages pooled TCP connections with capacity limits and reuse.
- **connection_pool_test.go** – verifies connection reuse, pool exhaustion and proper closure against a test server.
- **consensus.go** – SynnergyConsensus – hybrid PoH + PoS sub‑blocks, aggregated under PoW main block.
- **consensus_adaptive_management.go** – ConsensusAdaptiveManager monitors recent ledger activity and stake
- **consensus_difficulty.go** – ConsensusStatus exposes high level consensus metrics such as the current
- **consensus_network_adapter.go** – nodeNetworkAdapter adapts Node to the consensus engine's minimal
- **consensus_specific_node.go** – ConsensusSpecificNode implements a node tuned for a particular consensus algorithm.
- **consensus_start.go** – Start launches the consensus engine, spinning up proposer and block
- **consensus_validator_management.go** – ValidatorInfo represents a consensus validator and its staked amount.
- **content_node.go** – ContentNetworkNode mirrors Nodes.ContentNode information for registry.
- **content_node_impl.go** – ContentNode provides specialised handling for large encrypted content.
- **content_types.go** – ContentMeta describes stored content pinned by a content node.
- **contract_management.go** – ContractManager provides administrative lifecycle operations for
- **contract_vm_test.go** – TestHeavyVMInvokeWithReceipt compiles a sample contract, deploys it to the
- **contracts.go** – Smart‑Contract Runtime & Registry for Synnergy Network.
- **contracts_opcodes.go** – Opcode constants for contract-related actions.
- **cross_chain.go** – Bridge defines parameters for a cross-chain bridge
- **cross_chain_agnostic_protocols.go** – CrossChainProtocol defines a generic cross-chain integration profile.
- **cross_chain_bridge.go** – BridgeTransfer records a cross-chain transfer locked on this chain.
- **cross_chain_connection.go** – ChainConnection represents an active cross-chain connection between
- **cross_chain_contracts.go** – ContractMapping links a local contract address to a remote chain address.
- **cross_chain_transactions.go** – CrossChainTx records a cross-chain asset movement initiated via LockAndMint
- **cross_consensus_scaling_networks.go** – CCSNetwork represents a bridge between two independent consensus systems.
- **custodial_node.go** – CustodialConfig bundles network and ledger configuration for CustodialNode.
- **dao.go** – DAO represents a decentralised autonomous organisation managed on chain.
- **dao_access_control.go** – DAORole represents a simple role within the DAO access list.
- **dao_proposal.go** – DAOProposal represents an on-chain proposal within a DAO.
- **dao_quadratic_voting.go** – QuadraticVoteRecord stores an individual quadratic vote.
- **dao_staking.go** – DAOStaking manages member-only token staking for governance participation.
- **dao_token.go** – DAO2500Membership is persisted in the ledger for SYN2500 tokens.
- **data.go** – Opcode identifiers for CDN module
- **data_distribution.go** – DataSet represents a piece of content offered through the
- **data_operations.go** – DataFeed holds structured data referenced on chain.
- **data_resource_management.go** – DataResourceManager combines simple data storage helpers with
- **defi.go** – DeFiManager exposes basic decentralised finance helpers. The
- **devnet.go** – StartDevNet spins up a number of in-memory nodes listening on sequential ports.
- **disaster_recovery_node.go** – DisasterRecoveryConfig wires networking, ledger and backup parameters for a
- **distributed_network_coordination.go** – DistributedCoordinator orchestrates coordination tasks between nodes.
- **distribution.go** – Distribution module provides utilities for bulk token transfers and airdrops.
- **dynamic_consensus_hopping.go** – dynamic_consensus_hopping.go - runtime switch between PoW, PoS and PoH
- **ecommerce.go** – Ecommerce module provides a minimal marketplace allowing addresses to list
- **education_token.go** – go:build tokens
- **elected_authority_node.go** – ElectedAuthorityNode integrates the elected authority functionality with
- **employment.go** – EmploymentContract represents a simple on-chain employment agreement.
- **energy_efficiency.go** – energy_efficiency.go - Energy usage metrics and efficiency scoring
- **energy_efficient_node.go** – EnergyEfficientNode combines a network node with energy efficiency tracking.
- **energy_tokens.go** – EnergyAsset captures metadata about a renewable energy certificate or
- **environmental_monitoring_node.go** – EnvCondition evaluates sensor bytes and returns true when the action should trigger.
- **escrow.go** – EscrowParty defines a recipient and amount in an escrow agreement.
- **event_management.go** – Event represents a ledger anchored notification emitted by various modules.
- **execution_management.go** – ExecutionManager coordinates transaction execution against the ledger
- **experimental_node.go** – ExperimentalNode provides an isolated environment for testing new
- **external_sensor.go** – Sensor represents an external data source that can be polled or triggered
- **failover_recovery.go** – FailoverNode triggers a view change via the provided ViewChanger.
- **faucet.go** – FaucetAccount is the default funding account used by the faucet.
- **fault_tolerance.go** – fault_tolerance.go – Peer health‑checking and view‑change signaling for the
- **finalization_management.go** – FinalizationManager coordinates finalization of batches, channels and blocks.
- **firewall.go** – firewall.go - simple address/token/IP firewall for Synnergy Network
- **forum.go** – ForumEngine manages on-chain discussion threads and comments.
- **full_node.go** – FullNodeMode specifies the storage strategy of a full node.
- **gaming.go** – Game represents a simple on-chain gaming session. All funds are escrowed
- **gas_table.go** – SPDX-License-Identifier: BUSL-1.1
- **gateway_node.go** – GatewayConfig bundles dependencies required for a GatewayNode.
- **geolocation_network.go** – Location represents a geographic coordinate pair in decimal degrees.
- **geospatial_node.go** – GeoRecord stores a geospatial data point.
- **governance.go** – GovProposal represents a protocol parameter change proposal
- **governance_execution.go** – governance_execution.go - helpers for executing governance-specific smart contracts
- **governance_management.go** – GovernanceContract holds metadata about a governance smart contract.
- **governance_reputation_voting.go** – AddReputation mints SYN-REP tokens to the specified address.
- **governance_timelock.go** – TimelockEntry represents a queued governance proposal with its execution time.
- **governance_token_voting.go** – TokenVote allows a token weighted vote on a proposal.
- **government_authority_node.go** – GovernmentAuthorityNodeInterface defines operations for specialised nodes
- **green_technology.go** – green_technology.go – Carbon accounting, offset scoring, certificates & throttling.
- **healthcare.go** – HealthRecord stores a pointer to an off-chain medical record.
- **helpers.go** – InitLedger initialises the global ledger using OpenLedger at the given path.
- **high_availability.go** – HighAvailability provides failover helpers and ledger snapshot management.
- **historical_node.go** – HistoricalNode maintains a complete archive of all blocks and exposes
- **holographic.go** – Simple holographic data helpers used by HolographicNode.
- **identity_verification.go** – IdentityService manages verified addresses on the ledger.
- **idwallet_registration.go** – IDRegistry manages on-chain registration of wallets that
- **immutability_enforcement.go** – ImmutabilityEnforcer ensures the genesis block cannot be altered.
- **indexing_node.go** – IndexingNode provides fast query capabilities by indexing ledger data.
- **initialization_replication.go** – InitService orchestrates ledger bootstrap via the replication subsystem
- **intangible_assets.go** – IntangibleAsset models a non-physical asset tracked on the chain.
- **integration_node.go** – IntegrationNode extends a network node with facilities to track external APIs
- **integration_registry.go** – IntegrationRegistry manages external API and blockchain connections used by
- **ip_management.go** – IPMetadata captures basic information about an IP asset.
- **ipfs.go** – IPFSService provides high level helpers for interacting with an IPFS gateway.
- **kademlia.go** – Kademlia implements a minimal in-memory Kademlia DHT used for
- **ledger.go** – NewLedger initializes a ledger, replaying an existing WAL and optionally
- **ledger_test.go** – ------------------------------------------------------------
- **lightning_node.go** – LightningChannelID uniquely identifies a payment channel.
- **liquidity_pools.go** – Constant-product Automated Market Maker (AMM) for Synnergy Network.
- **liquidity_views.go** – PoolView exposes read-only information about a liquidity pool.
- **loanpool.go** – LoanPool – treasury that accumulates protocol income (10% of each transaction fee).
- **loanpool_apply.go** – LoanPoolApply implements a simplified application process that
- **loanpool_approval_process.go** – ApprovalRequest represents an off-chain approval workflow state.
- **loanpool_config.go** – LoanPoolConfig defines configuration parameters for LoanPool.
- **loanpool_grant_disbursement.go** – Grant represents a one-off payment from the loan pool treasury.
- **loanpool_management.go** – LoanPoolManager provides administrative helpers around LoanPool.
- **loanpool_proposal.go** – CancelProposal allows the creator to cancel a pending proposal before it is executed.
- **marketplace.go** – MarketListing represents a generic item listed for sale on chain.
- **master_node.go** – MasterNode provides enhanced transaction processing, privacy services and
- **merkle_tree_operations.go** – BuildMerkleTree returns the level-by-level nodes of a Merkle tree built from
- **messages.go** – MessageQueue is a concurrency safe FIFO queue for NetworkMessage items.
- **mining_node.go** – MiningNode bundles networking, ledger and consensus components for PoW mining.
- **mobile_mining_node.go** – MiningStats aggregates runtime statistics for a MobileMiningNode.
- **mobile_node.go** – MobileNode is a lightweight node designed for mobile devices. It wraps the
- **module_guide.md** – Documentation for module guide.
- **module_plugin.go** – OpcodeModule represents an external package that wishes to register additional
- **molecular_node.go** – MolecularNode operates at the molecular level combining networking with ledger
- **monomaniac_recovery.go** – Monomaniac account recovery module provides a 3-of-4 verification
- **syn1600.go** – MusicToken captures metadata and royalty splits for a SYN1600 music asset.
- **nat_traversal.go** – NATManager manages NAT traversal using NAT-PMP or UPnP.
- **network.go** – Package core implements P2P networking for Synnergy nodes.
- **network_test.go** – Implements network test functionality.
- **node.go** – NodeAdapter adapts Node to the minimal Nodes.NodeInterface.
- **offchain_wallet.go** – OffChainWallet wraps HDWallet for offline signing and storage utilities.
- **opcode_and_gas_guide.md** – Documentation for opcode and gas guide.
- **opcode_dispatcher.go** – SPDX-License-Identifier: BUSL-1.1
- **oracle_management.go** – OracleMetrics captures performance statistics for an oracle feed.
- **orphan/orphan_node.go** – OrphanNode manages orphan blocks and recycles their transactions back
- **partitioning_and_compression.go** – HorizontalPartition splits data into fixed-size chunks. The last chunk
- **peer_management.go** – PeerManagement implements PeerManager and provides discovery,
- **plasma.go** – SimplePlasmaDeposit represents a deposit into the Plasma chain.
- **plasma_management.go** – plasma_management.go - minimal plasma chain coordinator integrated with ledger and consensus.
- **plasma_operations.go** – PlasmaBlock represents a block reference on the plasma chain.
- **polls_management.go** – Poll represents a simple community poll stored in the global KV store.
- **private_transactions.go** – Private transactions helpers provide lightweight encryption and
- **quantum_resistant_node.go** – QuantumNodeConfig aggregates configuration for a quantum-resistant node.
- **quorum_tracker.go** – QuorumTracker tracks votes from validators or token holders and checks
- **real_estate.go** – Property represents a tokenised real estate asset registered on chain.
- **regulatory_management.go** – Regulator represents an approved regulatory authority
- **regulatory_node.go** – RegulatoryConfig aggregates config required to bootstrap a regulatory node.
- **rental_management.go** – RentalAgreement holds the on-chain details for a house rental.
- **replication.go** – Replication subsystem – decentralised block propagation & on-demand sync.
- **resource_allocation_management.go** – resourceKey returns the ledger state key for an address limit.
- **resource_allocator.go** – ResourceAllocator tracks per-address gas allowances.
- **resource_management.go** – ResourceQuota tracks allowed and consumed resources for an address.
- **resource_marketplace.go** – resource_marketplace.go - simple on-chain marketplace for compute resources.
- **rollup_management.go** – rollup_management.go - Administrative functions for controlling the roll-up aggregator.
- **rollups.go** – rollups.go – Layer‑2 Roll‑up framework for Synnergy Network.
- **rpc_webrtc.go** – RPCWebRTC bridges HTTP RPC calls with WebRTC data channels.
- **security.go** – SPDX-License-Identifier: Apache-2.0
- **sharding.go** – sharding.go – Horizontal ledger partitioning with cross‑shard messaging.
- **sidechain_ops.go** – sidechain_ops.go -- management helpers for sidechain lifecycle
- **sidechains.go** – sidechains.go – Trust‑minimised side‑chain bridge & header sync layer.
- **smart_legal_contracts.go** – SmartLegalRegistry manages Ricardian contracts and signer approvals.
- **stake_penalty.go** – StakePenaltyManager provides helper methods for adjusting validator stake
- **staking_node.go** – StakingNode combines networking with staking management for PoS consensus.
- **state_channel.go** – state_channel.go – Off‑chain payment/state channels for Synnergy Network.
- **state_channel_management.go** – Enterprise-grade state channel management operations.
- **storage.go** – core/storage.go
- **super_node.go** – SuperNode provides enhanced network capabilities combining networking,
- **supply_chain.go** – SupplyItem represents a tracked asset in the supply chain.
- **swarm.go** – Swarm orchestrates multiple network nodes that share a ledger and optional
- **syn10.go** – SYN10Engine manages a CBDC token pegged to fiat currency.
- **syn1155.go** – SYN1155Token implements a multi-asset token standard supporting both
- **syn11_token.go** – SYN11Token represents Central Bank Digital Gilts.
- **syn1300.go** – SupplyChainAsset holds metadata for a tracked asset.
- **syn131_token.go** – ValuationRecord captures historical valuations for SYN131 assets.
- **syn1401.go** – InvestmentRecord stores state for a SYN1401 investment token issuance.
- **syn1500.go** – ReputationEvent records an adjustment to a user's reputation score.
- **syn1700_token.go** – EventMetadata holds details about an event for ticketing purposes.
- **syn1800.go** – go:build tokens
- **syn20.go** – SYN20Token extends core.BaseToken with pause and freeze capabilities.
- **syn2100.go** – FinancialDocument captures metadata about an invoice or other trade finance instrument.
- **syn223_token.go** – SYN223Token implements the SYN223 safe-transfer token standard.
- **syn2400.go** – DataMarketplaceToken implements the SYN2400 standard.
- **syn2500_token.go** – Syn2500Member stores metadata about DAO membership.
- **syn2700.go** – VestingEntry defines a point in time when a portion becomes available.
- **syn2900.go** – TokenInsurancePolicy represents a blockchain based insurance policy.
- **syn3000_token.go** – go:build tokens
- **syn300_token.go** – GovernanceProposal represents a proposal created using SYN300 tokens.
- **syn3100.go** – EmploymentContractMeta stores employment token contract data
- **syn3200.go** – Bill represents a bill record managed by the SYN3200 standard.
- **syn3300_token.go** – Implements syn3300 token functionality.
- **syn3500_token.go** – SYN3500Token represents a currency or stablecoin token.
- **syn3600.go** – FuturesContract defines the essential metadata of a futures contract.
- **syn3700_token.go** – IndexComponent defines a single asset within an index token.
- **syn3800.go** – GrantRecord captures metadata for a SYN3800 grant token.
- **syn3900.go** – BenefitRecord holds metadata for a government benefit token issuance.
- **syn4200_token.go** – Implements syn4200 token functionality.
- **syn4700.go** – LegalToken represents a SYN4700 legal token which is tied to a legal document
- **syn500.go** – ServiceTier defines access tiers for SYN500 utility tokens.
- **syn5000.go** – BetRecord stores betting activity for SYN5000 tokens.
- **syn5000_index.go** – GamblingToken exposes methods of the SYN5000 token.
- **syn700.go** – SYN700Token implements the Intellectual Property token standard.
- **syn721_token.go** – SYN721Metadata holds metadata for a single NFT token
- **syn800_token.go** – AssetMetadata describes a real world asset backing a SYN800 token.
- **system_health_logging.go** – Metrics captures a snapshot of network and node health statistics.
- **tangible_assets.go** – TangibleAsset represents a tokenised physical asset tracked on chain.
- **time_locked_node.go** – TimeLockRecord represents a pending transfer with a release time.
- **token_management.go** – token_management.go - helper functions for working with tokens.
- **token_management_syn1000.go** – go:build tokens
- **token_syn130.go** – Syn130SaleRecord tracks sale price history for tangible assets
- **token_syn4900.go** – AgriculturalAsset holds detailed metadata for SYN4900 tokens.
- **token_syn600.go** – SYN600Token implements reward token mechanics including staking and
- **tokens.go** – tokens.go - central token registry and basic token implementation.
- **tokens_syn1000.go** – go:build tokens
- **tokens_syn1000_helpers.go** – go:build tokens
- **tokens_syn1000_opcodes.go** – Implements tokens syn1000 opcodes functionality.
- **tokens_syn1200.go** – Implements tokens syn1200 functionality.
- **tokens_syn900.go** – IdentityDetails stores personal information associated with an identity token.
- **tokens_syn900_index.go** – IdentityTokenAPI defines the exposed operations for the SYN900 token.
- **transaction_distribution.go** – TxDistributor splits transaction fees between network stakeholders.
- **transactionreversal.go** – ReverseTransactionFeeBps defines the fee charged on a reversal (2.5%).
- **transactions.go** – go:build tokens
- **tx_types.go** – go:build tokens
- **txpool_addtx.go** – go:build ignore
- **txpool_snapshot.go** – go:build !tokens
- **txpool_stub.go** – go:build ignore
- **user_feedback_system.go** – user_feedback_system.go -- user feedback collection and reward engine
- **utility_functions.go** – Short returns a shortened hex version of the hash (e.g. first 4 + last 4).
- **validator_node.go** – ValidatorNode bundles networking, ledger access and consensus participation.
- **virtual_machine.go** – Synnergy Network – virtual_machine.go
- **vm_sandbox_management.go** – SandboxInfo holds runtime limits and state for a single sandboxed contract
- **wallet.go** – Wallet implementation for the Synnergy Network blockchain.
- **wallet_management.go** – WalletManager wraps Ledger and HDWallet helpers to perform high level wallet operations.
- **warehouse.go** – WarehouseItem represents an item stored on-chain for supply chain tracking.
- **warfare_node.go** – Enforces signed command envelopes, records logistics/tactical updates and exposes replayable event streams for CLI/UI subscribers.
- **watchtower_node.go** – Emits start/stop/fork alerts, streams periodic health metrics and raises integrity sweeps against observed nodes.
- **workflow_integrations.go** – Workflow represents a sequence of opcode names executed in order.
- **zero_trust_data_channels.go** – Coordinates encrypted channels with participant governance, retention policies, key rotation and event feeds.
- **zkp_node.go** – ZKPNodeConfig aggregates network and ledger configuration for a zero-knowledge proof node.


## Working With Modules

1. **Keep dependencies minimal.**  Import from `common_structs.go` and other
   utilities rather than coupling modules directly.
2. **Document behaviour.**  Exported types and functions should have GoDoc style
   comments.  This guide summarises intent but the source is authoritative.
3. **Add tests.**  Create or update a matching file under
   `synnergy-network/tests` and ensure `scripts/run_tests.sh` passes.
4. **Update documentation.**  If your change affects public APIs or gas costs,
   also update `docs/developer-guide.md` or `core/opcode_and_gas_guide.md` as
   appropriate.

## Testing and Validation

The repository provides a helper script to run the entire test suite:

```bash
scripts/run_tests.sh
```

This script invokes `go test` across the workspace and performs additional lint
checks.  Run it before submitting a pull request.

## Additional Resources

- `docs/developer-guide.md` – broader project conventions and workflow.
- `core/opcode_and_gas_guide.md` – reference for virtual machine opcodes and gas
  pricing.
- `synnergy-network/tests` – examples of how modules are exercised in practice.

## Contributing

Changes to modules should be accompanied by clear commit messages and tests.
Open a pull request describing the rationale and reference this guide when
adding new modules or altering existing ones.
