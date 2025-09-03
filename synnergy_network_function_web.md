# Synnergy Network Function Web

This document provides a high-level function web for the Synnergy network, outlining major modules and their core functions. Stage 8 expands cross-chain registries and bridge managers with CLI access and gas‑priced opcodes. Stage 9 adds governance primitives, custodial nodes and cross-consensus network management. Stage 14 introduces a unified `internal/nodes` package with light, watchtower and logistics implementations.

## Diagram

```mermaid
graph TD
    subgraph BiometricSecurity
        BA[NewBiometricsAuth] --> BE[Enroll]
        BE --> BV[Verify]
        BV --> BR[Remove]
        BSN[NewBiometricSecurityNode] --> BSA[Authenticate]
        BSA --> BSE[SecureExecute]
    end

    subgraph WarfareNode
        WN[NewWarfareNode] --> WC[SecureCommand]
        WC --> WL[TrackLogistics]
        WL --> WT[ShareTactical]
    end

    subgraph Wallet
        NW[NewWallet] --> WS[Sign]
        WS --> WV[VerifySignature]
    end

    subgraph Watchtower
        WTN[NewWatchtowerNode] --> WTS[Start]
        WTS --> WTMetrics[Metrics]
    end

    subgraph GeospatialNode
        GN[NewGeospatialNode] --> GR[Record]
        GR --> GH[History]
    end

    subgraph CrossChain
        CCM[NewCrossChainManager] --> CCR[NewCrossChainRegistry]
        CCM --> CCT[NewCrossChainTxManager]
        CCBTM[NewBridgeTransferManager] --> CCBTD[Deposit]
        CCBTM --> CCBTC[Claim]
        CCBTM --> CCBTGT[GetTransfer]
        CCBTM --> CCBTLT[ListTransfers]
    end

    subgraph PrivateTransactions
        PTM[NewPrivateTxManager] --> PTS[Send]
        PTM --> PTL[List]
    end

    subgraph CentralBank
        CBN[NewCentralBankingNode] --> CBM[Mint]
    end

    subgraph Charity
        CP[NewCharityPool] --> CR[Register]
        CP --> CV[Vote]
    end

    subgraph Synchronization
        SM[NewSyncManager] --> SOnce[Once]
        SM --> SStart[Start]
    end

    subgraph Staking
        SN[NewStakingNode] --> SS[Stake]
        SN --> SU[Unstake]
        SN --> SB[Balance]
        SN --> ST[TotalStaked]
    end

    subgraph ZeroTrust
        ZE[NewZeroTrustEngine] --> ZO[OpenChannel]
        ZO --> ZS[Send]
        ZS --> ZR[Receive]
    end

    subgraph Regulatory
        RN[NewRegulatoryNode] --> RA[ApproveTransaction]
        RN --> RF[FlagEntity]
        RN --> RL[Logs]
    end

    subgraph Compliance
        CS[NewComplianceService] --> CK[ValidateKYC]
        CS --> CF[RecordFraud]
        CMG[NewComplianceManager] --> CSu[Suspend]
    end

    subgraph ConnectionPool
        CPoo[NewConnectionPool] --> CA[Acquire]
        CPoo --> CR[Release]
    end

    subgraph VirtualMachine
        VM[NewLightVM] --> VMExec[Execute]
        SMgr[NewSandboxManager] --> SDel[DeleteSandbox]
        SMgr --> SPurge[PurgeInactive]
    end

    subgraph Consensus
        CH[NewConsensusHopper] --> CM[Mode]
        AM[NewAdaptiveManager] --> Adj[Adjust]
        DM[NewDifficultyManager] --> DS[AddSample]
        CSvc[NewConsensusService] --> CSStart[Start]
        VMg[NewValidatorManager] --> VAdd[Add]
        Cmg[NewContractManager] --> CTran[Transfer]
    end

    BiometricSecurity --> Consensus
    WarfareNode --> Consensus
    Wallet --> Consensus
    Watchtower --> Consensus
    GeospatialNode --> CrossChain
    CrossChain --> Consensus
    PrivateTransactions --> Consensus
    Staking --> Consensus
    ZeroTrust --> Compliance
    Regulatory --> CrossChain
    Compliance --> Regulatory
    Compliance --> Consensus
    ConnectionPool --> Consensus
    VirtualMachine --> Consensus
    CentralBank --> Consensus
    Charity --> Consensus
    Synchronization --> Consensus
```

## Key Relationships

- **BiometricSecurity** functions protect node operations and feed into the overall consensus processes.
- **WarfareNode**, **Watchtower**, and **GeospatialNode** modules provide specialized data and monitoring that flows into consensus and cross-chain operations.
- **CrossChain** functions manage bridging and transaction management across ledgers.
- **PrivateTransactions**, **Staking**, and **Regulatory** modules interact with consensus for secure and compliant network activity.
- **Wallet** functionality signs transactions that ultimately feed into consensus.

This visualization can be rendered using any Mermaid-compatible Markdown viewer.

