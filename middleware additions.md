# Middleware Additions

| Function | Location | Middleware Needed |
| --- | --- | --- |
| TestCrosschainagnosticprotocolsPlaceholder | cross_chain_agnostic_protocols_test.go:5 | None |
| TestCrosschaincontractsPlaceholder | cross_chain_contracts_test.go:5 | None |
| SNVMOpcodeByName | snvm._opcodes.go:1289 | None |
| TestStakingnodePlaceholder | staking_node_test.go:5 | None |
| NewMiningNode | mining_node.go:26 | None |
| ID | mining_node.go:31 | None |
| Start | mining_node.go:37 | None |
| Stop | mining_node.go:50 | None |
| IsRunning | mining_node.go:57 | None |
| HashRate | mining_node.go:64 | None |
| SubmitBlock | mining_node.go:72 | None |
| LastBlock | mining_node.go:79 | None |
| mineLoop | mining_node.go:87 | None |
| randomHash | mining_node.go:108 | None |
| MineBlock | mining_node.go:118 | None |
| Encrypt | private_transactions.go:13 | None |
| Decrypt | private_transactions.go:31 | None |
| NewPrivateTxManager | private_transactions.go:62 | None |
| Send | private_transactions.go:67 | None |
| List | private_transactions.go:74 | None |
| TestStakepenaltyPlaceholder | stake_penalty_test.go:5 | None |
| NewProtocolRegistry | cross_chain_agnostic_protocols.go:23 | None |
| RegisterProtocol | cross_chain_agnostic_protocols.go:28 | None |
| ListProtocols | cross_chain_agnostic_protocols.go:37 | None |
| GetProtocol | cross_chain_agnostic_protocols.go:48 | None |
| TestIDRegistry | idwallet_registration_test.go:5 | None |
| TestCrosschainconnectionPlaceholder | cross_chain_connection_test.go:5 | None |
| NewZeroTrustEngine | zero_trust_data_channels.go:31 | None |
| OpenChannel | zero_trust_data_channels.go:36 | None |
| Send | zero_trust_data_channels.go:51 | None |
| Messages | zero_trust_data_channels.go:70 | None |
| Receive | zero_trust_data_channels.go:83 | None |
| CloseChannel | zero_trust_data_channels.go:105 | None |
| NewCrossChainManager | cross_chain.go:26 | None |
| RegisterBridge | cross_chain.go:34 | None |
| ListBridges | cross_chain.go:52 | None |
| GetBridge | cross_chain.go:63 | None |
| AuthorizeRelayer | cross_chain.go:71 | None |
| RevokeRelayer | cross_chain.go:78 | None |
| IsRelayerAuthorized | cross_chain.go:85 | None |
| NewMobileMiningNode | mobile_mining_node.go:20 | None |
| UpdateBattery | mobile_mining_node.go:25 | None |
| Start | mobile_mining_node.go:35 | None |
| Battery | mobile_mining_node.go:47 | None |
| Threshold | mobile_mining_node.go:54 | None |
| SetThreshold | mobile_mining_node.go:61 | None |
| TestComplianceServiceKYCAndRisk | compliance_test.go:5 | None |
| TestComplianceServiceMonitorTransaction | compliance_test.go:37 | None |
| TestAIContractRegistry | ai_enhanced_contract_test.go:5 | None |
| TestBiometricSecurityNode | biometric_security_node_test.go:11 | None |
| NewStakePenaltyManager | stake_penalty.go:27 | None |
| AdjustStake | stake_penalty.go:37 | None |
| Penalize | stake_penalty.go:48 | None |
| Info | stake_penalty.go:58 | None |
| NewAnomalyDetector | anomaly_detection.go:20 | None |
| Update | anomaly_detection.go:28 | None |
| IsAnomalous | anomaly_detection.go:38 | None |
| TestDriftMonitorBaseline | ai_drift_monitor_test.go:5 | None |
| NewFailoverManager | high_availability.go:20 | None |
| RegisterBackup | high_availability.go:29 | None |
| Heartbeat | high_availability.go:36 | None |
| Active | high_availability.go:45 | None |
| TestMainPlaceholder | walletserver/main_test.go:5 | Rate limiting, validation, auth, integrity, user, crypto |
| NewSimpleVM | virtual_machine.go:43 | None |
| Start | virtual_machine.go:82 | None |
| Stop | virtual_machine.go:93 | None |
| Status | virtual_machine.go:104 | None |
| RegisterOpcode | virtual_machine.go:112 | None |
| ExecuteContext | virtual_machine.go:124 | None |
| Execute | virtual_machine.go:192 | None |
| newServer | walletserver/handlers.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| healthHandler | walletserver/handlers.go:15 | Rate limiting, validation, auth, integrity, user, crypto |
| newWalletHandler | walletserver/handlers.go:20 | Rate limiting, validation, auth, integrity, user, crypto |
| NewDriftMonitor | ai_drift_monitor.go:15 | None |
| UpdateBaseline | ai_drift_monitor.go:20 | None |
| HasDrift | ai_drift_monitor.go:27 | None |
| TestNewWalletHandler | walletserver/handlers_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestHealthHandler | walletserver/handlers_test.go:27 | Rate limiting, validation, auth, integrity, user, crypto |
| TestAnomalyDetectorBasic | anomaly_detection_test.go:5 | None |
| TestAnomalyDetectorDefaultThreshold | anomaly_detection_test.go:19 | None |
| main | walletserver/main.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSandboxManager | vm_sandbox_management_test.go:5 | None |
| NewIdentityService | identity_verification.go:30 | None |
| Register | identity_verification.go:38 | None |
| Verify | identity_verification.go:49 | None |
| Info | identity_verification.go:61 | None |
| Logs | identity_verification.go:69 | None |
| TestDatadistributionPlaceholder | data_distribution_test.go:5 | None |
| NewStakingNode | staking_node.go:15 | None |
| Stake | staking_node.go:20 | None |
| Unstake | staking_node.go:28 | None |
| Balance | staking_node.go:40 | None |
| TotalStaked | staking_node.go:47 | None |
| TestSplitAndReconstructHolographic | holographic_test.go:5 | None |
| TestMobileMiningNodeBatteryThreshold | mobile_mining_node_test.go:7 | None |
| TestMobileMiningNodeStartError | mobile_mining_node_test.go:42 | None |
| TestRegulatoryManager | regulatory_management_test.go:5 | None |
| TestMiningNodeLifecycle | mining_node_test.go:9 | None |
| TestMiningNodeMineBlock | mining_node_test.go:31 | None |
| TestMiningNodeSubmitBlock | mining_node_test.go:45 | None |
| NewAIService | ai.go:42 | None |
| PredictFraud | ai.go:51 | None |
| OptimiseBaseFee | ai.go:61 | None |
| ForecastVolume | ai.go:71 | None |
| PublishModel | ai.go:81 | None |
| FetchModel | ai.go:96 | None |
| ListModel | ai.go:104 | None |
| BuyModel | ai.go:109 | None |
| RentModel | ai.go:131 | None |
| ReleaseEscrow | ai.go:154 | None |
| TestVersionPlaceholder | pkg/version/version_test.go:5 | None |
| NewSandboxManager | vm_sandbox_management.go:29 | None |
| StartSandbox | vm_sandbox_management.go:34 | None |
| StopSandbox | vm_sandbox_management.go:54 | None |
| DeleteSandbox | vm_sandbox_management.go:66 | None |
| ResetSandbox | vm_sandbox_management.go:77 | None |
| SandboxStatus | vm_sandbox_management.go:89 | None |
| ListSandboxes | vm_sandbox_management.go:97 | None |
| TestFailoverManager | high_availability_test.go:8 | None |
| NewSecureStorage | ai_secure_storage.go:19 | None |
| Store | ai_secure_storage.go:24 | None |
| Retrieve | ai_secure_storage.go:48 | None |
| TestEnvironmentalMonitoringNodeTrigger | environmental_monitoring_node_test.go:5 | None |
| NewWarfareNode | warfare_node.go:21 | None |
| GetID | warfare_node.go:26 | None |
| SecureCommand | warfare_node.go:30 | None |
| TrackLogistics | warfare_node.go:38 | None |
| ShareTactical | warfare_node.go:52 | None |
| Logistics | warfare_node.go:57 | None |
| LogisticsByAsset | warfare_node.go:67 | None |
| TestCrossChainManager | cross_chain_stage18_test.go:5 | None |
| TestProtocolRegistry | cross_chain_stage18_test.go:31 | None |
| TestBridgeTransferManager | cross_chain_stage18_test.go:46 | None |
| TestConnectionManager | cross_chain_stage18_test.go:63 | None |
| TestXContractRegistry | cross_chain_stage18_test.go:77 | None |
| TestTransactionManager | cross_chain_stage18_test.go:93 | None |
| TestCrosschaintransactionsPlaceholder | cross_chain_transactions_test.go:5 | None |
| NewDataFeed | data_operations.go:19 | None |
| Update | data_operations.go:24 | None |
| Get | data_operations.go:32 | None |
| Delete | data_operations.go:40 | None |
| Keys | data_operations.go:50 | None |
| Snapshot | data_operations.go:61 | None |
| LastUpdated | data_operations.go:72 | None |
| NewTrainingManager | ai_training.go:28 | None |
| Start | ai_training.go:34 | None |
| Status | ai_training.go:53 | None |
| List | ai_training.go:61 | None |
| Cancel | ai_training.go:72 | None |
| Complete | ai_training.go:89 | None |
| TestModelMarketplaceAddAndGet | ai_model_management_test.go:9 | None |
| TestModelMarketplaceList | ai_model_management_test.go:32 | None |
| TestModelMarketplaceUpdate | ai_model_management_test.go:54 | None |
| TestModelMarketplaceRemove | ai_model_management_test.go:72 | None |
| TestModelMarketplaceConcurrency | ai_model_management_test.go:93 | None |
| TestContractOpcodesValues | contracts_opcodes_test.go:7 | None |
| TestContractOpcodesUnique | contracts_opcodes_test.go:41 | None |
| TestContractManager | contract_management_test.go:5 | None |
| TestMovingAverageModel | financial_prediction_test.go:8 | None |
| TestLinearRegressionModel | financial_prediction_test.go:23 | None |
| TestAR1Model | financial_prediction_test.go:41 | None |
| TestForecastSeriesDefault | financial_prediction_test.go:56 | None |
| NewWatchtowerNode | watchtower_node.go:27 | None |
| ID | watchtower_node.go:37 | None |
| Start | watchtower_node.go:40 | None |
| monitorLoop | watchtower_node.go:54 | None |
| Stop | watchtower_node.go:71 | None |
| ReportFork | watchtower_node.go:83 | None |
| Metrics | watchtower_node.go:90 | None |
| Firewall | watchtower_node.go:95 | None |
| TestSecureStorageStoreRetrieve | ai_secure_storage_test.go:9 | None |
| TestSecureStorageBadKey | ai_secure_storage_test.go:28 | None |
| NewComplianceService | compliance.go:52 | None |
| ValidateKYC | compliance.go:62 | None |
| EraseKYC | compliance.go:77 | None |
| RecordFraud | compliance.go:93 | None |
| RiskScore | compliance.go:111 | None |
| AuditTrail | compliance.go:119 | None |
| MonitorTransaction | compliance.go:131 | None |
| VerifyZKP | compliance.go:143 | None |
| appendAudit | compliance.go:148 | None |
| TestInferenceEngineRun | ai_inference_analysis_test.go:5 | None |
| resetLanguages | contract_language_compatibility_test.go:9 | None |
| TestLanguageRegistry | contract_language_compatibility_test.go:23 | None |
| NewComplianceManager | compliance_management.go:23 | None |
| Suspend | compliance_management.go:32 | None |
| Resume | compliance_management.go:47 | None |
| Whitelist | compliance_management.go:59 | None |
| Unwhitelist | compliance_management.go:74 | None |
| Status | compliance_management.go:85 | None |
| ReviewTransaction | compliance_management.go:94 | None |
| TestWarfarenodePlaceholder | warfare_node_test.go:5 | None |
| TestGeospatialNodeRecordAndHistory | geospatial_node_test.go:5 | None |
| IsZeroAddress | address_zero.go:9 | None |
| TestFirewall | firewall_test.go:5 | None |
| TestContractRegistry | contracts_test.go:5 | None |
| NewModelMarketplace | ai_model_management.go:27 | None |
| AddListing | ai_model_management.go:32 | None |
| Get | ai_model_management.go:49 | None |
| List | ai_model_management.go:57 | None |
| Update | ai_model_management.go:70 | None |
| Remove | ai_model_management.go:84 | None |
| NewContentNetworkNode | content_node.go:20 | None |
| Register | content_node.go:30 | None |
| Unregister | content_node.go:45 | None |
| Content | content_node.go:56 | None |
| List | content_node.go:64 | None |
| TestEnergyefficiencyPlaceholder | energy_efficiency_test.go:5 | None |
| TestSynthronCoinPlaceholder | docs/Whitepaper_detailed/Synthron Coin_test.go:5 | None |
| SynthronCoin | docs/Whitepaper_detailed/Synthron Coin.go:4 | None |
| TestGovernanceSecurity | Security assessments & Benchmark assessments/Security assessments/Governance security_test.go:5 | None |
| TestNetworkSecurity | Security assessments & Benchmark assessments/Security assessments/Network security_test.go:5 | None |
| TestAuthorityNodeSecurity | Security assessments & Benchmark assessments/Security assessments/Authority node security_test.go:5 | None |
| TestTransactionSecurity | Security assessments & Benchmark assessments/Security assessments/Transaction security_test.go:5 | None |
| TestVmSecurity | Security assessments & Benchmark assessments/Security assessments/Vm security_test.go:5 | None |
| TestLoanpoolSecurity | Security assessments & Benchmark assessments/Security assessments/Loanpool security_test.go:5 | None |
| TestOpcodeSecurity | Security assessments & Benchmark assessments/Security assessments/Opcode security_test.go:5 | None |
| TestAiSecurity | Security assessments & Benchmark assessments/Security assessments/Ai security_test.go:5 | None |
| TestSpeedSecurity | Security assessments & Benchmark assessments/Security assessments/Speed security_test.go:5 | None |
| validateSecurityAssessment | Security assessments & Benchmark assessments/Security assessments/helpers_test.go:5 | None |
| TestTokenStandardsSecurity | Security assessments & Benchmark assessments/Security assessments/Token standards security_test.go:5 | None |
| TestBlockSecurity | Security assessments & Benchmark assessments/Security assessments/Block security_test.go:5 | None |
| TestContractSecurity | Security assessments & Benchmark assessments/Security assessments/Contract security_test.go:5 | None |
| TestGasSecurity | Security assessments & Benchmark assessments/Security assessments/Gas security_test.go:5 | None |
| NewGeospatialNode | geospatial_node.go:23 | None |
| Record | geospatial_node.go:28 | None |
| History | geospatial_node.go:40 | None |
| TestConsensusHopperInitial | dynamic_consensus_hopping_test.go:10 | None |
| TestConsensusHopperEvaluate | dynamic_consensus_hopping_test.go:23 | None |
| TestConsensusHopperConcurrency | dynamic_consensus_hopping_test.go:58 | None |
| NewRegulatoryNode | regulatory_node.go:17 | None |
| ApproveTransaction | regulatory_node.go:26 | None |
| FlagEntity | regulatory_node.go:36 | None |
| Logs | regulatory_node.go:43 | None |
| NewFaucet | faucet.go:21 | None |
| Request | faucet.go:27 | None |
| Balance | faucet.go:42 | None |
| Configure | faucet.go:49 | None |
| SplitHolographic | holographic.go:10 | None |
| ReconstructHolographic | holographic.go:30 | None |
| TestSimpleVMLifecycle | virtual_machine_test.go:9 | None |
| TestSNVMOpcodeExecution | virtual_machine_test.go:28 | None |
| TestVMContextCancel | virtual_machine_test.go:48 | None |
| NewEnergyEfficiencyTracker | energy_efficiency.go:18 | None |
| Record | energy_efficiency.go:23 | None |
| Efficiency | energy_efficiency.go:34 | None |
| NetworkAverage | energy_efficiency.go:45 | None |
| Stats | energy_efficiency.go:62 | None |
| Reset | energy_efficiency.go:70 | None |
| TestContentNetworkNode | stage12_content_data_test.go:9 | None |
| TestContentNode | stage12_content_data_test.go:50 | None |
| TestContentNodeErrorsAndMeta | stage12_content_data_test.go:73 | None |
| TestDataDistribution | stage12_content_data_test.go:100 | None |
| TestDataFeed | stage12_content_data_test.go:125 | None |
| TestDataFeedLastUpdated | stage12_content_data_test.go:147 | None |
| TestDataResourceManager | stage12_content_data_test.go:165 | None |
| TestDataResourceManagerIsolation | stage12_content_data_test.go:187 | None |
| TestIndexingNode | stage12_content_data_test.go:202 | None |
| TestIndexingNodeEdgeCases | stage12_content_data_test.go:223 | None |
| TestBiometricsauthPlaceholder | biometrics_auth_test.go:5 | None |
| TestNewContentMetaFromData | content_types_test.go:9 | None |
| TestContentMetaValidateErrors | content_types_test.go:30 | None |
| main | Security assessments & Benchmark assessments/Security assessments/cmd/update_docs.go:22 | None |
| NewAIContractRegistry | ai_enhanced_contract.go:19 | None |
| DeployAIContract | ai_enhanced_contract.go:28 | None |
| InvokeAIContract | ai_enhanced_contract.go:41 | None |
| ModelHash | ai_enhanced_contract.go:49 | None |
| TestSystemHealthLogger | system_health_logging_test.go:5 | None |
| NewContentMeta | content_types.go:28 | None |
| NewContentMetaFromData | content_types.go:41 | None |
| Validate | content_types.go:55 | None |
| TestHealthEndpoint | cmd/api-gateway/main_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| newHandler | cmd/api-gateway/main.go:8 | Rate limiting, validation, auth, integrity, user, crypto |
| main | cmd/api-gateway/main.go:17 | Rate limiting, validation, auth, integrity, user, crypto |
| NewXContractRegistry | cross_chain_contracts.go:19 | None |
| RegisterMapping | cross_chain_contracts.go:24 | None |
| ListMappings | cross_chain_contracts.go:35 | None |
| GetMapping | cross_chain_contracts.go:46 | None |
| RemoveMapping | cross_chain_contracts.go:54 | None |
| NewContentNode | content_node_impl.go:24 | None |
| StoreContent | content_node_impl.go:36 | None |
| RetrieveContent | content_node_impl.go:64 | None |
| Meta | content_node_impl.go:93 | None |
| DeleteContent | content_node_impl.go:101 | None |
| TestGenerateTables | cmd/opcodegen/main_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| NewAccessController | access_control.go:12 | None |
| Grant | access_control.go:17 | None |
| Revoke | access_control.go:27 | None |
| HasRole | access_control.go:39 | None |
| List | access_control.go:51 | None |
| generateTables | cmd/opcodegen/main.go:16 | Rate limiting, validation, auth, integrity, user, crypto |
| main | cmd/opcodegen/main.go:70 | Rate limiting, validation, auth, integrity, user, crypto |
| TestAccessController | access_control_test.go:5 | None |
| TestRegulatoryNode | regulatory_node_test.go:5 | None |
| NewIndexingNode | indexing_node.go:13 | None |
| Index | indexing_node.go:18 | None |
| Query | indexing_node.go:26 | None |
| Remove | indexing_node.go:39 | None |
| Keys | indexing_node.go:46 | None |
| Count | indexing_node.go:57 | None |
| loadGasTable | gas_table.go:56 | None |
| LoadGasTable | gas_table.go:97 | None |
| GasCost | gas_table.go:106 | None |
| HasOpcode | gas_table.go:115 | None |
| RegisterGasCost | gas_table.go:123 | None |
| ResetGasTable | gas_table.go:134 | None |
| captureOutput | cmd/governance/main_test.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| TestMainRunsGovernmentCommand | cmd/governance/main_test.go:25 | Rate limiting, validation, auth, integrity, user, crypto |
| NewBiometricSecurityNode | biometric_security_node.go:17 | None |
| GetID | biometric_security_node.go:25 | None |
| Enroll | biometric_security_node.go:28 | None |
| Remove | biometric_security_node.go:33 | None |
| Authenticate | biometric_security_node.go:38 | None |
| SecureExecute | biometric_security_node.go:43 | None |
| main | cmd/governance/main.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| TestMainPlaceholder | cmd/watchtower/main_test.go:5 | Rate limiting, validation, auth, integrity, user, crypto |
| main | cmd/watchtower/main.go:17 | Rate limiting, validation, auth, integrity, user, crypto |
| captureOutput | cmd/firewall/main_test.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| TestMainRunsFirewallCommand | cmd/firewall/main_test.go:25 | Rate limiting, validation, auth, integrity, user, crypto |
| main | cmd/firewall/main.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| TestRunAddAndListPeers | cmd/p2p-node/main_test.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| runWithManager | cmd/p2p-node/main.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| run | cmd/p2p-node/main.go:59 | Rate limiting, validation, auth, integrity, user, crypto |
| main | cmd/p2p-node/main.go:63 | Rate limiting, validation, auth, integrity, user, crypto |
| TestMetricsEndpoint | cmd/monitoring/main_test.go:16 | Rate limiting, validation, auth, integrity, user, crypto |
| TestLedgerSecurity | Security assessments & Benchmark assessments/Security assessments/Ledger security_test.go:5 | None |
| newHandler | cmd/monitoring/main.go:15 | Rate limiting, validation, auth, integrity, user, crypto |
| main | cmd/monitoring/main.go:26 | Rate limiting, validation, auth, integrity, user, crypto |
| AnalyzeAssessment | Security assessments & Benchmark assessments/Security assessments/assessment.go:18 | None |
| UpdateTestingSection | Security assessments & Benchmark assessments/Security assessments/assessment.go:71 | None |
| TestStorageSecurity | Security assessments & Benchmark assessments/Security assessments/Storage security_test.go:5 | None |
| TestMainPlaceholder | cmd/synnergy/main_test.go:5 | Rate limiting, validation, auth, integrity, user, crypto |
| TestConsensusSecurity | Security assessments & Benchmark assessments/Security assessments/Consensus security_test.go:5 | None |
| main | cmd/synnergy/main.go:19 | Rate limiting, validation, auth, integrity, user, crypto |
| TestTreasurySecurity | Security assessments & Benchmark assessments/Security assessments/Treasury security_test.go:5 | None |
| TestNodeSecurity | Security assessments & Benchmark assessments/Security assessments/Node security_test.go:5 | None |
| TestDocgenCreatesGuide | cmd/docgen/main_test.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| TestCharitySecurity | Security assessments & Benchmark assessments/Security assessments/Charity security_test.go:5 | None |
| main | cmd/docgen/main.go:15 | Rate limiting, validation, auth, integrity, user, crypto |
| TestComplianceSecurity | Security assessments & Benchmark assessments/Security assessments/Compliance security_test.go:5 | None |
| captureOutput | cmd/secrets-manager/main_test.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| run | cmd/secrets-manager/main_test.go:23 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSetGet | cmd/secrets-manager/main_test.go:31 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSubBlocksSecurity | Security assessments & Benchmark assessments/Security assessments/Sub blocks security_test.go:5 | None |
| main | cmd/secrets-manager/main.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| TestGasTableIncludesNewOpcodes | gas_table_test.go:5 | None |
| TestEnergyefficientnodePlaceholder | energy_efficient_node_test.go:5 | None |
| NewRegulatoryManager | regulatory_management.go:23 | None |
| AddRegulation | regulatory_management.go:28 | None |
| RemoveRegulation | regulatory_management.go:39 | None |
| GetRegulation | regulatory_management.go:46 | None |
| ListRegulations | regulatory_management.go:54 | None |
| EvaluateTransaction | regulatory_management.go:65 | None |
| TestSNVMOpcodeLookup | snvm._opcodes_test.go:5 | None |
| TestWatchtowerNode | watchtower_node_test.go:10 | None |
| BenchmarkWalletBenchmarks | Security assessments & Benchmark assessments/Benchmarks/Wallet benchmarks_test.go:5 | None |
| BenchmarkHighAvailabilityBenchmarks | Security assessments & Benchmark assessments/Benchmarks/High availability benchmarks_test.go:5 | None |
| BenchmarkLedgerBenchmarks | Security assessments & Benchmark assessments/Benchmarks/Ledger benchmarks_test.go:5 | None |
| BenchmarkConsensusBenchmarks | Security assessments & Benchmark assessments/Benchmarks/Consensus_benchmarks_test.go:5 | None |
| BenchmarkStorageBenchmarks | Security assessments & Benchmark assessments/Benchmarks/Storage benchmarks_test.go:5 | None |
| BenchmarkAuthorityNodeBenchmarks | Security assessments & Benchmark assessments/Benchmarks/Authority node benchmarks_test.go:5 | None |
| performComputation | Security assessments & Benchmark assessments/Benchmarks/benchmark_util.go:17 | None |
| repairAdvice | Security assessments & Benchmark assessments/Benchmarks/benchmark_util.go:41 | None |
| parseBenchmarks | Security assessments & Benchmark assessments/Benchmarks/benchmark_util.go:64 | None |
| GenerateBenchmarkReport | Security assessments & Benchmark assessments/Benchmarks/benchmark_util.go:98 | None |
| BenchmarkAllTokenStandardBenchmarks | Security assessments & Benchmark assessments/Benchmarks/All Token standard Benchmarks_test.go:5 | None |
| BenchmarkGovernanceBenchmarks | Security assessments & Benchmark assessments/Benchmarks/Governance benchmarks_test.go:5 | None |
| BenchmarkCharityBenchmarks | Security assessments & Benchmark assessments/Benchmarks/Charity benchmarks_test.go:5 | None |
| BenchmarkNodeBenchmarks | Security assessments & Benchmark assessments/Benchmarks/Node benchmarks_test.go:5 | None |
| BenchmarkContractBenchmarks | Security assessments & Benchmark assessments/Benchmarks/Contract benchmarks_test.go:5 | None |
| NewLifePolicyRegistry | internal/tokens/syn2800.go:40 | None |
| IssuePolicy | internal/tokens/syn2800.go:46 | None |
| PayPremium | internal/tokens/syn2800.go:75 | None |
| FileClaim | internal/tokens/syn2800.go:91 | None |
| GetPolicy | internal/tokens/syn2800.go:110 | None |
| ListPolicies | internal/tokens/syn2800.go:123 | None |
| Deactivate | internal/tokens/syn2800.go:136 | None |
| main | Security assessments & Benchmark assessments/Benchmarks/cmd/runbench/main.go:10 | None |
| NewSYN223Token | internal/tokens/syn223_token.go:20 | None |
| AddToWhitelist | internal/tokens/syn223_token.go:34 | None |
| RemoveFromWhitelist | internal/tokens/syn223_token.go:41 | None |
| AddToBlacklist | internal/tokens/syn223_token.go:48 | None |
| RemoveFromBlacklist | internal/tokens/syn223_token.go:55 | None |
| Transfer | internal/tokens/syn223_token.go:62 | None |
| BalanceOf | internal/tokens/syn223_token.go:82 | None |
| TestSyn10Placeholder | internal/tokens/syn10_test.go:5 | None |
| BenchmarkValidationBenchmarks | Security assessments & Benchmark assessments/Benchmarks/Validation benchmarks_test.go:5 | None |
| TestSyn3600Placeholder | internal/tokens/syn3600_test.go:5 | None |
| BenchmarkSpeedBenchmarks | Security assessments & Benchmark assessments/Benchmarks/Speed Benchmarks_test.go:5 | None |
| NewCarbonRegistry | internal/tokens/syn200.go:39 | None |
| Register | internal/tokens/syn200.go:44 | None |
| Issue | internal/tokens/syn200.go:62 | None |
| Retire | internal/tokens/syn200.go:78 | None |
| AddVerification | internal/tokens/syn200.go:96 | None |
| Verifications | internal/tokens/syn200.go:109 | None |
| ProjectInfo | internal/tokens/syn200.go:120 | None |
| ListProjects | internal/tokens/syn200.go:137 | None |
| BenchmarkFullReportAndAssessment | Security assessments & Benchmark assessments/Benchmarks/Benchmarks_full_report_and_assessment_test.go:5 | None |
| NewSYN500Token | internal/tokens/syn500.go:20 | None |
| Mint | internal/tokens/syn500.go:25 | None |
| Redeem | internal/tokens/syn500.go:32 | None |
| TestSyn3500tokenPlaceholder | internal/tokens/syn3500_token_test.go:5 | None |
| BenchmarkSecurityBenchmarks | Security assessments & Benchmark assessments/Benchmarks/Security benchmarks_test.go:5 | None |
| NewSYN3200Token | internal/tokens/syn3200.go:13 | None |
| Convert | internal/tokens/syn3200.go:18 | None |
| SetRatio | internal/tokens/syn3200.go:25 | None |
| BenchmarkVmBenchmarks | Security assessments & Benchmark assessments/Benchmarks/VM benchmarks_test.go:5 | None |
| NewDebtRegistry | internal/tokens/syn845.go:39 | None |
| CreateToken | internal/tokens/syn845.go:44 | None |
| IssueDebt | internal/tokens/syn845.go:55 | None |
| RecordPayment | internal/tokens/syn845.go:70 | None |
| GetDebt | internal/tokens/syn845.go:86 | None |
| NewSYN3500Token | internal/tokens/syn3500_token.go:19 | None |
| SetRate | internal/tokens/syn3500_token.go:30 | None |
| Info | internal/tokens/syn3500_token.go:37 | None |
| Mint | internal/tokens/syn3500_token.go:44 | None |
| Redeem | internal/tokens/syn3500_token.go:51 | None |
| BalanceOf | internal/tokens/syn3500_token.go:63 | None |
| NewSYN70Token | internal/tokens/syn70.go:26 | None |
| RegisterAsset | internal/tokens/syn70.go:34 | None |
| TransferAsset | internal/tokens/syn70.go:45 | None |
| SetAttribute | internal/tokens/syn70.go:64 | None |
| AddAchievement | internal/tokens/syn70.go:76 | None |
| AssetInfo | internal/tokens/syn70.go:88 | None |
| ListAssets | internal/tokens/syn70.go:105 | None |
| NewSyn2500Member | internal/tokens/syn2500_token.go:18 | None |
| UpdateVotingPower | internal/tokens/syn2500_token.go:33 | None |
| NewSyn2500Registry | internal/tokens/syn2500_token.go:44 | None |
| AddMember | internal/tokens/syn2500_token.go:49 | None |
| GetMember | internal/tokens/syn2500_token.go:56 | None |
| RemoveMember | internal/tokens/syn2500_token.go:64 | None |
| ListMembers | internal/tokens/syn2500_token.go:71 | None |
| NewSYN3800Token | internal/tokens/syn3800.go:16 | None |
| Mint | internal/tokens/syn3800.go:21 | None |
| Burn | internal/tokens/syn3800.go:32 | None |
| Supply | internal/tokens/syn3800.go:43 | None |
| TestSyn3400Placeholder | internal/tokens/syn3400_test.go:5 | None |
| TestSYN1000Index | internal/tokens/syn1000_index_test.go:9 | None |
| TestSyn3900Placeholder | internal/tokens/syn3900_test.go:5 | None |
| TestBaseTokenMintTransferBurn | internal/tokens/base_test.go:12 | None |
| TestBaseTokenAllowance | internal/tokens/base_test.go:31 | None |
| TestSYN10Info | internal/tokens/base_test.go:50 | None |
| TestSYN1000ReserveValue | internal/tokens/base_test.go:61 | None |
| TestSYN12Metadata | internal/tokens/base_test.go:76 | None |
| TestSYN20PauseFreeze | internal/tokens/base_test.go:84 | None |
| TestSYN70AssetLifecycle | internal/tokens/base_test.go:104 | None |
| TestRegistryInfoList | internal/tokens/base_test.go:130 | None |
| TestSYN1100Access | internal/tokens/base_test.go:143 | None |
| TestBaseTokenTransferErrors | internal/tokens/base_test.go:163 | None |
| TestBaseTokenConcurrent | internal/tokens/base_test.go:173 | None |
| TestTemplatesExist | smart-contracts/templates_test.go:9 | None |
| TestSyn20Placeholder | internal/tokens/syn20_test.go:5 | None |
| TestSyn4200tokenPlaceholder | internal/tokens/syn4200_token_test.go:5 | None |
| TestSyn845Placeholder | internal/tokens/syn845_test.go:5 | None |
| TestSyn500Placeholder | internal/tokens/syn500_test.go:5 | None |
| NewForexRegistry | internal/tokens/syn3400.go:27 | None |
| Register | internal/tokens/syn3400.go:32 | None |
| UpdateRate | internal/tokens/syn3400.go:43 | None |
| Get | internal/tokens/syn3400.go:56 | None |
| List | internal/tokens/syn3400.go:68 | None |
| NewSYN1000Token | internal/tokens/syn1000.go:25 | None |
| AddReserve | internal/tokens/syn1000.go:34 | None |
| SetReservePrice | internal/tokens/syn1000.go:46 | None |
| TotalReserveValue | internal/tokens/syn1000.go:59 | None |
| TestSYN2700Distribute | internal/tokens/token_extensions_test.go:8 | None |
| TestSYN3200Convert | internal/tokens/token_extensions_test.go:18 | None |
| TestSYN3600Weight | internal/tokens/token_extensions_test.go:25 | None |
| TestSYN3800Cap | internal/tokens/token_extensions_test.go:33 | None |
| TestSYN3900Vesting | internal/tokens/token_extensions_test.go:43 | None |
| TestSYN500Points | internal/tokens/token_extensions_test.go:52 | None |
| TestSYN5000Transfer | internal/tokens/token_extensions_test.go:61 | None |
| TestSyn300tokenPlaceholder | internal/tokens/syn300_token_test.go:5 | None |
| NewSYN4200Token | internal/tokens/syn4200_token.go:21 | None |
| Donate | internal/tokens/syn4200_token.go:26 | None |
| CampaignProgress | internal/tokens/syn4200_token.go:43 | None |
| Campaign | internal/tokens/syn4200_token.go:54 | None |
| TestSyn2500tokenPlaceholder | internal/tokens/syn2500_token_test.go:5 | None |
| TestSyn3700tokenPlaceholder | internal/tokens/syn3700_token_test.go:5 | None |
| TestSyn2369Placeholder | internal/tokens/syn2369_test.go:5 | None |
| NewInsuranceRegistry | internal/tokens/syn2900.go:42 | None |
| IssuePolicy | internal/tokens/syn2900.go:48 | None |
| FileClaim | internal/tokens/syn2900.go:79 | None |
| GetPolicy | internal/tokens/syn2900.go:98 | None |
| ListPolicies | internal/tokens/syn2900.go:111 | None |
| Deactivate | internal/tokens/syn2900.go:124 | None |
| TestSYN223TokenTransfer | internal/tokens/dao_tokens_test.go:8 | None |
| TestSYN300TokenGovernanceFlow | internal/tokens/dao_tokens_test.go:22 | None |
| TestSyn2500Registry | internal/tokens/dao_tokens_test.go:38 | None |
| TestSYN3500TokenMintRedeem | internal/tokens/dao_tokens_test.go:56 | None |
| TestSYN3700TokenValue | internal/tokens/dao_tokens_test.go:67 | None |
| TestSYN4200TokenDonations | internal/tokens/dao_tokens_test.go:77 | None |
| TestLegalTokenWorkflow | internal/tokens/dao_tokens_test.go:89 | None |
| NewLegalToken | internal/tokens/syn4700.go:44 | None |
| Sign | internal/tokens/syn4700.go:63 | None |
| RevokeSignature | internal/tokens/syn4700.go:74 | None |
| UpdateStatus | internal/tokens/syn4700.go:81 | None |
| Dispute | internal/tokens/syn4700.go:88 | None |
| partyExists | internal/tokens/syn4700.go:95 | None |
| NewLegalTokenRegistry | internal/tokens/syn4700.go:111 | None |
| Add | internal/tokens/syn4700.go:116 | None |
| Get | internal/tokens/syn4700.go:123 | None |
| Remove | internal/tokens/syn4700.go:131 | None |
| List | internal/tokens/syn4700.go:138 | None |
| NewSYN2700Token | internal/tokens/syn2700.go:14 | None |
| AddHolder | internal/tokens/syn2700.go:20 | None |
| Distribute | internal/tokens/syn2700.go:30 | None |
| TestSyn2700Placeholder | internal/tokens/syn2700_test.go:5 | None |
| NewSYN3900Token | internal/tokens/syn3900.go:21 | None |
| Grant | internal/tokens/syn3900.go:26 | None |
| Release | internal/tokens/syn3900.go:34 | None |
| TestSyn70Placeholder | internal/tokens/syn70_test.go:5 | None |
| NewSYN300Token | internal/tokens/syn300_token.go:30 | None |
| Delegate | internal/tokens/syn300_token.go:44 | None |
| RevokeDelegation | internal/tokens/syn300_token.go:55 | None |
| VotingPower | internal/tokens/syn300_token.go:62 | None |
| votingPowerLocked | internal/tokens/syn300_token.go:68 | None |
| CreateProposal | internal/tokens/syn300_token.go:79 | None |
| Vote | internal/tokens/syn300_token.go:96 | None |
| Execute | internal/tokens/syn300_token.go:117 | None |
| ProposalStatus | internal/tokens/syn300_token.go:139 | None |
| ListProposals | internal/tokens/syn300_token.go:159 | None |
| NewBaseToken | internal/tokens/base.go:50 | None |
| ID | internal/tokens/base.go:62 | None |
| Name | internal/tokens/base.go:69 | None |
| Symbol | internal/tokens/base.go:76 | None |
| Decimals | internal/tokens/base.go:83 | None |
| TotalSupply | internal/tokens/base.go:90 | None |
| BalanceOf | internal/tokens/base.go:97 | None |
| Transfer | internal/tokens/base.go:104 | None |
| TransferFrom | internal/tokens/base.go:116 | None |
| Mint | internal/tokens/base.go:132 | None |
| Burn | internal/tokens/base.go:141 | None |
| Approve | internal/tokens/base.go:153 | None |
| Allowance | internal/tokens/base.go:164 | None |
| TestSyn2900Placeholder | internal/tokens/syn2900_test.go:5 | None |
| TestSyn12Placeholder | internal/tokens/syn12_test.go:5 | None |
| NewSYN3700Token | internal/tokens/syn3700_token.go:23 | None |
| AddComponent | internal/tokens/syn3700_token.go:28 | None |
| RemoveComponent | internal/tokens/syn3700_token.go:35 | None |
| ListComponents | internal/tokens/syn3700_token.go:48 | None |
| Value | internal/tokens/syn3700_token.go:57 | None |
| TestSyn2600Placeholder | internal/tokens/syn2600_test.go:5 | None |
| TestSYN200CarbonRegistry | internal/tokens/advanced_tokens_test.go:8 | None |
| TestSYN2600InvestorRegistry | internal/tokens/advanced_tokens_test.go:40 | None |
| TestSYN2800LifePolicyRegistry | internal/tokens/advanced_tokens_test.go:66 | None |
| TestSYN2900InsuranceRegistry | internal/tokens/advanced_tokens_test.go:93 | None |
| TestSYN3400ForexRegistry | internal/tokens/advanced_tokens_test.go:117 | None |
| TestSYN845DebtRegistry | internal/tokens/advanced_tokens_test.go:133 | None |
| TestSYN2369ItemRegistry | internal/tokens/advanced_tokens_test.go:149 | None |
| TestSYN223Token | internal/tokens/advanced_tokens_test.go:168 | None |
| TestSyn2500RegistryOperations | internal/tokens/advanced_tokens_test.go:191 | None |
| TestSYN300Token | internal/tokens/advanced_tokens_test.go:214 | None |
| TestSYN3500Token | internal/tokens/advanced_tokens_test.go:248 | None |
| TestSYN3700Token | internal/tokens/advanced_tokens_test.go:267 | None |
| TestSYN4200Token | internal/tokens/advanced_tokens_test.go:287 | None |
| TestSYN4700LegalToken | internal/tokens/advanced_tokens_test.go:301 | None |
| TestSYN70Token | internal/tokens/advanced_tokens_test.go:332 | None |
| TestSyn2800Placeholder | internal/tokens/syn2800_test.go:5 | None |
| NewSYN3600Token | internal/tokens/syn3600.go:12 | None |
| SetWeight | internal/tokens/syn3600.go:17 | None |
| Weight | internal/tokens/syn3600.go:24 | None |
| TestSyn3200Placeholder | internal/tokens/syn3200_test.go:5 | None |
| TestSyn223tokenPlaceholder | internal/tokens/syn223_token_test.go:5 | None |
| TestSYN1000TokenReserveValue | internal/tokens/syn1000_test.go:11 | None |
| TestSYN1000TokenConcurrent | internal/tokens/syn1000_test.go:26 | None |
| TestSYN20ConcurrentPause | internal/tokens/standard_tokens_concurrency_test.go:11 | None |
| TestSYN70ConcurrentRegister | internal/tokens/standard_tokens_concurrency_test.go:30 | None |
| TestSyn200Placeholder | internal/tokens/syn200_test.go:5 | None |
| NewSYN1100Token | internal/tokens/syn1100.go:18 | None |
| AddRecord | internal/tokens/syn1100.go:23 | None |
| GrantAccess | internal/tokens/syn1100.go:32 | None |
| RevokeAccess | internal/tokens/syn1100.go:42 | None |
| GetRecord | internal/tokens/syn1100.go:52 | None |
| NewInvestorRegistry | internal/tokens/syn2600.go:36 | None |
| Issue | internal/tokens/syn2600.go:44 | None |
| Transfer | internal/tokens/syn2600.go:72 | None |
| RecordReturn | internal/tokens/syn2600.go:88 | None |
| Deactivate | internal/tokens/syn2600.go:104 | None |
| Get | internal/tokens/syn2600.go:116 | None |
| List | internal/tokens/syn2600.go:129 | None |
| TestSyn5000Placeholder | internal/tokens/syn5000_test.go:5 | None |
| NewSYN5000Token | internal/tokens/syn5000.go:16 | None |
| Mint | internal/tokens/syn5000.go:21 | None |
| Transfer | internal/tokens/syn5000.go:31 | None |
| Balance | internal/tokens/syn5000.go:47 | None |
| NewSYN12Token | internal/tokens/syn12.go:22 | None |
| NewSYN1000Index | internal/tokens/syn1000_index.go:18 | None |
| Create | internal/tokens/syn1000_index.go:23 | None |
| Token | internal/tokens/syn1000_index.go:33 | None |
| AddReserve | internal/tokens/syn1000_index.go:44 | None |
| SetReservePrice | internal/tokens/syn1000_index.go:54 | None |
| TotalValue | internal/tokens/syn1000_index.go:64 | None |
| TestContentNode_StoreRetrieve | content_node_impl_test.go:5 | None |
| TestContentNode_KeyLength | content_node_impl_test.go:25 | None |
| TestSyn4700Placeholder | internal/tokens/syn4700_test.go:5 | None |
| TestIndexingNodeBasic | indexing_node_test.go:9 | None |
| TestIndexingNodeQueryReturnsCopy | indexing_node_test.go:46 | None |
| TestIndexingNodeConcurrentAccess | indexing_node_test.go:62 | None |
| NewSYN20Token | internal/tokens/syn20.go:17 | None |
| Pause | internal/tokens/syn20.go:25 | None |
| Unpause | internal/tokens/syn20.go:32 | None |
| Freeze | internal/tokens/syn20.go:39 | None |
| Unfreeze | internal/tokens/syn20.go:46 | None |
| Transfer | internal/tokens/syn20.go:53 | None |
| Mint | internal/tokens/syn20.go:69 | None |
| Burn | internal/tokens/syn20.go:81 | None |
| TestComplianceManager | compliance_management_test.go:5 | None |
| NewItemRegistry | internal/tokens/syn2369.go:28 | None |
| CreateItem | internal/tokens/syn2369.go:33 | None |
| TransferItem | internal/tokens/syn2369.go:47 | None |
| UpdateAttributes | internal/tokens/syn2369.go:59 | None |
| GetItem | internal/tokens/syn2369.go:73 | None |
| ListItems | internal/tokens/syn2369.go:89 | None |
| TestContentNetworkNode_RegisterUnregister | content_node_test.go:5 | None |
| TestContentNetworkNode_ContentAndList | content_node_test.go:22 | None |
| NewRegistry | internal/tokens/index.go:14 | None |
| NextID | internal/tokens/index.go:19 | None |
| Register | internal/tokens/index.go:27 | None |
| Get | internal/tokens/index.go:34 | None |
| GetBySymbol | internal/tokens/index.go:42 | None |
| Info | internal/tokens/index.go:63 | None |
| InfoBySymbol | internal/tokens/index.go:82 | None |
| List | internal/tokens/index.go:100 | None |
| NewSYN10Token | internal/tokens/syn10.go:16 | None |
| SetExchangeRate | internal/tokens/syn10.go:25 | None |
| Info | internal/tokens/syn10.go:41 | None |
| TestIndexPlaceholder | internal/tokens/index_test.go:5 | None |
| TestSyn3800Placeholder | internal/tokens/syn3800_test.go:5 | None |
| BenchmarkBaseTokenTransfer | internal/tokens/base_benchmark_test.go:5 | None |
| BenchmarkRegistryInfo | internal/tokens/base_benchmark_test.go:14 | None |
| TestSyn1100Placeholder | internal/tokens/syn1100_test.go:5 | None |
| TestAuditLog | internal/governance/audit_log_test.go:5 | None |
| TestReplayProtectorSeen | internal/governance/replay_protection_test.go:5 | None |
| NewAuditLog | internal/governance/audit_log.go:12 | None |
| Append | internal/governance/audit_log.go:15 | None |
| Entries | internal/governance/audit_log.go:22 | None |
| NewReplayProtector | internal/governance/replay_protection.go:12 | None |
| Seen | internal/governance/replay_protection.go:17 | None |
| NewKeyRotator | internal/p2p/key_rotation.go:11 | None |
| Rotate | internal/p2p/key_rotation.go:16 | None |
| TestPFSChannel | internal/p2p/pfs_test.go:5 | None |
| generateCert | internal/p2p/tls_transport_test.go:19 | None |
| TestTLSTransportRoundTrip | internal/p2p/tls_transport_test.go:55 | None |
| TestKeyRotatorRotate | internal/p2p/key_rotation_test.go:8 | None |
| NewNoiseTransport | internal/p2p/noise_transport.go:27 | None |
| Dial | internal/p2p/noise_transport.go:37 | None |
| Listen | internal/p2p/noise_transport.go:51 | None |
| Accept | internal/p2p/noise_transport.go:67 | None |
| Write | internal/p2p/noise_transport.go:83 | None |
| Read | internal/p2p/noise_transport.go:99 | None |
| handshake | internal/p2p/noise_transport.go:118 | None |
| TestManager | internal/p2p/peer_test.go:5 | None |
| TestDiscoveryService | internal/p2p/peer_test.go:21 | None |
| TestNoiseHandshake | internal/p2p/noise_transport_test.go:9 | None |
| NewManager | internal/p2p/peer.go:19 | None |
| AddPeer | internal/p2p/peer.go:24 | None |
| RemovePeer | internal/p2p/peer.go:31 | None |
| GetPeer | internal/p2p/peer.go:38 | None |
| ListPeers | internal/p2p/peer.go:46 | None |
| NewDiscoveryService | internal/p2p/discovery.go:10 | None |
| DiscoverPeers | internal/p2p/discovery.go:15 | None |
| TestDiscoveryPlaceholder | internal/p2p/discovery_test.go:5 | None |
| NewPFSChannel | internal/p2p/pfs.go:9 | None |
| Encrypt | internal/p2p/pfs.go:12 | None |
| Decrypt | internal/p2p/pfs.go:20 | None |
| NewTLSTransport | internal/p2p/tls_transport.go:17 | None |
| Dial | internal/p2p/tls_transport.go:33 | None |
| Listen | internal/p2p/tls_transport.go:42 | None |
| BenchmarkComplianceBenchmarks | Security assessments & Benchmark assessments/Benchmarks/Compliance benchmarks_test.go:5 | None |
| TestDefaultConfigPath | internal/config/default_default_test.go:7 | None |
| TestDefaultprodPlaceholder | internal/config/default_prod_test.go:5 | None |
| Load | internal/config/config.go:39 | None |
| TestDefaultdevPlaceholder | internal/config/default_dev_test.go:5 | None |
| writeTempFile | internal/config/config_test.go:10 | None |
| TestLoadYAML | internal/config/config_test.go:19 | None |
| TestEnvOverrideAndJSON | internal/config/config_test.go:44 | None |
| TestValidation | internal/config/config_test.go:63 | None |
| Tracer | internal/telemetry/telemetry.go:9 | None |
| TestTracerNotNil | internal/telemetry/telemetry_test.go:5 | None |
| BenchmarkAiBenchmarks | Security assessments & Benchmark assessments/Benchmarks/ai benchmarks_test.go:5 | None |
| Error | internal/errors/errors.go:25 | None |
| Unwrap | internal/errors/errors.go:33 | None |
| New | internal/errors/errors.go:36 | None |
| Wrap | internal/errors/errors.go:39 | None |
| IsCode | internal/errors/errors.go:44 | None |
| TestErrorWrapping | internal/errors/errors_test.go:8 | None |
| BenchmarkCoinBenchmarks | Security assessments & Benchmark assessments/Benchmarks/Coin benchmarks_test.go:5 | None |
| TestParseBenchmarks | Security assessments & Benchmark assessments/Benchmarks/benchmark_util_test.go:8 | None |
| TestRepairAdvice | Security assessments & Benchmark assessments/Benchmarks/benchmark_util_test.go:30 | None |
| BenchmarkLoanpoolBenchmarks | Security assessments & Benchmark assessments/Benchmarks/Loanpool benchmarks_test.go:5 | None |
| BenchmarkOpcodeBenchmarks | Security assessments & Benchmark assessments/Benchmarks/Opcode benchmarks_test.go:5 | None |
| BenchmarkNetworkBenchmarks | Security assessments & Benchmark assessments/Benchmarks/Network benchmarks_test.go:5 | None |
| NewTracer | internal/monitoring/tracing.go:7 | None |
| StartSpan | internal/monitoring/tracing.go:10 | None |
| BenchmarkTransactionsBenchmarks | Security assessments & Benchmark assessments/Benchmarks/transactions benchmarks_test.go:5 | None |
| TestTracerStartSpan | internal/monitoring/tracing_test.go:5 | None |
| TestIdentityService | identity_verification_test.go:5 | None |
| TestMetrics | internal/monitoring/metrics_test.go:5 | None |
| NewContractRegistry | contracts.go:39 | None |
| CompileWASM | contracts.go:49 | None |
| Deploy | contracts.go:60 | None |
| Invoke | contracts.go:84 | None |
| List | contracts.go:101 | None |
| Get | contracts.go:112 | None |
| NewMetrics | internal/monitoring/metrics.go:12 | None |
| Inc | internal/monitoring/metrics.go:17 | None |
| Get | internal/monitoring/metrics.go:24 | None |
| NewInferenceEngine | ai_inference_analysis.go:22 | None |
| LoadModel | ai_inference_analysis.go:27 | None |
| Run | ai_inference_analysis.go:34 | None |
| Analyse | ai_inference_analysis.go:46 | None |
| TestAlerterAlert | internal/monitoring/alerting_test.go:5 | None |
| Forecast | financial_prediction.go:18 | None |
| Forecast | financial_prediction.go:44 | None |
| Forecast | financial_prediction.go:76 | None |
| ForecastSeries | financial_prediction.go:103 | None |
| NewAlerter | internal/monitoring/alerting.go:9 | None |
| Subscribe | internal/monitoring/alerting.go:12 | None |
| Alert | internal/monitoring/alerting.go:17 | None |
| TestPolicyEnforcerAuthorize | internal/auth/rbac_test.go:10 | None |
| NewRBAC | internal/auth/rbac.go:25 | None |
| AddRole | internal/auth/rbac.go:33 | None |
| AddPermissionToRole | internal/auth/rbac.go:43 | None |
| AssignRole | internal/auth/rbac.go:55 | None |
| hasPermission | internal/auth/rbac.go:69 | None |
| NewPolicyEnforcer | internal/auth/rbac.go:95 | None |
| Authorize | internal/auth/rbac.go:100 | None |
| TestAuditPlaceholder | internal/auth/audit_test.go:5 | None |
| NewStdAuditLogger | internal/auth/audit.go:21 | None |
| Log | internal/auth/audit.go:26 | None |
| TestExperimentalnodePlaceholder | internal/nodes/experimental_node_test.go:5 | None |
| TestGeospatialNodeRecordHistory | internal/nodes/geospatial_test.go:8 | None |
| TestHolographicNodeStoreRetrieve | internal/nodes/holographic_node_test.go:11 | None |
| TestHolographicNodePeers | internal/nodes/holographic_node_test.go:21 | None |
| TestHolographicNodeIDAndLifecycle | internal/nodes/holographic_node_test.go:31 | None |
| TestHolographicNodeRetrieveMissing | internal/nodes/holographic_node_test.go:44 | None |
| TestHolographicNodeConcurrentStoreRetrieve | internal/nodes/holographic_node_test.go:51 | None |
| NewForensicNode | internal/nodes/forensic_node.go:52 | None |
| RecordTransaction | internal/nodes/forensic_node.go:57 | None |
| RecordNetworkTrace | internal/nodes/forensic_node.go:65 | None |
| Transactions | internal/nodes/forensic_node.go:73 | None |
| NetworkTraces | internal/nodes/forensic_node.go:82 | None |
| NewGeospatialNode | internal/nodes/geospatial.go:36 | None |
| Record | internal/nodes/geospatial.go:41 | None |
| History | internal/nodes/geospatial.go:50 | None |
| NewHistoricalNode | internal/nodes/historical_node.go:40 | None |
| ArchiveBlock | internal/nodes/historical_node.go:45 | None |
| GetBlockByHeight | internal/nodes/historical_node.go:54 | None |
| GetBlockByHash | internal/nodes/historical_node.go:62 | None |
| TotalBlocks | internal/nodes/historical_node.go:70 | None |
| NewSimpleWarfareNode | internal/nodes/military_nodes/index.go:52 | None |
| GetID | internal/nodes/military_nodes/index.go:57 | None |
| SecureCommand | internal/nodes/military_nodes/index.go:61 | None |
| TrackLogistics | internal/nodes/military_nodes/index.go:67 | None |
| ShareTactical | internal/nodes/military_nodes/index.go:79 | None |
| Logistics | internal/nodes/military_nodes/index.go:82 | None |
| TestSimpleWarfareNodeLogistics | internal/nodes/military_nodes/index_test.go:5 | None |
| TestTypesPlaceholder | internal/nodes/types_test.go:5 | None |
| TestElectedAuthorityNode | internal/nodes/elected_authority_node_test.go:8 | None |
| NewBasicWatchtower | internal/nodes/watchtower/index.go:59 | None |
| ID | internal/nodes/watchtower/index.go:64 | None |
| Start | internal/nodes/watchtower/index.go:67 | None |
| Stop | internal/nodes/watchtower/index.go:79 | None |
| ReportFork | internal/nodes/watchtower/index.go:90 | None |
| Metrics | internal/nodes/watchtower/index.go:100 | None |
| TestProtocolRegistry | core/cross_chain_agnostic_protocols_test.go:5 | None |
| TestBasicWatchtowerLifecycle | internal/nodes/watchtower/index_test.go:8 | None |
| NewConnectionManager | cross_chain_connection.go:27 | None |
| OpenConnection | cross_chain_connection.go:32 | None |
| CloseConnection | cross_chain_connection.go:46 | None |
| GetConnection | cross_chain_connection.go:62 | None |
| ListConnections | cross_chain_connection.go:70 | None |
| TestDefaultGasTable | core/gas_test.go:5 | None |
| TestFaucetPlaceholder | faucet_test.go:5 | None |
| TestQuadraticWeight | core/dao_quadratic_voting_test.go:5 | None |
| TestCastQuadraticVoteZeroTokens | core/dao_quadratic_voting_test.go:11 | None |
| TestCastQuadraticVoteRequiresMembership | core/dao_quadratic_voting_test.go:28 | None |
| TestCastQuadraticVoteSuccess | core/dao_quadratic_voting_test.go:45 | None |
| TestCrossChainRegistry | core/cross_chain_contracts_test.go:5 | None |
| TestDataresourcemanagementPlaceholder | data_resource_management_test.go:5 | None |
| NewIDRegistry | idwallet_registration.go:15 | None |
| Register | idwallet_registration.go:20 | None |
| Info | idwallet_registration.go:31 | None |
| IsRegistered | idwallet_registration.go:39 | None |
| NewFirewall | firewall.go:16 | None |
| BlockAddress | firewall.go:25 | None |
| UnblockAddress | firewall.go:32 | None |
| IsAddressBlocked | firewall.go:39 | None |
| BlockToken | firewall.go:47 | None |
| UnblockToken | firewall.go:54 | None |
| IsTokenBlocked | firewall.go:61 | None |
| BlockIP | firewall.go:69 | None |
| UnblockIP | firewall.go:76 | None |
| IsIPBlocked | firewall.go:83 | None |
| Rules | firewall.go:91 | None |
| NewSynnergyConsensus | core/consensus.go:40 | None |
| Threshold | core/consensus.go:57 | None |
| AdjustWeights | core/consensus.go:65 | None |
| Tload | core/consensus.go:92 | None |
| Tsecurity | core/consensus.go:100 | None |
| Tstake | core/consensus.go:108 | None |
| TransitionThreshold | core/consensus.go:117 | None |
| DifficultyAdjust | core/consensus.go:123 | None |
| SetAvailability | core/consensus.go:131 | None |
| SetPoWRewards | core/consensus.go:139 | None |
| SelectValidator | core/consensus.go:146 | None |
| ValidateSubBlock | core/consensus.go:183 | None |
| MineBlock | core/consensus.go:193 | None |
| clamp | core/consensus.go:209 | None |
| ValidateBlock | core/consensus.go:222 | None |
| SupportedContractLanguages | contract_language_compatibility.go:25 | None |
| AddSupportedLanguage | contract_language_compatibility.go:37 | None |
| RemoveSupportedLanguage | contract_language_compatibility.go:49 | None |
| IsLanguageSupported | contract_language_compatibility.go:62 | None |
| TestAdaptiveManagerWindowAndReset | core/consensus_adaptive_management_test.go:8 | None |
| TestGeospatialPlaceholder | node_ext/geospatial_test.go:5 | None |
| ID | internal/nodes/consensus_specific_test.go:6 | None |
| Start | internal/nodes/consensus_specific_test.go:7 | None |
| Stop | internal/nodes/consensus_specific_test.go:8 | None |
| IsRunning | internal/nodes/consensus_specific_test.go:9 | None |
| Peers | internal/nodes/consensus_specific_test.go:10 | None |
| DialSeed | internal/nodes/consensus_specific_test.go:11 | None |
| ConsensusType | internal/nodes/consensus_specific_test.go:12 | None |
| TestDataoperationsPlaceholder | data_operations_test.go:5 | None |
| NewCustodialNode | core/custodial_node.go:18 | None |
| Custody | core/custodial_node.go:27 | None |
| AuthorizeRelayer | core/custodial_node.go:35 | None |
| RevokeRelayer | core/custodial_node.go:42 | None |
| IsRelayerAuthorized | core/custodial_node.go:49 | None |
| Release | core/custodial_node.go:59 | None |
| Balance | core/custodial_node.go:74 | None |
| NewHolographicNode | internal/nodes/holographic_node.go:17 | None |
| Store | internal/nodes/holographic_node.go:25 | None |
| Retrieve | internal/nodes/holographic_node.go:33 | None |
| NewCrossChainRegistry | core/cross_chain_contracts.go:23 | None |
| AuthorizeRelayer | core/cross_chain_contracts.go:31 | None |
| RevokeRelayer | core/cross_chain_contracts.go:38 | None |
| IsRelayerAuthorized | core/cross_chain_contracts.go:45 | None |
| RegisterMapping | core/cross_chain_contracts.go:52 | None |
| GetMapping | core/cross_chain_contracts.go:63 | None |
| ListMappings | core/cross_chain_contracts.go:71 | None |
| RemoveMapping | core/cross_chain_contracts.go:82 | None |
| NewIndex | internal/nodes/authority_nodes/index.go:19 | None |
| Add | internal/nodes/authority_nodes/index.go:24 | None |
| Get | internal/nodes/authority_nodes/index.go:34 | None |
| Remove | internal/nodes/authority_nodes/index.go:42 | None |
| List | internal/nodes/authority_nodes/index.go:49 | None |
| NewAgriculturalRegistry | core/token_syn4900.go:34 | None |
| Register | core/token_syn4900.go:39 | None |
| Transfer | core/token_syn4900.go:49 | None |
| UpdateStatus | core/token_syn4900.go:60 | None |
| Get | core/token_syn4900.go:71 | None |
| NewSYN223Token | core/syn223_token.go:20 | None |
| AddToWhitelist | core/syn223_token.go:34 | None |
| RemoveFromWhitelist | core/syn223_token.go:41 | None |
| AddToBlacklist | core/syn223_token.go:48 | None |
| RemoveFromBlacklist | core/syn223_token.go:55 | None |
| Transfer | core/syn223_token.go:62 | None |
| BalanceOf | core/syn223_token.go:82 | None |
| NewNATManager | core/nat_traversal.go:14 | None |
| MapPort | core/nat_traversal.go:19 | None |
| GetPort | core/nat_traversal.go:26 | None |
| RemoveMapping | core/nat_traversal.go:34 | None |
| Map | core/nat_traversal.go:41 | None |
| Unmap | core/nat_traversal.go:44 | None |
| SetExternalIP | core/nat_traversal.go:48 | None |
| ExternalIP | core/nat_traversal.go:55 | None |
| TestIndex | internal/nodes/authority_nodes/index_test.go:5 | None |
| NewAccessController | core/access_control.go:19 | None |
| Grant | core/access_control.go:24 | None |
| Revoke | core/access_control.go:41 | None |
| Has | core/access_control.go:57 | None |
| List | core/access_control.go:69 | None |
| GrantRole | core/access_control.go:87 | None |
| RevokeRole | core/access_control.go:96 | None |
| HasRole | core/access_control.go:105 | None |
| ListRoles | core/access_control.go:117 | None |
| TestLightnodePlaceholder | internal/nodes/light_node_test.go:5 | None |
| TestStorageMarketplaceLifecycle | core/storage_marketplace_test.go:12 | None |
| TestFuturesContract | core/syn3600_test.go:8 | None |
| TestHistoricalNodeArchiveRetrieve | internal/nodes/historical_node_test.go:8 | None |
| TestGeospatialPlaceholder | internal/nodes/extra/geospatial_test.go:5 | None |
| TestHolographicnodePlaceholder | node_ext/holographic_node_test.go:5 | None |
| TestHolographicnodePlaceholder | internal/nodes/extra/holographic_node_test.go:5 | None |
| TestAccessController | core/access_control_test.go:5 | None |
| TestRegulatoryNode | core/regulatory_node_test.go:5 | None |
| DefaultGasTable | core/gas_table.go:22 | None |
| parseGasGuide | core/gas_table.go:38 | None |
| SetGasCost | core/gas_table.go:72 | None |
| GasTableSnapshot | core/gas_table.go:80 | None |
| GasTableSnapshotJSON | core/gas_table.go:95 | None |
| WriteGasTableSnapshot | core/gas_table.go:121 | None |
| GasCostByName | core/gas_table.go:132 | None |
| TestIndexPlaceholder | node_ext/military_nodes/index_test.go:5 | None |
| TestIndexPlaceholder | internal/nodes/extra/military_nodes/index_test.go:5 | None |
| NewBiometricSecurityNode | core/biometric_security_node.go:17 | None |
| GetID | core/biometric_security_node.go:25 | None |
| Enroll | core/biometric_security_node.go:33 | None |
| Remove | core/biometric_security_node.go:38 | None |
| Authenticate | core/biometric_security_node.go:43 | None |
| SecureAddTransaction | core/biometric_security_node.go:49 | None |
| SecureExecute | core/biometric_security_node.go:59 | None |
| TestMineBlockFeeDistribution | core/node_test.go:5 | None |
| NewNodeAdapter | core/node_adapter.go:12 | None |
| TestIndexPlaceholder | internal/nodes/extra/watchtower/index_test.go:5 | None |
| TestIndexPlaceholder | node_ext/watchtower/index_test.go:5 | None |
| NewLoanPool | core/loanpool.go:18 | None |
| SubmitProposal | core/loanpool.go:27 | None |
| VoteProposal | core/loanpool.go:38 | None |
| Tick | core/loanpool.go:51 | None |
| Disburse | core/loanpool.go:61 | None |
| GetProposal | core/loanpool.go:78 | None |
| ListProposals | core/loanpool.go:84 | None |
| CancelProposal | core/loanpool.go:94 | None |
| ExtendProposal | core/loanpool.go:107 | None |
| TestPlasmaManagement | core/plasma_management_test.go:5 | None |
| NewHolographicNode | internal/nodes/extra/holographic_node.go:17 | None |
| ID | internal/nodes/extra/holographic_node.go:25 | None |
| Start | internal/nodes/extra/holographic_node.go:29 | None |
| Stop | internal/nodes/extra/holographic_node.go:33 | None |
| Store | internal/nodes/extra/holographic_node.go:36 | None |
| Retrieve | internal/nodes/extra/holographic_node.go:44 | None |
| NewHolographicNode | node_ext/holographic_node.go:17 | None |
| ID | node_ext/holographic_node.go:25 | None |
| Start | node_ext/holographic_node.go:29 | None |
| Stop | node_ext/holographic_node.go:33 | None |
| Store | node_ext/holographic_node.go:36 | None |
| Retrieve | node_ext/holographic_node.go:44 | None |
| TestHistoricalnodePlaceholder | internal/nodes/extra/historical_node_test.go:5 | None |
| TestHistoricalnodePlaceholder | node_ext/historical_node_test.go:5 | None |
| TestForensicnodePlaceholder | node_ext/forensic_node_test.go:5 | None |
| TestForensicnodePlaceholder | internal/nodes/extra/forensic_node_test.go:5 | None |
| NewSYN500Token | core/syn500.go:23 | None |
| Grant | core/syn500.go:28 | None |
| Use | core/syn500.go:33 | None |
| TestForensicNodeRecords | internal/nodes/forensic_node_test.go:10 | None |
| TestForensicNodeConcurrent | internal/nodes/forensic_node_test.go:28 | None |
| TestOptimizationPlaceholder | internal/nodes/extra/optimization_nodes/optimization_test.go:5 | None |
| Optimize | internal/nodes/extra/optimization_nodes/optimization.go:25 | None |
| tx | internal/nodes/optimization_nodes/optimization_test.go:10 | None |
| TestFeeOptimizerSortsByDensity | internal/nodes/optimization_nodes/optimization_test.go:14 | None |
| TestFeeOptimizerStableSort | internal/nodes/optimization_nodes/optimization_test.go:37 | None |
| TestFeeOptimizerZeroSize | internal/nodes/optimization_nodes/optimization_test.go:51 | None |
| TestFeeOptimizerEmpty | internal/nodes/optimization_nodes/optimization_test.go:67 | None |
| TestFeeOptimizerConcurrent | internal/nodes/optimization_nodes/optimization_test.go:77 | None |
| TestSyn3500tokenPlaceholder | core/syn3500_token_test.go:5 | None |
| TestIndexPlaceholder | internal/nodes/extra/optimization_nodes/index_test.go:5 | None |
| Optimize | internal/nodes/optimization_nodes/optimization.go:25 | None |
| NewBillRegistry | core/syn3200.go:32 | None |
| Create | core/syn3200.go:37 | None |
| Pay | core/syn3200.go:47 | None |
| Adjust | core/syn3200.go:65 | None |
| Get | core/syn3200.go:75 | None |
| Select | internal/nodes/optimization_nodes/index.go:20 | None |
| TestIndexPlaceholder | internal/nodes/extra/index_test.go:5 | None |
| TestComplianceManager | core/compliance_management_test.go:5 | None |
| TestIndexPlaceholder | internal/nodes/optimization_nodes/index_test.go:5 | None |
| NewSYN131Registry | core/syn131_token.go:20 | None |
| Create | core/syn131_token.go:25 | None |
| UpdateValuation | core/syn131_token.go:35 | None |
| Get | core/syn131_token.go:45 | None |
| NewExperimentalNode | internal/nodes/experimental_node.go:15 | None |
| ID | internal/nodes/experimental_node.go:20 | None |
| Start | internal/nodes/experimental_node.go:23 | None |
| Stop | internal/nodes/experimental_node.go:29 | None |
| IsRunning | internal/nodes/experimental_node.go:35 | None |
| Peers | internal/nodes/experimental_node.go:38 | None |
| DialSeed | internal/nodes/experimental_node.go:41 | None |
| NewElectedAuthorityNode | internal/nodes/elected_authority_node.go:16 | None |
| IsActive | internal/nodes/elected_authority_node.go:22 | None |
| NewNetwork | core/network.go:24 | None |
| Start | core/network.go:36 | None |
| Stop | core/network.go:50 | None |
| AddNode | core/network.go:66 | None |
| AddRelay | core/network.go:73 | None |
| Peers | core/network.go:80 | None |
| EnqueueTransaction | core/network.go:94 | None |
| Broadcast | core/network.go:107 | None |
| Subscribe | core/network.go:117 | None |
| Publish | core/network.go:127 | None |
| processQueue | core/network.go:142 | None |
| broadcast | core/network.go:157 | None |
| TestOptimizationPlaceholder | node_ext/optimization_nodes/optimization_test.go:5 | None |
| Optimize | node_ext/optimization_nodes/optimization.go:25 | None |
| TestIndexPlaceholder | node_ext/optimization_nodes/index_test.go:5 | None |
| TestIndexPlaceholder | node_ext/index_test.go:5 | None |
| TestContentNetworkNode_RegisterUnregister | core/content_node_test.go:5 | None |
| TestContentNetworkNode_ContentAndList | core/content_node_test.go:22 | None |
| TestTrainingLifecycle | ai_training_test.go:5 | None |
| TestSyn700Placeholder | core/syn700_test.go:5 | None |
| TestAIContractRegistry | core/ai_enhanced_contract_test.go:9 | None |
| NewSYN3500Token | core/syn3500_token.go:19 | None |
| SetRate | core/syn3500_token.go:30 | None |
| Info | core/syn3500_token.go:37 | None |
| Mint | core/syn3500_token.go:44 | None |
| Redeem | core/syn3500_token.go:51 | None |
| BalanceOf | core/syn3500_token.go:63 | None |
| TestEncryptor | internal/security/encryption_test.go:5 | None |
| ID | internal/nodes/bank_nodes/index_test.go:13 | None |
| Start | internal/nodes/bank_nodes/index_test.go:14 | None |
| Stop | internal/nodes/bank_nodes/index_test.go:15 | None |
| IsRunning | internal/nodes/bank_nodes/index_test.go:16 | None |
| Peers | internal/nodes/bank_nodes/index_test.go:17 | None |
| DialSeed | internal/nodes/bank_nodes/index_test.go:18 | None |
| ID | internal/nodes/bank_nodes/index_test.go:22 | None |
| Start | internal/nodes/bank_nodes/index_test.go:23 | None |
| Stop | internal/nodes/bank_nodes/index_test.go:24 | None |
| IsRunning | internal/nodes/bank_nodes/index_test.go:25 | None |
| Peers | internal/nodes/bank_nodes/index_test.go:26 | None |
| DialSeed | internal/nodes/bank_nodes/index_test.go:27 | None |
| MintCBDC | internal/nodes/bank_nodes/index_test.go:28 | None |
| ID | internal/nodes/bank_nodes/index_test.go:34 | None |
| Start | internal/nodes/bank_nodes/index_test.go:35 | None |
| Stop | internal/nodes/bank_nodes/index_test.go:36 | None |
| IsRunning | internal/nodes/bank_nodes/index_test.go:37 | None |
| Peers | internal/nodes/bank_nodes/index_test.go:38 | None |
| DialSeed | internal/nodes/bank_nodes/index_test.go:39 | None |
| TestInterfaceCompliance | internal/nodes/bank_nodes/index_test.go:42 | None |
| TestReplicator | core/replication_test.go:5 | None |
| NewBasicNode | internal/nodes/index.go:42 | None |
| ID | internal/nodes/index.go:47 | None |
| Start | internal/nodes/index.go:53 | None |
| Stop | internal/nodes/index.go:65 | None |
| IsRunning | internal/nodes/index.go:76 | None |
| Peers | internal/nodes/index.go:83 | None |
| DialSeed | internal/nodes/index.go:94 | None |
| NewSyn2500Member | core/syn2500_token.go:18 | None |
| UpdateVotingPower | core/syn2500_token.go:33 | None |
| NewSyn2500Registry | core/syn2500_token.go:44 | None |
| AddMember | core/syn2500_token.go:49 | None |
| GetMember | core/syn2500_token.go:56 | None |
| RemoveMember | core/syn2500_token.go:64 | None |
| ListMembers | core/syn2500_token.go:71 | None |
| TestBasicNodeLifecycle | internal/nodes/index_test.go:5 | None |
| NewGrantRegistry | core/syn3800.go:26 | None |
| CreateGrant | core/syn3800.go:31 | None |
| Disburse | core/syn3800.go:41 | None |
| GetGrant | core/syn3800.go:59 | None |
| ListGrants | core/syn3800.go:71 | None |
| TestIsZeroAddress | core/address_zero_test.go:5 | None |
| NewEncryptor | internal/security/encryption.go:9 | None |
| Encrypt | internal/security/encryption.go:12 | None |
| Decrypt | internal/security/encryption.go:21 | None |
| TestBiometricService | core/biometric_test.go:10 | None |
| TestSupplyChainRegistryRegisterGet | core/syn1300_test.go:7 | None |
| TestSupplyChainRegistryDuplicate | core/syn1300_test.go:25 | None |
| TestSupplyChainRegistryUpdate | core/syn1300_test.go:35 | None |
| TestSupplyChainRegistryUpdateNonexistent | core/syn1300_test.go:56 | None |
| NewInitService | core/initialization_replication.go:14 | None |
| Start | core/initialization_replication.go:19 | None |
| Stop | core/initialization_replication.go:30 | None |
| TestWebRTCRPC | core/rpc_webrtc_test.go:5 | None |
| NewPatchManager | internal/security/patch_manager.go:9 | None |
| Apply | internal/security/patch_manager.go:12 | None |
| Applied | internal/security/patch_manager.go:15 | None |
| NewDAOManager | core/dao.go:28 | None |
| AuthorizeRelayer | core/dao.go:33 | None |
| RevokeRelayer | core/dao.go:40 | None |
| IsRelayerAuthorized | core/dao.go:47 | None |
| Create | core/dao.go:55 | None |
| Join | core/dao.go:72 | None |
| Leave | core/dao.go:92 | None |
| Info | core/dao.go:112 | None |
| List | core/dao.go:123 | None |
| NewSecretsManager | internal/security/secrets_manager.go:11 | None |
| Store | internal/security/secrets_manager.go:17 | None |
| Retrieve | internal/security/secrets_manager.go:30 | None |
| TestEncryptDecrypt | core/private_transactions_test.go:8 | None |
| TestPrivateTxManager | core/private_transactions_test.go:24 | None |
| NewRateLimiter | internal/security/rate_limiter.go:12 | None |
| Allow | internal/security/rate_limiter.go:17 | None |
| TestPatchManager | internal/security/patch_manager_test.go:5 | None |
| TestNATManager | core/nat_traversal_test.go:5 | None |
| NewKeyManager | internal/security/key_management.go:12 | None |
| Rotate | internal/security/key_management.go:17 | None |
| Key | internal/security/key_management.go:24 | None |
| NewConsensusSwitcher | core/consensus_specific.go:21 | None |
| Evaluate | core/consensus_specific.go:27 | None |
| Mode | core/consensus_specific.go:49 | None |
| NewSystemHealthLogger | core/system_health_logging.go:19 | None |
| Collect | core/system_health_logging.go:25 | None |
| Snapshot | core/system_health_logging.go:44 | None |
| TestNewNodeAdapter | core/node_adapter_test.go:11 | None |
| TestRateLimiterAllow | internal/security/rate_limiter_test.go:8 | None |
| TestEventTickets | core/syn1700_token_test.go:5 | None |
| NewLightNode | internal/nodes/light_node.go:19 | None |
| AddHeader | internal/nodes/light_node.go:24 | None |
| Headers | internal/nodes/light_node.go:29 | None |
| LatestHeader | internal/nodes/light_node.go:37 | None |
| NewAuditManager | core/audit_management.go:35 | None |
| NewAuditManagerFromKey | core/audit_management.go:44 | None |
| Log | core/audit_management.go:53 | None |
| List | core/audit_management.go:83 | None |
| Verify | core/audit_management.go:94 | None |
| PublicKey | core/audit_management.go:108 | None |
| NewZeroTrustEngine | core/zero_trust_data_channels.go:31 | None |
| OpenChannel | core/zero_trust_data_channels.go:36 | None |
| Send | core/zero_trust_data_channels.go:51 | None |
| Messages | core/zero_trust_data_channels.go:70 | None |
| Receive | core/zero_trust_data_channels.go:83 | None |
| CloseChannel | core/zero_trust_data_channels.go:105 | None |
| NewConsensusHopper | dynamic_consensus_hopping.go:32 | None |
| Mode | dynamic_consensus_hopping.go:37 | None |
| SetMode | dynamic_consensus_hopping.go:44 | None |
| LastMetrics | dynamic_consensus_hopping.go:51 | None |
| Evaluate | dynamic_consensus_hopping.go:58 | None |
| TestIsZeroAddress | address_zero_test.go:5 | None |
| TestSecretsManager | internal/security/secrets_manager_test.go:5 | None |
| NewLightNode | core/light_node.go:12 | None |
| AddHeader | core/light_node.go:17 | None |
| LatestHeader | core/light_node.go:20 | None |
| Headers | core/light_node.go:28 | None |
| NewMiningNode | core/mining_node.go:19 | None |
| Start | core/mining_node.go:24 | None |
| Stop | core/mining_node.go:31 | None |
| IsMining | core/mining_node.go:38 | None |
| Mine | core/mining_node.go:45 | None |
| HashRateHint | core/mining_node.go:60 | None |
| NewNode | core/node.go:23 | None |
| AddTransaction | core/node.go:38 | None |
| ValidateTransaction | core/node.go:48 | None |
| MineBlock | core/node.go:59 | None |
| SetStake | core/node.go:99 | None |
| eligibleStakes | core/node.go:107 | None |
| ReportDoubleSign | core/node.go:118 | None |
| ReportDowntime | core/node.go:123 | None |
| Rehabilitate | core/node.go:128 | None |
| slash | core/node.go:132 | None |
| TestLoanpoolapplyPlaceholder | core/loanpool_apply_test.go:5 | None |
| Encrypt | core/private_transactions.go:13 | None |
| Decrypt | core/private_transactions.go:31 | None |
| NewPrivateTxManager | core/private_transactions.go:62 | None |
| Send | core/private_transactions.go:67 | None |
| List | core/private_transactions.go:74 | None |
| BenchmarkTransactionManagerLockAndMint | cross_chain_transactions_benchmark_test.go:5 | None |
| BenchmarkTransactionManagerBurnAndRelease | cross_chain_transactions_benchmark_test.go:12 | None |
| BenchmarkTransactionManagerListTransactions | cross_chain_transactions_benchmark_test.go:19 | None |
| BenchmarkTransactionManagerGetTransaction | cross_chain_transactions_benchmark_test.go:30 | None |
| TestDDoSMitigator | internal/security/ddos_mitigation_test.go:5 | None |
| NewDDoSMitigator | internal/security/ddos_mitigation.go:9 | None |
| Block | internal/security/ddos_mitigation.go:14 | None |
| IsBlocked | internal/security/ddos_mitigation.go:17 | None |
| TestKeyManager | internal/security/key_management_test.go:5 | None |
| TestDeployContractNoArgs | tests/scripts/deploy_contract_test.go:11 | None |
| TestDeployContractMissingFile | tests/scripts/deploy_contract_test.go:23 | None |
| TestDeployContractMissingBinary | tests/scripts/deploy_contract_test.go:35 | None |
| TestGateway | internal/api/api_test.go:5 | None |
| TestAuthMiddleware | internal/api/api_test.go:12 | None |
| TestRateLimiter | internal/api/api_test.go:19 | None |
| NewRateLimiter | internal/api/rate_limiter.go:13 | None |
| Allow | internal/api/rate_limiter.go:18 | None |
| TestAuthenticate | internal/api/auth_middleware_test.go:5 | None |
| TestGatewayStart | internal/api/gateway_test.go:5 | None |
| execCLI | tests/contracts/faucet_test.go:16 | None |
| TestTokenFaucetTemplate | tests/contracts/faucet_test.go:45 | None |
| TestRateLimiterAllow | internal/api/rate_limiter_test.go:5 | None |
| Authenticate | internal/api/auth_middleware.go:7 | None |
| NewGateway | internal/api/gateway.go:7 | None |
| Start | internal/api/gateway.go:10 | None |
| FuzzVM | tests/fuzz/vm_fuzz_test.go:5 | None |
| FuzzNetwork | tests/fuzz/network_fuzz_test.go:5 | None |
| FuzzCrypto | tests/fuzz/crypto_fuzz_test.go:5 | None |
| TestLogPlaceholder | internal/log/log_test.go:5 | None |
| execCLI | tests/gui_wallet_test.go:20 | None |
| TestGUIWalletIntegration | tests/gui_wallet_test.go:42 | None |
| log | internal/log/log.go:17 | None |
| Info | internal/log/log.go:37 | None |
| Error | internal/log/log.go:40 | None |
| TestDataPlaceholder | data_test.go:5 | None |
| TestCrosschainbridgePlaceholder | cross_chain_bridge_test.go:5 | None |
| TestContractsFormalVerification | tests/formal/contracts_verification_test.go:5 | None |
| TestAIServicePredictAndPublish | ai_test.go:8 | None |
| TestZeroTrustEngineBasicFlow | zero_trust_data_channels_test.go:12 | None |
| TestZeroTrustEngineDuplicateOpen | zero_trust_data_channels_test.go:47 | None |
| TestZeroTrustEngineSendErrors | zero_trust_data_channels_test.go:59 | None |
| TestZeroTrustEngineMessagesIsolation | zero_trust_data_channels_test.go:78 | None |
| TestZeroTrustEngineCloseChannel | zero_trust_data_channels_test.go:100 | None |
| TestZeroTrustEngineConcurrentSend | zero_trust_data_channels_test.go:108 | None |
| TestZeroTrustEngineMultipleChannels | zero_trust_data_channels_test.go:133 | None |
| TestZeroTrustEngineMessagesUnknown | zero_trust_data_channels_test.go:155 | None |
| execCLI | tests/e2e/network_harness_test.go:23 | None |
| TestNetworkHarness | tests/e2e/network_harness_test.go:55 | None |
| NewBiometricsAuth | biometrics_auth.go:21 | None |
| Enroll | biometrics_auth.go:26 | None |
| Verify | biometrics_auth.go:33 | None |
| Remove | biometrics_auth.go:48 | None |
| Enrolled | biometrics_auth.go:55 | None |
| List | biometrics_auth.go:64 | None |
| TestStakePenaltyManager | core/stake_penalty_test.go:5 | None |
| NewSyncManager | core/blockchain_synchronization.go:19 | None |
| Start | core/blockchain_synchronization.go:24 | None |
| Stop | core/blockchain_synchronization.go:31 | None |
| Status | core/blockchain_synchronization.go:38 | None |
| Once | core/blockchain_synchronization.go:45 | None |
| execCommand | tests/cli_integration_test.go:18 | None |
| TestCLIIntegration | tests/cli_integration_test.go:39 | None |
| TestEncryptDecrypt | private_transactions_test.go:8 | None |
| TestPrivateTxManager | private_transactions_test.go:24 | None |
| validRole | core/dao_access_control.go:18 | None |
| AddMember | core/dao_access_control.go:23 | None |
| UpdateMemberRole | core/dao_access_control.go:40 | None |
| RemoveMember | core/dao_access_control.go:57 | None |
| MemberRole | core/dao_access_control.go:68 | None |
| IsMember | core/dao_access_control.go:76 | None |
| IsAdmin | core/dao_access_control.go:84 | None |
| MembersList | core/dao_access_control.go:91 | None |
| NewSystemHealthLogger | system_health_logging.go:19 | None |
| Collect | system_health_logging.go:25 | None |
| Snapshot | system_health_logging.go:44 | None |
| NewContractManager | core/contract_management.go:16 | None |
| Transfer | core/contract_management.go:21 | None |
| Pause | core/contract_management.go:36 | None |
| Resume | core/contract_management.go:51 | None |
| Upgrade | core/contract_management.go:66 | None |
| Info | core/contract_management.go:87 | None |
| TestGovernmentAuthorityNode | core/government_authority_node_test.go:5 | None |
| QuadraticWeight | core/dao_quadratic_voting.go:9 | None |
| CastQuadraticVote | core/dao_quadratic_voting.go:15 | None |
| NewForensicNode | core/forensic_node.go:21 | None |
| RecordTransaction | core/forensic_node.go:26 | None |
| RecordNetworkTrace | core/forensic_node.go:34 | None |
| Transactions | core/forensic_node.go:42 | None |
| NetworkTraces | core/forensic_node.go:51 | None |
| TestRollupAggregator | core/rollups_test.go:5 | None |
| NewBridgeManager | core/cross_chain_bridge.go:39 | None |
| RegisterBridge | core/cross_chain_bridge.go:48 | None |
| ListBridges | core/cross_chain_bridge.go:62 | None |
| GetBridge | core/cross_chain_bridge.go:73 | None |
| AuthorizeRelayer | core/cross_chain_bridge.go:84 | None |
| RevokeRelayer | core/cross_chain_bridge.go:96 | None |
| IsRelayerAuthorized | core/cross_chain_bridge.go:109 | None |
| RemoveBridge | core/cross_chain_bridge.go:122 | None |
| Deposit | core/cross_chain_bridge.go:133 | None |
| Claim | core/cross_chain_bridge.go:154 | None |
| GetTransfer | core/cross_chain_bridge.go:180 | None |
| ListTransfers | core/cross_chain_bridge.go:191 | None |
| NewBridgeTransferManager | core/cross_chain_bridge.go:220 | None |
| Deposit | core/cross_chain_bridge.go:225 | None |
| Claim | core/cross_chain_bridge.go:244 | None |
| GetTransfer | core/cross_chain_bridge.go:259 | None |
| ListTransfers | core/cross_chain_bridge.go:267 | None |
| NewProtocolRegistry | core/cross_chain_agnostic_protocols.go:23 | None |
| AuthorizeRelayer | core/cross_chain_agnostic_protocols.go:31 | None |
| RevokeRelayer | core/cross_chain_agnostic_protocols.go:38 | None |
| IsRelayerAuthorized | core/cross_chain_agnostic_protocols.go:45 | None |
| Register | core/cross_chain_agnostic_protocols.go:53 | None |
| Remove | core/cross_chain_agnostic_protocols.go:66 | None |
| List | core/cross_chain_agnostic_protocols.go:80 | None |
| Get | core/cross_chain_agnostic_protocols.go:91 | None |
| NewTransactionManager | cross_chain_transactions.go:29 | None |
| LockAndMint | cross_chain_transactions.go:34 | None |
| BurnAndRelease | cross_chain_transactions.go:51 | None |
| ListTransactions | cross_chain_transactions.go:68 | None |
| GetTransaction | cross_chain_transactions.go:79 | None |
| TestCrosschainPlaceholder | cross_chain_test.go:5 | None |
| TestIDRegistry | core/idwallet_registration_test.go:5 | None |
| testScriptHelp | scripts/scripts_test.go:9 | None |
| TestPackageReleaseHelp | scripts/scripts_test.go:20 | None |
| TestGenerateDocsHelp | scripts/scripts_test.go:24 | None |
| TestCIScriptHelp | scripts/scripts_test.go:28 | None |
| TestBackupLedgerHelp | scripts/scripts_test.go:32 | None |
| NewContractManager | contract_management.go:11 | None |
| Transfer | contract_management.go:16 | None |
| Pause | contract_management.go:28 | None |
| Resume | contract_management.go:40 | None |
| Upgrade | contract_management.go:52 | None |
| Info | contract_management.go:70 | None |
| Catalogue | core/opcode.go:73 | None |
| Opcodes | core/opcode.go:85 | None |
| Register | core/opcode.go:97 | None |
| Dispatch | core/opcode.go:107 | None |
| wrap | core/opcode.go:123 | None |
| init | core/opcode.go:1630 | None |
| Hex | core/opcode.go:1663 | None |
| Bytes | core/opcode.go:1666 | None |
| String | core/opcode.go:1675 | None |
| ParseOpcode | core/opcode.go:1678 | None |
| MustParseOpcode | core/opcode.go:1686 | None |
| DebugDump | core/opcode.go:1696 | None |
| ToBytecode | core/opcode.go:1715 | None |
| HexDump | core/opcode.go:1724 | None |
| Evaluate | environmental_monitoring_node.go:17 | None |
| NewEnvironmentalMonitoringNode | environmental_monitoring_node.go:45 | None |
| SetCondition | environmental_monitoring_node.go:50 | None |
| Trigger | environmental_monitoring_node.go:57 | None |
| NewBridgeTransferManager | cross_chain_bridge.go:31 | None |
| Deposit | cross_chain_bridge.go:36 | None |
| Claim | cross_chain_bridge.go:53 | None |
| GetTransfer | cross_chain_bridge.go:70 | None |
| ListTransfers | cross_chain_bridge.go:78 | None |
| NewTransaction | core/transaction.go:42 | None |
| Hash | core/transaction.go:58 | None |
| Verify | core/transaction.go:75 | None |
| AttachBiometric | core/transaction.go:84 | None |
| NewDataDistribution | data_distribution.go:19 | None |
| Offer | data_distribution.go:24 | None |
| Revoke | data_distribution.go:37 | None |
| Meta | data_distribution.go:49 | None |
| Locations | data_distribution.go:60 | None |
| init | cli/watchtower.go:21 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn5000indexPlaceholder | core/syn5000_index_test.go:5 | None |
| TestGrantRegistry | core/syn3800_test.go:5 | None |
| NewDataResourceManager | data_resource_management.go:15 | None |
| Put | data_resource_management.go:20 | None |
| Get | data_resource_management.go:32 | None |
| Delete | data_resource_management.go:45 | None |
| Keys | data_resource_management.go:55 | None |
| Usage | data_resource_management.go:66 | None |
| TestOpcodePlaceholder | core/opcode_test.go:5 | None |
| TestChainConnectionManager | core/cross_chain_connection_test.go:5 | None |
| NewLedger | core/ledger.go:37 | None |
| replayWAL | core/ledger.go:54 | None |
| appendWAL | core/ledger.go:74 | None |
| Head | core/ledger.go:87 | None |
| GetBlock | core/ledger.go:98 | None |
| AddBlock | core/ledger.go:108 | None |
| GetBalance | core/ledger.go:116 | None |
| GetUTXOs | core/ledger.go:123 | None |
| Credit | core/ledger.go:135 | None |
| Mint | core/ledger.go:144 | None |
| Transfer | core/ledger.go:150 | None |
| ApplyTransaction | core/ledger.go:158 | None |
| AddToPool | core/ledger.go:173 | None |
| Pool | core/ledger.go:180 | None |
| newUTXO | core/ledger.go:188 | None |
| updateUTXO | core/ledger.go:194 | None |
| TestDAOManager | core/dao_test.go:5 | None |
| TestBenefitRegistry | core/syn3900_test.go:5 | None |
| Pause | core/plasma_management.go:4 | None |
| Resume | core/plasma_management.go:11 | None |
| Status | core/plasma_management.go:18 | None |
| opcodeByName | contracts_opcodes.go:1428 | None |
| opcodeByName | core/contracts_opcodes.go:22 | None |
| NewSidechainOps | core/sidechain_ops.go:11 | None |
| Deposit | core/sidechain_ops.go:16 | None |
| Withdraw | core/sidechain_ops.go:28 | None |
| EscrowBalance | core/sidechain_ops.go:47 | None |
| BlockReward | core/coin.go:21 | None |
| CirculatingSupply | core/coin.go:32 | None |
| RemainingSupply | core/coin.go:62 | None |
| InitialPrice | core/coin.go:72 | None |
| AlphaFactor | core/coin.go:78 | None |
| MinimumStake | core/coin.go:85 | None |
| LockupDuration | core/coin.go:94 | None |
| PriceToSupplyRatio | core/coin.go:103 | None |
| NewEnergyEfficientNode | energy_efficient_node.go:26 | None |
| ID | energy_efficient_node.go:31 | None |
| RecordUsage | energy_efficient_node.go:34 | None |
| AddOffset | energy_efficient_node.go:39 | None |
| OffsetCredits | energy_efficient_node.go:46 | None |
| Certify | energy_efficient_node.go:53 | None |
| Certificate | energy_efficient_node.go:68 | None |
| ShouldThrottle | energy_efficient_node.go:75 | None |
| TestConsensusNetworkManager | core/cross_consensus_scaling_networks_test.go:5 | None |
| GasCost | core/gas.go:16 | None |
| initGasTable | core/gas.go:24 | None |
| TestHistoricalNode_ArchiveAndRetrieve | core/historical_node_test.go:9 | None |
| TestHistoricalNode_Duplicate | core/historical_node_test.go:26 | None |
| TestNewCustodialNode | core/custodial_node_test.go:6 | None |
| TestCustodialNodeCustody | core/custodial_node_test.go:24 | None |
| TestCustodialNodeRelease | core/custodial_node_test.go:35 | None |
| TestLedgerApplyTransaction | core/ledger_test.go:5 | None |
| TestLedgerUTXOAndPool | core/ledger_test.go:24 | None |
| TotalVotes | core/authority_nodes.go:22 | None |
| MarshalJSON | core/authority_nodes.go:25 | None |
| NewAuthorityNodeRegistry | core/authority_nodes.go:43 | None |
| Register | core/authority_nodes.go:48 | None |
| Vote | core/authority_nodes.go:61 | None |
| RemoveVote | core/authority_nodes.go:80 | None |
| Electorate | core/authority_nodes.go:89 | None |
| IsAuthorityNode | core/authority_nodes.go:109 | None |
| Info | core/authority_nodes.go:117 | None |
| List | core/authority_nodes.go:128 | None |
| Deregister | core/authority_nodes.go:135 | None |
| TestContractRegistry | core/contracts_test.go:5 | None |
| NewGatewayNode | core/gateway_node.go:22 | None |
| RegisterEndpoint | core/gateway_node.go:31 | None |
| Handle | core/gateway_node.go:36 | None |
| RemoveEndpoint | core/gateway_node.go:44 | None |
| Endpoints | core/gateway_node.go:49 | None |
| TestStakingNodeStakeAndUnstake | core/staking_node_test.go:5 | None |
| TestStakingNodeTotal | core/staking_node_test.go:21 | None |
| TestInitServiceStartStop | core/initialization_replication_test.go:5 | None |
| TestLedgerCompressionRoundTrip | core/blockchain_compression_test.go:5 | None |
| TestBiometricSecurityNode | core/biometric_security_node_test.go:10 | None |
| NewStakePenaltyManager | core/stake_penalty.go:7 | None |
| Slash | core/stake_penalty.go:10 | None |
| Reward | core/stake_penalty.go:15 | None |
| String | core/charity.go:24 | None |
| NewCharityPool | core/charity.go:72 | None |
| Deposit | core/charity.go:77 | None |
| Register | core/charity.go:91 | None |
| Vote | core/charity.go:105 | None |
| Tick | core/charity.go:117 | None |
| Winners | core/charity.go:121 | None |
| GetRegistration | core/charity.go:126 | None |
| mustJSON | core/charity.go:140 | None |
| TestNewConsensusSpecificNode | core/consensus_specific_node_test.go:7 | None |
| TestImmutabilityEnforcer | core/immutability_enforcement_test.go:5 | None |
| TestConnectionPool | core/connection_pool_test.go:10 | None |
| NewSubBlock | core/block.go:22 | None |
| Hash | core/block.go:31 | None |
| VerifySignature | core/block.go:42 | None |
| Validate | core/block.go:50 | None |
| NewBlock | core/block.go:93 | None |
| HeaderHash | core/block.go:99 | None |
| signSubBlock | core/block.go:109 | None |
| Validate | core/block.go:117 | None |
| NewAdaptiveManager | core/consensus_adaptive_management.go:20 | None |
| record | core/consensus_adaptive_management.go:29 | None |
| averages | core/consensus_adaptive_management.go:36 | None |
| Adjust | core/consensus_adaptive_management.go:50 | None |
| Threshold | core/consensus_adaptive_management.go:64 | None |
| Weights | core/consensus_adaptive_management.go:76 | None |
| RecordMetrics | core/consensus_adaptive_management.go:87 | None |
| Reset | core/consensus_adaptive_management.go:94 | None |
| NewReplicator | core/replication.go:15 | None |
| Start | core/replication.go:20 | None |
| Stop | core/replication.go:27 | None |
| Status | core/replication.go:34 | None |
| ReplicateBlock | core/replication.go:42 | None |
| TestForensicNode_RecordAndRetrieve | core/forensic_node_test.go:9 | None |
| TestFullNodeModes | core/full_node_test.go:9 | None |
| TestMobileMiningNode | core/mobile_mining_node_test.go:5 | None |
| TestSidechainOps | core/sidechain_ops_test.go:5 | None |
| NewHistoricalNode | core/historical_node.go:21 | None |
| ArchiveBlock | core/historical_node.go:29 | None |
| GetBlockByHeight | core/historical_node.go:44 | None |
| GetBlockByHash | core/historical_node.go:52 | None |
| TotalBlocks | core/historical_node.go:60 | None |
| NewChainConnectionManager | core/cross_chain_connection.go:25 | None |
| Open | core/cross_chain_connection.go:30 | None |
| Close | core/cross_chain_connection.go:44 | None |
| Get | core/cross_chain_connection.go:59 | None |
| List | core/cross_chain_connection.go:70 | None |
| AuthorizeRelayer | core/cross_chain_connection.go:81 | None |
| RevokeRelayer | core/cross_chain_connection.go:96 | None |
| IsRelayerAuthorized | core/cross_chain_connection.go:109 | None |
| Remove | core/cross_chain_connection.go:122 | None |
| TestDAOProposal | core/dao_proposal_test.go:5 | None |
| NewElectedAuthorityNode | core/elected_authority_node.go:19 | None |
| IsActive | core/elected_authority_node.go:25 | None |
| RenewTerm | core/elected_authority_node.go:34 | None |
| TestRegulatoryManager | core/regulatory_management_test.go:5 | None |
| NewFirewall | core/firewall.go:14 | None |
| Allow | core/firewall.go:23 | None |
| Block | core/firewall.go:32 | None |
| IsAllowed | core/firewall.go:42 | None |
| Rules | core/firewall.go:56 | None |
| NewKademlia | core/kademlia.go:19 | None |
| Store | core/kademlia.go:24 | None |
| FindValue | core/kademlia.go:31 | None |
| Distance | core/kademlia.go:42 | None |
| Closest | core/kademlia.go:61 | None |
| TestMiningNode | core/mining_node_test.go:5 | None |
| TestDAOAccessControl | core/dao_access_control_test.go:5 | None |
| NewBridgeRegistry | core/cross_chain.go:24 | None |
| RegisterBridge | core/cross_chain.go:29 | None |
| ListBridges | core/cross_chain.go:45 | None |
| GetBridge | core/cross_chain.go:56 | None |
| AuthorizeRelayer | core/cross_chain.go:64 | None |
| RevokeRelayer | core/cross_chain.go:79 | None |
| IsRelayerAuthorized | core/cross_chain.go:92 | None |
| RemoveBridge | core/cross_chain.go:104 | None |
| TestQuorumTracker | core/quorum_tracker_test.go:5 | None |
| NewAuthorityApplicationManager | core/authority_apply.go:35 | None |
| Submit | core/authority_apply.go:45 | None |
| Vote | core/authority_apply.go:68 | None |
| Finalize | core/authority_apply.go:97 | None |
| Tick | core/authority_apply.go:118 | None |
| Get | core/authority_apply.go:129 | None |
| List | core/authority_apply.go:140 | None |
| MarshalJSON | core/authority_apply.go:151 | None |
| TestSYN4200Token | core/syn4200_token_test.go:6 | None |
| TestCampaignNotFound | core/syn4200_token_test.go:58 | None |
| BenchmarkTransactionHash | core/core_benchmarks_test.go:7 | None |
| BenchmarkLedgerApplyTransaction | core/core_benchmarks_test.go:25 | None |
| TestSyn500Placeholder | core/syn500_test.go:5 | None |
| NewCentralBankingNode | core/central_banking_node.go:19 | None |
| UpdatePolicy | core/central_banking_node.go:28 | None |
| MintCBDC | core/central_banking_node.go:32 | None |
| Mint | core/central_banking_node.go:42 | None |
| NewBiometricService | core/biometric.go:32 | None |
| Enroll | core/biometric.go:39 | None |
| Verify | core/biometric.go:61 | None |
| Remove | core/biometric.go:76 | None |
| Enrolled | core/biometric.go:83 | None |
| List | core/biometric.go:91 | None |
| NewMobileMiningNode | core/mobile_mining_node.go:13 | None |
| Start | core/mobile_mining_node.go:18 | None |
| Stop | core/mobile_mining_node.go:25 | None |
| IsMining | core/mobile_mining_node.go:28 | None |
| Mine | core/mobile_mining_node.go:31 | None |
| SetPowerLimit | core/mobile_mining_node.go:41 | None |
| PowerLimit | core/mobile_mining_node.go:48 | None |
| TestSwarmBroadcast | core/swarm_test.go:6 | None |
| TestShardManager | core/sharding_test.go:5 | None |
| TestFailoverManagerKeepsPrimaryWhenHealthy | core/high_availability_test.go:10 | None |
| TestFailoverManagerFailoverToLatestBackup | core/high_availability_test.go:31 | None |
| TestFailoverManagerHeartbeatAndRegister | core/high_availability_test.go:56 | None |
| NewEvent | core/syn1700_token.go:32 | None |
| IssueTicket | core/syn1700_token.go:45 | None |
| TransferTicket | core/syn1700_token.go:58 | None |
| VerifyTicket | core/syn1700_token.go:70 | None |
| NewFullNode | core/full_node.go:22 | None |
| SetMode | core/full_node.go:27 | None |
| CurrentMode | core/full_node.go:32 | None |
| IsArchive | core/full_node.go:37 | None |
| TestTokensyn4900Placeholder | core/token_syn4900_test.go:5 | None |
| TestIdentityService | core/identity_verification_test.go:5 | None |
| TestProtocolRegistryJSON | cli/cross_chain_agnostic_protocols_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| NewLiquidityPool | core/liquidity_pools.go:23 | None |
| AddLiquidity | core/liquidity_pools.go:34 | None |
| RemoveLiquidity | core/liquidity_pools.go:61 | None |
| Swap | core/liquidity_pools.go:78 | None |
| sqrt | core/liquidity_pools.go:110 | None |
| NewLiquidityPoolRegistry | core/liquidity_pools.go:126 | None |
| Create | core/liquidity_pools.go:131 | None |
| Get | core/liquidity_pools.go:143 | None |
| List | core/liquidity_pools.go:151 | None |
| TestGasPlaceholder | cli/gas_test.go:5 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSYN223TokenTransfer | core/syn223_token_test.go:5 | None |
| TestSYN223TokenWhitelistBlacklist | core/syn223_token_test.go:19 | None |
| TestSYN223TokenInsufficientBalance | core/syn223_token_test.go:54 | None |
| TestNewBankInstitutionalNode | core/bank_institutional_node_test.go:11 | None |
| TestBankInstitutionalNodeRegistration | core/bank_institutional_node_test.go:34 | None |
| TestBankInstitutionalNodeIsolation | core/bank_institutional_node_test.go:66 | None |
| TestRegisterEmptyInstitution | core/bank_institutional_node_test.go:88 | None |
| TestRemoveInstitution | core/bank_institutional_node_test.go:106 | None |
| TestRegisterInstitutionInvalidAddr | core/bank_institutional_node_test.go:125 | None |
| TestRegisterInstitutionInvalidSignature | core/bank_institutional_node_test.go:137 | None |
| TestRemoveInstitutionInvalidSignature | core/bank_institutional_node_test.go:150 | None |
| TestDAOQuadraticWeightJSON | cli/dao_quadratic_voting_test.go:18 | Rate limiting, validation, auth, integrity, user, crypto |
| TestDAOQuadraticVoteRequiresMembership | cli/dao_quadratic_voting_test.go:36 | Rate limiting, validation, auth, integrity, user, crypto |
| NewRollupManager | core/rollup_management.go:9 | None |
| Pause | core/rollup_management.go:14 | None |
| Resume | core/rollup_management.go:19 | None |
| Status | core/rollup_management.go:24 | None |
| TestXContractJSON | cli/cross_chain_contracts_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/consensus.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| TestComplianceServiceKYCAndRisk | core/compliance_test.go:9 | None |
| TestComplianceServiceMonitorTransaction | core/compliance_test.go:41 | None |
| TestComplianceServiceVerifyZKP | core/compliance_test.go:49 | None |
| TestRollupManager | core/rollup_management_test.go:5 | None |
| TestConsensusAdaptiveWeights | cli/consensus_adaptive_management_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| NewInvestmentRegistry | core/syn1401.go:25 | None |
| Issue | core/syn1401.go:30 | None |
| Accrue | core/syn1401.go:40 | None |
| Redeem | core/syn1401.go:57 | None |
| Get | core/syn1401.go:77 | None |
| TestInitGenesis | core/genesis_block_test.go:5 | None |
| TestBankNodeIndex | core/bank_nodes_index_test.go:8 | None |
| TestBankNodeIndexJSON | core/bank_nodes_index_test.go:21 | None |
| NewLoanProposalView | core/loanpool_views.go:20 | None |
| ProposalInfo | core/loanpool_views.go:36 | None |
| ProposalViews | core/loanpool_views.go:45 | None |
| NewLoanApplicationView | core/loanpool_views.go:67 | None |
| ApplicationInfo | core/loanpool_views.go:81 | None |
| ApplicationViews | core/loanpool_views.go:90 | None |
| NewLoanPoolApply | core/loanpool_apply.go:25 | None |
| Submit | core/loanpool_apply.go:34 | None |
| Vote | core/loanpool_apply.go:49 | None |
| Process | core/loanpool_apply.go:59 | None |
| Disburse | core/loanpool_apply.go:68 | None |
| Get | core/loanpool_apply.go:85 | None |
| List | core/loanpool_apply.go:91 | None |
| TestExperimentalnodePlaceholder | cli/experimental_node_test.go:5 | Rate limiting, validation, auth, integrity, user, crypto |
| InitGenesis | core/genesis_block.go:22 | None |
| NewImmutabilityEnforcer | core/immutability_enforcement.go:12 | None |
| CheckLedger | core/immutability_enforcement.go:22 | None |
| TestContractManager | core/contract_management_test.go:8 | None |
| init | cli/custodial_node.go:16 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/custodial_node.go:20 | Rate limiting, validation, auth, integrity, user, crypto |
| NewValidatorManager | core/consensus_validator_management.go:20 | None |
| Add | core/consensus_validator_management.go:29 | None |
| Remove | core/consensus_validator_management.go:44 | None |
| Slash | core/consensus_validator_management.go:55 | None |
| Eligible | core/consensus_validator_management.go:68 | None |
| Stake | core/consensus_validator_management.go:81 | None |
| TestLoanpoolviewsPlaceholder | core/loanpool_views_test.go:5 | None |
| init | cli/cross_chain_contracts.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSNVMArithmetic | core/snvm_test.go:5 | None |
| TestSNVMDivideByZero | core/snvm_test.go:21 | None |
| NewDAOStaking | core/dao_staking.go:21 | None |
| Stake | core/dao_staking.go:26 | None |
| Unstake | core/dao_staking.go:44 | None |
| Balance | core/dao_staking.go:63 | None |
| TotalStaked | core/dao_staking.go:70 | None |
| TestPeerManagerBasicOperations | core/peer_management_test.go:10 | None |
| TestPeerManagerConnect | core/peer_management_test.go:37 | None |
| TestPeerManagerAdvertiseAndDiscover | core/peer_management_test.go:49 | None |
| TestPeerManagerConcurrentAccess | core/peer_management_test.go:69 | None |
| contains | core/peer_management_test.go:97 | None |
| TestInstructionPlaceholder | core/instruction_test.go:5 | None |
| NewIPRegistry | core/syn700.go:42 | None |
| Register | core/syn700.go:47 | None |
| CreateLicense | core/syn700.go:57 | None |
| RecordRoyalty | core/syn700.go:70 | None |
| Get | core/syn700.go:83 | None |
| init | cli/syn2800.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSYN5000Token | core/syn5000_test.go:5 | None |
| TestBaseNodeLifecycle | core/base_node_test.go:12 | None |
| NewComplianceManager | core/compliance_management.go:18 | None |
| Suspend | core/compliance_management.go:27 | None |
| Resume | core/compliance_management.go:43 | None |
| Whitelist | core/compliance_management.go:56 | None |
| Unwhitelist | core/compliance_management.go:72 | None |
| Status | core/compliance_management.go:84 | None |
| ReviewTransaction | core/compliance_management.go:94 | None |
| IsZeroAddress | core/address_zero.go:8 | None |
| TestLightNodeHeaders | core/light_node_test.go:9 | None |
| NewSidechainRegistry | core/sidechains.go:25 | None |
| Register | core/sidechains.go:30 | None |
| SubmitHeader | core/sidechains.go:42 | None |
| GetHeader | core/sidechains.go:54 | None |
| Meta | core/sidechains.go:65 | None |
| List | core/sidechains.go:76 | None |
| Pause | core/sidechains.go:87 | None |
| Resume | core/sidechains.go:99 | None |
| UpdateValidators | core/sidechains.go:111 | None |
| Remove | core/sidechains.go:123 | None |
| NewSYN5000Token | core/syn5000.go:31 | None |
| PlaceBet | core/syn5000.go:36 | None |
| ResolveBet | core/syn5000.go:46 | None |
| GetBet | core/syn5000.go:65 | None |
| TestAuthorityApplication | core/authority_apply_test.go:11 | None |
| TestAuthorityApplicationJSON | core/authority_apply_test.go:32 | None |
| TestDAOTokenLedger | core/dao_token_test.go:5 | None |
| TestLiquidityPoolLifecycle | core/liquidity_pools_test.go:5 | None |
| TestWarfareNode | core/warfare_node_test.go:8 | None |
| TestSyn4700Placeholder | core/syn4700_test.go:5 | None |
| CompressLedger | core/blockchain_compression.go:21 | None |
| DecompressLedger | core/blockchain_compression.go:38 | None |
| SaveCompressedSnapshot | core/blockchain_compression.go:67 | None |
| LoadCompressedSnapshot | core/blockchain_compression.go:77 | None |
| TestFirewallPlaceholder | core/firewall_test.go:5 | None |
| NewBankInstitutionalNode | core/bank_institutional_node.go:19 | None |
| RegisterInstitution | core/bank_institutional_node.go:28 | None |
| RemoveInstitution | core/bank_institutional_node.go:47 | None |
| ListInstitutions | core/bank_institutional_node.go:62 | None |
| IsRegistered | core/bank_institutional_node.go:73 | None |
| MarshalJSON | core/bank_institutional_node.go:80 | None |
| NewTradeFinanceToken | core/syn2100.go:30 | None |
| RegisterDocument | core/syn2100.go:38 | None |
| FinanceDocument | core/syn2100.go:53 | None |
| GetDocument | core/syn2100.go:69 | None |
| ListDocuments | core/syn2100.go:81 | None |
| AddLiquidity | core/syn2100.go:93 | None |
| RemoveLiquidity | core/syn2100.go:100 | None |
| NewStorageMarketplace | core/storage_marketplace.go:42 | None |
| CreateListing | core/storage_marketplace.go:55 | None |
| ListListings | core/storage_marketplace.go:72 | None |
| OpenDeal | core/storage_marketplace.go:85 | None |
| CloseDeal | core/storage_marketplace.go:104 | None |
| GetListing | core/storage_marketplace.go:119 | None |
| GetDeal | core/storage_marketplace.go:129 | None |
| ListDeals | core/storage_marketplace.go:139 | None |
| MarshalListings | core/storage_marketplace.go:152 | None |
| MarshalDeals | core/storage_marketplace.go:157 | None |
| TestContractsOpcodes | core/contracts_opcodes_test.go:5 | None |
| NewSwarm | core/swarm.go:13 | None |
| Join | core/swarm.go:18 | None |
| Leave | core/swarm.go:25 | None |
| Members | core/swarm.go:32 | None |
| Broadcast | core/swarm.go:44 | None |
| Peers | core/swarm.go:51 | None |
| StartConsensus | core/swarm.go:63 | None |
| init | cli/token_syn4900.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| NewWatchtowerNode | core/watchtower_node.go:27 | None |
| ID | core/watchtower_node.go:37 | None |
| Start | core/watchtower_node.go:40 | None |
| monitorLoop | core/watchtower_node.go:54 | None |
| Stop | core/watchtower_node.go:71 | None |
| ReportFork | core/watchtower_node.go:83 | None |
| Metrics | core/watchtower_node.go:90 | None |
| Firewall | core/watchtower_node.go:95 | None |
| TestTokenInsurancePolicy | core/syn2900_test.go:8 | None |
| parseMap | cli/syn223_token.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/syn223_token.go:31 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn10InfoJSON | cli/syn10_test.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| NewComplianceService | core/compliance.go:46 | None |
| ValidateKYC | core/compliance.go:56 | None |
| EraseKYC | core/compliance.go:73 | None |
| RecordFraud | core/compliance.go:92 | None |
| RiskScore | core/compliance.go:113 | None |
| AuditTrail | core/compliance.go:122 | None |
| MonitorTransaction | core/compliance.go:134 | None |
| VerifyZKP | core/compliance.go:148 | None |
| appendAudit | core/compliance.go:155 | None |
| NewFaucet | core/faucet.go:19 | None |
| Request | core/faucet.go:29 | None |
| Balance | core/faucet.go:45 | None |
| UpdateConfig | core/faucet.go:52 | None |
| NewSupplyChainRegistry | core/syn1300.go:32 | None |
| Register | core/syn1300.go:37 | None |
| Update | core/syn1300.go:48 | None |
| Get | core/syn1300.go:60 | None |
| TestValidatorManagerLifecycle | core/consensus_validator_management_test.go:12 | None |
| TestStorageMarketplaceListEmpty | cli/storage_marketplace_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| NewWebRTCRPC | core/rpc_webrtc.go:13 | None |
| Connect | core/rpc_webrtc.go:19 | None |
| Send | core/rpc_webrtc.go:28 | None |
| Disconnect | core/rpc_webrtc.go:44 | None |
| NewPlasmaBridge | core/plasma.go:23 | None |
| NewSYN3700Token | core/syn3700_token.go:23 | None |
| AddComponent | core/syn3700_token.go:28 | None |
| RemoveComponent | core/syn3700_token.go:35 | None |
| ListComponents | core/syn3700_token.go:48 | None |
| Value | core/syn3700_token.go:57 | None |
| NewConsensusNetworkManager | core/cross_consensus_scaling_networks.go:24 | None |
| AuthorizeRelayer | core/cross_consensus_scaling_networks.go:29 | None |
| RevokeRelayer | core/cross_consensus_scaling_networks.go:36 | None |
| IsRelayerAuthorized | core/cross_consensus_scaling_networks.go:43 | None |
| RegisterNetwork | core/cross_consensus_scaling_networks.go:51 | None |
| RemoveNetwork | core/cross_consensus_scaling_networks.go:64 | None |
| ListNetworks | core/cross_consensus_scaling_networks.go:78 | None |
| GetNetwork | core/cross_consensus_scaling_networks.go:89 | None |
| TestComplianceMgmtStatusDefault | cli/compliance_mgmt_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| NewSNVM | core/snvm.go:13 | None |
| Execute | core/snvm.go:19 | None |
| TestSyn3600Lifecycle | cli/syn3600_test.go:8 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn3600Validation | cli/syn3600_test.go:30 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSimpleVM | core/virtual_machine_test.go:11 | None |
| TestVMVariants | core/virtual_machine_test.go:35 | None |
| TestVMContextCancel | core/virtual_machine_test.go:53 | None |
| TestVMCustomHandler | core/virtual_machine_test.go:65 | None |
| TestDifficultyManagerAdjustAndRetrieve | core/consensus_difficulty_test.go:7 | None |
| TestDifficultyManagerSlidingWindow | core/consensus_difficulty_test.go:24 | None |
| TestDifficultyManagerNilEngine | core/consensus_difficulty_test.go:42 | None |
| TestDifficultyManagerWindowFloor | core/consensus_difficulty_test.go:55 | None |
| Deposit | core/plasma_operations.go:7 | None |
| StartExit | core/plasma_operations.go:16 | None |
| FinalizeExit | core/plasma_operations.go:29 | None |
| GetExit | core/plasma_operations.go:41 | None |
| ListExits | core/plasma_operations.go:52 | None |
| TestBiometricsAuth | core/biometrics_auth_test.go:10 | None |
| TestValidatorNode | core/validator_node_test.go:5 | None |
| NewLoanPoolManager | core/loanpool_management.go:17 | None |
| Pause | core/loanpool_management.go:22 | None |
| Resume | core/loanpool_management.go:27 | None |
| Stats | core/loanpool_management.go:32 | None |
| init | cli/syn200.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| NewAIContractRegistry | core/ai_enhanced_contract.go:31 | None |
| DeployAIContract | core/ai_enhanced_contract.go:40 | None |
| InvokeAIContract | core/ai_enhanced_contract.go:60 | None |
| ModelHash | core/ai_enhanced_contract.go:72 | None |
| NewDAOTokenLedger | core/dao_token.go:18 | None |
| Mint | core/dao_token.go:23 | None |
| Transfer | core/dao_token.go:44 | None |
| Balance | core/dao_token.go:63 | None |
| Burn | core/dao_token.go:74 | None |
| TestPlasmaBridgeOperations | core/plasma_operations_test.go:5 | None |
| Start | core/audit_node_test.go:7 | None |
| TestAuditNode_StartLogAndVerify | core/audit_node_test.go:12 | None |
| TestSystemhealthloggingPlaceholder | core/system_health_logging_test.go:5 | None |
| NewGovernmentAuthorityNode | core/government_authority_node.go:13 | None |
| MintSYN | core/government_authority_node.go:20 | None |
| UpdateMonetaryPolicy | core/government_authority_node.go:26 | None |
| NewAuthorityNodeIndex | core/authority_node_index.go:15 | None |
| Add | core/authority_node_index.go:20 | None |
| Get | core/authority_node_index.go:30 | None |
| Remove | core/authority_node_index.go:38 | None |
| List | core/authority_node_index.go:45 | None |
| Snapshot | core/authority_node_index.go:56 | None |
| MarshalJSON | core/authority_node_index.go:67 | None |
| TestRegNodeCLI | cli/regulatory_node_test.go:8 | Rate limiting, validation, auth, integrity, user, crypto |
| TestThreshold | core/consensus_test.go:8 | None |
| TestAdjustWeightsAndAvailability | core/consensus_test.go:15 | None |
| TestTransitionThreshold | core/consensus_test.go:28 | None |
| TestDifficultyAdjust | core/consensus_test.go:37 | None |
| TestSelectValidator | core/consensus_test.go:44 | None |
| TestSelectValidatorMajorityStake | core/consensus_test.go:56 | None |
| TestValidateSubBlock | core/consensus_test.go:64 | None |
| TestValidateBlock | core/consensus_test.go:77 | None |
| NewFuturesContract | core/syn3600.go:19 | None |
| IsExpired | core/syn3600.go:24 | None |
| Settle | core/syn3600.go:29 | None |
| TestTokensyn130Placeholder | core/token_syn130_test.go:5 | None |
| NewBaseNode | core/base_node.go:21 | None |
| ID | core/base_node.go:29 | None |
| Start | core/base_node.go:32 | None |
| Stop | core/base_node.go:43 | None |
| IsRunning | core/base_node.go:54 | None |
| Peers | core/base_node.go:61 | None |
| DialSeed | core/base_node.go:72 | None |
| DialSeedSigned | core/base_node.go:84 | None |
| NewShardManager | core/sharding.go:17 | None |
| GetLeader | core/sharding.go:27 | None |
| SetLeader | core/sharding.go:35 | None |
| LeaderMap | core/sharding.go:42 | None |
| SubmitCrossShardTx | core/sharding.go:53 | None |
| PullReceipts | core/sharding.go:61 | None |
| Reshard | core/sharding.go:71 | None |
| Rebalance | core/sharding.go:78 | None |
| ShardCount | core/sharding.go:91 | None |
| TestLoanpoolproposalPlaceholder | core/loanpool_proposal_test.go:5 | None |
| TestNewTransactionAndHash | core/transaction_test.go:10 | None |
| TestAttachBiometric | core/transaction_test.go:25 | None |
| gasListPath | core/gas_table_test.go:13 | None |
| TestParseGasGuide | core/gas_table_test.go:24 | None |
| TestDefaultGasTableOverrides | core/gas_table_test.go:45 | None |
| TestSetGasCostAndSnapshot | core/gas_table_test.go:82 | None |
| TestAccessControlGasCosts | core/gas_table_test.go:104 | None |
| TestGasCostByName | core/gas_table_test.go:122 | None |
| TestGasTableSnapshotJSONDeterministic | core/gas_table_test.go:132 | None |
| TestWriteGasTableSnapshot | core/gas_table_test.go:154 | None |
| NewDifficultyManager | core/consensus_difficulty.go:17 | None |
| AddSample | core/consensus_difficulty.go:26 | None |
| Difficulty | core/consensus_difficulty.go:45 | None |
| NewLiquidityPoolView | core/liquidity_views.go:14 | None |
| PoolInfo | core/liquidity_views.go:26 | None |
| PoolViews | core/liquidity_views.go:35 | None |
| Liquidity_Pool | core/liquidity_views.go:45 | None |
| Liquidity_Pools | core/liquidity_views.go:50 | None |
| TestSyncManagerLifecycle | core/blockchain_synchronization_test.go:5 | None |
| NewRegulatoryManager | core/regulatory_management.go:23 | None |
| AddRegulation | core/regulatory_management.go:28 | None |
| RemoveRegulation | core/regulatory_management.go:39 | None |
| GetRegulation | core/regulatory_management.go:46 | None |
| ListRegulations | core/regulatory_management.go:54 | None |
| EvaluateTransaction | core/regulatory_management.go:65 | None |
| init | cli/gas_table.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| TestLoanPoolProposalLifecycle | core/loanpool_test.go:5 | None |
| TestSyn3200Placeholder | core/syn3200_test.go:5 | None |
| IsIDTokenHolder | core/charity_test.go:12 | None |
| newTestState | core/charity_test.go:19 | None |
| Transfer | core/charity_test.go:23 | None |
| SetState | core/charity_test.go:28 | None |
| GetState | core/charity_test.go:29 | None |
| HasState | core/charity_test.go:30 | None |
| PrefixIterator | core/charity_test.go:31 | None |
| BalanceOf | core/charity_test.go:32 | None |
| TestCharityPool | core/charity_test.go:34 | None |
| StringToAddress | core/address.go:16 | None |
| Hex | core/address.go:28 | None |
| Bytes | core/address.go:32 | None |
| Short | core/address.go:42 | None |
| String | core/address.go:51 | None |
| IsZero | core/address.go:54 | None |
| NewPeerManager | core/peer_management.go:14 | None |
| AddPeer | core/peer_management.go:22 | None |
| RemovePeer | core/peer_management.go:29 | None |
| GetPeer | core/peer_management.go:36 | None |
| ListPeers | core/peer_management.go:44 | None |
| Connect | core/peer_management.go:57 | None |
| Advertise | core/peer_management.go:64 | None |
| Discover | core/peer_management.go:71 | None |
| bsnOutput | cli/biometric_security_node.go:22 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/biometric_security_node.go:33 | Rate limiting, validation, auth, integrity, user, crypto |
| NewValidatorNode | core/validator_node.go:14 | None |
| AddValidator | core/validator_node.go:24 | None |
| RemoveValidator | core/validator_node.go:33 | None |
| SlashValidator | core/validator_node.go:39 | None |
| HasQuorum | core/validator_node.go:46 | None |
| TestTradeFinanceToken | core/syn2100_test.go:8 | None |
| parseOpcode | cli/instruction.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/instruction.go:24 | Rate limiting, validation, auth, integrity, user, crypto |
| NewFailoverManager | core/high_availability.go:20 | None |
| RegisterBackup | core/high_availability.go:29 | None |
| Heartbeat | core/high_availability.go:36 | None |
| Active | core/high_availability.go:45 | None |
| TestNodeCLICommands | cli/node_test.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| execNodeCLI | cli/node_test.go:60 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/node_adapter.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| Call | core/virtual_machine.go:43 | None |
| Gas | core/virtual_machine.go:47 | None |
| NewSimpleVM | core/virtual_machine.go:52 | None |
| RegisterHandler | core/virtual_machine.go:87 | None |
| Concurrency | core/virtual_machine.go:99 | None |
| Start | core/virtual_machine.go:102 | None |
| Stop | core/virtual_machine.go:113 | None |
| Status | core/virtual_machine.go:124 | None |
| ExecuteContext | core/virtual_machine.go:131 | None |
| Execute | core/virtual_machine.go:200 | None |
| init | cli/loanpool.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSandboxManager | core/vm_sandbox_management_test.go:8 | None |
| NewContractRegistry | core/contracts.go:39 | None |
| CompileWASM | core/contracts.go:49 | None |
| Deploy | core/contracts.go:60 | None |
| Invoke | core/contracts.go:84 | None |
| List | core/contracts.go:101 | None |
| Get | core/contracts.go:112 | None |
| TestLoanpoolmanagementPlaceholder | core/loanpool_management_test.go:5 | None |
| TestPlasmamanagementPlaceholder | cli/plasma_management_test.go:5 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSYN131RegistryCreateGet | core/syn131_token_test.go:5 | None |
| TestSYN131RegistryDuplicate | core/syn131_token_test.go:20 | None |
| TestSYN131RegistryUpdateValuation | core/syn131_token_test.go:30 | None |
| TestSYN131RegistryUpdateNonexistent | core/syn131_token_test.go:44 | None |
| init | cli/syn500.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| NewWallet | core/wallet.go:27 | None |
| Sign | core/wallet.go:45 | None |
| VerifySignature | core/wallet.go:61 | None |
| Save | core/wallet.go:73 | None |
| LoadWallet | core/wallet.go:112 | None |
| TestSyn3500Lifecycle | cli/syn3500_token_test.go:7 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn3500Validation | cli/syn3500_token_test.go:37 | Rate limiting, validation, auth, integrity, user, crypto |
| TestDefaultGenesisWalletsDeterministic | core/genesis_wallets_test.go:5 | None |
| TestAllocateToGenesisWallets | core/genesis_wallets_test.go:16 | None |
| init | cli/syn3200.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| TestFeeForTransfer | core/fees_test.go:5 | None |
| TestDistributeFees | core/fees_test.go:12 | None |
| TestApplyFeeCapFloor | core/fees_test.go:23 | None |
| TestFeePolicyEnforce | core/fees_test.go:32 | None |
| TestAdjustFeeRates | core/fees_test.go:45 | None |
| TestEstimateFee | core/fees_test.go:52 | None |
| TestShareProportional | core/fees_test.go:59 | None |
| TestAdjustForBlockUtilization | core/fees_test.go:67 | None |
| TestTxControlSchedule | cli/tx_control_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestBankNodeIndexAdd | cli/bank_node_index_test.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| TestStaterwPlaceholder | core/state_rw_test.go:5 | None |
| NewSYN300Token | core/syn300_token.go:30 | None |
| Delegate | core/syn300_token.go:44 | None |
| RevokeDelegation | core/syn300_token.go:55 | None |
| VotingPower | core/syn300_token.go:62 | None |
| votingPowerLocked | core/syn300_token.go:68 | None |
| CreateProposal | core/syn300_token.go:79 | None |
| Vote | core/syn300_token.go:96 | None |
| Execute | core/syn300_token.go:117 | None |
| ProposalStatus | core/syn300_token.go:139 | None |
| ListProposals | core/syn300_token.go:160 | None |
| init | cli/syn845.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| NewRegulatoryNode | core/regulatory_node.go:17 | None |
| ApproveTransaction | core/regulatory_node.go:26 | None |
| FlagEntity | core/regulatory_node.go:36 | None |
| Logs | core/regulatory_node.go:43 | None |
| ScheduleTransaction | core/transaction_control.go:24 | None |
| CancelTransaction | core/transaction_control.go:30 | None |
| ReverseTransaction | core/transaction_control.go:40 | None |
| RequestReversal | core/transaction_control.go:62 | None |
| Vote | core/transaction_control.go:75 | None |
| FinalizeReversal | core/transaction_control.go:84 | None |
| RejectReversal | core/transaction_control.go:108 | None |
| ConvertToPrivate | core/transaction_control.go:117 | None |
| Decrypt | core/transaction_control.go:140 | None |
| GenerateReceipt | core/transaction_control.go:169 | None |
| NewReceiptStore | core/transaction_control.go:180 | None |
| Store | core/transaction_control.go:185 | None |
| Get | core/transaction_control.go:192 | None |
| Search | core/transaction_control.go:200 | None |
| TestSyn800tokenPlaceholder | core/syn800_token_test.go:5 | None |
| NewBenefitRegistry | core/syn3900.go:25 | None |
| RegisterBenefit | core/syn3900.go:30 | None |
| Claim | core/syn3900.go:40 | None |
| GetBenefit | core/syn3900.go:55 | None |
| NewBankNodeIndex | core/bank_nodes_index.go:35 | None |
| Add | core/bank_nodes_index.go:40 | None |
| Get | core/bank_nodes_index.go:47 | None |
| Remove | core/bank_nodes_index.go:55 | None |
| List | core/bank_nodes_index.go:62 | None |
| Snapshot | core/bank_nodes_index.go:73 | None |
| MarshalJSON | core/bank_nodes_index.go:84 | None |
| init | cli/syn131_token.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/network.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn700Workflow | cli/syn700_test.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| TestGasSnapshotJSON | cli/gas_table_cli_test.go:16 | Rate limiting, validation, auth, integrity, user, crypto |
| TestGasSnapshotWrite | cli/gas_table_cli_test.go:37 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSmartContractMarketplaceDeployAndTrade | core/smart_contract_marketplace_test.go:12 | None |
| init | cli/syn3500_token.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| NewProposalManager | core/dao_proposal.go:29 | None |
| CreateProposal | core/dao_proposal.go:41 | None |
| Vote | core/dao_proposal.go:59 | None |
| Results | core/dao_proposal.go:85 | None |
| Execute | core/dao_proposal.go:104 | None |
| Get | core/dao_proposal.go:124 | None |
| List | core/dao_proposal.go:135 | None |
| NewAuditNode | core/audit_node.go:20 | None |
| Start | core/audit_node.go:25 | None |
| LogEvent | core/audit_node.go:33 | None |
| ListEvents | core/audit_node.go:44 | None |
| VerifyEvent | core/audit_node.go:52 | None |
| TestAuditManager_LogAndList | core/audit_management_test.go:5 | None |
| TestAuditManager_Invalid | core/audit_management_test.go:29 | None |
| NewBiometricsAuth | core/biometrics_auth.go:27 | None |
| Enroll | core/biometrics_auth.go:33 | None |
| Verify | core/biometrics_auth.go:53 | None |
| Remove | core/biometrics_auth.go:68 | None |
| Enrolled | core/biometrics_auth.go:75 | None |
| List | core/biometrics_auth.go:84 | None |
| init | cli/syn70.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| TestFaucet | core/faucet_test.go:8 | None |
| TestVestingSchedule | core/syn2700_test.go:8 | None |
| TestAuthorityNodeIndex | core/authority_node_index_test.go:5 | None |
| TestAuthorityNodeIndexJSON | core/authority_node_index_test.go:18 | None |
| cmOutput | cli/compliance_mgmt.go:17 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/compliance_mgmt.go:28 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn2500tokenPlaceholder | core/syn2500_token_test.go:5 | None |
| TestScheduleAndCancel | core/transaction_control_test.go:8 | None |
| TestReverseTransaction | core/transaction_control_test.go:20 | None |
| TestAuthorityMediatedReversal | core/transaction_control_test.go:35 | None |
| TestReversalRejection | core/transaction_control_test.go:57 | None |
| TestConvertToPrivate | core/transaction_control_test.go:79 | None |
| TestReceiptStore | core/transaction_control_test.go:95 | None |
| TestReplicationCLI | cli/replication_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestElectedAuthorityNode | core/elected_authority_node_test.go:8 | None |
| parseMeta | cli/syn2500_token.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/syn2500_token.go:28 | Rate limiting, validation, auth, integrity, user, crypto |
| TestNFTMarketplace | core/nft_marketplace_test.go:11 | None |
| TestNFTMarketplaceDuplicate | core/nft_marketplace_test.go:35 | None |
| BenchmarkNFTMarketplaceMint | core/nft_marketplace_test.go:45 | None |
| TestBankNodes | core/bank_nodes_test.go:11 | None |
| init | cli/syn3800.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| NewIDRegistry | core/idwallet_registration.go:15 | None |
| Register | core/idwallet_registration.go:20 | None |
| Info | core/idwallet_registration.go:31 | None |
| TestAddressZeroHelpers | cli/address_zero_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| newTestLiquidityPool | core/liquidity_views_test.go:6 | None |
| TestNewLiquidityPoolViewSnapshot | core/liquidity_views_test.go:15 | None |
| TestPoolInfo | core/liquidity_views_test.go:31 | None |
| TestPoolViews | core/liquidity_views_test.go:49 | None |
| TestBiometricVerifyNoEnrollment | cli/biometric_test.go:6 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSidechainRegistry | core/sidechains_test.go:5 | None |
| TestSyn1300RegisterRequiresFlags | cli/syn1300_test.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn1300Workflow | cli/syn1300_test.go:22 | Rate limiting, validation, auth, integrity, user, crypto |
| TestBridgeManager | core/cross_chain_bridge_test.go:5 | None |
| TestInvestmentRedeem | core/syn1401_test.go:8 | None |
| init | cli/initrep.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| TestVerifySignature | cli/ecdsa_util_test.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| msgHash | cli/ecdsa_util_test.go:48 | Rate limiting, validation, auth, integrity, user, crypto |
| TestRpcwebrtcPlaceholder | cli/rpc_webrtc_test.go:5 | Rate limiting, validation, auth, integrity, user, crypto |
| TestNewWatchtowerNode | core/watchtower_node_test.go:13 | None |
| TestWatchtowerStartStop | core/watchtower_node_test.go:34 | None |
| TestWatchtowerMetrics | core/watchtower_node_test.go:58 | None |
| TestWatchtowerReportFork | core/watchtower_node_test.go:70 | None |
| init | cli/dao.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn3400Lifecycle | cli/syn3400_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn3400InvalidRegister | cli/syn3400_test.go:32 | Rate limiting, validation, auth, integrity, user, crypto |
| TestStringToAddress | core/address_test.go:9 | None |
| TestAddressHexBytesAndString | core/address_test.go:62 | None |
| TestAddressBytesInvalid | core/address_test.go:83 | None |
| TestAddressShort | core/address_test.go:90 | None |
| CalculateBaseFee | core/fees.go:33 | None |
| CalculateVariableFee | core/fees.go:49 | None |
| CalculatePriorityFee | core/fees.go:54 | None |
| FeeForTransfer | core/fees.go:57 | None |
| FeeForPurchase | core/fees.go:64 | None |
| FeeForTokenUsage | core/fees.go:71 | None |
| FeeForContract | core/fees.go:78 | None |
| FeeForWalletVerification | core/fees.go:85 | None |
| FeeForValidatedTransfer | core/fees.go:94 | None |
| DistributeFees | core/fees.go:116 | None |
| ApplyFeeCapFloor | core/fees.go:131 | None |
| Enforce | core/fees.go:149 | None |
| AdjustFeeRates | core/fees.go:163 | None |
| EstimateFee | core/fees.go:176 | None |
| ShareProportional | core/fees.go:196 | None |
| NewFeeDistributionContract | core/fees.go:229 | None |
| Distribute | core/fees.go:234 | None |
| AdjustForBlockUtilization | core/fees.go:242 | None |
| TestPrivatetransactionsPlaceholder | cli/private_transactions_test.go:5 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/system_health_logging.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn300tokenPlaceholder | core/syn300_token_test.go:5 | None |
| TestBaseTokenMint | cli/base_token_test.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| NewCrossChainTxManager | core/cross_chain_transactions.go:40 | None |
| AuthorizeRelayer | core/cross_chain_transactions.go:45 | None |
| RevokeRelayer | core/cross_chain_transactions.go:52 | None |
| IsRelayerAuthorized | core/cross_chain_transactions.go:59 | None |
| LockMint | core/cross_chain_transactions.go:66 | None |
| BurnRelease | core/cross_chain_transactions.go:87 | None |
| ListTransfers | core/cross_chain_transactions.go:108 | None |
| GetTransfer | core/cross_chain_transactions.go:119 | None |
| TestBridgeRegistry | core/cross_chain_test.go:5 | None |
| NewTangibleAssetRegistry | core/token_syn130.go:40 | None |
| Register | core/token_syn130.go:45 | None |
| UpdateValuation | core/token_syn130.go:55 | None |
| RecordSale | core/token_syn130.go:65 | None |
| StartLease | core/token_syn130.go:76 | None |
| EndLease | core/token_syn130.go:86 | None |
| Get | core/token_syn130.go:98 | None |
| NewMusicToken | core/syn1600.go:19 | None |
| Info | core/syn1600.go:29 | None |
| Update | core/syn1600.go:36 | None |
| SetRoyaltyShare | core/syn1600.go:51 | None |
| Distribute | core/syn1600.go:60 | None |
| TestCentralBankingMintCBDC | core/central_banking_node_test.go:9 | None |
| TestConsensusServiceStartStop | core/consensus_start_test.go:10 | None |
| NewIdentityService | core/identity_verification.go:30 | None |
| Register | core/identity_verification.go:38 | None |
| Verify | core/identity_verification.go:49 | None |
| Info | core/identity_verification.go:61 | None |
| Logs | core/identity_verification.go:69 | None |
| TestPlasmaPlaceholder | core/plasma_test.go:5 | None |
| NewTokenInsurancePolicy | core/syn2900.go:23 | None |
| IsActive | core/syn2900.go:38 | None |
| Claim | core/syn2900.go:43 | None |
| hashAddress | core/genesis_wallets.go:23 | None |
| DefaultGenesisWallets | core/genesis_wallets.go:29 | None |
| AllocateToGenesisWallets | core/genesis_wallets.go:46 | None |
| NewStakingNode | core/staking_node.go:12 | None |
| Stake | core/staking_node.go:17 | None |
| Unstake | core/staking_node.go:25 | None |
| Balance | core/staking_node.go:37 | None |
| TotalStaked | core/staking_node.go:44 | None |
| NewQuorumTracker | core/quorum_tracker.go:14 | None |
| Join | core/quorum_tracker.go:22 | None |
| Leave | core/quorum_tracker.go:29 | None |
| Count | core/quorum_tracker.go:36 | None |
| Reached | core/quorum_tracker.go:43 | None |
| TestNodeAdapterInfo | cli/node_adapter_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| NewLegalToken | core/syn4700.go:44 | None |
| Sign | core/syn4700.go:63 | None |
| RevokeSignature | core/syn4700.go:74 | None |
| UpdateStatus | core/syn4700.go:81 | None |
| Dispute | core/syn4700.go:88 | None |
| partyExists | core/syn4700.go:95 | None |
| NewLegalTokenRegistry | core/syn4700.go:111 | None |
| Add | core/syn4700.go:116 | None |
| Get | core/syn4700.go:123 | None |
| Remove | core/syn4700.go:131 | None |
| List | core/syn4700.go:138 | None |
| TestMusicToken | core/syn1600_test.go:5 | None |
| TestWalletSignAndVerify | core/wallet_test.go:8 | None |
| TestWalletSaveLoad | core/wallet_test.go:39 | None |
| TestAuthorityNodeRegistry | core/authority_nodes_test.go:9 | None |
| TestAuthorityNodeJSONAndRemoveVote | core/authority_nodes_test.go:33 | None |
| TestSyn1700InitRequiresFlags | cli/syn1700_token_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn1700Workflow | cli/syn1700_token_test.go:20 | Rate limiting, validation, auth, integrity, user, crypto |
| NewConnectionPool | core/connection_pool.go:26 | None |
| Acquire | core/connection_pool.go:35 | None |
| Release | core/connection_pool.go:65 | None |
| Size | core/connection_pool.go:79 | None |
| Dial | core/connection_pool.go:90 | None |
| Close | core/connection_pool.go:99 | None |
| Stats | core/connection_pool.go:119 | None |
| init | cli/node.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSnvmopcodesPlaceholder | core/snvm_opcodes_test.go:5 | None |
| NewSandboxManager | core/vm_sandbox_management.go:32 | None |
| StartSandbox | core/vm_sandbox_management.go:41 | None |
| StopSandbox | core/vm_sandbox_management.go:61 | None |
| ResetSandbox | core/vm_sandbox_management.go:73 | None |
| DeleteSandbox | core/vm_sandbox_management.go:85 | None |
| PurgeInactive | core/vm_sandbox_management.go:98 | None |
| SandboxStatus | core/vm_sandbox_management.go:113 | None |
| ListSandboxes | core/vm_sandbox_management.go:121 | None |
| init | cli/dao_access_control.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestConsensusSwitcher | core/consensus_specific_test.go:5 | None |
| TestGeospatialRecordHistory | cli/geospatial_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSetStakeEnforcesMinimum | core/security_test.go:5 | None |
| TestSlashingAndRehabilitation | core/security_test.go:16 | None |
| TestEligibleStakesExcludesSlashed | core/security_test.go:33 | None |
| TestSubBlockSignature | core/security_test.go:48 | None |
| init | cli/contract_management.go:17 | Rate limiting, validation, auth, integrity, user, crypto |
| printErr | cli/contract_management.go:109 | Rate limiting, validation, auth, integrity, user, crypto |
| NewSYN4200Token | core/syn4200_token.go:23 | None |
| Donate | core/syn4200_token.go:28 | None |
| CampaignProgress | core/syn4200_token.go:45 | None |
| Campaign | core/syn4200_token.go:56 | None |
| TestDAOStaking | core/dao_staking_test.go:7 | None |
| init | cli/dao_quadratic_voting.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn3700tokenPlaceholder | core/syn3700_token_test.go:5 | None |
| TestBlockRewardHalving | core/coin_test.go:8 | None |
| TestCirculatingAndRemainingSupply | core/coin_test.go:20 | None |
| TestCirculatingSupplyCapped | core/coin_test.go:31 | None |
| TestEconomicHelpers | core/coin_test.go:37 | None |
| TestKademliaStoreFind | core/kademlia_test.go:8 | None |
| TestDistance | core/kademlia_test.go:17 | None |
| NewAssetRegistry | core/syn800_token.go:25 | None |
| Register | core/syn800_token.go:30 | None |
| UpdateValuation | core/syn800_token.go:40 | None |
| Get | core/syn800_token.go:51 | None |
| NewConsensusSpecificNode | core/consensus_specific_node.go:10 | None |
| configure | core/consensus_specific_node.go:19 | None |
| TestZeroTrustEngine | core/zero_trust_data_channels_test.go:5 | None |
| NewLoanProposal | core/loanpool_proposal.go:20 | None |
| Vote | core/loanpool_proposal.go:34 | None |
| VoteCount | core/loanpool_proposal.go:39 | None |
| IsExpired | core/loanpool_proposal.go:44 | None |
| TestConsensusSpecificNodeCreateInfo | cli/consensus_specific_node_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| TestHolographicNodeWorkflow | cli/holographic_node_test.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/block.go:16 | Rate limiting, validation, auth, integrity, user, crypto |
| TestStakePenaltySlashJSON | cli/stake_penalty_test.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/cross_chain.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/forensic_node.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/genesis.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/mobile_mining_node.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/geospatial.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/syn10.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| execCommand | cli/cli_core_test.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| TestAddressParse | cli/cli_core_test.go:31 | Rate limiting, validation, auth, integrity, user, crypto |
| TestRootHelp | cli/cli_core_test.go:42 | Rate limiting, validation, auth, integrity, user, crypto |
| TestNetworkStartStop | cli/cli_core_test.go:52 | Rate limiting, validation, auth, integrity, user, crypto |
| TestPeerDiscoverEmpty | cli/cli_core_test.go:62 | Rate limiting, validation, auth, integrity, user, crypto |
| TestConsensusMineGas | cli/cli_core_test.go:74 | Rate limiting, validation, auth, integrity, user, crypto |
| TestDAOCreationGas | cli/cli_core_test.go:84 | Rate limiting, validation, auth, integrity, user, crypto |
| TestRollupsCLI | cli/rollups_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestShardingCLI | cli/sharding_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/cross_chain_agnostic_protocols.go:15 | Rate limiting, validation, auth, integrity, user, crypto |
| resetSyn4900 | cli/token_syn4900_test.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| TestTokenSyn4900RegisterInfo | cli/token_syn4900_test.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/cross_chain_bridge.go:15 | Rate limiting, validation, auth, integrity, user, crypto |
| TestBankInstitutionalRegister | cli/bank_institutional_node_test.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/transaction.go:22 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn5000Index | cli/syn5000_index_test.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| TestComplianceRiskDefault | cli/compliance_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn3800Lifecycle | cli/syn3800_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn3800Validation | cli/syn3800_test.go:46 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/ledger.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| TestGenesisShow | cli/genesis_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| TestGenesisAllocate | cli/genesis_test.go:23 | Rate limiting, validation, auth, integrity, user, crypto |
| TestTokenregistryPlaceholder | cli/token_registry_test.go:5 | Rate limiting, validation, auth, integrity, user, crypto |
| TestCrossChainCLI | cli/cross_chain_cli_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| TestPlasmaMgmtCLI | cli/cross_chain_cli_test.go:26 | Rate limiting, validation, auth, integrity, user, crypto |
| TestBridgeDepositCLI | cli/cross_chain_cli_test.go:42 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/dao_staking.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| NewSmartContractMarketplace | core/smart_contract_marketplace.go:24 | None |
| DeployContract | core/smart_contract_marketplace.go:36 | None |
| TradeContract | core/smart_contract_marketplace.go:50 | None |
| init | cli/light_node.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| NewWarfareNode | core/warfare_node.go:19 | None |
| GetID | core/warfare_node.go:24 | None |
| SecureCommand | core/warfare_node.go:33 | None |
| TrackLogistics | core/warfare_node.go:41 | None |
| ShareTactical | core/warfare_node.go:55 | None |
| Logistics | core/warfare_node.go:61 | None |
| LogisticsByAsset | core/warfare_node.go:71 | None |
| init | cli/mining_node.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| TestIdentityWorkflow | cli/identity_test.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSubBlockCreationAndVerification | core/block_test.go:12 | None |
| TestBlockHeaderHash | core/block_test.go:20 | None |
| TestSubBlockValidateRejectsDuplicateTransactions | core/block_test.go:41 | None |
| TestBlockValidateRejectsFutureTimestamp | core/block_test.go:49 | None |
| TestBlockValidateRejectsSubBlockAfterBlockTimestamp | core/block_test.go:61 | None |
| TestBlockValidateRequiresHash | core/block_test.go:73 | None |
| TestBlockValidateDetectsHashMismatch | core/block_test.go:83 | None |
| init | cli/bank_node_index.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| TestCrossChainTxManager | core/cross_chain_transactions_test.go:5 | None |
| TestLoanpoolApplySubmitList | cli/loanpool_apply_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| NewVestingSchedule | core/syn2700.go:18 | None |
| Claim | core/syn2700.go:23 | None |
| Pending | core/syn2700.go:36 | None |
| init | cli/private_transactions.go:15 | Rate limiting, validation, auth, integrity, user, crypto |
| NewNFTMarketplace | core/nft_marketplace.go:27 | None |
| Mint | core/nft_marketplace.go:33 | None |
| List | core/nft_marketplace.go:53 | None |
| Buy | core/nft_marketplace.go:65 | None |
| ListAll | core/nft_marketplace.go:85 | None |
| init | cli/syn12.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn1000IndexValueZeroJSON | cli/syn1000_index_test.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| TestNetworkBroadcast | core/network_test.go:13 | None |
| TestAuthorityApplySubmit | cli/authority_apply_test.go:16 | Rate limiting, validation, auth, integrity, user, crypto |
| TestGatewayNodeEndpoints | core/gateway_node_test.go:9 | None |
| TestLiquidityPoolsCreateList | cli/liquidity_pools_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| NewConsensusService | core/consensus_start.go:19 | None |
| Start | core/consensus_start.go:25 | None |
| Stop | core/consensus_start.go:48 | None |
| Info | core/consensus_start.go:56 | None |
| TestSyn3900Lifecycle | cli/syn3900_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn3900Validation | cli/syn3900_test.go:35 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/syn1000_index.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| NewRollupAggregator | core/rollups.go:25 | None |
| SubmitBatch | core/rollups.go:30 | None |
| ChallengeBatch | core/rollups.go:43 | None |
| FinalizeBatch | core/rollups.go:59 | None |
| BatchInfo | core/rollups.go:75 | None |
| ListBatches | core/rollups.go:86 | None |
| BatchTransactions | core/rollups.go:97 | None |
| Pause | core/rollups.go:108 | None |
| Resume | core/rollups.go:115 | None |
| Status | core/rollups.go:122 | None |
| init | cli/plasma_management.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn4700Lifecycle | cli/syn4700_test.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn4700Validation | cli/syn4700_test.go:46 | Rate limiting, validation, auth, integrity, user, crypto |
| TestFirewallAllowCheck | cli/firewall_test.go:8 | Rate limiting, validation, auth, integrity, user, crypto |
| TestLedgerMintBalance | cli/ledger_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/centralbank.go:18 | Rate limiting, validation, auth, integrity, user, crypto |
| TestNATCLIMapAndUnmap | cli/nat_test.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| TestNATCLIInvalidPort | cli/nat_test.go:34 | Rate limiting, validation, auth, integrity, user, crypto |
| TestContractsList | cli/contracts_test.go:5 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/audit.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/optimization_node.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| TestConnectionManagerJSON | cli/cross_chain_connection_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| IsIDTokenHolder | cli/charity.go:18 | Rate limiting, validation, auth, integrity, user, crypto |
| newMemState | cli/charity.go:25 | Rate limiting, validation, auth, integrity, user, crypto |
| Transfer | cli/charity.go:29 | Rate limiting, validation, auth, integrity, user, crypto |
| SetState | cli/charity.go:35 | Rate limiting, validation, auth, integrity, user, crypto |
| GetState | cli/charity.go:36 | Rate limiting, validation, auth, integrity, user, crypto |
| HasState | cli/charity.go:37 | Rate limiting, validation, auth, integrity, user, crypto |
| BalanceOf | cli/charity.go:38 | Rate limiting, validation, auth, integrity, user, crypto |
| PrefixIterator | cli/charity.go:46 | Rate limiting, validation, auth, integrity, user, crypto |
| Next | cli/charity.go:57 | Rate limiting, validation, auth, integrity, user, crypto |
| Value | cli/charity.go:62 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/charity.go:68 | Rate limiting, validation, auth, integrity, user, crypto |
| TestDAOCLIFlow | cli/dao_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/historical.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn1100AddRequiresFlags | cli/syn1100_test.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn1100Workflow | cli/syn1100_test.go:23 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/consensus_adaptive_management.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/sidechain_ops.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestForensicRecordTx | cli/forensic_node_test.go:8 | Rate limiting, validation, auth, integrity, user, crypto |
| TestFullNodeLifecycle | cli/full_node_test.go:8 | Rate limiting, validation, auth, integrity, user, crypto |
| TestOutputPlaceholder | cli/output_test.go:5 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/syn20.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSidechainopsPlaceholder | cli/sidechain_ops_test.go:5 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/firewall.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/gateway.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| TestRootCommand | cli/root_test.go:7 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/cross_chain_connection.go:15 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/authority_nodes.go:16 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/experimental_node.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| TestStakingNodeTotalZero | cli/staking_node_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| TestGasPrint | cli/gas_print_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| parseAttrs | cli/syn2369.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/syn2369.go:28 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/elected_authority_node.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| TestBSNAuthNoEnrollment | cli/biometric_security_node_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| run223 | cli/syn223_token_test.go:7 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn223InitTransfer | cli/syn223_token_test.go:16 | Rate limiting, validation, auth, integrity, user, crypto |
| TestAIContractModelHashNotFound | cli/ai_contract_test.go:6 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/rollup_management.go:5 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/stake_penalty.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| TestRollupManagerCLI | cli/rollup_management_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn200RegisterMissingFlags | cli/syn200_test.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn200RegisterWorkflow | cli/syn200_test.go:23 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/access.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestAuditLogAndList | cli/audit_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| TestAccessListEmpty | cli/access_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/contracts_opcodes.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/zero_trust_data_channels.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| TestImmutabilityWorkflow | cli/immutability_test.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| TestLiquidityViewsList | cli/liquidity_views_cli_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/coin.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| coinOutput | cli/coin.go:151 | Rate limiting, validation, auth, integrity, user, crypto |
| TestAIServiceEndToEnd | ai_modules_test.go:8 | None |
| TestModelMarketplace | ai_modules_test.go:56 | None |
| TestTrainingManager | ai_modules_test.go:71 | None |
| TestInferenceEngine | ai_modules_test.go:92 | None |
| TestDriftMonitor | ai_modules_test.go:108 | None |
| TestSecureStorage | ai_modules_test.go:119 | None |
| TestAnomalyDetector | ai_modules_test.go:131 | None |
| TestSyn1000ValueZeroJSON | cli/syn1000_test.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| TestHistoricalWorkflow | cli/historical_test.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/syn2600.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| TestGovernmentNew | cli/government_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestAIContractListEmpty | cli/ai_contract_cli_test.go:7 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/base_token.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/immutability.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn5000Workflow | cli/syn5000_test.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/syn1401.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| TestConsensusNetworkJSON | cli/cross_consensus_scaling_networks_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/address_zero.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestCentralbankInfo | cli/centralbank_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| TestContentNodeCLI | cli/content_node_mgmt_test.go:8 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/output.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| printOutput | cli/output.go:19 | Rate limiting, validation, auth, integrity, user, crypto |
| TestBankNodeTypes | cli/bank_nodes_index_test.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| captureOutput | cli/bank_nodes_index_test.go:25 | Rate limiting, validation, auth, integrity, user, crypto |
| TestLightNodeHeader | cli/light_node_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/nat.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/syn5000.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/loanpool_apply.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/gas.go:18 | Rate limiting, validation, auth, integrity, user, crypto |
| TestWarfareNodeLogistics | cli/warfare_node_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| TestCustodialNodeJSON | cli/custodial_node_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| TestContractManagerInfoError | cli/contract_management_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/bank_institutional_node.go:16 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/replication.go:7 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSnvmExecJSON | cli/snvm_test.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/connpool.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| TestMobileMiningCLI | cli/mobile_mining_node_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| cnOutput | cli/content_node_mgmt.go:17 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/content_node_mgmt.go:28 | Rate limiting, validation, auth, integrity, user, crypto |
| TestDAOProposalCreateJSON | cli/dao_proposal_test.go:17 | Rate limiting, validation, auth, integrity, user, crypto |
| TestPeerConnectGas | cli/peer_management_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/node_types.go:8 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/sidechains.go:16 | Rate limiting, validation, auth, integrity, user, crypto |
| TestInstructionNew | cli/instruction_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestInstructionList | cli/instruction_test.go:31 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn20InitRequiresFlags | cli/syn20_test.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn20MintWorkflow | cli/syn20_test.go:24 | Rate limiting, validation, auth, integrity, user, crypto |
| TestDAOTokenMintJSON | cli/dao_token_test.go:17 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSynchronizationPlaceholder | cli/synchronization_test.go:5 | Rate limiting, validation, auth, integrity, user, crypto |
| run2600 | cli/syn2600_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn2600IssueAndTransfer | cli/syn2600_test.go:16 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn2600IssueMissingFields | cli/syn2600_test.go:36 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/syn700.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| TestRegulatorCLI | cli/regulatory_management_test.go:8 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSimpleVMStatus | cli/virtual_machine_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| TestGasTableSnapshot | cli/gas_table_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/plasma_operations.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/liquidity_views.go:7 | Rate limiting, validation, auth, integrity, user, crypto |
| TestBioAuthListEmpty | cli/biometrics_auth_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| resetValidatorNode | cli/validator_node_test.go:8 | Rate limiting, validation, auth, integrity, user, crypto |
| TestValidatorNodeCreateQuorum | cli/validator_node_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/regulatory_management.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/syn1100.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/loanpool_management.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/dao_token.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/kademlia.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/authority_node_index.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/consensus_service.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| run2800 | cli/syn2800_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn2800IssueAndPay | cli/syn2800_test.go:16 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn2800IssueMissingFields | cli/syn2800_test.go:33 | Rate limiting, validation, auth, integrity, user, crypto |
| TestMiningNodeCLI | cli/mining_node_test.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| execMiningCLI | cli/mining_node_test.go:39 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/syn3600.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| execLoanCLI | cli/loanpool_proposal_test.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| TestLoanProposalCLI | cli/loanpool_proposal_test.go:30 | Rate limiting, validation, auth, integrity, user, crypto |
| TestDAOMemberAddJSON | cli/dao_access_control_test.go:17 | Rate limiting, validation, auth, integrity, user, crypto |
| TestDAOMemberUpdateRole | cli/dao_access_control_test.go:64 | Rate limiting, validation, auth, integrity, user, crypto |
| TestDAOMemberUpdateRoleUnauthorized | cli/dao_access_control_test.go:112 | Rate limiting, validation, auth, integrity, user, crypto |
| gasPrint | cli/gas_print.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| TestQuorumtrackerPlaceholder | cli/quorum_tracker_test.go:5 | Rate limiting, validation, auth, integrity, user, crypto |
| TestWalletNewCLI | cli/wallet_cli_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSystemhealthloggingPlaceholder | cli/system_health_logging_test.go:5 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/authority_apply.go:20 | Rate limiting, validation, auth, integrity, user, crypto |
| TestOpcodesHex | cli/opcodes_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/validator_management.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| TestConnPoolLifecycle | cli/connpool_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn4200TokenLifecycle | cli/syn4200_token_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn4200TokenValidation | cli/syn4200_token_test.go:23 | Rate limiting, validation, auth, integrity, user, crypto |
| TestConsensusServiceStartStop | cli/consensus_service_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn845Workflow | cli/syn845_test.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/sharding.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/ai_contract.go:21 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn500Lifecycle | cli/syn500_test.go:5 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn500Validation | cli/syn500_test.go:28 | Rate limiting, validation, auth, integrity, user, crypto |
| TestFullNodeCreateJSON | cli/node_commands_test.go:8 | Rate limiting, validation, auth, integrity, user, crypto |
| TestStakingNodeStake | cli/node_commands_test.go:18 | Rate limiting, validation, auth, integrity, user, crypto |
| TestTransactionVariableFee | cli/transaction_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/biometric.go:16 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/consensus_difficulty.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSwarmPeersEmpty | cli/swarm_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/peer_management.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestHighAvailabilityWorkflow | cli/high_availability_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| run2700 | cli/syn2700_test.go:8 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn2700CreateAndClaim | cli/syn2700_test.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn2700CreateMissingEntries | cli/syn2700_test.go:25 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/full_node.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| Next | cli/state_rw.go:23 | Rate limiting, validation, auth, integrity, user, crypto |
| Value | cli/state_rw.go:31 | Rate limiting, validation, auth, integrity, user, crypto |
| newStateStore | cli/state_rw.go:35 | Rate limiting, validation, auth, integrity, user, crypto |
| Transfer | cli/state_rw.go:39 | Rate limiting, validation, auth, integrity, user, crypto |
| SetState | cli/state_rw.go:48 | Rate limiting, validation, auth, integrity, user, crypto |
| GetState | cli/state_rw.go:52 | Rate limiting, validation, auth, integrity, user, crypto |
| HasState | cli/state_rw.go:60 | Rate limiting, validation, auth, integrity, user, crypto |
| PrefixIterator | cli/state_rw.go:65 | Rate limiting, validation, auth, integrity, user, crypto |
| BalanceOf | cli/state_rw.go:75 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/state_rw.go:81 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/token_registry.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSmartContractMarketplaceDeployJSON | cli/smart_contract_marketplace_test.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| TestBaseNodeLifecycle | cli/base_node_test.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/dao_proposal.go:15 | Rate limiting, validation, auth, integrity, user, crypto |
| parseTime | cli/syn2100.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/syn2100.go:25 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn800TokenWorkflow | cli/syn800_token_test.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/validator_node.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/syn3900.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/bank_nodes_index.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/syn3400.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| execLoanpoolCLI | cli/loanpool_test.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| TestLoanpoolCLIFlow | cli/loanpool_test.go:30 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn3200Lifecycle | cli/syn3200_test.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn3200InvalidDue | cli/syn3200_test.go:29 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/contracts.go:18 | Rate limiting, validation, auth, integrity, user, crypto |
| TestCharityRegistrationJSON | cli/charity_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| TestCharityBalancesJSON | cli/charity_test.go:23 | Rate limiting, validation, auth, integrity, user, crypto |
| TestStateBalanceZero | cli/state_rw_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| TestAuditNodeStart | cli/audit_node_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| TestAuditNodeLogAndList | cli/audit_node_test.go:17 | Rate limiting, validation, auth, integrity, user, crypto |
| TestInitrepWorkflow | cli/initrep_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestGenesisInitBlock | cli/genesis_cli_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestGenesisInitOnce | cli/genesis_cli_test.go:29 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/opcodes.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/address.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/government.go:8 | Rate limiting, validation, auth, integrity, user, crypto |
| run2100 | cli/syn2100_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn2100RegisterAndFinance | cli/syn2100_test.go:18 | Rate limiting, validation, auth, integrity, user, crypto |
| Start | cli/audit_node.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/audit_node.go:21 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/high_availability.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| TestFaucetRequest | cli/faucet_test.go:8 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/quorum_tracker.go:16 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/virtual_machine.go:15 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/syn1000.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/syn4700.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| resetSandboxMgr | cli/vm_sandbox_management_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSandboxStartStatus | cli/vm_sandbox_management_test.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| TestAuthorityNodeIndexAdd | cli/authority_node_index_test.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| TestWatchtowerLifecycle | cli/watchtower_test.go:5 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/syn1700_token.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| execNFTCLI | cli/nft_marketplace_test.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| TestNFTMarketplaceCLI | cli/nft_marketplace_test.go:30 | Rate limiting, validation, auth, integrity, user, crypto |
| TestNFTMarketplaceInvalidPrice | cli/nft_marketplace_test.go:59 | Rate limiting, validation, auth, integrity, user, crypto |
| TestNodeAddrCLI | cli/node_types_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/liquidity_pools.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/synchronization.go:15 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn1600InitRequiresFlags | cli/syn1600_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn1600Workflow | cli/syn1600_test.go:20 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSidechainListEmpty | cli/sidechains_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/swarm.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| TestWalletNew | cli/wallet_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn1401IssueRequiresFlags | cli/syn1401_test.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn1401Workflow | cli/syn1401_test.go:25 | Rate limiting, validation, auth, integrity, user, crypto |
| run2900 | cli/syn2900_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn2900IssueAndClaim | cli/syn2900_test.go:16 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn2900IssueMissingFields | cli/syn2900_test.go:33 | Rate limiting, validation, auth, integrity, user, crypto |
| TestAuthorityRegister | cli/authority_nodes_test.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| TestAuthorityVote | cli/authority_nodes_test.go:28 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn12InitRequiresFlags | cli/syn12_test.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn12MintWorkflow | cli/syn12_test.go:25 | Rate limiting, validation, auth, integrity, user, crypto |
| TestWatchtowerNodeLifecycle | cli/watchtower_node_test.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/faucet.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/tx_control.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/fees.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/rpc_webrtc.go:15 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/vm_sandbox_management.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| run300 | cli/syn300_token_test.go:5 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn300InitAndProposal | cli/syn300_token_test.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn300InitMissingBalances | cli/syn300_token_test.go:27 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/plasma.go:18 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/idwallet.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/syn4200_token.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/smart_contract_marketplace.go:15 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/syn3700_token.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/root.go:30 | Rate limiting, validation, auth, integrity, user, crypto |
| Execute | cli/root.go:36 | Rate limiting, validation, auth, integrity, user, crypto |
| RootCmd | cli/root.go:39 | Rate limiting, validation, auth, integrity, user, crypto |
| TestCoinInfoJSON | cli/coin_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| TestCoinRewardValidation | cli/coin_test.go:20 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/cross_consensus_scaling_networks.go:15 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/warfare_node.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| TestKademliaStoreGet | cli/kademlia_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestKademliaClosest | cli/kademlia_test.go:49 | Rate limiting, validation, auth, integrity, user, crypto |
| TestLoanmgrPauseResume | cli/loanpool_management_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/consensus_specific_node.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| run2500 | cli/syn2500_token_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn2500AddAndUpdate | cli/syn2500_token_test.go:15 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn2500AddMissingFields | cli/syn2500_token_test.go:33 | Rate limiting, validation, auth, integrity, user, crypto |
| TestElectedNodeCreateJSON | cli/elected_authority_node_test.go:16 | Rate limiting, validation, auth, integrity, user, crypto |
| TestPoolViews | cli/liquidity_views_test.go:6 | Rate limiting, validation, auth, integrity, user, crypto |
| TestBlockCreateAndHeader | cli/block_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| TestIDWalletWorkflow | cli/idwallet_test.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| TestCrossChainBridgeDepositJSON | cli/cross_chain_bridge_test.go:8 | Rate limiting, validation, auth, integrity, user, crypto |
| TestConsensusModeSetShow | cli/consensus_mode_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| TestCrossChainTransactionsCommands | cli/cross_chain_transactions_test.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| TestAddressUtilities | cli/address_test.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn70Workflow | cli/syn70_test.go:11 | Rate limiting, validation, auth, integrity, user, crypto |
| parsePubKey | cli/ecdsa_util.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| decodeSig | cli/ecdsa_util.go:27 | Rate limiting, validation, auth, integrity, user, crypto |
| VerifySignature | cli/ecdsa_util.go:35 | Rate limiting, validation, auth, integrity, user, crypto |
| TestOptimizationnodePlaceholder | cli/optimization_node_test.go:5 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn131CreateRequiresFlags | cli/syn131_token_test.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn131Workflow | cli/syn131_token_test.go:22 | Rate limiting, validation, auth, integrity, user, crypto |
| parseEntries | cli/syn2700.go:15 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/syn2700.go:37 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/syn5000_index.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/wallet.go:8 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/nft_marketplace.go:16 | Rate limiting, validation, auth, integrity, user, crypto |
| TestDAOStakingStakeJSON | cli/dao_staking_test.go:17 | Rate limiting, validation, auth, integrity, user, crypto |
| TestFeesEstimate | cli/fees_test.go:8 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn3700Lifecycle | cli/syn3700_token_test.go:5 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn3700Validation | cli/syn3700_token_test.go:43 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/syn800_token.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| TestNetworkCLIStartStop | cli/network_test.go:15 | Rate limiting, validation, auth, integrity, user, crypto |
| TestNetworkCLIPeersAndBroadcast | cli/network_test.go:37 | Rate limiting, validation, auth, integrity, user, crypto |
| execNetCLI | cli/network_test.go:66 | Rate limiting, validation, auth, integrity, user, crypto |
| TestZeroTrustDataChannelsFlow | cli/zero_trust_data_channels_test.go:5 | Rate limiting, validation, auth, integrity, user, crypto |
| compOut | cli/compression.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/compression.go:23 | Rate limiting, validation, auth, integrity, user, crypto |
| parseBalances | cli/syn300_token.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/syn300_token.go:31 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/loanpool_proposal.go:17 | Rate limiting, validation, auth, integrity, user, crypto |
| parseMode | cli/consensus_mode.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/consensus_mode.go:26 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/regulatory_node.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| bioOutput | cli/biometrics_auth.go:20 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/biometrics_auth.go:31 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/rollups.go:15 | Rate limiting, validation, auth, integrity, user, crypto |
| TestCompressionSaveLoad | cli/compression_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/cross_chain_transactions.go:15 | Rate limiting, validation, auth, integrity, user, crypto |
| TestConsensusDifficultySample | cli/consensus_difficulty_test.go:6 | Rate limiting, validation, auth, integrity, user, crypto |
| execCLI | cli/cross_chain_test.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| TestCrossChainBridgeCommands | cli/cross_chain_test.go:27 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/token_syn130.go:18 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/syn1600.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| TestPlasmaoperationsPlaceholder | cli/plasma_operations_test.go:5 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/storage_marketplace.go:15 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/identity.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestConsensusWeights | cli/consensus_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestContractOpcodesGas | cli/contracts_opcodes_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| TestGatewayPlaceholder | cli/gateway_test.go:5 | Rate limiting, validation, auth, integrity, user, crypto |
| resetSyn130 | cli/token_syn130_test.go:10 | Rate limiting, validation, auth, integrity, user, crypto |
| TestTokenSyn130RegisterAndList | cli/token_syn130_test.go:15 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/watchtower_node.go:16 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/snvm.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| compOutput | cli/compliance.go:19 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/compliance.go:30 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/syn1300.go:12 | Rate limiting, validation, auth, integrity, user, crypto |
| TestPlasmaPlaceholder | cli/plasma_test.go:5 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/staking_node.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/syn2900.go:14 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/base_node.go:15 | Rate limiting, validation, auth, integrity, user, crypto |
| run2369 | cli/syn2369_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn2369CreateAndTransfer | cli/syn2369_test.go:15 | Rate limiting, validation, auth, integrity, user, crypto |
| TestSyn2369MissingFields | cli/syn2369_test.go:33 | Rate limiting, validation, auth, integrity, user, crypto |
| TestValidatorAddStake | cli/validator_management_test.go:9 | Rate limiting, validation, auth, integrity, user, crypto |
| init | cli/holographic_node.go:13 | Rate limiting, validation, auth, integrity, user, crypto |
