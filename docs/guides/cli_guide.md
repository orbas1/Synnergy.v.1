
# Synnergy Command Line Guide

This short guide summarises the CLI entry points found in `cmd/cli`.  Each Go file wires a set of commands using the [Cobra](https://github.com/spf13/cobra) framework.  Commands are grouped by module and can be imported individually into the root program.

Most commands require environment variables or a configuration file to be present.  Refer to inline comments for a full list of options.

## Available Command Groups

The following command groups expose the same functionality available in the core modules. Each can be mounted on a root [`cobra.Command`](https://github.com/spf13/cobra).

- **ai** – Tools for publishing ML models and running anomaly detection jobs via gRPC to the AI service. Useful for training pipelines and on‑chain inference.
- **ai_contract** – Deploy and interact with AI-enhanced smart contracts.
- **ai-train** – Manage on-chain AI model training jobs.
- **ai_mgmt** – Manage marketplace listings for AI models.
- **ai_infer** – Advanced inference and batch analysis utilities.
- **amm** – Swap tokens and manage liquidity pools. Includes helpers to quote routes and add/remove liquidity.
- **authority_node** – Register new validators, vote on authority proposals and list the active electorate.
- **access** – Manage role based access permissions.
- **authority_apply** – Submit and vote on authority node applications.
- **charity_pool** – Query the community charity fund and trigger payouts for the current cycle.
- **charity_mgmt** – Donate to and withdraw from the charity pool.
- **identity** – Register and verify user identities.
- **coin** – Mint the base coin, transfer balances and inspect supply metrics.
 - **compliance_management** – Manage suspensions and whitelists for addresses.
- **compliance** – Run KYC/AML checks on addresses and export audit reports.
- **audit** – Manage on-chain audit logs.
- **consensus_hop** – Switch between consensus modes based on network metrics.
- **adaptive** – Manage adaptive consensus weights.
- **stake** – Adjust validator stakes and record penalties.
- **contracts** – Deploy, upgrade and invoke smart contracts stored on chain.
- **contractops** – Administrative operations like pausing and transferring ownership.
- **cross_chain** – Bridge assets to or from other chains using lock and release commands.
- **ccsn** – Manage cross-consensus scaling networks.
- **xcontract** – Register and query cross-chain contract mappings.
- **cross_tx** – Execute cross-chain lock/mint and burn/release transfers.
- **cross_chain_connection** – Create and monitor links between chains.
- **cross_chain_agnostic_protocols** – Register cross-chain protocols.
- **cross_chain_bridge** – Manage cross-chain token transfers.
- **data** – Inspect raw key/value pairs in the underlying data store for debugging.
- **fork** – Manage chain forks and resolve competing branches.
- **messages** – Queue, process and broadcast network messages.
- **partition** – Partition and compress data sets.
- **data_ops** – Manage and transform on-chain data feeds.
- **anomaly_detection** – Run anomaly analysis on transactions and list flagged hashes.
- **resource** – Manage stored data and VM gas allocations.
- **immutability** – Verify the chain against the genesis block.
- **fault_tolerance** – Inject faults, simulate network partitions and test recovery procedures.
- **plasma** – Manage deposits and exits on the plasma bridge.
- **resource_allocation** – Manage per-contract gas limits.
- **failover** – Manage ledger snapshots and coordinate recovery.
- **employment** – Manage on-chain employment contracts and salaries.
- **governance** – Create proposals, cast votes and check DAO parameters.
- **token_vote** – Cast token weighted governance votes.
- **qvote** – Submit quadratic votes and view weighted results.
- **polls_management** – Create and vote on community polls.
- **governance_management** – Manage governance contracts on chain.
- **reputation_voting** – Reputation weighted governance commands.
- **timelock** – Manage delayed proposal execution.
- **dao** – Manage DAO creation and membership.
- **green_technology** – View energy metrics and toggle any experimental sustainability features.
- **resource_management** – Manage VM resource quotas and usage.
- **carbon_credit_system** – Manage carbon offset projects and credits.
- **energy_efficiency** – Record transaction counts and compute efficiency scores.
- **ledger** – Inspect blocks, query balances and perform administrative token operations via the ledger daemon.
- **account** – manage accounts and balances
- **network** – Manage peer connections and print networking statistics.
- **bootstrap** – Start a dedicated bootstrap node for new peers.
- **connpool** – Manage reusable outbound connections.
- **nat** – Manage router port mappings for inbound connectivity.
- **peer** – Discover peers, connect to them and advertise this node.
 - **replication** – Trigger snapshot creation and replicate the ledger to new nodes.
 - **high_availability** – Manage standby nodes and promote backups.
 - **rollups** – Create rollup batches or inspect existing ones.
- **plasma** – Deposit into and withdraw from the Plasma chain.
- **replication** – Trigger snapshot creation and replicate the ledger to new nodes.
- **coordination** – Coordinate distributed nodes and broadcast ledger state.
 - **rollups** – Create rollup batches, inspect existing ones and control the aggregator state.
- **initrep** – Bootstrap a ledger via peer replication.
- **synchronization** – Coordinate block download and verification.
- **rollups** – Create rollup batches or inspect existing ones.
- **compression** – Save and load compressed ledger snapshots.
- **security** – Key generation, signing utilities and password helpers.
- **firewall** – Manage address, token and IP block lists.
- **biometrics** – Manage biometric authentication templates.
- **sharding** – Migrate data between shards and check shard status.
 - **sidechain** – Launch, manage and interact with remote side‑chain nodes.
- **state_channel** – Open, close and settle payment channels.
- **loanpool** – Submit loan proposals and manage disbursements.
- **grant** – Create and release grants from the loan pool.
- **plasma** – Manage plasma deposits and submit block roots.
- **state_channel_mgmt** – Pause, resume and force-close channels.
- **zero_trust_data_channels** – Manage encrypted data channels with escrow.
- **swarm** – Manage groups of nodes running together.
- **storage** – Configure the backing key/value store and inspect content.
- **legal** – Manage Ricardian contracts and sign agreements.
- **resource** – Manage compute resource rentals.
- **staking** – Stake and unstake tokens for DAO governance.
- **dao_access** – Manage DAO membership roles.
- **sensor** – Manage external sensor inputs and webhooks.
- **real_estate** – Manage tokenised real estate assets.
- **healthcare** – Manage healthcare records and permissions.
- **warehouse** – Manage on-chain inventory records.
- **tokens** – Register new token types and move balances between accounts.
- **defi** – Insurance policies and other DeFi utilities.
- **event_management** – Emit and query custom events stored on chain.
- **gaming** – Manage simple on-chain games.
- **transactions** – Build raw transactions, sign them and broadcast to the network.
- **private_tx** – Encrypt data and submit private transactions.
- **transactionreversal** – Reverse confirmed transactions with authority approval.
- **transaction_distribution** – Distribute transaction fees between stakeholders.
- **utility_functions** – Miscellaneous helpers shared by other command groups.
- **geolocation** – Manage node geolocation data.
- **distribution** – Bulk token distribution and airdrop helpers.
- **finalization_management** – Finalize blocks, batches and channels.
- **quorum** – Manage quorum trackers for proposals or validation.
- **virtual_machine** – Execute scripts in the built‑in VM for testing.
- **sandbox** – Manage VM sandbox environments.
- **supply** – Manage supply chain records.
- **wallet** – Generate mnemonics, derive addresses and sign transactions.
- **execution** – Manage block execution and transaction pipelines.
- **binarytree** – Manage ledger-backed binary search trees.
- **regulator** – Manage on-chain regulators and rule checks.
- **feedback** – Submit and review on‑chain user feedback.
- **system_health** – Monitor node metrics and emit log entries.
- **idwallet** – Register ID-token wallets and verify status.
- **offwallet** – Offline wallet utilities.
- **recovery** – Manage account recovery registration and execution.
- **workflow** – Build on-chain workflows using triggers and webhooks.
- **wallet_mgmt** – Manage wallets and submit ledger transfers.
- **devnet** – Launch a local multi-node developer network.
- **testnet** – Start an ephemeral test network from a YAML config.
- **faucet** – Dispense test funds with rate limiting.

# Newly Added Command Groups

The CLI has grown considerably. The modules below expose functionality found in
recently added Go files under `cmd/cli`.

- **agriculture** – Manage agricultural campaigns and crop investments.
- **ai_enhanced_node** – Start or stop an AI-assisted node with volume prediction.
- **api_node** – Run the API gateway service.
- **archival_witness_node** – Notarize transactions and witness blocks.
- **audit_node** – Perform audit logging on-chain.
- **autonomous_agent_node** – Launch autonomous agent services.
- **bank_institutional_node** – Operate regulated banking nodes.
- **biometric_security_node** – Nodes using biometric admin controls.
- **central_banking_node** – Issue and manage central bank currency.
- **charity_token** – Donate to charity campaigns and query progress.
- **consensus_specific_node** – Nodes tuned for a particular consensus mode.
- **content_node** – Serve decentralised content from storage.
- **custodial_node** – Provide custodial wallet services.
- **dao_token** – Manage DAO membership tokens.
- **disaster_recovery_node** – Coordinate disaster recovery replicas.
- **elected_authority_node** – Operate an elected authority node.
- **employment_token** – Issue employment contracts and payroll tokens.
- **energy_efficient_node** – Node implementing energy saving features.
- **energy_tokens** – Register energy credits and record usage.
- **environmental_monitoring_node** – Aggregate external sensor data.
- **event_ticket** – Sell and transfer event tickets.
- **experimental_node** – Run experimental consensus or features.
- **forensic_node** – Inspect logs and network traces.
- **forex_token** – Manage foreign exchange pegged tokens.
- **full_node** – Operate a standard validating node.
- **gateway_node** – Bridge data sources and chain connections.
- **geospatial_node** – Track location based data.
- **grant_tokens** – Manage grant disbursement tokens.
- **historical_node** – Serve historical chain data.
- **holographic_node** – Node for holographic/VR content.
- **identity_token** – Issue identity verified tokens.
- **index** – Manage search indexes for ledger data.
- **indexing_node** – Node specialised for indexing the chain.
- **insurance_token** – Manage insurance policies as tokens.
- **integration_node** – Integration test harness node.
- **iptoken** – Register intellectual property tokens.
- **legal_token** – Track legal documents on chain.
- **life_insurance** – Issue life insurance policies.
- **lightning_node** – Operate a Lightning Network bridge.
- **master_node** – Configure and monitor masternodes.
- **military_node** – Run a military grade node.
- **mining_node** – Execute mining operations.
- **mobile_mining_node** – Mining node for mobile devices.
- **mobile_node** – Lightweight mobile node.
- **molecular_node** – Node used for molecular simulations.
- **optimization** – Optimisation utilities for tokens.
- **orphan_node** – Maintain orphaned blocks or offline nodes.
- **pension_tokens** – Manage pension fund tokens.
- **quantum_resistant_node** – Use quantum resistant cryptography.
- **regulatory_node** – Regulatory oversight commands.
- **rental_token** – Tokenised rentals for assets.
- **reputation_tokens** – Issue reputation based tokens.
- **staking_node** – Dedicated staking node commands.
- **super_node** – High performance super node.
- **syn10** – Manage SYN10 governance tokens.
- **syn1000** – Stablecoin token operations.
- **syn11** – Legacy SYN11 token utilities.
- **syn1100** – Healthcare record tokens.
- **syn1155** – ERC1155 style multi‑token commands.
- **syn1200** – Advanced NFT management.
- **syn130** – Commodity token operations.
- **syn1300_token** – SYN1300 asset token utilities.
- **syn131** – Multisig token management.
- **syn1401** – Regulated security tokens.
- **syn1600** – Cross jurisdiction asset tokens.
- **syn1800** – Derivatives token utilities.
- **syn1900** – Weighted index token commands.
- **syn1967** – Historical asset tokens.
- **syn200** – Simple utility token management.
- **syn2100** – Bond issuance tokens.
- **syn2200** – Real estate backed tokens.
- **syn223** – Security token operations.
- **syn2400** – Perpetual futures tokens.
- **syn300** – Charity funding token standard.
- **syn3200** – Prediction market tokens.
- **syn3300** – Regulated commodity tokens.
- **syn3500** – Fiat currency tokens.
- **syn500** – Reward point token standard.
- **syn5000** – Fixed income token standard.
- **syn600** – Loyalty points token commands.
- **syn70** – Generic token operations.
- **syn721** – NFT standard token commands.
- **syn800** – Option contract token utilities.
- **syn845** – Debt instrument tokens.
- **time_locked_node** – Nodes enforcing time locked execution.
- **tokens** – Inspect and administer any token type.
- **validator_node** – Manage validator registration.
- **watchtower_node** – Monitor network health and detect forks.
- **zkp_node** – Zero knowledge proof node operations.


To use these groups, import the corresponding command constructor (e.g. `ledger.NewLedgerCommand()`) in your main program and attach it to the root `cobra.Command`.

If you want to enable **all** CLI modules with a single call, use `cli.RegisterRoutes(rootCmd)` from the `cli` package. This helper mounts every exported command group so routes can be invoked like:

```bash
$ synnergy ~network ~start
```

## Command Reference

The sections below list each root command and its available sub‑commands. Every
command maps directly to logic in `synnergy-network/core` and can be composed as
needed in custom tooling.

### ai

| Sub-command | Description |
|-------------|-------------|
| `predict <tx.json>` | Predict fraud probability for a transaction. |
| `optimise <stats.json>` | Suggest an optimal base fee for the next block. |
| `volume <stats.json>` | Forecast upcoming transaction volume. |
| `publish <cid>` | Publish a model hash with optional royalty basis points. |
| `fetch <model-hash>` | Fetch metadata for a published model. |
| `list <price> <cid>` | Create a marketplace listing for a model. |
| `buy <listing-id> <buyer-addr>` | Buy a listed model with escrow. |
| `rent <listing-id> <renter-addr> <hours>` | Rent a model for a period of time. |
| `release <escrow-id>` | Release funds from escrow to the seller. |

### ai-train

| Sub-command | Description |
|-------------|-------------|
| `start <datasetCID> <modelCID>` | Begin a new training job. |
| `status <jobID>` | Display status for a training job. |
| `list` | List all active training jobs. |
| `cancel <jobID>` | Cancel a running job. |
### ai_mgmt

| Sub-command | Description |
|-------------|-------------|
| `get <id>` | Fetch a marketplace listing. |
| `ls` | List all AI model listings. |
| `update <id> <price>` | Update the price of your listing. |
| `remove <id>` | Remove a listing you own. |


| Sub-command | Description |
|-------------|-------------|
| `ai_infer run <model-hash> <input-file>` | Execute model inference on input data. |
| `ai_infer analyse <txs.json>` | Analyse a batch of transactions for fraud risk. |

### amm

| Sub-command | Description |
|-------------|-------------|
| `init <fixture-file>` | Initialise pools from a JSON fixture. |
| `swap <tokenIn> <amtIn> <tokenOut> <minOut> [trader]` | Swap tokens via the router. |
| `add <poolID> <provider> <amtA> <amtB>` | Add liquidity to a pool. |
| `remove <poolID> <provider> <lpTokens>` | Remove liquidity from a pool. |
| `quote <tokenIn> <amtIn> <tokenOut>` | Estimate output amount without executing. |
| `pairs` | List all tradable token pairs. |

### authority_node

| Sub-command | Description |
|-------------|-------------|
| `register <addr> <role>` | Submit a new authority-node candidate. |
| `vote <voterAddr> <candidateAddr>` | Cast a vote for a candidate. |
| `electorate <size>` | Sample a weighted electorate of active nodes. |
| `is <addr>` | Check if an address is an active authority node. |
| `info <addr>` | Display details for an authority node. |
| `list` | List authority nodes. |
| `deregister <addr>` | Remove an authority node and its votes. |

### access

| Sub-command | Description |
|-------------|-------------|
| `grant <role> <addr>` | Grant a role to an address. |
| `revoke <role> <addr>` | Revoke a role from an address. |
| `check <role> <addr>` | Check whether an address has a role. |
| `list <addr>` | List all roles assigned to an address. |
### authority_apply

| Sub-command | Description |
|-------------|-------------|
| `submit <candidate> <role> <desc>` | Submit an authority node application. |
| `vote <voter> <id>` | Vote on an application. Use `--approve=false` to reject. |
| `finalize <id>` | Finalize and register the node if the vote passed. |
| `tick` | Check all pending applications for expiry. |
| `get <id>` | Display an application by ID. |
| `list` | List all applications. |

### charity_pool

| Sub-command | Description |
|-------------|-------------|
| `register <addr> <category> <name>` | Register a charity with the pool. |
| `vote <voterAddr> <charityAddr>` | Vote for a charity during the cycle. |
| `tick [timestamp]` | Manually trigger pool cron tasks. |
| `registration <addr> [cycle]` | Show registration info for a charity. |
| `winners [cycle]` | List winning charities for a cycle. |

### charity_mgmt

| Sub-command | Description |
|-------------|-------------|
| `donate <from> <amt>` | Donate tokens to the charity pool. |
| `withdraw <to> <amt>` | Withdraw internal charity funds. |
| `balances` | Show pool and internal balances. |

### coin

| Sub-command | Description |
|-------------|-------------|
| `mint <addr> <amt>` | Mint the base SYNN coin. |
| `supply` | Display total supply. |
| `balance <addr>` | Query balance for an address. |
| `transfer <from> <to> <amt>` | Transfer SYNN between accounts. |
| `burn <addr> <amt>` | Burn SYNN from an address. |

### compliance

| Sub-command | Description |
|-------------|-------------|
| `validate <kyc.json>` | Validate and store a KYC document commitment. |
| `erase <address>` | Remove a user's KYC data. |
| `fraud <address> <severity>` | Record a fraud signal. |
| `risk <address>` | Retrieve accumulated fraud risk score. |
| `audit <address>` | Display the audit trail for an address. |
| `monitor <tx.json> <threshold>` | Run anomaly detection on a transaction. |
| `verifyzkp <blob.bin> <commitmentHex> <proofHex>` | Verify a zero‑knowledge proof. |

### audit

| Sub-command | Description |
|-------------|-------------|
| `log <addr> <event> [meta.json]` | Record an audit event. |
| `list <addr>` | List audit events for an address. |
### compliance_management

| Sub-command | Description |
|-------------|-------------|
| `suspend <addr>` | Suspend an address from transfers. |
| `resume <addr>` | Lift an address suspension. |
| `whitelist <addr>` | Add an address to the whitelist. |
| `unwhitelist <addr>` | Remove an address from the whitelist. |
| `status <addr>` | Show suspension and whitelist status. |
| `review <tx.json>` | Check a transaction before broadcast. |
### anomaly_detection

| Sub-command | Description |
|-------------|-------------|
| `analyze <tx.json>` | Run anomaly detection on a transaction. |
| `list` | List flagged transactions. |

### consensus

| Sub-command | Description |
|-------------|-------------|
| `start` | Launch the consensus engine. |
| `stop` | Gracefully stop the consensus service. |
| `info` | Show consensus height and running status. |
| `weights <demand> <stake>` | Calculate dynamic consensus weights. |
| `threshold <demand> <stake>` | Compute the consensus switch threshold. |
| `set-weight-config <alpha> <beta> <gamma> <dmax> <smax>` | Update weight coefficients. |
| `get-weight-config` | Display current weight configuration. |

### consensus_hop

| Sub-command | Description |
|-------------|-------------|
| `eval <demand> <stake>` | Evaluate metrics and possibly switch consensus mode. |
| `mode` | Show the current consensus mode. |
### adaptive

| Sub-command | Description |
|-------------|-------------|
| `metrics` | Show current demand and stake levels. |
| `adjust` | Recompute consensus weights. |
| `set-config <alpha> <beta> <gamma> <dmax> <smax>` | Update weighting coefficients. |

### stake

| Sub-command | Description |
|-------------|-------------|
| `adjust <addr> <delta>` | Increase or decrease stake for a validator. |
| `penalize <addr> <points> [reason]` | Record penalty points against a validator. |
| `info <addr>` | Display stake and penalty totals. |

### contracts

| Sub-command | Description |
|-------------|-------------|
| `compile <src.wat|src.wasm>` | Compile WAT or WASM to deterministic bytecode. |
| `deploy --wasm <path> [--ric <file>] [--gas <limit>]` | Deploy compiled WASM. |
| `invoke <address>` | Invoke a contract method. |
| `list` | List deployed contracts. |
| `info <address>` | Show Ricardian manifest for a contract. |

### contractops

| Sub-command | Description |
|-------------|-------------|
| `transfer <addr> <newOwner>` | Transfer contract ownership. |
| `pause <addr>` | Pause contract execution. |
| `resume <addr>` | Resume a paused contract. |
| `upgrade <addr> <wasm>` | Replace contract bytecode. |
| `info <addr>` | Display owner and paused status. |

### cross_chain

| Sub-command | Description |
|-------------|-------------|
| `register <source_chain> <target_chain> <relayer_addr>` | Register a bridge. |
| `list` | List registered bridges. |
| `get <bridge_id>` | Retrieve a bridge configuration. |
| `authorize <relayer_addr>` | Whitelist a relayer address. |
| `revoke <relayer_addr>` | Remove a relayer from the whitelist. |
### cross_chain_agnostic_protocols

| Sub-command | Description |
|-------------|-------------|
| `register <name>` | Register a new protocol definition. |
| `list` | List registered protocols. |
| `get <id>` | Retrieve a protocol configuration. |


### cross_chain_bridge

| Sub-command | Description |
|-------------|-------------|
| `deposit <bridge_id> <from> <to> <amount> [tokenID]` | Lock assets for bridging. |
| `claim <transfer_id> <proof.json>` | Release assets using a proof. |
| `get <id>` | Show a transfer record. |
| `list` | List all transfers. |

### cross_chain_connection

| Sub-command | Description |
|-------------|-------------|
| `open <local_chain> <remote_chain>` | Establish a new connection. |
| `close <connection_id>` | Terminate a connection. |
| `get <connection_id>` | Retrieve connection details. |
| `list` | List active and historic connections. |

### cross_tx

| Sub-command | Description |
|-------------|-------------|
| `lockmint <bridge_id> <asset_id> <amount> <proof>` | Lock native assets and mint wrapped tokens. |
| `burnrelease <bridge_id> <to> <asset_id> <amount>` | Burn wrapped tokens and release native assets. |
| `list` | List cross-chain transfer records. |
| `get <tx_id>` | Retrieve a cross-chain transfer by ID. |

### xcontract

| Sub-command | Description |
|-------------|-------------|
| `register <local_addr> <remote_chain> <remote_addr>` | Register a contract mapping. |
| `list` | List registered mappings. |
| `get <local_addr>` | Retrieve mapping info. |
| `remove <local_addr>` | Delete a mapping. |

### ccsn

| Sub-command | Description |
|-------------|-------------|
| `register <source_consensus> <target_consensus>` | Register a cross-consensus network. |
| `list` | List configured networks. |
| `get <network_id>` | Retrieve a network configuration. |

### data

**Node operations**

| Sub-command | Description |
|-------------|-------------|
| `node register <address> <host:port> <capacityMB>` | Register a CDN node. |
| `node list` | List CDN nodes. |

**Asset operations**

| Sub-command | Description |
|-------------|-------------|
| `asset upload <filePath>` | Upload and pin an asset. |
| `asset retrieve <cid> [output]` | Retrieve an asset by CID. |

**Oracle feeds**

| Sub-command | Description |
|-------------|-------------|
| `oracle register <source>` | Register a new oracle feed. |
| `oracle push <oracleID> <value>` | Push a value to an oracle feed. |
| `oracle query <oracleID>` | Query the latest oracle value. |
| `oracle list` | List registered oracles. |

### messages

| Sub-command | Description |
|-------------|-------------|
| `enqueue <src> <dst> <topic> <type> <payload>` | Queue a message for processing. |
| `process` | Process the next queued message using the ledger. |
| `broadcast` | Broadcast the next message to peers. |
### distribution

| Sub-command | Description |
|-------------|-------------|
| `distribution create <owner> <cid> <price>` | Register a dataset for sale. |
| `distribution buy <datasetID> <buyer>` | Purchase dataset access. |
| `distribution info <datasetID>` | Show dataset metadata. |
| `distribution list` | List all datasets. |
**Oracle management**

| Sub-command | Description |
|-------------|-------------|
| `oracle_mgmt metrics <oracleID>` | Show performance metrics. |
| `oracle_mgmt request <oracleID>` | Fetch value and record latency. |
| `oracle_mgmt sync <oracleID>` | Sync local oracle data. |
| `oracle_mgmt update <oracleID> <source>` | Update oracle source. |
| `oracle_mgmt remove <oracleID>` | Remove oracle configuration. |
### data_ops

| Sub-command | Description |
|-------------|-------------|
| `create <desc> <v1,v2,..>` | Create a new data feed. |
| `query <id>` | Query a feed and print JSON. |
| `normalize <id>` | Normalize feed values. |
| `impute <id>` | Impute missing values using the mean. |
### resource

| Sub-command | Description |
|-------------|-------------|
| `store <owner> <key> <file> <gas>` | Store data and set a gas limit. |
| `load <owner> <key> [out|-]` | Load data for a key. |
| `delete <owner> <key>` | Remove stored data and reset the limit. |

### fault_tolerance
- **employment** – Manage on-chain employment contracts and salaries.

| Sub-command | Description |
|-------------|-------------|
| `snapshot` | Dump current peer statistics. |
| `add-peer <addr>` | Add a peer to the health-checker set. |
| `rm-peer <addr|id>` | Remove a peer from the set. |
| `view-change` | Force a leader rotation. |
| `backup` | Create a ledger backup snapshot. |
| `restore <file>` | Restore ledger state from a snapshot. |
| `failover <addr>` | Force failover of a node. |
| `predict <addr>` | Predict failure probability for a node. |

### resource_allocation

| Sub-command | Description |
|-------------|-------------|
| `set <addr> --limit=<n>` | Set gas limit for an address. |
| `get <addr>` | Display current gas limit. |
| `list` | List limits for all addresses. |
| `consume <addr> --amt=<n>` | Deduct gas from an address limit. |
| `transfer <from> <to> --amt=<n>` | Transfer limit between addresses. |
### failover

| Sub-command | Description |
|-------------|-------------|
| `backup <path>` | Create a ledger snapshot. |
| `restore <file>` | Restore ledger state from a snapshot file. |
| `verify <file>` | Verify a snapshot against the current ledger. |
| `node [reason]` | Trigger a view change. |

### governance

| Sub-command | Description |
|-------------|-------------|
| `propose` | Submit a new governance proposal. |
| `vote <proposal-id>` | Cast a vote on a proposal. |
| `execute <proposal-id>` | Execute a proposal after the deadline. |
| `get <proposal-id>` | Display a single proposal. |
| `list` | List all proposals. |

### token_vote

| Sub-command | Description |
|-------------|-------------|
| `cast <proposal-id> <voter> <token-id> <amount> [approve]` | Cast a token weighted vote on a proposal. |

### qvote

| Sub-command | Description |
|-------------|-------------|
| `cast` | Submit a quadratic vote on a proposal. |
| `results <proposal-id>` | Display aggregated quadratic weights. |
### dao_access

| Sub-command | Description |
|-------------|-------------|
| `add <addr> <role>` | Add a DAO member with role `member` or `admin`. |
| `remove <addr>` | Remove a DAO member. |
| `role <addr>` | Display the member role. |
| `list` | List all DAO members. |
### polls_management

| Sub-command | Description |
|-------------|-------------|
| `create` | Create a new poll. |
| `vote <id>` | Cast a vote on a poll. |
| `close <id>` | Close a poll immediately. |
| `get <id>` | Display a poll. |
| `list` | List existing polls. |
### governance_management

| Sub-command | Description |
|-------------|-------------|
| `contract:add <addr> <name>` | Register a governance contract. |
| `contract:enable <addr>` | Enable a contract for voting. |
| `contract:disable <addr>` | Disable a contract. |
| `contract:get <addr>` | Display contract information. |
| `contract:list` | List registered contracts. |
| `contract:rm <addr>` | Remove a contract from the registry. |
### reputation_voting

| Sub-command | Description |
|-------------|-------------|
| `propose` | Submit a new reputation proposal. |
| `vote <proposal-id>` | Cast a weighted vote using SYN-REP. |
| `execute <proposal-id>` | Execute a reputation proposal. |
| `get <proposal-id>` | Display a reputation proposal. |
| `list` | List all reputation proposals. |
| `balance <addr>` | Show reputation balance. |
### timelock

| Sub-command | Description |
|-------------|-------------|
| `queue <proposal-id>` | Queue a proposal with a delay. |
| `cancel <proposal-id>` | Remove a queued proposal. |
| `execute` | Execute all due proposals. |
| `list` | List queued proposals. |
### dao

| Sub-command | Description |
|-------------|-------------|
| `create <name> <creator>` | Create a new DAO. |
| `join <dao-id> <addr>` | Join an existing DAO. |
| `leave <dao-id> <addr>` | Leave a DAO. |
| `info <dao-id>` | Display DAO information. |
| `list` | List all DAOs. |

### green_technology

| Sub-command | Description |
|-------------|-------------|
| `usage <validator-addr>` | Record energy and carbon usage for a validator. |
| `offset <validator-addr>` | Record carbon offset credits. |
| `certify` | Recompute certificates immediately. |
| `cert <validator-addr>` | Show the sustainability certificate. |
| `throttle <validator-addr>` | Check if a validator should be throttled. |
| `list` | List certificates for all validators. |

### resource_management

| Sub-command | Description |
|-------------|-------------|
| `set <addr> <cpu> <mem> <store>` | Set resource quota for an address. |
| `show <addr>` | Display quota and current usage. |
| `charge <addr> <cpu> <mem> <store>` | Charge and record consumed resources. |

### carbon_credit_system

| Sub-command | Description |
|-------------|-------------|
| `register <owner> <name> <total>` | Register a carbon offset project. |
| `issue <projectID> <to> <amount>` | Issue credits from a project. |
| `retire <holder> <amount>` | Burn carbon credits permanently. |
| `info <projectID>` | Show details of a project. |
| `list` | List all projects. |
### energy_efficiency

| Sub-command | Description |
|-------------|-------------|
| `record <validator-addr>` | Record processed transactions and energy use. |
| `efficiency <validator-addr>` | Show tx per kWh for a validator. |
| `network` | Display the network average efficiency. |

### ledger

| Sub-command | Description |
|-------------|-------------|
| `head` | Show chain height and latest block hash. |
| `block <height>` | Fetch a block by height. |
| `balance <addr>` | Display token balances of an address. |
| `utxo <addr>` | List UTXOs for an address. |
| `pool` | List mem-pool transactions. |
| `mint <addr>` | Mint tokens to an address. |
| `transfer <from> <to>` | Transfer tokens between addresses. |

### fork

| Sub-command | Description |
|-------------|-------------|
| `list` | Show currently tracked forks. |
| `resolve` | Resolve forks extending the tip. |
### account

| Sub-command | Description |
|-------------|-------------|
| `create <addr>` | Create a new account. |
| `delete <addr>` | Delete an account. |
| `balance <addr>` | Show account balance. |
| `transfer` | Transfer between accounts. |

### liquidity_pools

| Sub-command | Description |
|-------------|-------------|
| `create <tokenA> <tokenB> [feeBps]` | Create a new liquidity pool. |
| `add <poolID> <provider> <amtA> <amtB>` | Add liquidity to a pool. |
| `swap <poolID> <trader> <tokenIn> <amtIn> <minOut>` | Swap tokens within a pool. |
| `remove <poolID> <provider> <lpTokens>` | Remove liquidity from a pool. |
| `info <poolID>` | Show pool state. |
| `list` | List all pools. |

### loanpool

| Sub-command | Description |
|-------------|-------------|
| `submit <creator> <recipient> <type> <amount> <desc>` | Submit a loan proposal. |
| `vote <voter> <id>` | Vote on a proposal. |
| `disburse <id>` | Disburse an approved loan. |
| `tick` | Process proposals and update cycles. |
| `get <id>` | Display a single proposal. |
| `list` | List proposals in the pool. |
| `cancel <creator> <id>` | Cancel an active proposal. |
| `extend <creator> <id> <hrs>` | Extend the voting deadline. |

### loanmgr

| Sub-command | Description |
|-------------|-------------|
| `pause` | Pause new proposals. |
| `resume` | Resume proposal submissions. |
| `stats` | Display treasury and proposal stats. |
### loanpool_apply

| Sub-command | Description |
|-------------|-------------|
| `submit <applicant> <amount> <termMonths> <purpose>` | Submit a loan application. |
| `vote <voter> <id>` | Vote on an application. |
| `process` | Finalise pending applications. |
| `disburse <id>` | Disburse an approved application. |
| `get <id>` | Display a single application. |
| `list` | List loan applications. |

### grant

| Sub-command | Description |
|-------------|-------------|
| `create <recipient> <amount>` | Create a new grant funded from the loan pool. |
| `release <id>` | Release funds for a grant. |
| `get <id>` | Display a single grant record. |

### network

| Sub-command | Description |
|-------------|-------------|
| `start` | Start the networking stack. |
| `stop` | Stop network services. |
| `peers` | List connected peers. |
| `broadcast <topic> <data>` | Publish data on the network. |
| `subscribe <topic>` | Subscribe to a topic. |

### bootstrap

| Sub-command | Description |
|-------------|-------------|
| `start` | Start the bootstrap node. |
| `stop` | Stop the bootstrap node. |
| `peers` | List peers connected to the bootstrap node. |
### connpool

| Sub-command | Description |
|-------------|-------------|
| `stats` | Show pool statistics. |
| `dial <addr>` | Dial an address using the pool. |
| `close` | Close the pool. |
### nat

| Sub-command | Description |
|-------------|-------------|
| `map <port>` | Open a port on the router via UPnP/NAT-PMP. |
| `unmap` | Remove the current port mapping. |
| `ip` | Show the discovered external IP address. |
### peer

| Sub-command | Description |
|-------------|-------------|
| `discover` | List peers discovered via mDNS. |
| `connect <addr>` | Connect to a peer by multi-address. |
| `advertise [topic]` | Broadcast this node ID on a topic. |

### replication

| Sub-command | Description |
|-------------|-------------|
| `start` | Launch replication goroutines. |
| `stop` | Stop the replication subsystem. |
| `status` | Show replication status. |
| `replicate <block-hash>` | Gossip a known block. |
| `request <block-hash>` | Request a block from peers. |
| `sync` | Synchronize blocks from peers. |

### coordination

| Sub-command | Description |
|-------------|-------------|
| `start` | Start coordination background tasks. |
| `stop` | Stop coordination tasks. |
| `broadcast` | Broadcast the current ledger height. |
| `mint <addr> <token> <amount>` | Mint tokens via the coordinator. |
### initrep

| Sub-command | Description |
|-------------|-------------|
| `start` | Bootstrap the ledger and start replication. |
| `stop` | Stop the initialization service. |
### synchronization

| Sub-command | Description |
|-------------|-------------|
| `start` | Start the sync manager. |
| `stop` | Stop the sync manager. |
| `status` | Show sync progress. |
| `once` | Perform one synchronization round. |

### high_availability

| Sub-command | Description |
|-------------|-------------|
| `add <addr>` | Register a standby node. |
| `remove <addr>` | Remove a standby node. |
| `list` | List registered standby nodes. |
| `promote <addr>` | Promote a standby to leader via view change. |
| `snapshot [path]` | Write a ledger snapshot to disk. |

### rollups

| Sub-command | Description |
|-------------|-------------|
| `submit` | Submit a new rollup batch. |
| `challenge <batchID> <txIdx> <proof...>` | Submit a fraud proof for a batch. |
| `finalize <batchID>` | Finalize or revert a batch. |
| `info <batchID>` | Display batch header and state. |
| `list` | List recent batches. |
| `txs <batchID>` | List transactions in a batch. |
| `pause` | Pause the rollup aggregator. |
| `resume` | Resume the rollup aggregator. |
| `status` | Show current aggregator status. |

### compression

| Sub-command | Description |
|-------------|-------------|
| `save <file>` | Write a compressed ledger snapshot. |
| `load <file>` | Load a compressed snapshot and display the height. |

### security

| Sub-command | Description |
|-------------|-------------|
| `sign` | Sign a message with a private key. |
| `verify` | Verify a signature. |
| `aggregate <sig1,sig2,...>` | Aggregate BLS signatures. |
| `encrypt` | Encrypt data using XChacha20‑Poly1305. |
| `decrypt` | Decrypt an encrypted blob. |
| `merkle <leaf1,leaf2,...>` | Compute a double-SHA256 Merkle root. |
| `dilithium-gen` | Generate a Dilithium3 key pair. |
| `dilithium-sign` | Sign a message with a Dilithium key. |
| `dilithium-verify` | Verify a Dilithium signature. |
| `anomaly-score` | Compute an anomaly z-score from data. |

### firewall

| Sub-command | Description |
|-------------|-------------|
| `block-address <addr>` | Block a wallet address. |
| `unblock-address <addr>` | Remove an address from the block list. |
| `block-token <id>` | Block transfers of a token id. |
| `unblock-token <id>` | Allow transfers of a token id. |
| `block-ip <ip>` | Block a peer IP address. |
| `unblock-ip <ip>` | Unblock a peer IP address. |
| `list` | Display current firewall rules. |
### biometrics

| Sub-command | Description |
|-------------|-------------|
| `enroll <file>` | Enroll biometric data for an address. |
| `verify <file>` | Verify biometric data against an address. |
| `delete <addr>` | Remove stored biometric data. |

### sharding

| Sub-command | Description |
|-------------|-------------|
| `leader get <shardID>` | Show the leader for a shard. |
| `leader set <shardID> <addr>` | Set the leader address for a shard. |
| `map` | List shard-to-leader mappings. |
| `submit <fromShard> <toShard> <txHash>` | Submit a cross-shard transaction header. |
| `pull <shardID>` | Pull receipts for a shard. |
| `reshard <newBits>` | Increase the shard count. |
| `rebalance <threshold>` | List shards exceeding the load threshold. |

### sidechain

| Sub-command | Description |
|-------------|-------------|
| `register` | Register a new side-chain. |
| `header` | Submit a side-chain header. |
| `deposit` | Deposit tokens to a side-chain escrow. |
| `withdraw <proofHex>` | Verify a withdrawal proof. |
| `get-header` | Fetch a submitted side-chain header. |
| `meta <chainID>` | Display side-chain metadata. |
| `list` | List registered side-chains. |
| `pause <chainID>` | Pause a side-chain. |
| `resume <chainID>` | Resume a paused side-chain. |
| `update-validators` | Update side-chain validator set. |
| `remove <chainID>` | Remove a side-chain and all data. |

### plasma

| Sub-command | Description |
|-------------|-------------|
| `deposit` | Deposit funds into the Plasma chain. |
| `withdraw <nonce>` | Finalise a Plasma exit. |

### plasma

| Sub-command | Description |
|-------------|-------------|
| `deposit <from> <token> <amount>` | Deposit tokens into the plasma bridge. |
| `exit <owner> <token> <amount>` | Start an exit from the bridge. |
| `finalize <nonce>` | Finalize a pending exit. |
| `get <nonce>` | Get details about an exit. |
| `list <owner>` | List exits initiated by an address. |

### state_channel

| Sub-command | Description |
|-------------|-------------|
| `open` | Open a new payment/state channel. |
| `close` | Submit a signed state to start closing. |
| `challenge` | Challenge a closing state with a newer one. |
| `finalize` | Finalize and settle an expired channel. |
| `status` | Show the current channel state. |
| `list` | List all open channels. |

### plasma

| Sub-command | Description |
|-------------|-------------|
| `deposit` | Deposit tokens into the plasma chain. |
| `withdraw` | Withdraw a previously deposited amount. |
| `submit` | Submit a plasma block root. |
### state_channel_mgmt

| Sub-command | Description |
|-------------|-------------|
| `pause` | Pause a channel to block new updates. |
| `resume` | Resume a paused channel. |
| `cancel` | Cancel a pending close operation. |
| `force-close` | Immediately settle a channel with a signed state. |
### zero_trust_data_channels

| Sub-command | Description |
|-------------|-------------|
| `open` | Open a new zero trust data channel. |
| `send` | Send a hex encoded payload over the channel. |
| `close` | Close the channel and release escrow. |

### storage

| Sub-command | Description |
|-------------|-------------|
| `pin` | Pin a file or data blob to the gateway. |
| `get` | Retrieve data by CID. |
| `listing:create` | Create a storage listing. |
| `listing:get` | Get a storage listing by ID. |
| `listing:list` | List storage listings. |
| `deal:open` | Open a storage deal backed by escrow. |
| `deal:close` | Close a storage deal and release funds. |
| `deal:get` | Get details for a storage deal. |
| `deal:list` | List storage deals. |
### real_estate

| Sub-command | Description |
|-------------|-------------|
| `register` | Register a new property. |
| `transfer` | Transfer a property to another owner. |
| `get` | Get property details. |
| `list` | List properties, optionally by owner. |


### escrow

| Sub-command | Description |
|-------------|-------------|
| `create` | Create a new multi-party escrow |
| `deposit` | Deposit additional funds |
| `release` | Release funds to participants |
| `cancel` | Cancel an escrow and refund |
| `info` | Show escrow details |
| `list` | List all escrows |
### marketplace

| Sub-command | Description |
|-------------|-------------|
| `listing:create <price> <metaJSON>` | Create a marketplace listing. |
| `listing:get <id>` | Fetch a listing by ID. |
| `listing:list` | List marketplace listings. |
| `buy <id> <buyer>` | Purchase a listing via escrow. |
| `cancel <id>` | Cancel an unsold listing. |
| `release <escrow>` | Release escrow funds to seller. |
| `deal:get <id>` | Retrieve deal details. |
| `deal:list` | List marketplace deals. |

| Sub-command | Description |
|-------------|-------------|
| `register <addr>` | Register a patient address. |
| `grant <patient> <provider>` | Allow a provider to submit records. |
| `revoke <patient> <provider>` | Revoke provider access. |
| `add <patient> <provider> <cid>` | Add a record CID for a patient. |
| `list <patient>` | List stored record IDs for a patient. |
### warehouse

| Sub-command | Description |
|-------------|-------------|
| `add` | Add a new inventory item. |
| `remove` | Delete an existing item. |
| `move` | Transfer item ownership. |
| `list` | List all warehouse items. |

### staking

| Sub-command | Description |
|-------------|-------------|
| `stake <addr> <amt>` | Stake tokens to participate in governance. |
| `unstake <addr> <amt>` | Unstake previously locked tokens. |
| `balance <addr>` | Show staked balance of an address. |
| `total` | Display the total amount staked. |

### staking

| Sub-command | Description |
|-------------|-------------|
| `stake <addr> <amt>` | Stake tokens to participate in governance. |
| `unstake <addr> <amt>` | Unstake previously locked tokens. |
| `balance <addr>` | Show staked balance of an address. |
| `total` | Display the total amount staked. |

### resource

| Sub-command | Description |
|-------------|-------------|
| `listing:create` | Create a resource listing. |
| `listing:get` | Get a resource listing by ID. |
| `listing:list` | List resource listings. |
| `deal:open` | Open a resource deal. |
| `deal:close` | Close a resource deal. |
| `deal:get` | Get resource deal details. |
| `deal:list` | List resource deals. |

### ipfs

| Sub-command | Description |
|-------------|-------------|
| `add` | Add a file to the configured IPFS gateway. |
| `get` | Fetch a CID and write to stdout or a file. |
| `unpin` | Remove a CID from the gateway pinset. |

### legal

| Sub-command | Description |
|-------------|-------------|
| `register <addr> <json>` | Register a Ricardian contract. |
| `sign <contract> <party>` | Sign a contract as a party. |
| `revoke <contract> <party>` | Revoke a signature. |
| `info <addr>` | Show contract and signers. |
| `list` | List all registered legal contracts. |

### partition

| Sub-command | Description |
|-------------|-------------|
| `split <file>` | Split a file into equally sized chunks. |
| `compress <file>` | Compress a file and print base64 output. |
| `decompress <b64>` | Decompress base64 input and print bytes. |

### tokens

| Sub-command | Description |
|-------------|-------------|
| `list` | List registered tokens. |
| `info <id|symbol>` | Display token metadata. |
| `balance <tok> <addr>` | Query token balance of an address. |
| `transfer <tok>` | Transfer tokens between addresses. |
| `mint <tok>` | Mint new tokens. |
| `burn <tok>` | Burn tokens from an address. |
| `approve <tok>` | Approve a spender allowance. |
| `allowance <tok> <owner> <spender>` | Show current allowance. |

### defi

| Sub-command | Description |
|-------------|-------------|
| `insurance new <id> <holder> <premium> <payout>` | Create an insurance policy. |
| `insurance claim <id>` | Claim a payout. |
### event_management

| Sub-command | Description |
|-------------|-------------|
| `emit <type> <data>` | Emit a new event and broadcast it. |
| `list <type>` | List recent events of a given type. |
| `get <type> <id>` | Fetch a specific event by ID. |
### token_management

| Sub-command | Description |
|-------------|-------------|
| `create` | Create a new token. |
| `balance <id> <addr>` | Check balance for a token ID. |
| `transfer <id>` | Transfer tokens between addresses. |
### tangible

| Sub-command | Description |
|-------------|-------------|
| `register <id> <owner> <meta> <value>` | Register a new tangible asset. |
| `transfer <id> <owner>` | Transfer ownership of an asset. |
| `info <id>` | Display asset metadata. |
| `list` | List all tangible assets. |
### gaming

| Sub-command | Description |
|-------------|-------------|
| `create` | Create a new game. |
| `join <id>` | Join an existing game. |
| `finish <id>` | Finish a game and release funds. |
| `get <id>` | Display a game record. |
| `list` | List games. |

### transactions

| Sub-command | Description |
|-------------|-------------|
| `create` | Craft an unsigned transaction JSON. |
| `sign` | Sign a transaction JSON with a keystore key. |
| `verify` | Verify a signed transaction JSON. |
| `submit` | Submit a signed transaction to the network. |
| `pool` | List pending pool transaction hashes. |

### distribution

| Sub-command | Description |
|-------------|-------------|
| `airdrop` | Mint tokens to a list of recipients. |
| `batch` | Transfer tokens from one account to many. |
### private_tx

| Sub-command | Description |
|-------------|-------------|
| `encrypt` | Encrypt transaction payload bytes. |
| `decrypt` | Decrypt previously encrypted payload. |
| `send` | Submit an encrypted transaction JSON file. |
### transactionreversal

| Sub-command | Description |
|-------------|-------------|
| `reversal` | Reverse a confirmed transaction. Requires authority signatures. |

### utility_functions

| Sub-command | Description |
|-------------|-------------|
| `hash` | Compute a cryptographic hash. |
| `short-hash` | Shorten a 32-byte hash to first4..last4 format. |
| `bytes2addr` | Convert big-endian bytes to an address. |

### finalization_management

| Sub-command | Description |
|-------------|-------------|
| `block <file>` | Finalize a block from JSON. |
| `batch <batchID>` | Finalize a rollup batch. |
| `channel <channelID>` | Finalize a state channel. |
### quorum

| Sub-command | Description |
|-------------|-------------|
| `init <total> <threshold>` | Initialise a global quorum tracker. |
| `vote <address>` | Record a vote from an address. |
| `check` | Check if the configured quorum is reached. |
| `reset` | Clear all recorded votes. |
### supply

| Sub-command | Description |
|-------------|-------------|
| `register <id> <desc> <owner> <location>` | Register a new item on chain. |
| `update-location <id> <location>` | Update item location. |
| `status <id> <status>` | Update item status. |
| `get <id>` | Fetch item metadata. |

### virtual_machine

| Sub-command | Description |
|-------------|-------------|
| `start` | Start the VM HTTP daemon. |
| `stop` | Stop the VM daemon. |
| `status` | Show daemon status. |

### sandbox

| Sub-command | Description |
|-------------|-------------|
| `start` | Create a sandbox for a contract. |
| `stop` | Stop a running sandbox. |
| `reset` | Reset sandbox timers. |
| `status` | Display sandbox info. |
| `list` | List all sandboxes. |
### swarm

| Sub-command | Description |
|-------------|-------------|
| `add <id> <addr>` | Add a node to the swarm. |
| `remove <id>` | Remove a node from the swarm. |
| `broadcast <tx.json>` | Broadcast a transaction to all nodes. |
| `peers` | List nodes currently in the swarm. |
| `start` | Start consensus for the swarm. |
| `stop` | Stop all nodes and consensus. |

### wallet

| Sub-command | Description |
|-------------|-------------|
| `create` | Generate a new wallet and mnemonic. |
| `import` | Import an existing mnemonic. |
| `address` | Derive an address from a wallet. |
| `sign` | Sign a transaction JSON using the wallet. |

### execution

| Sub-command | Description |
|-------------|-------------|
| `begin` | Begin a new block at a given height. |
| `run <tx.json>` | Execute a transaction JSON file. |
| `finalize` | Finalize the current block and output it. |
### binarytree

| Sub-command | Description |
|-------------|-------------|
| `create <name>` | Create a new binary tree bound to the ledger. |
| `insert <tree> <key> <value>` | Insert or update a key. Value may be hex encoded. |
| `search <tree> <key>` | Retrieve a value by key. |
| `delete <tree> <key>` | Remove a key from the tree. |
| `list <tree>` | List all keys in order. |


### system_health

| Sub-command | Description |
|-------------|-------------|
| `snapshot` | Display current system metrics. |
| `log <level> <msg>` | Append a message to the system log. |

### idwallet

| Sub-command | Description |
|-------------|-------------|
| `register <address> <info>` | Register wallet and mint a SYN-ID token. |
| `check <address>` | Verify registration status. |
### offwallet

| Sub-command | Description |
|-------------|-------------|
| `create` | Create an offline wallet file. |
| `sign` | Sign a transaction offline using the wallet. |
### recovery

| Sub-command | Description |
|-------------|-------------|
| `register` | Register recovery credentials for an address. |
| `recover` | Restore an address by proving three credentials. |
### workflow

| Sub-command | Description |
|-------------|-------------|
| `new` | Create a new workflow by ID. |
| `add` | Append an opcode name to the workflow. |
| `trigger` | Set a cron expression for execution. |
| `webhook` | Register a webhook called after completion. |
| `run` | Execute the workflow immediately. |

### wallet_mgmt

| Sub-command | Description |
|-------------|-------------|
| `create` | Create a wallet and print the mnemonic. |
| `balance` | Show the SYNN balance for an address. |
| `transfer` | Send SYNN from a mnemonic to a target address. |
### devnet

| Sub-command | Description |
|-------------|-------------|
| `start [nodes]` | Start a local developer network with the given number of nodes. |

### testnet

| Sub-command | Description |
|-------------|-------------|
| `start <config.yaml>` | Launch a testnet using the node definitions in the YAML file. |
### faucet

| Sub-command | Description |
|-------------|-------------|
| `request <addr>` | Request faucet funds for an address. |
| `balance` | Display remaining faucet balance. |
| `config --amount <n> --cooldown <d>` | Update faucet parameters. |

### syn10

| Sub-command | Description |
|-------------|-------------|
| `set-rate <rate>` | Update exchange rate. |
| `info` | Display CBDC info. |
| `mint <to> <amt>` | Mint CBDC tokens. |
| `burn <from> <amt>` | Burn CBDC tokens. |

### syn11

| Sub-command | Description |
|-------------|-------------|
| `issue <address> <amount>` | Issue SYN11 gilts. |
| `redeem <address> <amount>` | Redeem gilts from an address. |
| `set-coupon <rate>` | Update coupon rate. |
| `pay-coupon` | Pay accrued coupons. |

### syn70

| Sub-command | Description |
|-------------|-------------|
| `register <id> <owner> <name> <game>` | Register an in-game asset. |
| `transfer <id> <newOwner>` | Transfer asset ownership. |
| `setattr <id> <key> <value>` | Set a custom asset attribute. |
| `achievement <id> <name>` | Record an achievement. |
| `info <id>` | Show asset details. |
| `list` | List registered assets. |

### syn200

| Sub-command | Description |
|-------------|-------------|
| `register <owner> <name> <total>` | Register a carbon project. |
| `issue <projectID> <to> <amount>` | Issue carbon credits. |
| `retire <holder> <amount>` | Retire credits from circulation. |
| `verify <projectID> <verifier> <verID> [status]` | Add a verification record. |
| `verifications <projectID>` | List project verifications. |
| `info <projectID>` | Show project info. |
| `list` | List all carbon projects. |

### syn2100

| Sub-command | Description |
|-------------|-------------|
| `register-document <token> <docID> <issuer> <recipient> <amount> <issue> <due> <desc>` | Register a financing document. |
| `finance <token> <docID> <financier>` | Finance a document. |
| `get-document <token> <docID>` | Fetch document details. |
| `list-documents <token>` | List registered documents. |
| `add-liquidity <token> --from <addr> --amt <n>` | Add liquidity to the pool. |
| `remove-liquidity <token> --to <addr> --amt <n>` | Remove liquidity from the pool. |

### syn2200

| Sub-command | Description |
|-------------|-------------|
| `create --name <n> --symbol <s> --owner <addr> [--dec <d>] [--supply <n>]` | Create a payment token. |
| `pay --id <token> --from <addr> --to <addr> --amt <n> --cur <code>` | Send a payment through the network. |

### syn223

| Sub-command | Description |
|-------------|-------------|
| `whitelist-add <token> <addr>` | Add an address to the whitelist. |
| `whitelist-remove <token> <addr>` | Remove an address from the whitelist. |
| `blacklist-add <token> <addr>` | Add an address to the blacklist. |
| `blacklist-remove <token> <addr>` | Remove an address from the blacklist. |
| `transfer <token> --from <addr> --to <addr> --amt <n>` | Perform a safe token transfer. |

### syn2400

| Sub-command | Description |
|-------------|-------------|
| `create --name <n> --symbol <s> --owner <addr> --hash <h> [--desc <d>] [--price <n>] [--supply <n>]` | Create a data token. |

### syn300

| Sub-command | Description |
|-------------|-------------|
| `delegate <owner> <delegate>` | Delegate voting power. |
| `revoke <owner>` | Revoke delegation. |
| `power <addr>` | Show voting power. |
| `propose <creator> <description>` | Create a governance proposal. |
| `vote <id> <voter> <approve>` | Cast a vote on a proposal. |
| `execute <id> <quorum>` | Execute an approved proposal. |
| `status <id>` | Get proposal status. |
| `list` | List proposals. |

### syn3200

| Sub-command | Description |
|-------------|-------------|
| `create --issuer <addr> --payer <addr> --amt <n> --due <time> [--meta <m>]` | Create a bill. |
| `pay --bill <id> --payer <addr> --amt <n>` | Pay part of a bill. |
| `adjust --bill <id> --amt <n>` | Adjust bill amount. |
| `info --bill <id>` | Show bill details. |

### syn3300

| Sub-command | Description |
|-------------|-------------|
| `info <id>` | Display ETF token info. |
| `update <id> --price <n>` | Update ETF price. |
| `mint <id> <to> --shares <n>` | Mint ETF shares. |
| `burn <id> <from> --shares <n>` | Burn ETF shares. |

### syn500

| Sub-command | Description |
|-------------|-------------|
| `create --name <n> --symbol <s> --owner <addr> [--dec <d>] [--supply <n>]` | Create a utility token. |
| `grant <addr> --tier <n> --max <n>` | Grant a usage tier. |
| `use <addr>` | Record token usage. |

### syn5000

| Sub-command | Description |
|-------------|-------------|
| `create --name <n> --symbol <s> [--dec <d>]` | Create a gambling token. |
| `bet <bettor> --id <token> --amt <n> --odds <r> --game <type>` | Place a bet. |
| `resolve --id <token> --bet <id> [--win]` | Resolve a bet. |

### syn600

| Sub-command | Description |
|-------------|-------------|
| `stake <addr> <amt> --days <d>` | Stake SYN600 tokens. |
| `unstake <addr>` | Unstake tokens. |
| `reward <addr> <amt>` | Mint staking rewards. |
| `engage <addr> <pts>` | Record engagement points. |
| `engagement <addr>` | Show engagement for an address. |

### syn800

| Sub-command | Description |
|-------------|-------------|
| `register <desc> <valuation> <loc> <type> <cert>` | Register an asset. |
| `update <valuation>` | Update asset valuation. |
| `info` | Display asset information. |

### syn845

| Sub-command | Description |
|-------------|-------------|
| `create --name <n> --symbol <s> --owner <addr> --supply <amt>` | Create a new debt token. |
| `issue <token> <debtID> <borrower> <principal> <rate> <penalty> <due>` | Issue a debt instrument. |
| `pay <token> <debtID> <amount>` | Record a payment toward a debt. |
| `info <token> <debtID>` | Display debt information. |

### syn1000

| Sub-command | Description |
|-------------|-------------|
| `create --name <n> --symbol <s> [--dec <d>]` | Create a SYN1000 stablecoin. |
| `reserve --id <token> --asset <name> --amt <n>` | Add a reserve backing asset. |
| `setprice --id <token> --asset <name> --price <p>` | Update reserve price. |
| `value --id <token>` | Display total reserve value. |

### syn1100

| Sub-command | Description |
|-------------|-------------|
| `add <id> <owner> <hexdata>` | Add a healthcare record. |
| `grant <id> <grantee>` | Grant access to a record. |
| `revoke <id> <grantee>` | Revoke record access. |
| `get <id> <caller>` | Fetch a record. |
| `transfer <id> <newowner>` | Transfer record ownership. |

### syn1155

| Sub-command | Description |
|-------------|-------------|
| `create --name <n> --symbol <s> --owner <addr>` | Create a multi-asset token. |
| `batch-transfer --id <token> --from <addr> --to <addr,...> --tokenids <ids> --amounts <amts>` | Batch transfer assets. |
| `approve-all --id <token> --owner <addr> --operator <addr> --approve` | Set or revoke operator approval. |

### syn1200

| Sub-command | Description |
|-------------|-------------|
| `add-bridge <token> <chain> <addr>` | Add a bridge address. |
| `swap <token> --id <swap> --chain <name> --from <addr> --to <addr> --amt <n>` | Start an atomic swap. |
| `status <token> <id>` | Check swap status. |

### syn130

| Sub-command | Description |
|-------------|-------------|
| `register <id> <owner> <meta> <value>` | Register a tangible asset. |
| `value <id> <val>` | Update asset valuation. |
| `sale <id> <buyer> <price>` | Record a sale. |
| `lease <id> <lessee> <pay> <start> <end>` | Start a lease. |
| `endlease <id>` | End a lease. |

### syn1300_token

| Sub-command | Description |
|-------------|-------------|
| `info` | Show SYN1300 token info. |

### syn131

| Sub-command | Description |
|-------------|-------------|
| `create --name <n> --symbol <s> --owner <addr>` | Create an intangible asset token. |
| `value --id <token> --val <n>` | Update token valuation. |

### syn1401

| Sub-command | Description |
|-------------|-------------|
| `issue --id <id> --owner <addr> --principal <n> --rate <r> --maturity <days>` | Issue an investment. |
| `accrue <id>` | Accrue interest on an investment. |
| `redeem <id> --to <addr>` | Redeem an investment. |
| `info <id>` | Show investment details. |

### syn1600

| Sub-command | Description |
|-------------|-------------|
| `info` | Show music metadata. |
| `update --title <t> --artist <a> --album <a>` | Update music info. |
| `distribute --amt <n>` | Distribute royalties. |

### syn1700

| Sub-command | Description |
|-------------|-------------|
| `create-event --name <n> --desc <d> --location <loc> --start <unix> --end <unix> --supply <n>` | Create an event. |
| `issue --event <id> --owner <addr> --class <c> --type <t> --price <p>` | Issue a ticket. |
| `transfer --ticket <id> --from <addr> --to <addr>` | Transfer a ticket. |
| `verify --ticket <id> --holder <addr>` | Verify ticket ownership. |

### syn1800

| Sub-command | Description |
|-------------|-------------|
| `emit <owner> <amount> <desc> <source>` | Record a carbon emission. |
| `offset <owner> <amount> <desc> <source>` | Record an offset action. |
| `balance <owner>` | Show net carbon balance. |
| `records <owner>` | List footprint records. |

### syn1900

| Sub-command | Description |
|-------------|-------------|
| `issue --id <token> --credit <id> --course <id> --cname <name> --issuer <addr> --recipient <addr> --value <n> --meta <m> --expiry <time>` | Issue an education credit. |
| `verify <creditID> --id <token>` | Verify a credit. |
| `revoke <creditID> --id <token>` | Revoke a credit. |
| `get <creditID> --id <token>` | Show credit info. |
| `list <recipient> --id <token>` | List credits for a recipient. |

### syn1967

| Sub-command | Description |
|-------------|-------------|
| `create --name <n> --symbol <s> --commodity <c> --unit <u> --price <n> --owner <addr> --supply <n>` | Create a commodity token. |
| `update-price <tok> <price>` | Update spot price. |
| `price <tok>` | Show current price. |
| `history <tok>` | Show price history. |

### syn3500

| Sub-command | Description |
|-------------|-------------|
| `set-rate <rate>` | Update the fiat exchange rate. |
| `info` | Display currency code, issuer and current rate. |
| `mint <to> <amt>` | Mint currency tokens. |
| `redeem <from> <amt>` | Redeem currency tokens for fiat. |

### charity_token

| Sub-command | Description |
|-------------|-------------|
| `donate <symbol> --from <addr> --amt <n> [--purpose <p>]` | Donate to a charity campaign. |
| `progress <symbol>` | Show campaign progress. |

### employmenttoken

| Sub-command | Description |
|-------------|-------------|
| `add <id> <employer> <employee> <salary> --position <title>` | Create an employment contract. |
| `pay <id>` | Pay salary for a contract. |
| `show <id>` | Display contract details. |

### energy-token

| Sub-command | Description |
|-------------|-------------|
| `register <owner> <type> <qty> <validUnix> <location> <cert>` | Register an energy asset. |
| `transfer <assetID> <to>` | Transfer asset ownership. |
| `record <assetID> <info>` | Record sustainability info. |
| `info <assetID>` | Show asset info. |
| `list` | List energy assets. |

### forex_token

| Sub-command | Description |
|-------------|-------------|
| `rate` | Display current forex rate. |
| `update <rate>` | Update the exchange rate. |

### granttoken

| Sub-command | Description |
|-------------|-------------|
| `create <beneficiary> <name> <amount>` | Create a grant. |
| `disburse <id> <amount> [note]` | Disburse grant funds. |
| `info <id>` | Show grant details. |
| `list` | List all grants. |

### idtoken

| Sub-command | Description |
|-------------|-------------|
| `register <addr> --name <n> --dob <date> --nat <c>` | Register identity details. |
| `verify <addr> --method <m>` | Record a verification method. |
| `info <addr>` | Retrieve identity info. |
| `logs <addr>` | Show verification logs. |

### insurance_token

| Sub-command | Description |
|-------------|-------------|
| `issue <holder> --coverage <desc> --premium <n> --payout <n> --deductible <n> --limit <n> --start <time> --end <time>` | Issue a policy. |
| `claim <policyID>` | Claim an insurance policy. |
| `info <policyID>` | Show policy info. |

### iptoken

| Sub-command | Description |
|-------------|-------------|
| `register <token> --id <id> --title <t> --desc <d> --creator <c> --owner <addr>` | Register an IP asset. |
| `license <token> --id <id> --type <t> --licensee <addr> --royalty <n>` | Create a license. |
| `royalty <token> --id <id> --licensee <addr> --amount <n>` | Record a royalty payment. |

### legal_token

| Sub-command | Description |
|-------------|-------------|
| `create --name <n> --symbol <s> --doctype <t> --hash <h> --expiry <time> --owner <addr> [--supply <n>] --party <addr>...` | Create a legal document token. |
| `sign <id> <party> <sig>` | Add a party signature. |
| `revoke <id> <party>` | Revoke a signature. |
| `status <id> <status>` | Update token status. |
| `dispute <id> <action> [result]` | Start or resolve a dispute. |

### life

| Sub-command | Description |
|-------------|-------------|
| `example` | Show example life insurance policy metadata. |

### pension

| Sub-command | Description |
|-------------|-------------|
| `register <owner> <name> <maturity> [schedule]` | Register a pension plan. |
| `contribute <planID> <to> <amount>` | Contribute to a plan. |
| `withdraw <planID> <holder> <amount>` | Withdraw vested tokens. |
| `info <planID>` | Show plan info. |
| `list` | List pension plans. |

### rental_token

| Sub-command | Description |
|-------------|-------------|
| `register --token <id> --property <id> --tenant <addr> --landlord <addr> --rent <n> --deposit <n> --start <time> --end <time>` | Register a rental agreement. |
| `pay <agreementID> <amount>` | Pay rent. |
| `terminate <agreementID>` | Terminate an agreement. |

### reputation

| Sub-command | Description |
|-------------|-------------|
| `add <addr> <points> [desc]` | Add reputation activity. |
| `penalize <addr> <points> [reason]` | Penalize reputation. |
| `score <addr>` | Show reputation score. |
| `history <addr>` | Show reputation events. |
