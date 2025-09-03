# Synnergy Opcode and Gas Guide

This guide describes how Synnergy assigns opcodes, prices execution and enforces gas usage.  It is derived from the canonical tables in `gas_table.go` and the dispatcher in `opcode_dispatcher.go`.

## Opcode Structure

Every exported function in the core packages is mapped to a unique 24‑bit opcode.  The format `0xCCNNNN` splits the value into a one‑byte **category** `CC` and a two‑byte **index** `NNNN`.  Categories correspond to major modules such as the ledger, AMM, state channels or the virtual machine.  The catalogue is generated automatically; the dispatcher resolves the opcode at runtime and invokes the appropriate handler.

## Operand Semantics

Opcodes follow a stack‑based calling convention:

* Operands are pushed on the VM stack in big‑endian word order.
* Handlers pop only the words they require.  Left‑over values remain on the stack for subsequent calls.
* Words are 32 bytes.  Smaller integers should be packed into a single word.  When strings or byte slices are used, the first word contains the length followed by contiguous data words.
* Address operands are 20 byte values right‑aligned within a 32 byte word.

Using the smallest possible data types, re‑using memory buffers and avoiding redundant stack reads lowers overall gas consumption.

### String and byte considerations

Strings are encoded as UTF‑8.  For fixed identifiers prefer `bytes32` to avoid dynamic‑length penalties.  Large blobs should be stored off‑chain and referred to by hash.

### Expansion techniques

Complex logic can be decomposed into multiple small opcodes or combined into a bulk operation that touches state once.  Expansion through macros or higher level language features (e.g. Solidity libraries, Rust traits, Yul functions, Go packages or JavaScript modules) keeps bytecode compact while enabling reuse.

## Gas Accounting and Fees

Gas prices are defined in [`gas_table.go`](gas_table.go).  Each opcode has a deterministic base cost reflecting CPU, storage and network impact.  Missing entries are charged `DefaultGasCost`.

The VM charges gas **before** executing an opcode.  Dynamic portions – such as per‑word memory fees or refunds for resource release – are handled by the VM's gas meter.

```go
// gas_table.go
const DefaultGasCost uint64 = 1
```

### Fee calculation

`helpers.go` exposes the `DynamicGasCalculator` used by tooling and the CLI to estimate fees.  It walks a payload, decodes 3‑byte opcodes and sums their base cost.

```go
// helpers.go
func (d *DynamicGasCalculator) Estimate(payload []byte) (uint64, error) {
        ...
        total += GasCost(op)
        ...
}
```

`GasCost` performs a concurrent lookup in the in‑memory table and falls back to the default when an opcode is unknown. Stage 6 added structured logging so missing entries are reported for auditing.

```go
// gas_table.go
func GasCost(op Opcode) uint64 {
        if cost, ok := gasTable[op]; ok {
                return cost
        }
        log.Printf("gas_table: missing cost for opcode %d – charging default", op)
        return DefaultGasCost
}
```

`GasCostByName` lets applications fetch prices using a function's exported name
without manually converting it to an opcode.  This is useful for CLI and GUI
tooling that operates on high‑level function identifiers.

```go
// gas_table.go
cost := GasCostByName("Add")
```

Biometric authentication opcodes added in Stage 4 follow the same pricing
model and are listed in `gas_table_list.md` alongside other functions.

## Ledger synchronization and compression

Stage 5 introduces opcodes for `CompressLedger`, `DecompressLedger` and the
`SyncManager` (`Start`, `Stop`, `Status`, `Once`). Central bank and charity
modules (`Mint`, `UpdatePolicy`, `Register`, `Vote`) also receive dedicated
codes. Their base gas costs are defined in `gas_table_list.md` ensuring the VM
and CLI can accurately meter usage.

Stage 7 expands the catalogue with consensus management primitives. Validator registration, contract administration and the consensus service publish opcodes with explicit gas prices and emit OpenTelemetry traces for observability.

Lowering fees is therefore a matter of choosing cheaper opcodes, batching writes and avoiding default‑priced unknown operations.

## Efficiency Patterns

* **Batch operations** – use bulk transfer, liquidity or governance methods to amortise fixed costs.
* **Read before write** – query state with cheap `Get*` calls before attempting expensive updates.
* **Free refunds** – deleting storage or self‑destructing contracts may refund part of the gas.
* **Externalise heavy data** – large strings or arrays should live in IPFS or another storage system and be referenced by hash.

## Language‑specific tips

### Solidity
Use libraries and events to keep runtime bytecode small.  Prefer `bytes32` for identifiers, pack structs and leverage `view` functions for off‑chain reads.

### Rust
When using `ink!` or other Rust‑based frameworks, mark functions as `#[payable]` only when necessary and reuse `Vec` allocations.  Traits allow opcode wrappers to be composed efficiently.

### Yul
Yul gives fine grained control over stack layout.  Preallocate memory and avoid leaving junk on the stack.  Macro‑based expansion can bundle repeated opcode sequences.

### JavaScript / TypeScript
Client side scripts typically interact with the VM via RPC.  Batch RPC calls and cache results.  When generating bytecode, reuse the dispatcher catalogue to look up opcodes instead of hard coding numbers.

### Go
Go contracts compiled to WASM should keep exported functions minimal and reuse byte slices.  Use the `DynamicGasCalculator` during development to estimate costs before deployment.

## Complete Gas Catalogue

The following sections list every opcode grouped by category together with its base gas cost.  Values are in gas units and represent the base fee charged before any dynamic component.

### AI

Operations related to ai.


| Opcode | Gas Cost |
|---|---|
| `InitAI` | `5000` |
| `AI` | `4000` |
| `PredictAnomaly` | `3500` |
| `OptimizeFees` | `2500` |
| `PublishModel` | `4500` |
| `FetchModel` | `1500` |
| `ListModel` | `800` |
| `ValidateKYC` | `100` |
| `BuyModel` | `3000` |
| `RentModel` | `2000` |
| `ReleaseEscrow` | `1200` |
| `PredictVolume` | `1500` |
| `DeployAIContract` | `5000` |
| `InvokeAIContract` | `750` |
| `UpdateAIModel` | `2000` |
| `GetAIModel` | `200` |
| `StartTraining` | `5000` |
| `TrainingStatus` | `500` |
| `ListTrainingJobs` | `800` |
| `CancelTraining` | `1000` |
| `GetModelListing` | `100` |
| `ListModelListings` | `200` |
| `UpdateListingPrice` | `200` |
| `RemoveListing` | `200` |
| `InferModel` | `3000` |
| `AnalyseTransactions` | `3500` |


### Automated-Market-Maker

Operations related to automated-market-maker.


| Opcode | Gas Cost |
|---|---|
| `SwapExactIn` | `450` |
| `AddLiquidity` | `500` |
| `RemoveLiquidity` | `500` |
| `Quote` | `250` |
| `AllPairs` | `200` |
| `InitPoolsFromFile` | `600` |


### Authority / Validator-Set

Operations related to authority / validator-set.


| Opcode | Gas Cost |
|---|---|
| `NewAuthoritySet` | `2000` |
| `RecordVote` | `300` |
| `RemoveVote` | `150` |
| `RegisterCandidate` | `800` |
| `RandomElectorate` | `400` |
| `IsAuthority` | `80` |
| `GetAuthority` | `100` |
| `ListAuthorities` | `200` |
| `DeregisterAuthority` | `600` |
| `NewAuthorityApplier` | `2000` |
| `SubmitApplication` | `400` |
| `VoteApplication` | `300` |
| `FinalizeApplication` | `500` |
| `GetApplication` | `100` |
| `ListApplications` | `200` |
| `ElectedAuth_RecordVote` | `300` |
| `ElectedAuth_Report` | `300` |
| `ElectedAuth_ValidateTx` | `500` |
| `ElectedAuth_CreateBlock` | `1000` |
| `ElectedAuth_ReverseTx` | `600` |
| `ElectedAuth_ViewPrivateTx` | `300` |
| `ElectedAuth_ApproveLoan` | `400` |


### Charity Pool

Operations related to charity pool.


| Opcode | Gas Cost |
|---|---|
| `NewCharityPool` | `1000` |
| `Deposit` | `210` |
| `Charity_Register` | `0` |
| `Vote` | `300` |
| `Tick` | `100` |
| `GetRegistration` | `80` |
| `Winners` | `80` |
| `Charity_Donate` | `0` |
| `Charity_WithdrawInternal` | `0` |
| `Charity_Balances` | `0` |


### Coin

Operations related to coin.


| Opcode | Gas Cost |
|---|---|
| `NewCoin` | `1200` |
| `Mint` | `210` |
| `TotalSupply` | `80` |
| `BalanceOf` | `40` |


### Compliance

Operations related to compliance.


| Opcode | Gas Cost |
|---|---|
| `InitCompliance` | `800` |
| `EraseData` | `500` |
| `RecordFraudSignal` | `700` |
| `Compliance_LogAudit` | `0` |
| `Compliance_AuditTrail` | `0` |
| `Compliance_MonitorTx` | `0` |
| `Compliance_VerifyZKP` | `0` |
| `Audit_Init` | `0` |
| `Audit_Log` | `100` |
| `Audit_Events` | `20` |
| `Audit_Close` | `0` |
| `InitComplianceManager` | `1000` |
| `SuspendAccount` | `400` |
| `ResumeAccount` | `400` |
| `IsSuspended` | `50` |
| `WhitelistAccount` | `300` |
| `RemoveWhitelist` | `300` |
| `IsWhitelisted` | `50` |
| `Compliance_ReviewTx` | `0` |
| `AnalyzeAnomaly` | `600` |
| `FlagAnomalyTx` | `250` |


### Consensus Core

Operations related to consensus core.


| Opcode | Gas Cost |
|---|---|
| `Pick` | `200` |
| `Broadcast` | `500` |
| `Subscribe` | `150` |
| `Sign` | `300` |
| `Verify` | `350` |
| `ValidatorPubKey` | `80` |
| `StakeOf` | `100` |
| `LoanPoolAddress` | `80` |
| `Hash` | `60` |
| `SerializeWithoutNonce` | `120` |
| `NewConsensus` | `2500` |
| `Start` | `500` |
| `ProposeSubBlock` | `1500` |
| `ValidatePoH` | `2000` |
| `SealMainBlockPOW` | `6000` |
| `DistributeRewards` | `1000` |
| `CalculateWeights` | `800` |
| `ComputeThreshold` | `600` |
| `NewConsensusAdaptiveManager` | `1000` |
| `ComputeDemand` | `200` |
| `ComputeStakeConcentration` | `200` |
| `AdjustConsensus` | `500` |
| `HopConsensus` | `400` |
| `CurrentConsensus` | `50` |
| `Status` | `100` |
| `SetDifficulty` | `200` |
| `AdjustStake` | `300` |
| `PenalizeValidator` | `400` |
| `RegisterValidator` | `800` |
| `DeregisterValidator` | `600` |
| `StakeValidator` | `200` |
| `UnstakeValidator` | `200` |
| `SlashValidator` | `300` |
| `GetValidator` | `100` |
| `ListValidators` | `200` |
| `IsValidator` | `80` |
| `StartValidatorNode` | `1000` |
| `StopValidatorNode` | `800` |
| `ProposeBlock` | `1200` |
| `VoteBlock` | `400` |
| `ConsensusNode_Start` | `500` |
| `ConsensusNode_Stop` | `500` |
| `ConsensusNode_SubmitBlock` | `800` |
| `ConsensusNode_ProcessTx` | `400` |


### Contracts (WASM / EVM‐compat)

Operations related to contracts (wasm / evm‐compat).


| Opcode | Gas Cost |
|---|---|
| `InitContracts` | `1500` |
| `CompileWASM` | `4500` |
| `Invoke` | `700` |
| `Deploy` | `2500` |
| `TransferOwnership` | `500` |
| `PauseContract` | `300` |
| `ResumeContract` | `300` |
| `UpgradeContract` | `2000` |
| `ContractInfo` | `100` |


### Cross-Chain

Operations related to cross-chain.


| Opcode | Gas Cost |
|---|---|
| `RegisterBridge` | `2000` |
| `AuthorizeRelayer` | `500` |
| `RevokeRelayer` | `500` |
| `LockMint` | `3000` |
| `BurnRelease` | `3000` |
| `GetBridge` | `100` |
| `ListBridges` | `120` |
| `RegisterMapping` | `2200` |
| `GetMapping` | `100` |
| `ListMappings` | `120` |
| `RemoveMapping` | `500` |
| `GetTransfer` | `100` |
| `ListTransfers` | `200` |
| `OpenConnection` | `1000` |
| `CloseConnection` | `500` |
| `GetConnection` | `100` |
| `ListConnections` | `200` |
| `RegisterProtocol` | `2000` |
| `ListProtocols` | `200` |
| `GetProtocol` | `100` |


### Cross-Consensus Scaling Networks

Operations related to cross-consensus scaling networks.


| Opcode | Gas Cost |
|---|---|
| `RegisterCCSNetwork` | `2000` |
| `ListCCSNetworks` | `500` |
| `GetCCSNetwork` | `100` |
| `CCSLockAndTransfer` | `3000` |
| `CCSBurnAndRelease` | `3000` |


### Data / Oracle / IPFS Integration

Operations related to data / oracle / ipfs integration.


| Opcode | Gas Cost |
|---|---|
| `RegisterNode` | `1000` |
| `UploadAsset` | `3000` |
| `Pin` | `500` |
| `Retrieve` | `400` |
| `RetrieveAsset` | `400` |
| `RegisterOracle` | `1000` |
| `PushFeed` | `300` |
| `QueryOracle` | `300` |
| `ListCDNNodes` | `300` |
| `RegisterContentNode` | `1000` |
| `UploadContent` | `3000` |
| `RetrieveContent` | `400` |
| `ListContentNodes` | `300` |
| `ListOracles` | `300` |
| `PushFeedSigned` | `400` |
| `UpdateOracleSource` | `400` |
| `RemoveOracle` | `400` |
| `GetOracleMetrics` | `200` |
| `RequestOracleData` | `300` |
| `SyncOracle` | `500` |
| `CreateDataSet` | `800` |
| `PurchaseDataSet` | `500` |
| `GetDataSet` | `100` |
| `ListDataSets` | `200` |
| `HasAccess` | `100` |
| `CreateDataFeed` | `600` |
| `QueryDataFeed` | `300` |
| `ManageDataFeed` | `500` |
| `ImputeMissing` | `400` |
| `NormalizeFeed` | `400` |
| `AddProvenance` | `200` |
| `SampleFeed` | `300` |
| `ScaleFeed` | `300` |
| `TransformFeed` | `400` |
| `VerifyFeedTrust` | `300` |
| `ZTDC_Open` | `0` |
| `ZTDC_Send` | `0` |
| `ZTDC_Close` | `0` |
| `StoreManagedData` | `800` |
| `LoadManagedData` | `300` |
| `DeleteManagedData` | `200` |


### External Sensors

Operations related to external sensors.


| Opcode | Gas Cost |
|---|---|
| `RegisterSensor` | `1000` |
| `GetSensor` | `100` |
| `ListSensors` | `200` |
| `UpdateSensorValue` | `150` |
| `PollSensor` | `500` |
| `TriggerWebhook` | `500` |


### Environmental Monitoring Node

Operations related to environmental monitoring node.


| Opcode | Gas Cost |
|---|---|
| `NewEnvironmentalNode` | `2000` |
| `EnvNode_AddTrigger` | `500` |
| `EnvNode_RemoveTrigger` | `300` |
| `EnvNode_Start` | `300` |
| `EnvNode_Stop` | `300` |
| `EnvNode_ListSensors` | `100` |


### Fault-Tolerance / Health-Checker

Operations related to fault-tolerance / health-checker.


| Opcode | Gas Cost |
|---|---|
| `NewHealthChecker` | `800` |
| `AddPeer` | `150` |
| `RemovePeer` | `150` |
| `FT_Snapshot` | `0` |
| `Recon` | `800` |
| `Ping` | `30` |
| `SendPing` | `30` |
| `AwaitPong` | `30` |
| `BackupSnapshot` | `1000` |
| `RestoreSnapshot` | `1200` |
| `VerifyBackup` | `600` |
| `FailoverNode` | `800` |
| `PredictFailure` | `100` |
| `AdjustResources` | `150` |
| `InitResourceManager` | `500` |
| `SetLimit` | `100` |
| `GetLimit` | `50` |
| `ConsumeLimit` | `80` |
| `TransferLimit` | `120` |
| `ListLimits` | `70` |
| `HA_Register` | `0` |
| `HA_Remove` | `0` |
| `HA_List` | `0` |
| `HA_Sync` | `0` |
| `HA_Promote` | `0` |


### Disaster Recovery Node

Operations related to disaster recovery node.


| Opcode | Gas Cost |
|---|---|
| `DR_Start` | `0` |
| `DR_Stop` | `0` |
| `DR_BackupNow` | `1000` |
| `DR_Restore` | `1200` |
| `DR_Verify` | `600` |


### Governance

Operations related to governance.


| Opcode | Gas Cost |
|---|---|
| `UpdateParam` | `500` |
| `ProposeChange` | `1000` |
| `VoteChange` | `300` |
| `EnactChange` | `800` |
| `SubmitProposal` | `1000` |
| `BalanceOfAsset` | `60` |
| `CastVote` | `300` |
| `ExecuteProposal` | `1500` |
| `GetProposal` | `100` |
| `ListProposals` | `200` |
| `NewQuorumTracker` | `100` |
| `QuorumAddVote` | `0` |
| `QuorumHasQuorum` | `0` |
| `QuorumReset` | `0` |
| `SubmitQuadraticVote` | `350` |
| `QuadraticResults` | `200` |
| `QuadraticWeight` | `5` |
| `RegisterGovContract` | `800` |
| `GetGovContract` | `100` |
| `ListGovContracts` | `200` |
| `EnableGovContract` | `100` |
| `DeleteGovContract` | `100` |
| `DeployGovContract` | `2500` |
| `InvokeGovContract` | `700` |
| `AddReputation` | `200` |
| `SubtractReputation` | `200` |
| `ReputationOf` | `50` |
| `SubmitRepGovProposal` | `1000` |
| `CastRepGovVote` | `300` |
| `ExecuteRepGovProposal` | `1500` |
| `GetRepGovProposal` | `100` |
| `ListRepGovProposals` | `200` |
| `CastTokenVote` | `400` |
| `DAO_Stake` | `0` |
| `DAO_Unstake` | `0` |
| `DAO_Staked` | `0` |
| `DAO_TotalStaked` | `0` |
| `AddDAOMember` | `120` |
| `RemoveDAOMember` | `120` |
| `RoleOfMember` | `50` |
| `ListDAOMembers` | `100` |
| `NewTimelock` | `400` |
| `QueueProposal` | `300` |
| `CancelProposal` | `300` |
| `ExecuteReady` | `500` |
| `ListTimelocks` | `100` |
| `SYN300_Delegate` | `40` |
| `SYN300_RevokeDelegate` | `40` |
| `SYN300_VotingPower` | `30` |
| `SYN300_CreateProposal` | `100` |
| `SYN300_Vote` | `80` |
| `SYN300_ExecuteProposal` | `150` |
| `SYN300_ProposalStatus` | `30` |
| `SYN300_ListProposals` | `50` |
| `CreateDAO` | `1000` |
| `JoinDAO` | `300` |
| `LeaveDAO` | `200` |
| `DAOInfo` | `100` |
| `ListDAOs` | `200` |
| `UpdateParam` | `500` |
| `ProposeChange` | `1000` |
| `VoteChange` | `300` |
| `EnactChange` | `800` |
| `SubmitProposal` | `1000` |
| `BalanceOfAsset` | `60` |
| `CastVote` | `300` |
| `ExecuteProposal` | `1500` |
| `GetProposal` | `100` |
| `ListProposals` | `200` |
| `NewQuorumTracker` | `100` |
| `QuorumAddVote` | `0` |
| `QuorumHasQuorum` | `0` |
| `QuorumReset` | `0` |
| `SubmitQuadraticVote` | `350` |
| `QuadraticResults` | `200` |
| `QuadraticWeight` | `5` |
| `RegisterGovContract` | `800` |
| `GetGovContract` | `100` |
| `ListGovContracts` | `200` |
| `EnableGovContract` | `100` |
| `DeleteGovContract` | `100` |
| `DeployGovContract` | `2500` |
| `InvokeGovContract` | `700` |
| `AddReputation` | `200` |
| `SubtractReputation` | `200` |
| `ReputationOf` | `50` |
| `SubmitRepGovProposal` | `1000` |
| `CastRepGovVote` | `300` |
| `ExecuteRepGovProposal` | `1500` |
| `GetRepGovProposal` | `100` |
| `ListRepGovProposals` | `200` |
| `RepAddActivity` | `200` |
| `RepEndorse` | `200` |
| `RepPenalize` | `200` |
| `RepScore` | `50` |
| `RepLevel` | `50` |
| `RepHistory` | `100` |
| `CastTokenVote` | `400` |
| `DAO_Stake` | `0` |
| `DAO_Unstake` | `0` |
| `DAO_Staked` | `0` |
| `DAO_TotalStaked` | `0` |
| `AddDAOMember` | `120` |
| `RemoveDAOMember` | `120` |
| `RoleOfMember` | `50` |
| `ListDAOMembers` | `100` |
| `AddSYN2500Member` | `120` |
| `RemoveSYN2500Member` | `120` |
| `DelegateSYN2500Vote` | `80` |
| `SYN2500VotingPower` | `50` |
| `CastSYN2500Vote` | `150` |
| `SYN2500MemberInfo` | `50` |
| `ListSYN2500Members` | `100` |
| `NewTimelock` | `400` |
| `QueueProposal` | `300` |
| `CancelProposal` | `300` |
| `ExecuteReady` | `500` |
| `ListTimelocks` | `100` |
| `CreateDAO` | `1000` |
| `JoinDAO` | `300` |
| `LeaveDAO` | `200` |
| `DAOInfo` | `100` |
| `ListDAOs` | `200` |


### Green Technology

Operations related to green technology.


| Opcode | Gas Cost |
|---|---|
| `InitGreenTech` | `800` |
| `Green` | `200` |
| `RecordUsage` | `300` |
| `RecordOffset` | `300` |
| `Certify` | `700` |
| `CertificateOf` | `50` |
| `ShouldThrottle` | `20` |
| `ListCertificates` | `100` |


### Energy Efficiency

Operations related to energy efficiency.


| Opcode | Gas Cost |
|---|---|
| `InitEnergyEfficiency` | `800` |
| `EnergyEff` | `200` |
| `RecordStats` | `300` |
| `EfficiencyOf` | `50` |
| `NetworkAverage` | `100` |
| `ListEfficiency` | `100` |
| `NewEnergyNode` | `2000` |
| `EnergyNodeStart` | `50` |
| `EnergyNodeStop` | `50` |
| `EnergyNodeRecord` | `300` |
| `EnergyNodeEfficiency` | `100` |
| `EnergyNodeNetworkAvg` | `100` |


### Ledger / UTXO / Account-Model

Operations related to ledger / utxo / account-model.


| Opcode | Gas Cost |
|---|---|
| `NewLedger` | `5000` |
| `GetPendingSubBlocks` | `200` |
| `LastBlockHash` | `60` |
| `AppendBlock` | `5000` |
| `MintBig` | `220` |
| `EmitApproval` | `120` |
| `EmitTransfer` | `120` |
| `DeductGas` | `210` |
| `WithinBlock` | `100` |
| `IsIDTokenHolder` | `40` |
| `TokenBalance` | `40` |
| `AddBlock` | `4000` |
| `GetBlock` | `200` |
| `GetUTXO` | `150` |
| `AddToPool` | `100` |
| `ListPool` | `80` |
| `GetContract` | `100` |
| `Snapshot` | `300` |
| `MintToken` | `200` |
| `LastSubBlockHeight` | `50` |
| `LastBlockHeight` | `50` |
| `RecordPoSVote` | `300` |
| `AppendSubBlock` | `800` |
| `Transfer` | `210` |
| `Burn` | `210` |
| `InitForkManager` | `500` |
| `AddForkBlock` | `700` |
| `ResolveForks` | `1200` |
| `ListForks` | `200` |
| `Account_Create` | `0` |
| `Account_Delete` | `0` |
| `Account_Balance` | `0` |
| `Account_Transfer` | `0` |


### Liquidity Manager (high-level AMM façade)

Operations related to liquidity manager (high-level amm façade).


| Opcode | Gas Cost |
|---|---|
| `InitAMM` | `800` |
| `Manager` | `100` |
| `CreatePool` | `1000` |
| `Swap` | `450` |


### AddLiquidity & RemoveLiquidity already defined above

Operations related to addliquidity & removeliquidity already defined above.


| Opcode | Gas Cost |
|---|---|
| `Pool` | `150` |
| `Pools` | `200` |


### Loan-Pool

Operations related to loan-pool.


| Opcode | Gas Cost |
|---|---|
| `NewLoanPool` | `2000` |
| `Submit` | `300` |
| `Disburse` | `800` |
| `Loanpool_GetProposal` | `0` |
| `Loanpool_ListProposals` | `0` |
| `Redistribute` | `500` |
| `Loanpool_CancelProposal` | `0` |
| `Loanpool_ExtendProposal` | `0` |
| `Loanpool_RequestApproval` | `0` |
| `Loanpool_ApproveRequest` | `0` |
| `Loanpool_RejectRequest` | `0` |
| `Loanpool_CreateGrant` | `0` |
| `Loanpool_ReleaseGrant` | `0` |
| `Loanpool_GetGrant` | `0` |
| `NewLoanPoolManager` | `1000` |
| `Loanpool_Pause` | `0` |
| `Loanpool_Resume` | `0` |
| `Loanpool_IsPaused` | `0` |
| `Loanpool_Stats` | `0` |
| `NewLoanPoolApply` | `2000` |
| `LoanApply_Submit` | `0` |
| `LoanApply_Vote` | `0` |
| `LoanApply_Process` | `0` |
| `LoanApply_Disburse` | `0` |
| `LoanApply_Get` | `0` |
| `LoanApply_List` | `0` |


### Networking

Operations related to networking.


| Opcode | Gas Cost |
|---|---|
| `NewNode` | `1800` |
| `HandlePeerFound` | `150` |
| `DialSeed` | `200` |
| `ListenAndServe` | `800` |
| `Close` | `50` |
| `Peers` | `40` |
| `NewDialer` | `200` |
| `Dial` | `200` |
| `SetBroadcaster` | `50` |
| `GlobalBroadcast` | `100` |
| `NewBootstrapNode` | `2000` |
| `Bootstrap_Start` | `0` |
| `Bootstrap_Stop` | `0` |
| `Bootstrap_Peers` | `0` |
| `Bootstrap_DialSeed` | `0` |
| `NewConnPool` | `800` |
| `AcquireConn` | `50` |
| `ReleaseConn` | `20` |
| `ClosePool` | `40` |
| `PoolStats` | `10` |
| `NewGatewayNode` | `2500` |
| `Gateway_Start` | `800` |
| `Gateway_Stop` | `400` |
| `Gateway_AddSource` | `200` |
| `Gateway_RemoveSource` | `200` |
| `Gateway_ListSources` | `50` |
| `Gateway_ConnectChain` | `1000` |
| `Gateway_DisconnectChain` | `500` |
| `Gateway_ListConnections` | `200` |
| `Gateway_PushExternalData` | `400` |
| `Gateway_QueryExternalData` | `300` |
| `NewNATManager` | `500` |
| `NAT_Map` | `0` |
| `NAT_Unmap` | `0` |
| `NAT_ExternalIP` | `0` |
| `DiscoverPeers` | `100` |
| `Connect` | `150` |
| `Disconnect` | `100` |
| `AdvertiseSelf` | `80` |
| `StartDevNet` | `5000` |
| `StartTestNet` | `6000` |
| `NewNode` | `1800` |
| `HandlePeerFound` | `150` |
| `DialSeed` | `200` |
| `ListenAndServe` | `800` |
| `Close` | `50` |
| `Peers` | `40` |
| `NewDialer` | `200` |
| `Dial` | `200` |
| `SetBroadcaster` | `50` |
| `GlobalBroadcast` | `100` |
| `NewBootstrapNode` | `2000` |
| `Bootstrap_Start` | `0` |
| `Bootstrap_Stop` | `0` |
| `Bootstrap_Peers` | `0` |
| `Bootstrap_DialSeed` | `0` |
| `NewConnPool` | `800` |
| `AcquireConn` | `50` |
| `ReleaseConn` | `20` |
| `ClosePool` | `40` |
| `PoolStats` | `10` |
| `NewNATManager` | `500` |
| `NAT_Map` | `0` |
| `NAT_Unmap` | `0` |
| `NAT_ExternalIP` | `0` |
| `DiscoverPeers` | `100` |
| `Connect` | `150` |
| `Disconnect` | `100` |
| `AdvertiseSelf` | `80` |
| `StartDevNet` | `5000` |
| `StartTestNet` | `6000` |
| `NewStakingNode` | `2500` |
| `Staking_Start` | `800` |
| `Staking_Stop` | `400` |
| `Staking_Stake` | `300` |
| `Staking_Unstake` | `300` |
| `Staking_ProposeBlock` | `1000` |
| `Staking_ValidateBlock` | `500` |
| `Staking_Status` | `10` |
| `NewSuperNode` | `2500` |
| `Super_Start` | `1000` |
| `Super_Stop` | `500` |
| `Super_Peers` | `100` |
| `Super_DialSeed` | `200` |
| `Super_ExecuteContract` | `3000` |
| `BroadcastOrphanBlock` | `150` |
| `SubscribeOrphanBlocks` | `150` |
| `NewOrphanNode` | `1000` |
| `Orphan_Process` | `200` |
| `Orphan_Detect` | `50` |
| `Orphan_Analyse` | `80` |
| `Orphan_Recycle` | `100` |
| `Orphan_Archive` | `70` |
| `NewNode` | `1800` |
| `HandlePeerFound` | `150` |
| `DialSeed` | `200` |
| `ListenAndServe` | `800` |
| `Close` | `50` |
| `Peers` | `40` |
| `NewDialer` | `200` |
| `Dial` | `200` |
| `SetBroadcaster` | `50` |
| `GlobalBroadcast` | `100` |
| `NewBootstrapNode` | `2000` |
| `Bootstrap_Start` | `0` |
| `Bootstrap_Stop` | `0` |
| `Bootstrap_Peers` | `0` |
| `Bootstrap_DialSeed` | `0` |
| `NewConnPool` | `800` |
| `AcquireConn` | `50` |
| `ReleaseConn` | `20` |
| `ClosePool` | `40` |
| `PoolStats` | `10` |
| `NewNATManager` | `500` |
| `NAT_Map` | `0` |
| `NAT_Unmap` | `0` |
| `NAT_ExternalIP` | `0` |
| `DiscoverPeers` | `100` |
| `Connect` | `150` |
| `Disconnect` | `100` |
| `AdvertiseSelf` | `80` |
| `StartDevNet` | `5000` |
| `StartTestNet` | `6000` |
| `NewQuantumResistantNode` | `2500` |
| `Quantum_Start` | `800` |
| `Quantum_Stop` | `400` |
| `Quantum_SecureBroadcast` | `300` |
| `Quantum_SecureSubscribe` | `300` |
| `Quantum_RotateKeys` | `200` |
| `NewNode` | `1800` |
| `HandlePeerFound` | `150` |
| `DialSeed` | `200` |
| `ListenAndServe` | `800` |
| `Close` | `50` |
| `Peers` | `40` |
| `NewDialer` | `200` |
| `Dial` | `200` |
| `SetBroadcaster` | `50` |
| `GlobalBroadcast` | `100` |
| `NewBootstrapNode` | `2000` |
| `Bootstrap_Start` | `0` |
| `Bootstrap_Stop` | `0` |
| `Bootstrap_Peers` | `0` |
| `Bootstrap_DialSeed` | `0` |
| `NewConnPool` | `800` |
| `AcquireConn` | `50` |
| `ReleaseConn` | `20` |
| `ClosePool` | `40` |
| `PoolStats` | `10` |
| `NewNATManager` | `500` |
| `NAT_Map` | `0` |
| `NAT_Unmap` | `0` |
| `NAT_ExternalIP` | `0` |
| `DiscoverPeers` | `100` |
| `Connect` | `150` |
| `Disconnect` | `100` |
| `AdvertiseSelf` | `80` |
| `StartDevNet` | `5000` |
| `StartTestNet` | `6000` |
| `NewExperimentalNode` | `2500` |
| `Exp_StartTesting` | `1000` |
| `Exp_StopTesting` | `500` |
| `Exp_DeployFeature` | `1500` |
| `Exp_RollbackFeature` | `1200` |
| `Exp_SimulateTx` | `800` |
| `Exp_TestContract` | `1000` |
| `NewNode` | `1800` |
| `HandlePeerFound` | `150` |
| `DialSeed` | `200` |
| `ListenAndServe` | `800` |
| `Close` | `50` |
| `Peers` | `40` |
| `NewDialer` | `200` |
| `Dial` | `200` |
| `SetBroadcaster` | `50` |
| `GlobalBroadcast` | `100` |
| `NewBootstrapNode` | `2000` |
| `Bootstrap_Start` | `0` |
| `Bootstrap_Stop` | `0` |
| `Bootstrap_Peers` | `0` |
| `Bootstrap_DialSeed` | `0` |
| `NewAutonomousAgentNode` | `2500` |
| `Autonomous_Start` | `0` |
| `Autonomous_Stop` | `0` |
| `Autonomous_AddRule` | `0` |
| `Autonomous_RemoveRule` | `0` |
| `NewConnPool` | `800` |
| `AcquireConn` | `50` |
| `ReleaseConn` | `20` |
| `ClosePool` | `40` |
| `PoolStats` | `10` |
| `NewNATManager` | `500` |
| `NAT_Map` | `0` |
| `NAT_Unmap` | `0` |
| `NAT_ExternalIP` | `0` |
| `DiscoverPeers` | `100` |
| `Connect` | `150` |
| `Disconnect` | `100` |
| `AdvertiseSelf` | `80` |
| `StartDevNet` | `5000` |
| `StartTestNet` | `6000` |
| `NewMasterNode` | `3000` |
| `Master_Start` | `500` |
| `Master_Stop` | `300` |
| `Master_ProcessTx` | `200` |
| `Master_HandlePrivateTx` | `300` |
| `Master_VoteProposal` | `100` |
| `NewNode` | `1800` |
| `HandlePeerFound` | `150` |
| `DialSeed` | `200` |
| `ListenAndServe` | `800` |
| `Close` | `50` |
| `Peers` | `40` |
| `NewDialer` | `200` |
| `Dial` | `200` |
| `SetBroadcaster` | `50` |
| `GlobalBroadcast` | `100` |
| `NewBootstrapNode` | `2000` |
| `Bootstrap_Start` | `0` |
| `Bootstrap_Stop` | `0` |
| `Bootstrap_Peers` | `0` |
| `Bootstrap_DialSeed` | `0` |
| `NewFullNode` | `2500` |
| `Full_Start` | `0` |
| `Full_Stop` | `0` |
| `Full_Peers` | `0` |
| `Full_DialSeed` | `0` |
| `NewAuditNode` | `0` |
| `AuditNode_Start` | `0` |
| `AuditNode_Stop` | `0` |
| `AuditNode_Log` | `100` |
| `AuditNode_Events` | `20` |
| `NewConnPool` | `800` |
| `AcquireConn` | `50` |
| `ReleaseConn` | `20` |
| `ClosePool` | `40` |
| `PoolStats` | `10` |
| `NewNATManager` | `500` |
| `NAT_Map` | `0` |
| `NAT_Unmap` | `0` |
| `NAT_ExternalIP` | `0` |
| `DiscoverPeers` | `100` |
| `Connect` | `150` |
| `Disconnect` | `100` |
| `AdvertiseSelf` | `80` |
| `StartDevNet` | `5000` |
| `StartTestNet` | `6000` |


### Mining Node

Operations related to mining node.


| Opcode | Gas Cost |
|---|---|
| `NewMiningNode` | `2000` |
| `StartMining` | `1000` |
| `StopMining` | `500` |
| `AddTransaction` | `200` |
| `SolvePuzzle` | `5000` |
| `NewAPINode` | `2500` |
| `APINode_Start` | `500` |
| `APINode_Stop` | `200` |
| `NewAIEnhancedNode` | `2500` |
| `AINode_Start` | `800` |
| `AINode_Stop` | `400` |
| `AINode_PredictLoad` | `200` |
| `AINode_AnalyseTx` | `300` |
| `NewMobileNode` | `1500` |
| `Mobile_Start` | `400` |
| `Mobile_Stop` | `200` |
| `Mobile_QueueTx` | `20` |
| `Mobile_FlushTxs` | `50` |
| `Mobile_SetOffline` | `10` |
| `Mobile_SyncLedger` | `200` |
| `NewZKPNode` | `2000` |
| `ZKP_Start` | `800` |
| `ZKP_Stop` | `400` |
| `ZKP_GenerateProof` | `500` |
| `ZKP_VerifyProof` | `500` |
| `ZKP_StoreProof` | `40` |
| `ZKP_GetProof` | `40` |
| `ZKP_SubmitTx` | `600` |
| `NewTimeLockedNode` | `2000` |
| `TL_Queue` | `200` |
| `TL_Cancel` | `100` |
| `TL_ExecuteDue` | `500` |
| `TL_List` | `50` |
| `NewNode` | `1800` |
| `HandlePeerFound` | `150` |
| `DialSeed` | `200` |
| `ListenAndServe` | `800` |
| `Close` | `50` |
| `Peers` | `40` |
| `NewDialer` | `200` |
| `Dial` | `200` |
| `SetBroadcaster` | `50` |
| `GlobalBroadcast` | `100` |
| `NewBootstrapNode` | `2000` |
| `Bootstrap_Start` | `0` |
| `Bootstrap_Stop` | `0` |
| `Bootstrap_Peers` | `0` |
| `Bootstrap_DialSeed` | `0` |
| `NewConnPool` | `800` |
| `AcquireConn` | `50` |
| `ReleaseConn` | `20` |
| `ClosePool` | `40` |
| `PoolStats` | `10` |
| `NewNATManager` | `500` |
| `NAT_Map` | `0` |
| `NAT_Unmap` | `0` |
| `NAT_ExternalIP` | `0` |
| `DiscoverPeers` | `100` |
| `Connect` | `150` |
| `Disconnect` | `100` |
| `AdvertiseSelf` | `80` |
| `StartDevNet` | `5000` |
| `StartTestNet` | `6000` |
| `NewHoloNode` | `2200` |
| `Holo_Start` | `800` |
| `Holo_Stop` | `50` |
| `Holo_EncodeStore` | `500` |
| `Holo_Retrieve` | `150` |
| `Holo_Sync` | `300` |
| `Holo_ProcessTx` | `400` |
| `Holo_ExecuteContract` | `600` |
| `Molecular_AtomicTx` | `2000` |
| `Molecular_EncodeData` | `1500` |
| `Molecular_Monitor` | `1000` |
| `Molecular_Control` | `1500` |
| `NewNode` | `1800` |
| `HandlePeerFound` | `150` |
| `DialSeed` | `200` |
| `ListenAndServe` | `800` |
| `Close` | `50` |
| `Peers` | `40` |
| `NewDialer` | `200` |
| `Dial` | `200` |
| `SetBroadcaster` | `50` |
| `GlobalBroadcast` | `100` |
| `NewBootstrapNode` | `2000` |
| `Bootstrap_Start` | `0` |
| `Bootstrap_Stop` | `0` |
| `Bootstrap_Peers` | `0` |
| `Bootstrap_DialSeed` | `0` |
| `NewConnPool` | `800` |
| `AcquireConn` | `50` |
| `ReleaseConn` | `20` |
| `ClosePool` | `40` |
| `PoolStats` | `10` |
| `NewNATManager` | `500` |
| `NAT_Map` | `0` |
| `NAT_Unmap` | `0` |
| `NAT_ExternalIP` | `0` |
| `DiscoverPeers` | `100` |
| `Connect` | `150` |
| `Disconnect` | `100` |
| `AdvertiseSelf` | `80` |
| `StartDevNet` | `5000` |
| `StartTestNet` | `6000` |
| `MobileMiner_Start` | `3000` |
| `MobileMiner_Stop` | `500` |
| `MobileMiner_Status` | `100` |
| `MobileMiner_SetIntensity` | `50` |


### Replication / Data Availability

Operations related to replication / data availability.


| Opcode | Gas Cost |
|---|---|
| `NewReplicator` | `1200` |
| `ReplicateBlock` | `3000` |
| `RequestMissing` | `400` |
| `Synchronize` | `2500` |
| `Stop` | `300` |
| `NewInitService` | `800` |
| `BootstrapLedger` | `2000` |
| `ShutdownInitService` | `300` |


### Distributed Coordination

Operations related to distributed coordination.


| Opcode | Gas Cost |
|---|---|
| `NewCoordinator` | `1000` |
| `StartCoordinator` | `500` |
| `StopCoordinator` | `500` |
| `BroadcastLedgerHeight` | `300` |
| `DistributeToken` | `500` |
| `NewSyncManager` | `1200` |
| `Sync_Start` | `0` |
| `Sync_Stop` | `0` |
| `Sync_Status` | `0` |
| `SyncOnce` | `800` |


### Roll-ups

Operations related to roll-ups.


| Opcode | Gas Cost |
|---|---|
| `NewAggregator` | `1500` |
| `SubmitBatch` | `1000` |
| `SubmitFraudProof` | `3000` |
| `FinalizeBatch` | `1000` |
| `BatchHeader` | `50` |
| `BatchState` | `30` |
| `BatchTransactions` | `100` |
| `ListBatches` | `200` |
| `PauseAggregator` | `50` |
| `ResumeAggregator` | `50` |
| `AggregatorStatus` | `20` |


### Security / Cryptography

Operations related to security / cryptography.


| Opcode | Gas Cost |
|---|---|
| `AggregateBLSSigs` | `700` |
| `VerifyAggregated` | `800` |
| `CombineShares` | `600` |
| `ComputeMerkleRoot` | `120` |
| `Encrypt` | `150` |
| `Decrypt` | `150` |
| `NewTLSConfig` | `500` |
| `DilithiumKeypair` | `600` |
| `DilithiumSign` | `500` |
| `DilithiumVerify` | `500` |
| `PredictRisk` | `200` |
| `AnomalyScore` | `200` |
| `BuildMerkleTree` | `150` |
| `MerkleProof` | `120` |
| `VerifyMerklePath` | `120` |


### Sharding

Operations related to sharding.


| Opcode | Gas Cost |
|---|---|
| `NewShardCoordinator` | `2000` |
| `SetLeader` | `100` |
| `Leader` | `80` |
| `SubmitCrossShard` | `1500` |
| `Send` | `200` |
| `PullReceipts` | `300` |
| `Reshard` | `3000` |
| `GossipTx` | `500` |
| `RebalanceShards` | `800` |
| `VerticalPartition` | `200` |
| `HorizontalPartition` | `200` |
| `CompressData` | `400` |
| `DecompressData` | `400` |


### Side-chains

Operations related to side-chains.


| Opcode | Gas Cost |
|---|---|
| `InitSidechains` | `1200` |
| `Sidechains` | `60` |
| `Sidechain_Register` | `0` |
| `SubmitHeader` | `800` |
| `VerifyWithdraw` | `400` |
| `VerifyAggregateSig` | `800` |
| `VerifyMerkleProof` | `120` |
| `GetSidechainMeta` | `100` |
| `ListSidechains` | `120` |
| `GetSidechainHeader` | `100` |
| `PauseSidechain` | `300` |
| `ResumeSidechain` | `300` |
| `UpdateSidechainValidators` | `500` |
| `RemoveSidechain` | `600` |


### State-Channels

Operations related to state-channels.


| Opcode | Gas Cost |
|---|---|
| `InitStateChannels` | `800` |
| `Channels` | `60` |
| `OpenChannel` | `1000` |
| `VerifyECDSASignature` | `200` |
| `InitiateClose` | `300` |
| `Challenge` | `400` |
| `Finalize` | `500` |
| `GetChannel` | `80` |
| `ListChannels` | `120` |
| `PauseChannel` | `150` |
| `ResumeChannel` | `150` |
| `CancelClose` | `300` |
| `ForceClose` | `600` |
| `Lightning_OpenChannel` | `800` |
| `Lightning_RoutePayment` | `50` |
| `Lightning_CloseChannel` | `400` |
| `Lightning_ListChannels` | `80` |


### Storage / Marketplace

Operations related to storage / marketplace.


| Opcode | Gas Cost |
|---|---|
| `NewStorage` | `1200` |
| `CreateListing` | `800` |
| `Exists` | `40` |
| `OpenDeal` | `500` |
| `Create` | `800` |
| `CloseDeal` | `500` |
| `Release` | `200` |
| `GetListing` | `100` |
| `ListListings` | `100` |
| `GetDeal` | `100` |
| `ListDeals` | `100` |
| `IPFS_Add` | `0` |
| `IPFS_Get` | `0` |
| `IPFS_Unpin` | `0` |


### General Marketplace

Operations related to general marketplace.


| Opcode | Gas Cost |
|---|---|
| `CreateMarketListing` | `800` |
| `PurchaseItem` | `600` |
| `CancelListing` | `300` |
| `ReleaseFunds` | `200` |
| `GetMarketListing` | `100` |
| `ListMarketListings` | `100` |
| `GetMarketDeal` | `100` |
| `ListMarketDeals` | `100` |


### Tangible assets

Operations related to tangible assets.


| Opcode | Gas Cost |
|---|---|
| `Assets_Register` | `0` |
| `Assets_Transfer` | `0` |
| `Assets_Get` | `0` |
| `Assets_List` | `0` |


### SYN1401 investment tokens

Operations related to syn1401 investment tokens.


| Opcode | Gas Cost |
|---|---|
| `SYN1401_Issue` | `0` |
| `SYN1401_Accrue` | `0` |
| `SYN1401_Redeem` | `0` |
| `SYN1401_Info` | `0` |


### Identity Verification

Operations related to identity verification.


| Opcode | Gas Cost |
|---|---|
| `RegisterIdentity` | `500` |
| `VerifyIdentity` | `100` |
| `RemoveIdentity` | `200` |
| `ListIdentities` | `200` |


### Resource Marketplace

Operations related to resource marketplace.


| Opcode | Gas Cost |
|---|---|
| `ListResource` | `800` |
| `OpenResourceDeal` | `500` |
| `CloseResourceDeal` | `500` |
| `GetResourceListing` | `100` |
| `ListResourceListings` | `100` |
| `GetResourceDeal` | `100` |
| `ListResourceDeals` | `100` |


### Token Standards (constants – zero-cost markers)

Operations related to token standards (constants – zero-cost markers).


| Opcode | Gas Cost |
|---|---|
| `StdSYN10` | `1` |
| `StdSYN11` | `11` |
| `StdSYN12` | `12` |
| `StdSYN20` | `2` |
| `StdSYN70` | `7` |
| `StdSYN130` | `13` |
| `StdSYN131` | `13` |
| `StdSYN200` | `20` |
| `StdSYN223` | `22` |
| `StdSYN300` | `30` |
| `StdSYN500` | `50` |
| `StdSYN600` | `60` |
| `StdSYN700` | `70` |
| `StdSYN721` | `72` |
| `StdSYN722` | `72` |
| `StdSYN800` | `80` |
| `StdSYN845` | `84` |
| `StdSYN900` | `90` |
| `StdSYN1000` | `100` |
| `StdSYN1100` | `110` |
| `StdSYN1155` | `115` |
| `StdSYN1200` | `120` |
| `StdSYN1300` | `130` |
| `StdSYN1401` | `140` |
| `StdSYN1500` | `150` |
| `StdSYN1600` | `160` |
| `StdSYN1700` | `170` |
| `StdSYN1800` | `180` |
| `StdSYN1900` | `190` |
| `StdSYN1967` | `196` |
| `StdSYN2100` | `210` |
| `StdSYN2200` | `220` |
| `StdSYN2369` | `236` |
| `StdSYN2400` | `240` |
| `StdSYN2500` | `250` |
| `StdSYN2600` | `260` |
| `StdSYN2700` | `270` |
| `StdSYN2800` | `280` |
| `StdSYN2900` | `290` |
| `StdSYN3000` | `300` |
| `StdSYN3100` | `310` |
| `StdSYN3200` | `320` |
| `StdSYN3300` | `330` |
| `StdSYN3400` | `340` |
| `StdSYN3500` | `350` |
| `StdSYN3600` | `360` |
| `StdSYN3700` | `370` |
| `StdSYN3800` | `380` |
| `StdSYN3900` | `390` |
| `StdSYN4200` | `420` |
| `StdSYN4300` | `430` |
| `StdSYN4700` | `470` |
| `StdSYN4900` | `490` |
| `StdSYN5000` | `500` |


### Token Utilities

Operations related to token utilities.


| Opcode | Gas Cost |
|---|---|
| `ID` | `40` |
| `Meta` | `40` |
| `Allowance` | `40` |
| `Approve` | `80` |
| `Add` | `60` |
| `Sub` | `60` |
| `Get` | `40` |
| `transfer` | `210` |
| `Calculate` | `80` |
| `RegisterToken` | `800` |
| `NewBalanceTable` | `500` |
| `Set` | `60` |
| `RefundGas` | `10` |
| `PopUint32` | `3` |
| `PopAddress` | `30` |
| `PopUint64` | `6` |
| `PushBool` | `30` |
| `Push` | `30` |
| `Len` | `20` |
| `InitTokens` | `800` |
| `GetRegistryTokens` | `40` |
| `TokenManager_Create` | `0` |
| `TokenManager_Transfer` | `0` |
| `TokenManager_Mint` | `0` |
| `TokenManager_Burn` | `0` |
| `TokenManager_Approve` | `0` |
| `TokenManager_BalanceOf` | `0` |
| `InsuranceToken_IssuePolicy` | `0` |
| `InsuranceToken_ClaimPolicy` | `0` |
| `InsuranceToken_UpdatePolicy` | `0` |
| `InsuranceToken_GetPolicy` | `0` |
| `InsuranceToken_CancelPolicy` | `0` |
| `ID` | `40` |
| `Meta` | `40` |
| `Allowance` | `40` |
| `Approve` | `80` |
| `Add` | `60` |
| `Sub` | `60` |
| `Get` | `40` |
| `transfer` | `210` |
| `Calculate` | `80` |
| `RegisterToken` | `800` |
| `NewBalanceTable` | `500` |
| `Set` | `60` |
| `RefundGas` | `10` |
| `PopUint32` | `3` |
| `PopAddress` | `30` |
| `PopUint64` | `6` |
| `PushBool` | `30` |
| `Push` | `30` |
| `Len` | `20` |
| `InitTokens` | `800` |
| `GetRegistryTokens` | `40` |
| `TokenManager_Create` | `0` |
| `TokenManager_Transfer` | `0` |
| `TokenManager_Mint` | `0` |
| `TokenManager_Burn` | `0` |
| `TokenManager_Approve` | `0` |
| `TokenManager_BalanceOf` | `0` |
| `SYN1155_BatchTransfer` | `400` |
| `SYN1155_BatchBalance` | `200` |
| `SYN1155_SetApprovalForAll` | `100` |
| `SYN1155_IsApprovedForAll` | `50` |
| `SYN1155_RegisterHook` | `50` |
| `ID` | `40` |
| `Meta` | `40` |
| `Allowance` | `40` |
| `Approve` | `80` |
| `Add` | `60` |
| `Sub` | `60` |
| `Get` | `40` |
| `transfer` | `210` |
| `Calculate` | `80` |
| `RegisterToken` | `800` |
| `NewBalanceTable` | `500` |
| `Set` | `60` |
| `RefundGas` | `10` |
| `PopUint32` | `3` |
| `PopAddress` | `30` |
| `PopUint64` | `6` |
| `PushBool` | `30` |
| `Push` | `30` |
| `Len` | `20` |
| `InitTokens` | `800` |
| `GetRegistryTokens` | `40` |
| `TokenManager_Create` | `0` |
| `TokenManager_Transfer` | `0` |
| `TokenManager_Mint` | `0` |
| `TokenManager_Burn` | `0` |
| `TokenManager_Approve` | `0` |
| `TokenManager_BalanceOf` | `0` |
| `EmpToken_CreateContract` | `0` |
| `EmpToken_PaySalary` | `0` |
| `EmpToken_UpdateBenefits` | `0` |
| `EmpToken_Terminate` | `0` |
| `EmpToken_GetContract` | `0` |
| `Tokens_Pause` | `80` |
| `Tokens_Unpause` | `80` |
| `Tokens_IsPaused` | `40` |
| `Tokens_BulkTransfer` | `300` |
| `Tokens_BulkApprove` | `200` |
| `Tokens_TransferWithMemo` | `220` |
| `ID` | `40` |
| `Meta` | `40` |
| `Allowance` | `40` |
| `Approve` | `80` |
| `Add` | `60` |
| `Sub` | `60` |
| `Get` | `40` |
| `transfer` | `210` |
| `Calculate` | `80` |
| `RegisterToken` | `800` |
| `NewBalanceTable` | `500` |
| `Set` | `60` |
| `RefundGas` | `10` |
| `PopUint32` | `3` |
| `PopAddress` | `30` |
| `PopUint64` | `6` |
| `PushBool` | `30` |
| `Push` | `30` |
| `Len` | `20` |
| `InitTokens` | `800` |
| `GetRegistryTokens` | `40` |
| `TokenManager_Create` | `0` |
| `TokenManager_Transfer` | `0` |
| `TokenManager_Mint` | `0` |
| `TokenManager_Burn` | `0` |
| `TokenManager_Approve` | `0` |
| `TokenManager_BalanceOf` | `0` |
| `SYN4900_RegisterAsset` | `0` |
| `SYN4900_UpdateStatus` | `0` |
| `SYN4900_TransferAsset` | `0` |
| `SYN4900_RecordInvestment` | `0` |
| `SYN4900_GetInvestment` | `0` |
| `ID` | `40` |
| `Meta` | `40` |
| `Allowance` | `40` |
| `Approve` | `80` |
| `Add` | `60` |
| `Sub` | `60` |
| `Get` | `40` |
| `transfer` | `210` |
| `Calculate` | `80` |
| `RegisterToken` | `800` |
| `NewBalanceTable` | `500` |
| `Set` | `60` |
| `RefundGas` | `10` |
| `PopUint32` | `3` |
| `PopAddress` | `30` |
| `PopUint64` | `6` |
| `PushBool` | `30` |
| `Push` | `30` |
| `Len` | `20` |
| `InitTokens` | `800` |
| `GetRegistryTokens` | `40` |
| `TokenManager_Create` | `0` |
| `TokenManager_Transfer` | `0` |
| `TokenManager_Mint` | `0` |
| `TokenManager_Burn` | `0` |
| `TokenManager_Approve` | `0` |
| `TokenManager_BalanceOf` | `0` |
| `Forex_UpdateRate` | `0` |
| `Forex_OpenPosition` | `0` |
| `Forex_ClosePosition` | `0` |
| `ETF_UpdatePrice` | `0` |
| `ETF_FractionalMint` | `0` |
| `ETF_FractionalBurn` | `0` |
| `ETF_Info` | `0` |
| `SYN3500_UpdateRate` | `0` |
| `SYN3500_Mint` | `0` |
| `SYN3500_Redeem` | `0` |
| `SYN3500_Info` | `0` |
| `Syn3200_CreateBill` | `0` |
| `Syn3200_PayFraction` | `0` |
| `Syn3200_AdjustAmount` | `0` |
| `Syn3200_GetBill` | `0` |
| `SYN2369_CreateItem` | `250` |
| `SYN2369_TransferItem` | `210` |
| `SYN2369_UpdateAttrs` | `150` |
| `Mint721` | `0` |
| `Transfer721` | `0` |
| `Burn721` | `0` |
| `Metadata721` | `0` |
| `UpdateMetadata721` | `0` |
| `SYN223_SafeTransfer` | `0` |
| `SYN223_AddWhitelist` | `0` |
| `SYN223_RemoveWhitelist` | `0` |
| `SYN223_AddBlacklist` | `0` |
| `SYN223_RemoveBlacklist` | `0` |
| `SYN223_SetRequiredSigs` | `0` |
| `SYN223_IsWhitelisted` | `0` |
| `SYN223_IsBlacklisted` | `0` |
| `SYN131UpdateValuation` | `0` |
| `SYN131RecordSale` | `0` |
| `SYN131AddRental` | `0` |
| `SYN131IssueLicense` | `0` |
| `SYN131TransferShare` | `0` |
| `SYN130_UpdateValuation` | `0` |
| `SYN130_RecordSale` | `0` |
| `SYN130_StartLease` | `0` |
| `SYN130_EndLease` | `0` |
| `Benefit_Issue` | `0` |
| `Benefit_Claim` | `0` |
| `Benefit_Record` | `0` |
| `Benefit_List` | `0` |
| `CharityToken_Donate` | `0` |
| `CharityToken_Release` | `0` |
| `CharityToken_Progress` | `0` |
| `SYN11_Issue` | `0` |
| `SYN11_Redeem` | `0` |
| `SYN11_UpdateCoupon` | `0` |
| `SYN11_PayCoupon` | `0` |
| `ID` | `40` |
| `Meta` | `40` |
| `Allowance` | `40` |
| `Approve` | `80` |
| `Add` | `60` |
| `Sub` | `60` |
| `Get` | `40` |
| `transfer` | `210` |
| `Calculate` | `80` |
| `RegisterToken` | `800` |
| `NewBalanceTable` | `500` |
| `Set` | `60` |
| `RefundGas` | `10` |
| `PopUint32` | `3` |
| `PopAddress` | `30` |
| `PopUint64` | `6` |
| `PushBool` | `30` |
| `Push` | `30` |
| `Len` | `20` |
| `InitTokens` | `800` |
| `GetRegistryTokens` | `40` |
| `TokenManager_Create` | `0` |
| `TokenManager_Transfer` | `0` |
| `TokenManager_Mint` | `0` |
| `TokenManager_Burn` | `0` |
| `TokenManager_Approve` | `0` |
| `TokenManager_BalanceOf` | `0` |
| `SYN1967_UpdatePrice` | `0` |
| `SYN1967_CurrentPrice` | `0` |
| `SYN1967_PriceHistory` | `0` |
| `SYN1967_AddCertification` | `0` |
| `SYN1967_AddTrace` | `0` |
| `TokenManager_CreateSYN1967` | `0` |
| `LegalToken_New` | `0` |
| `LegalToken_AddSignature` | `0` |
| `LegalToken_RevokeSignature` | `0` |
| `LegalToken_UpdateStatus` | `0` |
| `LegalToken_StartDispute` | `0` |
| `LegalToken_ResolveDispute` | `0` |
| `SYN1100_AddRecord` | `50` |
| `SYN1100_GrantAccess` | `30` |
| `SYN1100_RevokeAccess` | `20` |
| `SYN1100_GetRecord` | `40` |
| `SYN1100_TransferOwnership` | `50` |
| `ID` | `40` |
| `Meta` | `40` |
| `Allowance` | `40` |
| `Approve` | `80` |
| `Add` | `60` |
| `Sub` | `60` |
| `Get` | `40` |
| `transfer` | `210` |
| `Calculate` | `80` |
| `RegisterToken` | `800` |
| `NewBalanceTable` | `500` |
| `Set` | `60` |
| `RefundGas` | `10` |
| `PopUint32` | `3` |
| `PopAddress` | `30` |
| `PopUint64` | `6` |
| `PushBool` | `30` |
| `Push` | `30` |
| `Len` | `20` |
| `InitTokens` | `800` |
| `GetRegistryTokens` | `40` |
| `TokenManager_Create` | `0` |
| `TokenManager_Transfer` | `0` |
| `TokenManager_Mint` | `0` |
| `TokenManager_Burn` | `0` |
| `TokenManager_Approve` | `0` |
| `TokenManager_BalanceOf` | `0` |
| `SupplyChain_RegisterAsset` | `0` |
| `SupplyChain_UpdateLocation` | `0` |
| `SupplyChain_UpdateStatus` | `0` |
| `SupplyChain_TransferAsset` | `0` |
| `ID` | `40` |
| `Meta` | `40` |
| `Allowance` | `40` |
| `Approve` | `80` |
| `Add` | `60` |
| `Sub` | `60` |
| `Get` | `40` |
| `transfer` | `210` |
| `Calculate` | `80` |
| `RegisterToken` | `800` |
| `NewBalanceTable` | `500` |
| `Set` | `60` |
| `RefundGas` | `10` |
| `PopUint32` | `3` |
| `PopAddress` | `30` |
| `PopUint64` | `6` |
| `PushBool` | `30` |
| `Push` | `30` |
| `Len` | `20` |
| `InitTokens` | `800` |
| `GetRegistryTokens` | `40` |
| `TokenManager_Create` | `0` |
| `TokenManager_Transfer` | `0` |
| `TokenManager_Mint` | `0` |
| `TokenManager_Burn` | `0` |
| `TokenManager_Approve` | `0` |
| `TokenManager_BalanceOf` | `0` |
| `SYN70_RegisterAsset` | `0` |
| `SYN70_TransferAsset` | `0` |
| `SYN70_UpdateAttributes` | `0` |
| `SYN70_RecordAchievement` | `0` |
| `SYN70_GetAsset` | `0` |
| `SYN70_ListAssets` | `0` |
| `MusicRoyalty_AddRevenue` | `80` |
| `MusicRoyalty_Distribute` | `100` |
| `MusicRoyalty_UpdateInfo` | `50` |
| `ID` | `40` |
| `Meta` | `40` |
| `Allowance` | `40` |
| `Approve` | `80` |
| `Add` | `60` |
| `Sub` | `60` |
| `Get` | `40` |
| `transfer` | `210` |
| `Calculate` | `80` |
| `RegisterToken` | `800` |
| `NewBalanceTable` | `500` |
| `Set` | `60` |
| `RefundGas` | `10` |
| `PopUint32` | `3` |
| `PopAddress` | `30` |
| `PopUint64` | `6` |
| `PushBool` | `30` |
| `Push` | `30` |
| `Len` | `20` |
| `InitTokens` | `800` |
| `GetRegistryTokens` | `40` |
| `TokenManager_Create` | `0` |
| `TokenManager_Transfer` | `0` |
| `TokenManager_Mint` | `0` |
| `TokenManager_Burn` | `0` |
| `TokenManager_Approve` | `0` |
| `TokenManager_BalanceOf` | `0` |
| `SYN600_Stake` | `200` |
| `SYN600_Unstake` | `200` |
| `SYN600_AddEngagement` | `50` |
| `SYN600_EngagementOf` | `40` |
| `SYN600_DistributeRewards` | `500` |
| `SYN2100_RegisterDocument` | `0` |
| `SYN2100_FinanceDocument` | `0` |
| `SYN2100_GetDocument` | `0` |
| `SYN2100_ListDocuments` | `0` |
| `SYN2100_AddLiquidity` | `0` |
| `SYN2100_RemoveLiquidity` | `0` |
| `SYN2100_LiquidityOf` | `0` |
| `ID` | `40` |
| `Meta` | `40` |
| `Allowance` | `40` |
| `Approve` | `80` |
| `Add` | `60` |
| `Sub` | `60` |
| `Get` | `40` |
| `transfer` | `210` |
| `Calculate` | `80` |
| `RegisterToken` | `800` |
| `NewBalanceTable` | `500` |
| `Set` | `60` |
| `RefundGas` | `10` |
| `PopUint32` | `3` |
| `PopAddress` | `30` |
| `PopUint64` | `6` |
| `PushBool` | `30` |
| `Push` | `30` |
| `Len` | `20` |
| `InitTokens` | `800` |
| `GetRegistryTokens` | `40` |
| `TokenManager_Create` | `0` |
| `TokenManager_Transfer` | `0` |
| `TokenManager_Mint` | `0` |
| `TokenManager_Burn` | `0` |
| `TokenManager_Approve` | `0` |
| `TokenManager_BalanceOf` | `0` |
| `SYN500_GrantAccess` | `0` |
| `SYN500_UpdateAccess` | `0` |
| `SYN500_RevokeAccess` | `0` |
| `SYN500_RecordUsage` | `0` |
| `SYN500_RedeemReward` | `0` |
| `SYN500_RewardBalance` | `0` |
| `SYN500_Usage` | `0` |
| `SYN500_AccessInfo` | `0` |
| `SYN800_RegisterAsset` | `10` |
| `SYN800_UpdateValuation` | `10` |
| `SYN800_GetAsset` | `5` |
| `IDToken_Register` | `0` |
| `IDToken_Verify` | `0` |
| `IDToken_Get` | `0` |
| `IDToken_Logs` | `0` |
| `SYN1200_AddBridge` | `0` |
| `SYN1200_AtomicSwap` | `0` |
| `SYN1200_CompleteSwap` | `0` |
| `SYN1200_GetSwap` | `0` |
| `RegisterIPAsset` | `800` |
| `TransferIPOwnership` | `300` |
| `CreateLicense` | `400` |
| `RevokeLicense` | `200` |
| `RecordRoyalty` | `100` |
| `Event_Create` | `100` |
| `Event_IssueTicket` | `150` |
| `Event_Transfer` | `210` |
| `Event_Verify` | `50` |
| `Event_Use` | `40` |
| `Tokens_RecordEmission` | `0` |
| `Tokens_RecordOffset` | `0` |
| `Tokens_NetBalance` | `0` |
| `Tokens_ListRecords` | `0` |
| `Edu_RegisterCourse` | `0` |
| `Edu_IssueCredit` | `0` |
| `Edu_VerifyCredit` | `0` |
| `Edu_RevokeCredit` | `0` |
| `Edu_GetCredit` | `0` |
| `Edu_ListCredits` | `0` |
| `Tokens_CreateSYN2200` | `0` |
| `TokensCreateSYN1000` | `0` |
| `SYN1000_AddReserve` | `0` |
| `SYN1000_RemoveReserve` | `0` |
| `SYN1000_SetPrice` | `0` |
| `SYN1000_ReserveValue` | `0` |
| `Tokens_SendPayment` | `0` |
| `Tokens_GetPayment` | `0` |
| `DataToken_UpdateMeta` | `0` |
| `DataToken_SetPrice` | `0` |
| `DataToken_GrantAccess` | `0` |
| `DataToken_RevokeAccess` | `0` |
| `SYN845_IssueDebt` | `0` |
| `SYN845_RecordPayment` | `0` |
| `SYN845_AdjustInterest` | `0` |
| `SYN845_MarkDefault` | `0` |
| `SYN845_GetDebt` | `0` |
| `SYN845_ListDebts` | `0` |
| `SYN5000_PlaceBet` | `0` |
| `SYN5000_ResolveBet` | `0` |
| `SYN5000_BetInfo` | `0` |


### Transactions

Operations related to transactions.


| Opcode | Gas Cost |
|---|---|
| `VerifySig` | `350` |
| `ValidateTx` | `500` |
| `NewTxPool` | `1200` |
| `AddTx` | `600` |
| `PickTxs` | `150` |
| `TxPoolSnapshot` | `80` |
| `Exec_Begin` | `0` |
| `Exec_RunTx` | `0` |
| `Exec_Finalize` | `0` |
| `EncryptTxPayload` | `350` |
| `DecryptTxPayload` | `300` |
| `SubmitPrivateTx` | `650` |
| `EncodeEncryptedHex` | `30` |
| `ReverseTransaction` | `1000` |
| `NewTxDistributor` | `800` |
| `DistributeFees` | `150` |


### corrections are applied at run-time by the VM).

Operations related to corrections are applied at run-time by the vm)..


| Opcode | Gas Cost |
|---|---|
| `Short` | `0` |
| `BytesToAddress` | `0` |
| `Pop` | `0` |
| `opADD` | `0` |
| `opMUL` | `0` |
| `opSUB` | `0` |
| `OpDIV` | `0` |
| `opSDIV` | `0` |
| `opMOD` | `0` |
| `opSMOD` | `0` |
| `opADDMOD` | `0` |
| `opMULMOD` | `0` |
| `opEXP` | `1` |
| `opSIGNEXTEND` | `0` |
| `opLT` | `0` |
| `opGT` | `0` |
| `opSLT` | `0` |
| `opSGT` | `0` |
| `opEQ` | `0` |
| `opISZERO` | `0` |
| `opAND` | `0` |
| `opOR` | `0` |
| `opXOR` | `0` |
| `opNOT` | `0` |
| `opBYTE` | `0` |
| `opSHL` | `0` |
| `opSHR` | `0` |
| `opSAR` | `0` |
| `opECRECOVER` | `70` |
| `opEXTCODESIZE` | `70` |
| `opEXTCODECOPY` | `70` |
| `opEXTCODEHASH` | `70` |
| `opRETURNDATASIZE` | `0` |
| `opRETURNDATACOPY` | `70` |
| `opMLOAD` | `0` |
| `opMSTORE` | `0` |
| `opMSTORE8` | `0` |
| `opCALLDATALOAD` | `0` |
| `opCALLDATASIZE` | `0` |
| `opCALLDATACOPY` | `70` |
| `opCODESIZE` | `0` |
| `opCODECOPY` | `70` |
| `opJUMP` | `0` |
| `opJUMPI` | `1` |
| `opPC` | `0` |
| `opMSIZE` | `0` |
| `opGAS` | `0` |
| `opJUMPDEST` | `0` |
| `opSHA256` | `25` |
| `opKECCAK256` | `25` |
| `opRIPEMD160` | `16` |
| `opBLAKE2B256` | `0` |
| `opADDRESS` | `0` |
| `opCALLER` | `0` |
| `opORIGIN` | `0` |
| `opCALLVALUE` | `0` |
| `opGASPRICE` | `0` |
| `opNUMBER` | `0` |
| `opTIMESTAMP` | `0` |
| `opDIFFICULTY` | `0` |
| `opGASLIMIT` | `0` |
| `opCHAINID` | `0` |
| `opBLOCKHASH` | `2` |
| `opBALANCE` | `40` |
| `opSELFBALANCE` | `0` |
| `opLOG0` | `0` |
| `opLOG1` | `0` |
| `opLOG2` | `0` |
| `opLOG3` | `0` |
| `opLOG4` | `0` |
| `logN` | `200` |
| `opCREATE` | `3200` |
| `opCALL` | `70` |
| `opCALLCODE` | `70` |
| `opDELEGATECALL` | `70` |
| `opSTATICCALL` | `70` |
| `opRETURN` | `0` |
| `opREVERT` | `0` |
| `opSTOP` | `0` |
| `opSELFDESTRUCT` | `500` |


### Shared accounting ops

Operations related to shared accounting ops.


| Opcode | Gas Cost |
|---|---|
| `TransferVM` | `210` |


### Virtual Machine Internals

Operations related to virtual machine internals.


| Opcode | Gas Cost |
|---|---|
| `BurnVM` | `210` |
| `BurnLP` | `210` |
| `MintLP` | `210` |
| `NewInMemory` | `50` |
| `CallCode` | `70` |
| `CallContract` | `70` |
| `StaticCallVM` | `70` |
| `GetBalance` | `40` |
| `GetTokenBalance` | `40` |
| `SetTokenBalance` | `50` |
| `GetTokenSupply` | `50` |
| `SetBalance` | `50` |
| `DelegateCall` | `70` |
| `GetToken` | `40` |
| `NewMemory` | `50` |
| `Read` | `0` |
| `Write` | `0` |
| `LenVM` | `0` |
| `Call` | `70` |
| `SelectVM` | `100` |
| `CreateContract` | `3200` |
| `AddLog` | `37` |
| `GetCode` | `20` |
| `GetCodeHash` | `20` |
| `MintTokenVM` | `200` |
| `PrefixIterator` | `50` |
| `NonceOf` | `40` |
| `GetState` | `40` |
| `SetState` | `50` |
| `HasState` | `40` |
| `DeleteState` | `50` |
| `NewGasMeter` | `50` |
| `SelfDestructVM` | `500` |
| `Remaining` | `0` |
| `Consume` | `0` |
| `ExecuteVM` | `200` |
| `NewSuperLightVM` | `50` |
| `NewLightVM` | `80` |
| `NewHeavyVM` | `120` |
| `ExecuteSuperLight` | `100` |
| `ExecuteLight` | `150` |
| `ExecuteHeavy` | `200` |


### Sandbox management

Operations related to sandbox management.


| Opcode | Gas Cost |
|---|---|
| `VM_SandboxStart` | `0` |
| `VM_SandboxStop` | `0` |
| `VM_SandboxReset` | `0` |
| `VM_SandboxStatus` | `0` |
| `VM_SandboxList` | `0` |


### Smart Legal Contracts

Operations related to smart legal contracts.


| Opcode | Gas Cost |
|---|---|
| `Legal_Register` | `0` |
| `Legal_Sign` | `0` |
| `Legal_Revoke` | `0` |
| `Legal_Info` | `0` |
| `Legal_List` | `0` |


### Plasma

Operations related to plasma.


| Opcode | Gas Cost |
|---|---|
| `InitPlasma` | `800` |
| `Plasma_Deposit` | `0` |
| `Plasma_Withdraw` | `0` |
| `Plasma_StartExit` | `0` |
| `Plasma_FinalizeExit` | `0` |
| `Plasma_GetExit` | `0` |
| `Plasma_ListExits` | `0` |


### Gaming

Operations related to gaming.


| Opcode | Gas Cost |
|---|---|
| `CreateGame` | `800` |
| `JoinGame` | `400` |
| `FinishGame` | `600` |
| `GetGame` | `100` |
| `ListGames` | `200` |


### Messaging / Queue Management

Operations related to messaging / queue management.


| Opcode | Gas Cost |
|---|---|
| `NewMessageQueue` | `500` |
| `EnqueueMessage` | `50` |
| `DequeueMessage` | `50` |
| `BroadcastNextMessage` | `100` |
| `ProcessNextMessage` | `200` |
| `QueueLength` | `10` |
| `ClearQueue` | `20` |


### Wallet / Key-Management

Operations related to wallet / key-management.


| Opcode | Gas Cost |
|---|---|
| `NewRandomWallet` | `1000` |
| `WalletFromMnemonic` | `500` |
| `NewHDWalletFromSeed` | `600` |
| `PrivateKey` | `40` |
| `NewAddress` | `50` |
| `SignTx` | `300` |
| `RegisterIDWallet` | `800` |
| `IsIDWalletRegistered` | `50` |
| `NewOffChainWallet` | `800` |
| `OffChainWalletFromMnemonic` | `500` |
| `SignOffline` | `250` |
| `StoreSignedTx` | `30` |
| `LoadSignedTx` | `30` |
| `BroadcastSignedTx` | `100` |


### Access Control

Operations related to access control. Addresses are validated before roles are
assigned or revoked.


| Opcode | Gas Cost |
|---|---|
| `GrantRole` | `100` |
| `RevokeRole` | `100` |
| `HasRole` | `30` |
| `ListRoles` | `50` |


### Geolocation Network

Operations related to geolocation network.


| Opcode | Gas Cost |
|---|---|
| `RegisterLocation` | `200` |
| `GetLocation` | `50` |
| `ListLocations` | `100` |
| `NodesInRadius` | `150` |


### Geospatial Node

Operations related to geospatial node.


| Opcode | Gas Cost |
|---|---|
| `NewGeospatialNode` | `1000` |
| `RegisterGeoData` | `200` |
| `TransformCoordinates` | `300` |
| `AddGeofence` | `150` |
| `InGeofence` | `100` |
| `QueryGeoData` | `200` |


### Firewall

Operations related to firewall.


| Opcode | Gas Cost |
|---|---|
| `NewFirewall` | `400` |
| `Firewall_BlockAddress` | `0` |
| `Firewall_UnblockAddress` | `0` |
| `Firewall_IsAddressBlocked` | `0` |
| `Firewall_BlockToken` | `0` |
| `Firewall_UnblockToken` | `0` |
| `Firewall_IsTokenBlocked` | `0` |
| `Firewall_BlockIP` | `0` |
| `Firewall_UnblockIP` | `0` |
| `Firewall_IsIPBlocked` | `0` |
| `Firewall_ListRules` | `0` |
| `Firewall_CheckTx` | `0` |


### RPC / WebRTC

Operations related to rpc / webrtc.


| Opcode | Gas Cost |
|---|---|
| `NewRPCWebRTC` | `1000` |
| `RPC_Serve` | `0` |
| `RPC_Close` | `0` |
| `RPC_ConnectPeer` | `0` |
| `RPC_Broadcast` | `0` |


### Plasma Management

Operations related to plasma management.


| Opcode | Gas Cost |
|---|---|
| `Plasma_SubmitBlock` | `0` |
| `Plasma_GetBlock` | `0` |


### Resource Management

Operations related to resource management.


| Opcode | Gas Cost |
|---|---|
| `SetQuota` | `100` |
| `GetQuota` | `50` |
| `ChargeResources` | `200` |
| `ReleaseResources` | `100` |


### Distribution

Operations related to distribution.


| Opcode | Gas Cost |
|---|---|
| `NewDistributor` | `100` |
| `BatchTransfer` | `400` |
| `Airdrop` | `300` |
| `DistributeEven` | `200` |


### Carbon Credit System

Operations related to carbon credit system.


| Opcode | Gas Cost |
|---|---|
| `InitCarbonEngine` | `800` |
| `Carbon` | `200` |
| `RegisterProject` | `500` |
| `IssueCredits` | `500` |
| `RetireCredits` | `300` |
| `ProjectInfo` | `100` |
| `ListProjects` | `100` |


### Pension Token System

Operations related to pension token system.


| Opcode | Gas Cost |
|---|---|
| `InitPensionEngine` | `800` |
| `Pension` | `200` |
| `RegisterPlan` | `500` |
| `Contribute` | `500` |
| `Withdraw` | `300` |
| `PlanInfo` | `100` |
| `ListPlans` | `100` |


### Government Grant Tokens

Operations related to government grant tokens.


| Opcode | Gas Cost |
|---|---|
| `InitGrantEngine` | `800` |
| `GrantEngine` | `200` |
| `GrantToken_Create` | `500` |
| `GrantToken_Disburse` | `500` |
| `GrantToken_Info` | `100` |
| `GrantToken_List` | `100` |
| `InitSYN10` | `800` |
| `SYN10` | `200` |
| `SYN10_UpdateRate` | `100` |
| `SYN10_Info` | `50` |
| `SYN10_Mint` | `200` |
| `SYN10_Burn` | `200` |
| `SYN3500_UpdateRate` | `100` |
| `SYN3500_Info` | `50` |
| `SYN3500_Mint` | `200` |
| `SYN3500_Redeem` | `200` |
| `InitCarbonEngine` | `800` |
| `Carbon` | `200` |
| `RegisterProject` | `500` |
| `IssueCredits` | `500` |
| `RetireCredits` | `300` |
| `ProjectInfo` | `100` |
| `ListProjects` | `100` |
| `AddVerification` | `200` |
| `ListVerifications` | `100` |


### Energy Token System

Operations related to energy token system.


| Opcode | Gas Cost |
|---|---|
| `InitEnergyEngine` | `800` |
| `Energy` | `200` |
| `RegisterEnergyAsset` | `500` |
| `TransferEnergyAsset` | `300` |
| `RecordSustainability` | `100` |
| `EnergyAssetInfo` | `100` |
| `ListEnergyAssets` | `100` |


### Finalization Management

Operations related to finalization management.


| Opcode | Gas Cost |
|---|---|
| `NewFinalizationManager` | `800` |
| `FinalizeBlock` | `400` |
| `FinalizeBatchManaged` | `350` |
| `FinalizeChannelManaged` | `350` |
| `RegisterRecovery` | `500` |
| `RecoverAccount` | `800` |


### DeFi

Operations related to defi.


| Opcode | Gas Cost |
|---|---|
| `DeFi_CreateInsurance` | `0` |
| `DeFi_ClaimInsurance` | `0` |
| `DeFi_PlaceBet` | `0` |
| `DeFi_SettleBet` | `0` |
| `DeFi_StartCrowdfund` | `0` |
| `DeFi_Contribute` | `0` |
| `DeFi_FinalizeCrowdfund` | `0` |
| `DeFi_CreatePrediction` | `0` |
| `DeFi_VotePrediction` | `0` |
| `DeFi_ResolvePrediction` | `0` |
| `DeFi_RequestLoan` | `0` |
| `DeFi_RepayLoan` | `0` |
| `DeFi_StartYieldFarm` | `0` |
| `DeFi_Stake` | `0` |
| `DeFi_Unstake` | `0` |
| `DeFi_CreateSynthetic` | `0` |
| `DeFi_MintSynthetic` | `0` |
| `DeFi_BurnSynthetic` | `0` |


### Binary Tree Operations

Operations related to binary tree operations.


| Opcode | Gas Cost |
|---|---|
| `BinaryTreeNew` | `500` |
| `BinaryTreeInsert` | `600` |
| `BinaryTreeSearch` | `400` |
| `BinaryTreeDelete` | `600` |
| `BinaryTreeInOrder` | `300` |


### Regulatory Management

Operations related to regulatory management.


| Opcode | Gas Cost |
|---|---|
| `InitRegulatory` | `400` |
| `RegisterRegulator` | `600` |
| `GetRegulator` | `200` |
| `ListRegulators` | `200` |
| `EvaluateRuleSet` | `500` |
| `NewRegulatoryNode` | `2000` |
| `RegNode_Start` | `800` |
| `RegNode_Stop` | `400` |
| `RegNode_Peers` | `100` |
| `RegNode_DialSeed` | `200` |
| `RegNode_VerifyTx` | `500` |
| `RegNode_KYC` | `400` |
| `RegNode_EraseKYC` | `300` |
| `RegNode_RiskScore` | `200` |
| `RegNode_GenReport` | `600` |
| `InitRegulatory` | `400` |
| `RegisterRegulator` | `600` |
| `GetRegulator` | `200` |
| `ListRegulators` | `200` |
| `EvaluateRuleSet` | `500` |
| `NewGovAuthorityNode` | `1000` |
| `Gov_CheckCompliance` | `500` |
| `Gov_EnforceRegulation` | `700` |
| `Gov_InterfaceRegulator` | `400` |
| `Gov_UpdateLegalFramework` | `800` |
| `Gov_AuditTrail` | `300` |


### Polls Management

Operations related to polls management.


| Opcode | Gas Cost |
|---|---|
| `CreatePoll` | `800` |
| `VotePoll` | `300` |
| `ClosePoll` | `200` |
| `GetPoll` | `50` |
| `ListPolls` | `100` |


### Feedback System

Operations related to feedback system.


| Opcode | Gas Cost |
|---|---|
| `InitFeedback` | `800` |
| `Feedback_Submit` | `0` |
| `Feedback_Get` | `0` |
| `Feedback_List` | `0` |
| `Feedback_Reward` | `0` |


### Forum

Operations related to forum.


| Opcode | Gas Cost |
|---|---|
| `ForumCreateThread` | `0` |
| `ForumGetThread` | `0` |
| `ForumListThreads` | `0` |
| `ForumAddComment` | `0` |
| `ForumListComments` | `0` |


### Blockchain Compression

Operations related to blockchain compression.


| Opcode | Gas Cost |
|---|---|
| `CompressLedger` | `600` |
| `DecompressLedger` | `600` |
| `SaveCompressedSnapshot` | `800` |
| `LoadCompressedSnapshot` | `800` |


### Biometrics Authentication

Operations related to biometrics authentication.


| Opcode | Gas Cost |
|---|---|
| `Bio_Enroll` | `0` |
| `Bio_Verify` | `0` |
| `Bio_Delete` | `0` |
| `BSN_Register` | `500` |
| `BSN_VerifyTx` | `400` |
| `BSN_Remove` | `200` |


### System Health & Logging

Operations related to system health & logging.


| Opcode | Gas Cost |
|---|---|
| `NewHealthLogger` | `800` |
| `MetricsSnapshot` | `100` |
| `LogEvent` | `50` |
| `RotateLogs` | `400` |


### Workflow / Key-Management

Operations related to workflow / key-management.


| Opcode | Gas Cost |
|---|---|
| `NewWorkflow` | `1500` |
| `AddWorkflowAction` | `200` |
| `SetWorkflowTrigger` | `100` |
| `SetWebhook` | `100` |
| `ExecuteWorkflow` | `500` |
| `ListWorkflows` | `50` |


### Swarm

Operations related to swarm.


| Opcode | Gas Cost |
|---|---|
| `NewSwarm` | `1000` |
| `Swarm_AddNode` | `0` |
| `Swarm_RemoveNode` | `0` |
| `Swarm_BroadcastTx` | `0` |
| `Swarm_Start` | `0` |
| `Swarm_Stop` | `0` |
| `Swarm_Peers` | `0` |


### Real Estate

Operations related to real estate.


| Opcode | Gas Cost |
|---|---|
| `RegisterProperty` | `400` |
| `TransferProperty` | `350` |
| `GetProperty` | `100` |
| `ListProperties` | `150` |


### Rental Management

Operations related to rental management.


| Opcode | Gas Cost |
|---|---|
| `RegisterRentalAgreement` | `400` |
| `PayRent` | `200` |
| `TerminateRentalAgreement` | `300` |


### Event Management

Operations related to event management.


| Opcode | Gas Cost |
|---|---|
| `InitEvents` | `500` |
| `EmitEvent` | `40` |
| `GetEvent` | `80` |
| `ListEvents` | `100` |
| `CreateWallet` | `1000` |
| `ImportWallet` | `500` |
| `WalletBalance` | `40` |
| `WalletTransfer` | `210` |


### Employment Contracts

Operations related to employment contracts.


| Opcode | Gas Cost |
|---|---|
| `InitEmployment` | `1000` |
| `CreateJob` | `800` |
| `SignJob` | `300` |
| `RecordWork` | `100` |
| `PaySalary` | `800` |
| `GetJob` | `100` |


### Escrow Management

Operations related to escrow management.


| Opcode | Gas Cost |
|---|---|
| `EscrowCreate` | `0` |
| `EscrowDeposit` | `0` |
| `EscrowRelease` | `0` |
| `EscrowCancel` | `0` |
| `EscrowGet` | `0` |
| `EscrowList` | `0` |


### Faucet

Operations related to faucet.


| Opcode | Gas Cost |
|---|---|
| `NewFaucet` | `500` |
| `Faucet_Request` | `0` |
| `Faucet_Balance` | `0` |
| `Faucet_SetAmount` | `0` |
| `Faucet_SetCooldown` | `0` |


### Supply Chain

Operations related to supply chain.


| Opcode | Gas Cost |
|---|---|
| `GetItem` | `100` |
| `RegisterItem` | `1000` |
| `UpdateLocation` | `500` |
| `MarkStatus` | `500` |


### Healthcare Records

Operations related to healthcare records.


| Opcode | Gas Cost |
|---|---|
| `InitHealthcare` | `800` |
| `RegisterPatient` | `300` |
| `AddHealthRecord` | `400` |
| `GrantAccess` | `150` |
| `RevokeAccess` | `100` |
| `ListHealthRecords` | `200` |


### Warehouse Records

Operations related to warehouse records.


| Opcode | Gas Cost |
|---|---|
| `WarehouseNew` | `0` |
| `WarehouseAddItem` | `0` |
| `WarehouseRemoveItem` | `0` |
| `WarehouseMoveItem` | `0` |
| `WarehouseListItems` | `0` |
| `WarehouseGetItem` | `0` |


### Integration Node

Operations related to integration node.


| Opcode | Gas Cost |
|---|---|
| `IntRegisterAPI` | `1000` |
| `IntRemoveAPI` | `1000` |
| `IntListAPIs` | `200` |
| `IntConnectChain` | `1500` |
| `IntDisconnectChain` | `1500` |
| `IntListChains` | `200` |
| `IntRelayTx` | `1000` |


### Bank Institutional Node

Operations related to bank institutional node.


| Opcode | Gas Cost |
|---|---|
| `NewBankInstitutionalNode` | `1500` |
| `BankNode_Start` | `300` |
| `BankNode_Stop` | `200` |
| `BankNode_MonitorTx` | `250` |
| `BankNode_ComplianceReport` | `400` |
| `BankNode_ConnectFinNet` | `300` |
| `BankNode_UpdateRules` | `250` |
| `BankNode_SubmitTx` | `300` |
| `BankNode_RegisterInstitution` | `250` |
| `BankNode_RemoveInstitution` | `250` |
| `BankNode_ListInstitutions` | `200` |


### Custodial Node

Operations related to custodial node.


| Opcode | Gas Cost |
|---|---|
| `NewCustodialNode` | `2000` |
| `Custodial_Start` | `500` |
| `Custodial_Stop` | `500` |
| `Custodial_Register` | `300` |
| `Custodial_Deposit` | `500` |
| `Custodial_Withdraw` | `500` |
| `Custodial_Transfer` | `800` |
| `Custodial_Balance` | `200` |
| `Custodial_Audit` | `1000` |


### Forensic Node

Operations related to forensic node.


| Opcode | Gas Cost |
|---|---|
| `Forensic_Init` | `5000` |
| `Forensic_AnalyseTx` | `4000` |
| `Forensic_CheckCompliance` | `3000` |
| `Forensic_ThreatResponse` | `2000` |


### Indexing Node

Operations related to indexing node.


| Opcode | Gas Cost |
|---|---|
| `Indexing_Build` | `1000` |
| `Indexing_QueryTxHistory` | `200` |
| `Indexing_QueryState` | `200` |


### Immutability Enforcement

Operations related to immutability enforcement.


| Opcode | Gas Cost |
|---|---|
| `InitImmutability` | `800` |
| `VerifyChain` | `400` |
| `RestoreChain` | `600` |


### Watchtower Node

Operations related to watchtower node.


| Opcode | Gas Cost |
|---|---|
| `NewWatchtowerNode` | `5000` |
| `Watchtower_Start` | `1000` |
| `Watchtower_Stop` | `500` |
| `Watchtower_Log` | `200` |
| `Watchtower_Resolve` | `800` |


### Historical Node

Operations related to historical node.


| Opcode | Gas Cost |
|---|---|
| `NewHistoricalNode` | `5000` |
| `ArchiveBlock` | `500` |
| `BlockByHeight` | `400` |
| `RangeBlocks` | `800` |
| `SyncFromLedger` | `2000` |


### Optimization Node

Operations related to optimization node.


| Opcode | Gas Cost |
|---|---|
| `InitOptimization` | `1000` |
| `OptimizeTransactions` | `500` |
| `BalanceLoad` | `400` |


### Archival Witness Node

Operations related to archival witness node.


| Opcode | Gas Cost |
|---|---|
| `NewArchivalWitnessNode` | `1000` |
| `Witness_NotarizeTx` | `300` |
| `Witness_NotarizeBlock` | `500` |
| `Witness_GetTx` | `50` |
| `Witness_GetBlock` | `50` |


### Warfare / Military Nodes

Operations related to warfare / military nodes.


| Opcode | Gas Cost |
|---|---|
| `NewWarfareNode` | `2000` |
| `Warfare_SecureCommand` | `300` |
| `Warfare_TrackLogistics` | `300` |
| `Warfare_ShareTactical` | `200` |
