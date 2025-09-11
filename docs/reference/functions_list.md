# Function List

This catalogue lists exported functions across the repository for quick navigation. Regenerate with `go list` and `rg` when APIs change.

| File | Line | Signature |
|---|---|---|
| `biometrics_auth.go` | `15` | `func NewBiometricsAuth() *BiometricsAuth {` |
| `biometrics_auth.go` | `20` | `func (b *BiometricsAuth) Enroll(addr string, biometric []byte) {` |
| `biometrics_auth.go` | `27` | `func (b *BiometricsAuth) Verify(addr string, biometric []byte) bool {` |
| `biometrics_auth.go` | `38` | `func (b *BiometricsAuth) Remove(addr string) {` |
| `biometric_security_node.go` | `14` | `func NewBiometricSecurityNode(id string, auth *BiometricsAuth) *BiometricSecurityNode {` |
| `biometric_security_node.go` | `22` | `func (b *BiometricSecurityNode) GetID() string { return b.id }` |
| `biometric_security_node.go` | `25` | `func (b *BiometricSecurityNode) Enroll(addr string, biometric []byte) {` |
| `biometric_security_node.go` | `30` | `func (b *BiometricSecurityNode) Authenticate(addr string, biometric []byte) bool {` |
| `biometric_security_node.go` | `35` | `func (b *BiometricSecurityNode) SecureExecute(addr string, biometric []byte, fn func() error) error {` |
| `warfare_node.go` | `27` | `func NewWarfareNode(id string) *WarfareNode {` |
| `warfare_node.go` | `32` | `func (w *WarfareNode) GetID() string { return w.id }` |
| `warfare_node.go` | `36` | `func (w *WarfareNode) SecureCommand(cmd string) error {` |
| `warfare_node.go` | `44` | `func (w *WarfareNode) TrackLogistics(assetID, location, status string) {` |
| `warfare_node.go` | `58` | `func (w *WarfareNode) ShareTactical(info string) {` |
| `warfare_node.go` | `63` | `func (w *WarfareNode) Logistics() []LogisticsRecord {` |
| `geospatial_node.go` | `23` | `func NewGeospatialNode() *GeospatialNode {` |
| `geospatial_node.go` | `28` | `func (n *GeospatialNode) Record(subject string, lat, lon float64) {` |
| `geospatial_node.go` | `40` | `func (n *GeospatialNode) History(subject string) []GeoRecord {` |
| `dynamic_consensus_hopping.go` | `31` | `func NewConsensusHopper(initial ConsensusMode) *ConsensusHopper {` |
| `dynamic_consensus_hopping.go` | `36` | `func (h *ConsensusHopper) Mode() ConsensusMode {` |
| `dynamic_consensus_hopping.go` | `43` | `func (h *ConsensusHopper) Evaluate(m NetworkMetrics) ConsensusMode {` |
| `holographic.go` | `10` | `func SplitHolographic(id string, data []byte, n int) HolographicFrame {` |
| `holographic.go` | `30` | `func ReconstructHolographic(frame HolographicFrame) []byte {` |
| `energy_efficient_node.go` | `26` | `func NewEnergyEfficientNode(id string, tracker *EnergyEfficiencyTracker) *EnergyEfficientNode {` |
| `energy_efficient_node.go` | `31` | `func (n *EnergyEfficientNode) ID() string { return n.id }` |
| `energy_efficient_node.go` | `34` | `func (n *EnergyEfficientNode) RecordUsage(txProcessed int, energyKWh float64) {` |
| `energy_efficient_node.go` | `39` | `func (n *EnergyEfficientNode) AddOffset(credits float64) {` |
| `energy_efficient_node.go` | `46` | `func (n *EnergyEfficientNode) Certify() SustainabilityCertificate {` |
| `energy_efficient_node.go` | `61` | `func (n *EnergyEfficientNode) Certificate() SustainabilityCertificate {` |
| `energy_efficient_node.go` | `68` | `func (n *EnergyEfficientNode) ShouldThrottle(threshold float64) bool {` |
| `content_node.go` | `20` | `func NewContentNetworkNode(id, addr string) *ContentNetworkNode {` |
| `content_node.go` | `30` | `func (n *ContentNetworkNode) Register(meta ContentMeta) error {` |
| `content_node.go` | `45` | `func (n *ContentNetworkNode) Unregister(id string) error {` |
| `content_node.go` | `56` | `func (n *ContentNetworkNode) Content(id string) (ContentMeta, bool) {` |
| `content_node.go` | `64` | `func (n *ContentNetworkNode) List() []ContentMeta {` |
| `nodesextra/holographic_node.go` | `17` | `func NewHolographicNode(id string) *HolographicNode {` |
| `nodesextra/holographic_node.go` | `25` | `func (n *HolographicNode) ID() string { return n.id }` |
| `nodesextra/holographic_node.go` | `29` | `func (n *HolographicNode) Start() error { return nil }` |
| `nodesextra/holographic_node.go` | `33` | `func (n *HolographicNode) Stop() error { return nil }` |
| `nodesextra/holographic_node.go` | `36` | `func (n *HolographicNode) Store(frame synnergy.HolographicFrame) {` |
| `nodesextra/holographic_node.go` | `44` | `func (n *HolographicNode) Retrieve(id string) (synnergy.HolographicFrame, bool) {` |
| `core/replication_test.go` | `5` | `func TestReplicator(t *testing.T) {` |
| `core/syn5000_test.go` | `5` | `func TestSYN5000Token(t *testing.T) {` |
| `core/cross_consensus_scaling_networks.go` | `23` | `func NewConsensusNetworkManager() *ConsensusNetworkManager {` |
| `core/cross_consensus_scaling_networks.go` | `28` | `func (m *ConsensusNetworkManager) RegisterNetwork(source, target string) int {` |
| `core/cross_consensus_scaling_networks.go` | `38` | `func (m *ConsensusNetworkManager) ListNetworks() []ConsensusNetwork {` |
| `core/cross_consensus_scaling_networks.go` | `49` | `func (m *ConsensusNetworkManager) GetNetwork(id int) (ConsensusNetwork, error) {` |
| `core/cross_chain_transactions_test.go` | `5` | `func TestCrossChainTxManager(t *testing.T) {` |
| `nodesextra/optimization_nodes/optimization.go` | `25` | `func (o *SimpleOptimizer) Optimize(m Metrics) Suggestion {` |
| `core/block.go` | `20` | `func NewSubBlock(txs []*Transaction, validator string) *SubBlock {` |
| `core/block.go` | `29` | `func (sb *SubBlock) Hash() string {` |
| `core/block.go` | `40` | `func (sb *SubBlock) VerifySignature() bool {` |
| `core/block.go` | `55` | `func NewBlock(subBlocks []*SubBlock, prevHash string) *Block {` |
| `core/block.go` | `61` | `func (b *Block) HeaderHash(nonce uint64) string {` |
| `core/block.go` | `71` | `func signSubBlock(validator, msg string) string {` |
| `content_node_impl.go` | `24` | `func NewContentNode(key []byte) (*ContentNode, error) {` |
| `content_node_impl.go` | `36` | `func (n *ContentNode) StoreContent(name string, data []byte) (ContentMeta, error) {` |
| `content_node_impl.go` | `64` | `func (n *ContentNode) RetrieveContent(id string) ([]byte, bool, error) {` |
| `content_node_impl.go` | `93` | `func (n *ContentNode) Meta(id string) (ContentMeta, bool) {` |
| `core/syn2700.go` | `18` | `func NewVestingSchedule(entries []VestingEntry) *VestingSchedule {` |
| `core/syn2700.go` | `23` | `func (v *VestingSchedule) Claim(now time.Time) uint64 {` |
| `core/syn2700.go` | `36` | `func (v *VestingSchedule) Pending(now time.Time) uint64 {` |
| `indexing_node.go` | `13` | `func NewIndexingNode() *IndexingNode {` |
| `indexing_node.go` | `18` | `func (n *IndexingNode) Index(key string, value []byte) {` |
| `indexing_node.go` | `26` | `func (n *IndexingNode) Query(key string) ([]byte, bool) {` |
| `indexing_node.go` | `39` | `func (n *IndexingNode) Remove(key string) {` |
| `indexing_node.go` | `46` | `func (n *IndexingNode) Keys() []string {` |
| `core/dao_proposal_test.go` | `5` | `func TestDAOProposal(t *testing.T) {` |
| `core/syn1700_token.go` | `28` | `func NewEvent(name, desc, location string, start, end int64, supply uint64) *EventMetadata {` |
| `core/syn1700_token.go` | `41` | `func (e *EventMetadata) IssueTicket(owner, class, ticketType string, price uint64) (uint64, error) {` |
| `core/syn1700_token.go` | `52` | `func (e *EventMetadata) TransferTicket(id uint64, from, to string) error {` |
| `core/syn1700_token.go` | `62` | `func (e *EventMetadata) VerifyTicket(id uint64, holder string) bool {` |
| `core/syn3800_test.go` | `5` | `func TestGrantRegistry(t *testing.T) {` |
| `core/transaction.go` | `42` | `func NewTransaction(from, to string, amount, fee, nonce uint64) *Transaction {` |
| `core/transaction.go` | `58` | `func (t *Transaction) Hash() string {` |
| `core/transaction.go` | `75` | `func (t *Transaction) Verify(pub *ecdsa.PublicKey) bool {` |
| `core/transaction.go` | `84` | `func (t *Transaction) AttachBiometric(userID string, biometric []byte, svc *BiometricService) error {` |
| `ai_inference_analysis.go` | `22` | `func NewInferenceEngine() *InferenceEngine {` |
| `ai_inference_analysis.go` | `27` | `func (e *InferenceEngine) LoadModel(hash string, data []byte) {` |
| `ai_inference_analysis.go` | `34` | `func (e *InferenceEngine) Run(hash string, input []byte) ([]byte, error) {` |
| `ai_inference_analysis.go` | `46` | `func (e *InferenceEngine) Analyse(txIDs []string) []FraudResult {` |
| `core/syn2500_token.go` | `18` | `func NewSyn2500Member(id, addr string, power uint64, meta map[string]string) *Syn2500Member {` |
| `core/syn2500_token.go` | `33` | `func (m *Syn2500Member) UpdateVotingPower(power uint64) {` |
| `core/syn2500_token.go` | `44` | `func NewSyn2500Registry() *Syn2500Registry {` |
| `core/syn2500_token.go` | `49` | `func (r *Syn2500Registry) AddMember(m *Syn2500Member) {` |
| `core/syn2500_token.go` | `56` | `func (r *Syn2500Registry) GetMember(id string) (*Syn2500Member, bool) {` |
| `core/syn2500_token.go` | `64` | `func (r *Syn2500Registry) RemoveMember(id string) {` |
| `core/syn2500_token.go` | `71` | `func (r *Syn2500Registry) ListMembers() []*Syn2500Member {` |
| `ai_secure_storage.go` | `19` | `func NewSecureStorage() *SecureStorage {` |
| `ai_secure_storage.go` | `24` | `func (s *SecureStorage) Store(hash string, data, key []byte) error {` |
| `ai_secure_storage.go` | `45` | `func (s *SecureStorage) Retrieve(hash string, key []byte) ([]byte, error) {` |
| `core/syn3200.go` | `32` | `func NewBillRegistry() *BillRegistry {` |
| `core/syn3200.go` | `37` | `func (r *BillRegistry) Create(id, issuer, payer string, amt uint64, due time.Time, meta string) (*Bill, error) {` |
| `core/syn3200.go` | `47` | `func (r *BillRegistry) Pay(id, payer string, amt uint64) error {` |
| `core/syn3200.go` | `65` | `func (r *BillRegistry) Adjust(id string, amt uint64) error {` |
| `core/syn3200.go` | `75` | `func (r *BillRegistry) Get(id string) (*Bill, bool) {` |
| `core/syn3900_test.go` | `5` | `func TestBenefitRegistry(t *testing.T) {` |
| `core/syn2900.go` | `23` | `func NewTokenInsurancePolicy(id, holder, coverage string, premium, payout, deductible, limit uint64, start, end time.Time) *TokenInsurancePolicy {` |
| `core/syn2900.go` | `38` | `func (p *TokenInsurancePolicy) IsActive(now time.Time) bool {` |
| `core/syn2900.go` | `43` | `func (p *TokenInsurancePolicy) Claim(now time.Time) (uint64, error) {` |
| `core/authority_node_index_test.go` | `5` | `func TestAuthorityNodeIndex(t *testing.T) {` |
| `core/regulatory_management_test.go` | `5` | `func TestRegulatoryManager(t *testing.T) {` |
| `core/dao_token.go` | `18` | `func NewDAOTokenLedger(mgr *DAOManager) *DAOTokenLedger {` |
| `core/dao_token.go` | `23` | `func (l *DAOTokenLedger) Mint(daoID, admin, addr string, amount uint64) error {` |
| `core/dao_token.go` | `44` | `func (l *DAOTokenLedger) Transfer(daoID, from, to string, amount uint64) error {` |
| `core/dao_token.go` | `63` | `func (l *DAOTokenLedger) Balance(daoID, addr string) uint64 {` |
| `core/dao_token.go` | `74` | `func (l *DAOTokenLedger) Burn(daoID, admin, addr string, amount uint64) error {` |
| `core/syn5000.go` | `27` | `func NewSYN5000Token(name, symbol string, decimals uint8) *SYN5000Token {` |
| `core/syn5000.go` | `32` | `func (t *SYN5000Token) PlaceBet(bettor string, amount uint64, odds float64, game string) uint64 {` |
| `core/syn5000.go` | `40` | `func (t *SYN5000Token) ResolveBet(betID uint64, win bool) (uint64, error) {` |
| `core/syn5000.go` | `57` | `func (t *SYN5000Token) GetBet(betID uint64) (*BetRecord, bool) {` |
| `node_ext/holographic_node.go` | `17` | `func NewHolographicNode(id string) *HolographicNode {` |
| `node_ext/holographic_node.go` | `25` | `func (n *HolographicNode) ID() string { return n.id }` |
| `node_ext/holographic_node.go` | `29` | `func (n *HolographicNode) Start() error { return nil }` |
| `node_ext/holographic_node.go` | `33` | `func (n *HolographicNode) Stop() error { return nil }` |
| `node_ext/holographic_node.go` | `36` | `func (n *HolographicNode) Store(frame synnergy.HolographicFrame) {` |
| `node_ext/holographic_node.go` | `44` | `func (n *HolographicNode) Retrieve(id string) (synnergy.HolographicFrame, bool) {` |
| `core/cross_chain.go` | `24` | `func NewBridgeRegistry() *BridgeRegistry {` |
| `core/cross_chain.go` | `29` | `func (r *BridgeRegistry) RegisterBridge(source, target, relayer string) (*Bridge, error) {` |
| `core/cross_chain.go` | `45` | `func (r *BridgeRegistry) ListBridges() []*Bridge {` |
| `core/cross_chain.go` | `56` | `func (r *BridgeRegistry) GetBridge(id string) (*Bridge, bool) {` |
| `core/cross_chain.go` | `64` | `func (r *BridgeRegistry) AuthorizeRelayer(id, relayer string) error {` |
| `core/cross_chain.go` | `79` | `func (r *BridgeRegistry) RevokeRelayer(id, relayer string) error {` |
| `core/full_node.go` | `22` | `func NewFullNode(id nodes.Address, mode FullNodeMode) *FullNode {` |
| `core/full_node.go` | `27` | `func (f *FullNode) SetMode(m FullNodeMode) {` |
| `core/full_node.go` | `32` | `func (f *FullNode) CurrentMode() FullNodeMode {` |
| `core/full_node.go` | `37` | `func (f *FullNode) IsArchive() bool {` |
| `ai.go` | `42` | `func NewAIService() *AIService {` |
| `ai.go` | `51` | `func (s *AIService) PredictFraud(txJSON []byte) (float64, error) {` |
| `ai.go` | `61` | `func (s *AIService) OptimiseBaseFee(statsJSON []byte) (uint64, error) {` |
| `ai.go` | `71` | `func (s *AIService) ForecastVolume(statsJSON []byte) (uint64, error) {` |
| `ai.go` | `81` | `func (s *AIService) PublishModel(cid string, royaltyBps uint16) (string, error) {` |
| `ai.go` | `96` | `func (s *AIService) FetchModel(hash string) (AIModelMetadata, bool) {` |
| `ai.go` | `104` | `func (s *AIService) ListModel(hash, cid, seller string, price uint64) string {` |
| `ai.go` | `109` | `func (s *AIService) BuyModel(listingID, buyer string, amount uint64) (string, error) {` |
| `ai.go` | `131` | `func (s *AIService) RentModel(listingID, renter string, hours int, amount uint64) (string, error) {` |
| `ai.go` | `154` | `func (s *AIService) ReleaseEscrow(id string) error {` |
| `data_distribution.go` | `19` | `func NewDataDistribution() *DataDistribution {` |
| `data_distribution.go` | `24` | `func (d *DataDistribution) Offer(nodeID string, meta ContentMeta) {` |
| `data_distribution.go` | `36` | `func (d *DataDistribution) Locations(contentID string) []string {` |
| `node_ext/optimization_nodes/optimization.go` | `25` | `func (o *SimpleOptimizer) Optimize(m Metrics) Suggestion {` |
| `cmd/synnergy/main.go` | `9` | `func main() {` |
| `anomaly_detection.go` | `18` | `func NewAnomalyDetector(threshold float64) *AnomalyDetector {` |
| `anomaly_detection.go` | `23` | `func (a *AnomalyDetector) Update(v float64) {` |
| `anomaly_detection.go` | `33` | `func (a *AnomalyDetector) IsAnomalous(v float64) bool {` |
| `environmental_monitoring_node.go` | `17` | `func (c EnvCondition) Evaluate(data []byte) bool {` |
| `environmental_monitoring_node.go` | `45` | `func NewEnvironmentalMonitoringNode() *EnvironmentalMonitoringNode {` |
| `environmental_monitoring_node.go` | `50` | `func (n *EnvironmentalMonitoringNode) SetCondition(c EnvCondition) {` |
| `environmental_monitoring_node.go` | `57` | `func (n *EnvironmentalMonitoringNode) Trigger(sensorID string, data []byte) bool {` |
| `content_types.go` | `23` | `func NewContentMeta(id, name string, size int64, hash string) ContentMeta {` |
| `ai_training.go` | `28` | `func NewTrainingManager() *TrainingManager {` |
| `ai_training.go` | `33` | `func (m *TrainingManager) Start(datasetCID, modelCID string) string {` |
| `ai_training.go` | `49` | `func (m *TrainingManager) Status(id string) (TrainingJob, bool) {` |
| `ai_training.go` | `57` | `func (m *TrainingManager) List() []TrainingJob {` |
| `ai_training.go` | `68` | `func (m *TrainingManager) Cancel(id string) error {` |
| `Tokens/syn845.go` | `39` | `func NewDebtRegistry() *DebtRegistry {` |
| `Tokens/syn845.go` | `44` | `func (r *DebtRegistry) CreateToken(name, symbol, owner string, supply uint64) (string, *DebtToken) {` |
| `Tokens/syn845.go` | `55` | `func (r *DebtRegistry) IssueDebt(tokenID, debtID, borrower string, principal uint64, rate, penalty float64, due time.Time) error {` |
| `Tokens/syn845.go` | `70` | `func (r *DebtRegistry) RecordPayment(tokenID, debtID string, amount uint64) error {` |
| `Tokens/syn845.go` | `86` | `func (r *DebtRegistry) GetDebt(tokenID, debtID string) (*DebtRecord, error) {` |
| `energy_efficiency.go` | `18` | `func NewEnergyEfficiencyTracker() *EnergyEfficiencyTracker {` |
| `energy_efficiency.go` | `23` | `func (t *EnergyEfficiencyTracker) Record(validator string, txProcessed int, energyKWh float64) {` |
| `energy_efficiency.go` | `34` | `func (t *EnergyEfficiencyTracker) Efficiency(validator string) (float64, bool) {` |
| `energy_efficiency.go` | `45` | `func (t *EnergyEfficiencyTracker) NetworkAverage() float64 {` |
| `Tokens/syn12.go` | `22` | `func NewSYN12Token(id TokenID, name, symbol string, meta SYN12Metadata, decimals uint8) *SYN12Token {` |
| `ai_drift_monitor.go` | `15` | `func NewDriftMonitor() *DriftMonitor {` |
| `ai_drift_monitor.go` | `20` | `func (d *DriftMonitor) UpdateBaseline(modelHash string, metric float64) {` |
| `ai_drift_monitor.go` | `27` | `func (d *DriftMonitor) HasDrift(modelHash string, metric, threshold float64) bool {` |
| `Tokens/base.go` | `32` | `func NewBaseToken(id TokenID, name, symbol string, decimals uint8) *BaseToken {` |
| `Tokens/base.go` | `43` | `func (t *BaseToken) ID() TokenID { return t.id }` |
| `Tokens/base.go` | `46` | `func (t *BaseToken) Name() string { return t.name }` |
| `Tokens/base.go` | `49` | `func (t *BaseToken) Symbol() string { return t.symbol }` |
| `Tokens/base.go` | `52` | `func (t *BaseToken) Decimals() uint8 { return t.decimals }` |
| `Tokens/base.go` | `55` | `func (t *BaseToken) TotalSupply() uint64 { return t.supply }` |
| `Tokens/base.go` | `58` | `func (t *BaseToken) BalanceOf(addr string) uint64 {` |
| `Tokens/base.go` | `63` | `func (t *BaseToken) Transfer(from, to string, amount uint64) error {` |
| `Tokens/base.go` | `73` | `func (t *BaseToken) Mint(to string, amount uint64) error {` |
| `Tokens/base.go` | `80` | `func (t *BaseToken) Burn(from string, amount uint64) error {` |
| `Tokens/base_test.go` | `8` | `func TestBaseTokenMintTransferBurn(t *testing.T) {` |
| `Tokens/base_test.go` | `27` | `func TestSYN10Info(t *testing.T) {` |
| `Tokens/base_test.go` | `38` | `func TestSYN1000ReserveValue(t *testing.T) {` |
| `Tokens/base_test.go` | `53` | `func TestSYN1100Access(t *testing.T) {` |
| `core/connection_pool.go` | `21` | `func NewConnectionPool(max int) *ConnectionPool {` |
| `core/connection_pool.go` | `30` | `func (p *ConnectionPool) Acquire(id string) (*Connection, error) {` |
| `core/connection_pool.go` | `45` | `func (p *ConnectionPool) Release(id string) {` |
| `core/connection_pool.go` | `52` | `func (p *ConnectionPool) Size() int {` |
| `Tokens/syn70.go` | `19` | `func NewSYN70Token(id TokenID, name, symbol string, decimals uint8) *SYN70Token {` |
| `Tokens/syn70.go` | `27` | `func (t *SYN70Token) MintAsset(owner, assetID, metadata string) error {` |
| `Tokens/syn70.go` | `36` | `func (t *SYN70Token) TransferAsset(assetID, from, to string) error {` |
| `Tokens/syn70.go` | `53` | `func (t *SYN70Token) GetAsset(assetID string) (SYN70Asset, bool) {` |
| `core/biometric.go` | `18` | `func NewBiometricService() *BiometricService {` |
| `core/biometric.go` | `39` | `func (b *BiometricService) Enroll(userID string, biometric []byte, pub ed25519.PublicKey) error {` |
| `core/biometric.go` | `61` | `func (b *BiometricService) Verify(userID string, biometric []byte, sig []byte) bool {` |
| `Tokens/syn20.go` | `13` | `func NewSYN20Token(id TokenID, name, symbol string, decimals uint8) *SYN20Token {` |
| `Tokens/syn20.go` | `21` | `func (t *SYN20Token) Pause() { t.paused = true }` |
| `Tokens/syn20.go` | `24` | `func (t *SYN20Token) Unpause() { t.paused = false }` |
| `Tokens/syn20.go` | `27` | `func (t *SYN20Token) Freeze(addr string) { t.frozen[addr] = true }` |
| `Tokens/syn20.go` | `30` | `func (t *SYN20Token) Unfreeze(addr string) { delete(t.frozen, addr) }` |
| `Tokens/syn20.go` | `33` | `func (t *SYN20Token) Transfer(from, to string, amount uint64) error {` |
| `Tokens/syn20.go` | `44` | `func (t *SYN20Token) Mint(to string, amount uint64) error {` |
| `Tokens/syn20.go` | `52` | `func (t *SYN20Token) Burn(from string, amount uint64) error {` |
| `core/fees.go` | `33` | `func CalculateBaseFee(recent []uint64, adjustment float64) uint64 {` |
| `core/fees.go` | `49` | `func CalculateVariableFee(gasUnits, gasPrice uint64) uint64 {` |
| `core/fees.go` | `54` | `func CalculatePriorityFee(tip uint64) uint64 { return tip }` |
| `core/fees.go` | `57` | `func FeeForTransfer(dataSize, baseFee, variableRate, tip uint64) FeeBreakdown {` |
| `core/fees.go` | `64` | `func FeeForPurchase(calls, baseFee, variableRate, tip uint64) FeeBreakdown {` |
| `core/fees.go` | `71` | `func FeeForTokenUsage(computationUnits, baseFee, variableRate, tip uint64) FeeBreakdown {` |
| `core/fees.go` | `78` | `func FeeForContract(complexityFactor, baseFee, variableRate, tip uint64) FeeBreakdown {` |
| `core/fees.go` | `85` | `func FeeForWalletVerification(securityLevel, baseFee, variableRate, tip uint64) FeeBreakdown {` |
| `core/fees.go` | `94` | `func FeeForValidatedTransfer(dataSize, baseFee, variableRate, tip uint64, validated bool) FeeBreakdown {` |
| `core/fees.go` | `115` | `func DistributeFees(total uint64) FeeDistribution {` |
| `core/fees.go` | `129` | `func ApplyFeeCapFloor(fee, cap, floor uint64) uint64 {` |
| `core/fees.go` | `147` | `func (p FeePolicy) Enforce(fee uint64) (uint64, string) {` |
| `core/fees.go` | `161` | `func AdjustFeeRates(baseFee, variableRate uint64, load float64) (uint64, uint64) {` |
| `core/fees.go` | `174` | `func EstimateFee(txType TransactionType, units, baseFee, variableRate, tip uint64) FeeBreakdown {` |
| `core/fees.go` | `194` | `func ShareProportional(total uint64, weights map[string]uint64) map[string]uint64 {` |
| `core/fees.go` | `227` | `func NewFeeDistributionContract(l *Ledger) *FeeDistributionContract {` |
| `core/fees.go` | `232` | `func (f *FeeDistributionContract) Distribute(shares map[string]uint64) {` |
| `core/fees.go` | `240` | `func AdjustForBlockUtilization(pool uint64, used, capacity int) uint64 {` |
| `cli/biometric.go` | `14` | `func init() {` |
| `core/private_transactions_test.go` | `8` | `func TestEncryptDecrypt(t *testing.T) {` |
| `core/private_transactions_test.go` | `24` | `func TestPrivateTxManager(t *testing.T) {` |
| `core/cross_chain_contracts.go` | `22` | `func NewCrossChainRegistry() *CrossChainRegistry {` |
| `core/cross_chain_contracts.go` | `27` | `func (r *CrossChainRegistry) RegisterMapping(local, remoteChain, remoteAddr string) {` |
| `core/cross_chain_contracts.go` | `34` | `func (r *CrossChainRegistry) GetMapping(local string) (*ContractMapping, bool) {` |
| `core/cross_chain_contracts.go` | `42` | `func (r *CrossChainRegistry) ListMappings() []*ContractMapping {` |
| `core/cross_chain_contracts.go` | `53` | `func (r *CrossChainRegistry) RemoveMapping(local string) error {` |
| `cli/fees.go` | `11` | `func init() {` |
| `core/syn300_token.go` | `30` | `func NewSYN300Token(initial map[string]uint64) *SYN300Token {` |
| `core/syn300_token.go` | `44` | `func (t *SYN300Token) Delegate(owner, delegate string) {` |
| `core/syn300_token.go` | `55` | `func (t *SYN300Token) RevokeDelegation(owner string) {` |
| `core/syn300_token.go` | `62` | `func (t *SYN300Token) VotingPower(addr string) uint64 {` |
| `core/syn300_token.go` | `68` | `func (t *SYN300Token) votingPowerLocked(addr string) uint64 {` |
| `core/syn300_token.go` | `79` | `func (t *SYN300Token) CreateProposal(creator, description string) uint64 {` |
| `core/syn300_token.go` | `96` | `func (t *SYN300Token) Vote(id uint64, voter string, approve bool) error {` |
| `core/syn300_token.go` | `117` | `func (t *SYN300Token) Execute(id uint64, quorum uint64) error {` |
| `core/syn300_token.go` | `139` | `func (t *SYN300Token) ProposalStatus(id uint64) (*GovernanceProposal, error) {` |
| `core/syn300_token.go` | `160` | `func (t *SYN300Token) ListProposals() []*GovernanceProposal {` |
| `cli/snvm.go` | `13` | `func init() {` |
| `Tokens/syn200.go` | `39` | `func NewCarbonRegistry() *CarbonRegistry {` |
| `Tokens/syn200.go` | `44` | `func (r *CarbonRegistry) Register(owner, name string, total uint64) *CarbonProject {` |
| `Tokens/syn200.go` | `62` | `func (r *CarbonRegistry) Issue(projectID, holder string, amount uint64) error {` |
| `Tokens/syn200.go` | `78` | `func (r *CarbonRegistry) Retire(projectID, holder string, amount uint64) error {` |
| `Tokens/syn200.go` | `96` | `func (r *CarbonRegistry) AddVerification(projectID, verifier, verID, status string) error {` |
| `Tokens/syn200.go` | `109` | `func (r *CarbonRegistry) Verifications(projectID string) ([]Verification, bool) {` |
| `Tokens/syn200.go` | `120` | `func (r *CarbonRegistry) ProjectInfo(projectID string) (*CarbonProject, bool) {` |
| `Tokens/syn200.go` | `137` | `func (r *CarbonRegistry) ListProjects() []*CarbonProject {` |
| `core/swarm.go` | `13` | `func NewSwarm() *Swarm {` |
| `core/swarm.go` | `18` | `func (s *Swarm) Join(n *Node) {` |
| `core/swarm.go` | `25` | `func (s *Swarm) Leave(id string) {` |
| `core/swarm.go` | `32` | `func (s *Swarm) Members() []*Node {` |
| `core/swarm.go` | `44` | `func (s *Swarm) Broadcast(tx *Transaction) {` |
| `cli/tx_control.go` | `13` | `func init() {` |
| `Tokens/syn2369.go` | `28` | `func NewItemRegistry() *ItemRegistry {` |
| `Tokens/syn2369.go` | `33` | `func (r *ItemRegistry) CreateItem(owner, name, desc string, attrs map[string]string) *VirtualItem {` |
| `Tokens/syn2369.go` | `47` | `func (r *ItemRegistry) TransferItem(itemID, newOwner string) error {` |
| `Tokens/syn2369.go` | `59` | `func (r *ItemRegistry) UpdateAttributes(itemID string, attrs map[string]string) error {` |
| `Tokens/syn2369.go` | `73` | `func (r *ItemRegistry) GetItem(itemID string) (*VirtualItem, bool) {` |
| `Tokens/syn2369.go` | `89` | `func (r *ItemRegistry) ListItems() []*VirtualItem {` |
| `cli/root.go` | `12` | `func Execute() error { return rootCmd.Execute() }` |
| `core/node_test.go` | `5` | `func TestMineBlockFeeDistribution(t *testing.T) {` |
| `cli/coin.go` | `11` | `func init() {` |
| `Tokens/syn3400.go` | `27` | `func NewForexRegistry() *ForexRegistry {` |
| `Tokens/syn3400.go` | `32` | `func (r *ForexRegistry) Register(base, quote string, rate float64) *ForexPair {` |
| `Tokens/syn3400.go` | `43` | `func (r *ForexRegistry) UpdateRate(pairID string, rate float64) error {` |
| `Tokens/syn3400.go` | `56` | `func (r *ForexRegistry) Get(pairID string) (*ForexPair, bool) {` |
| `Tokens/syn3400.go` | `68` | `func (r *ForexRegistry) List() []*ForexPair {` |
| `core/gas_test.go` | `5` | `func TestDefaultGasTable(t *testing.T) {` |
| `cli/network.go` | `13` | `func init() {` |
| `Tokens/syn2600.go` | `36` | `func NewInvestorRegistry() *InvestorRegistry {` |
| `Tokens/syn2600.go` | `41` | `func (r *InvestorRegistry) Issue(asset, owner string, shares uint64, expiry time.Time) *InvestorTokenMeta {` |
| `Tokens/syn2600.go` | `60` | `func (r *InvestorRegistry) Transfer(tokenID, newOwner string) error {` |
| `Tokens/syn2600.go` | `72` | `func (r *InvestorRegistry) RecordReturn(tokenID string, amount uint64) error {` |
| `Tokens/syn2600.go` | `84` | `func (r *InvestorRegistry) Get(tokenID string) (*InvestorTokenMeta, bool) {` |
| `Tokens/syn2600.go` | `97` | `func (r *InvestorRegistry) List() []*InvestorTokenMeta {` |
| `core/consensus_adaptive_management_test.go` | `5` | `func TestAdaptiveManagerAdjust(t *testing.T) {` |
| `cli/wallet.go` | `10` | `func init() {` |
| `core/syn1300.go` | `32` | `func NewSupplyChainRegistry() *SupplyChainRegistry {` |
| `core/syn1300.go` | `37` | `func (r *SupplyChainRegistry) Register(id, desc, owner, location string) (*SupplyChainAsset, error) {` |
| `core/syn1300.go` | `48` | `func (r *SupplyChainRegistry) Update(id, location, status, note string) error {` |
| `core/syn1300.go` | `60` | `func (r *SupplyChainRegistry) Get(id string) (*SupplyChainAsset, bool) {` |
| `cli/opcodes.go` | `10` | `func init() {` |
| `core/biometric_security_node_test.go` | `5` | `func TestBiometricSecurityNode(t *testing.T) {` |
| `cli/consensus.go` | `13` | `func init() {` |
| `Tokens/syn1000_index.go` | `12` | `func NewSYN1000Index() *SYN1000Index {` |
| `Tokens/syn1000_index.go` | `17` | `func (i *SYN1000Index) Create(name, symbol string, decimals uint8) TokenID {` |
| `Tokens/syn1000_index.go` | `25` | `func (i *SYN1000Index) Token(id TokenID) (*SYN1000Token, error) {` |
| `Tokens/syn1000_index.go` | `34` | `func (i *SYN1000Index) AddReserve(id TokenID, asset string, amount float64) error {` |
| `Tokens/syn1000_index.go` | `44` | `func (i *SYN1000Index) SetReservePrice(id TokenID, asset string, price float64) error {` |
| `Tokens/syn1000_index.go` | `54` | `func (i *SYN1000Index) TotalValue(id TokenID) (float64, error) {` |
| `core/faucet_test.go` | `8` | `func TestFaucet(t *testing.T) {` |
| `cli/gas.go` | `12` | `func init() {` |
| `Tokens/index.go` | `10` | `func NewRegistry() *Registry {` |
| `Tokens/index.go` | `15` | `func (r *Registry) NextID() TokenID {` |
| `Tokens/index.go` | `21` | `func (r *Registry) Register(t Token) {` |
| `Tokens/index.go` | `26` | `func (r *Registry) Get(id TokenID) (Token, bool) {` |
| `core/contracts.go` | `39` | `func NewContractRegistry(vm VirtualMachine) *ContractRegistry {` |
| `core/contracts.go` | `49` | `func CompileWASM(src []byte) ([]byte, string, error) {` |
| `core/contracts.go` | `60` | `func (r *ContractRegistry) Deploy(wasm []byte, manifest string, gasLimit uint64, owner string) (string, error) {` |
| `core/contracts.go` | `84` | `func (r *ContractRegistry) Invoke(addr, method string, args []byte, gasLimit uint64) ([]byte, uint64, error) {` |
| `core/contracts.go` | `101` | `func (r *ContractRegistry) List() []*Contract {` |
| `core/contracts.go` | `112` | `func (r *ContractRegistry) Get(addr string) (*Contract, bool) {` |
| `Tokens/syn10.go` | `11` | `func NewSYN10Token(id TokenID, name, symbol, issuer string, rate float64, decimals uint8) *SYN10Token {` |
| `Tokens/syn10.go` | `20` | `func (t *SYN10Token) SetExchangeRate(rate float64) {` |
| `Tokens/syn10.go` | `34` | `func (t *SYN10Token) Info() SYN10Info {` |
| `cli/node.go` | `13` | `func init() {` |
| `core/token_syn4900.go` | `34` | `func NewAgriculturalRegistry() *AgriculturalRegistry {` |
| `core/token_syn4900.go` | `39` | `func (r *AgriculturalRegistry) Register(id, assetType, owner, origin string, qty uint64, harvest, expiry time.Time, cert string) (*AgriculturalAsset, error) {` |
| `core/token_syn4900.go` | `49` | `func (r *AgriculturalRegistry) Transfer(id, newOwner string) error {` |
| `core/token_syn4900.go` | `60` | `func (r *AgriculturalRegistry) UpdateStatus(id, status string) error {` |
| `core/token_syn4900.go` | `71` | `func (r *AgriculturalRegistry) Get(id string) (*AgriculturalAsset, bool) {` |
| `Tokens/syn1100.go` | `18` | `func NewSYN1100Token() *SYN1100Token {` |
| `Tokens/syn1100.go` | `23` | `func (t *SYN1100Token) AddRecord(id TokenID, owner string, data []byte) error {` |
| `Tokens/syn1100.go` | `32` | `func (t *SYN1100Token) GrantAccess(id TokenID, grantee string) error {` |
| `Tokens/syn1100.go` | `42` | `func (t *SYN1100Token) RevokeAccess(id TokenID, grantee string) error {` |
| `Tokens/syn1100.go` | `52` | `func (t *SYN1100Token) GetRecord(id TokenID, caller string) ([]byte, error) {` |
| `cli/ledger.go` | `13` | `func init() {` |
| `core/charity.go` | `24` | `func (c CharityCategory) String() string {` |
| `core/charity.go` | `73` | `func NewCharityPool(lg *logrus.Logger, led StateRW, el electorate, genesis time.Time) *CharityPool {` |
| `core/charity.go` | `78` | `func (cp *CharityPool) Deposit(from Address, amount uint64) error {` |
| `core/charity.go` | `86` | `func (cp *CharityPool) Register(addr Address, name string, cat CharityCategory) error {` |
| `core/charity.go` | `97` | `func (cp *CharityPool) Vote(voter, charity Address) error {` |
| `core/charity.go` | `106` | `func (cp *CharityPool) Tick(ts time.Time) {}` |
| `core/charity.go` | `110` | `func (cp *CharityPool) Winners(cycle uint64) ([]Address, error) {` |
| `core/charity.go` | `115` | `func (cp *CharityPool) GetRegistration(cycle uint64, addr Address) (CharityRegistration, bool, error) {` |
| `core/charity.go` | `129` | `func mustJSON(v interface{}) []byte {` |
| `cli/transaction.go` | `23` | `func init() {` |
| `core/syn3600.go` | `15` | `func NewFuturesContract(underlying string, quantity, price uint64, expiration time.Time) *FuturesContract {` |
| `core/syn3600.go` | `20` | `func (f *FuturesContract) IsExpired(now time.Time) bool {` |
| `core/syn3600.go` | `25` | `func (f *FuturesContract) Settle(marketPrice uint64) int64 {` |
| `ai_model_management.go` | `27` | `func NewModelMarketplace() *ModelMarketplace {` |
| `ai_model_management.go` | `32` | `func (m *ModelMarketplace) AddListing(hash, cid, seller string, price uint64) string {` |
| `ai_model_management.go` | `49` | `func (m *ModelMarketplace) Get(id string) (ModelListing, bool) {` |
| `ai_model_management.go` | `57` | `func (m *ModelMarketplace) List() []ModelListing {` |
| `ai_model_management.go` | `70` | `func (m *ModelMarketplace) Update(id string, price uint64) error {` |
| `ai_model_management.go` | `84` | `func (m *ModelMarketplace) Remove(id, seller string) error {` |
| `core/address_zero_test.go` | `5` | `func TestIsZeroAddress(t *testing.T) {` |
| `core/authority_apply.go` | `30` | `func NewAuthorityApplicationManager(reg *AuthorityNodeRegistry, ttl time.Duration) *AuthorityApplicationManager {` |
| `core/authority_apply.go` | `40` | `func (m *AuthorityApplicationManager) Submit(candidate, role, desc string) string {` |
| `core/authority_apply.go` | `57` | `func (m *AuthorityApplicationManager) Vote(voter, id string, approve bool) error {` |
| `core/authority_apply.go` | `76` | `func (m *AuthorityApplicationManager) Finalize(id string) error {` |
| `core/authority_apply.go` | `93` | `func (m *AuthorityApplicationManager) Tick(now time.Time) {` |
| `core/authority_apply.go` | `102` | `func (m *AuthorityApplicationManager) Get(id string) (*AuthorityApplication, error) {` |
| `core/authority_apply.go` | `111` | `func (m *AuthorityApplicationManager) List() []*AuthorityApplication {` |
| `Tokens/syn1000.go` | `16` | `func NewSYN1000Token(id TokenID, name, symbol string, decimals uint8) *SYN1000Token {` |
| `Tokens/syn1000.go` | `24` | `func (t *SYN1000Token) AddReserve(asset string, amount float64) {` |
| `Tokens/syn1000.go` | `31` | `func (t *SYN1000Token) SetReservePrice(asset string, price float64) {` |
| `Tokens/syn1000.go` | `38` | `func (t *SYN1000Token) TotalReserveValue() float64 {` |
| `core/staking_node_test.go` | `5` | `func TestStakingNodeStakeAndUnstake(t *testing.T) {` |
| `core/staking_node_test.go` | `21` | `func TestStakingNodeTotal(t *testing.T) {` |
| `core/compliance_test.go` | `5` | `func TestComplianceServiceKYCAndRisk(t *testing.T) {` |
| `core/compliance_test.go` | `20` | `func TestComplianceServiceMonitorTransaction(t *testing.T) {` |
| `core/syn4200_token.go` | `23` | `func NewSYN4200Token() *SYN4200Token {` |
| `core/syn4200_token.go` | `28` | `func (t *SYN4200Token) Donate(symbol, from string, amount uint64, purpose string) {` |
| `core/syn4200_token.go` | `45` | `func (t *SYN4200Token) CampaignProgress(symbol string) (uint64, bool) {` |
| `core/syn4200_token.go` | `56` | `func (t *SYN4200Token) Campaign(symbol string) (*CharityCampaign, bool) {` |
| `core/elected_authority_node_test.go` | `8` | `func TestElectedAuthorityNode(t *testing.T) {` |
| `core/firewall.go` | `14` | `func NewFirewall() *Firewall {` |
| `core/firewall.go` | `23` | `func (f *Firewall) Allow(ip string) {` |
| `core/firewall.go` | `32` | `func (f *Firewall) Block(ip string) {` |
| `core/firewall.go` | `42` | `func (f *Firewall) IsAllowed(ip string) bool {` |
| `core/firewall.go` | `56` | `func (f *Firewall) Rules() (allowed []string, blocked []string) {` |
| `core/forensic_node.go` | `21` | `func NewForensicNode() *ForensicNode {` |
| `core/forensic_node.go` | `26` | `func (f *ForensicNode) RecordTransaction(tx nodes.TransactionLite) error {` |
| `core/forensic_node.go` | `34` | `func (f *ForensicNode) RecordNetworkTrace(trace nodes.NetworkTrace) error {` |
| `core/forensic_node.go` | `42` | `func (f *ForensicNode) Transactions() []nodes.TransactionLite {` |
| `core/forensic_node.go` | `51` | `func (f *ForensicNode) NetworkTraces() []nodes.NetworkTrace {` |
| `core/idwallet_registration.go` | `15` | `func NewIDRegistry() *IDRegistry {` |
| `core/idwallet_registration.go` | `20` | `func (r *IDRegistry) Register(addr, info string) error {` |
| `core/idwallet_registration.go` | `31` | `func (r *IDRegistry) Info(addr string) (string, bool) {` |
| `Tokens/syn2800.go` | `40` | `func NewLifePolicyRegistry() *LifePolicyRegistry {` |
| `Tokens/syn2800.go` | `45` | `func (r *LifePolicyRegistry) IssuePolicy(insured, beneficiary string, coverage, premium uint64, start, end time.Time) *LifePolicy {` |
| `Tokens/syn2800.go` | `65` | `func (r *LifePolicyRegistry) PayPremium(policyID string, amount uint64) error {` |
| `Tokens/syn2800.go` | `77` | `func (r *LifePolicyRegistry) FileClaim(policyID string, amount uint64) (*Claim, error) {` |
| `Tokens/syn2800.go` | `92` | `func (r *LifePolicyRegistry) GetPolicy(policyID string) (*LifePolicy, bool) {` |
| `Tokens/syn2800.go` | `105` | `func (r *LifePolicyRegistry) ListPolicies() []*LifePolicy {` |
| `core/blockchain_compression.go` | `18` | `func CompressLedger(l *Ledger) ([]byte, error) {` |
| `core/blockchain_compression.go` | `35` | `func DecompressLedger(data []byte) (*Ledger, error) {` |
| `core/blockchain_compression.go` | `53` | `func SaveCompressedSnapshot(l *Ledger, path string) error {` |
| `core/blockchain_compression.go` | `62` | `func LoadCompressedSnapshot(path string) (*Ledger, error) {` |
| `core/idwallet_registration_test.go` | `5` | `func TestIDRegistry(t *testing.T) {` |
| `core/warfare_node.go` | `27` | `func NewWarfareNode(base *Node) *WarfareNode {` |
| `core/warfare_node.go` | `32` | `func (w *WarfareNode) GetID() string {` |
| `core/warfare_node.go` | `41` | `func (w *WarfareNode) SecureCommand(cmd string) error {` |
| `core/warfare_node.go` | `49` | `func (w *WarfareNode) TrackLogistics(assetID, location, status string) {` |
| `core/warfare_node.go` | `63` | `func (w *WarfareNode) ShareTactical(info string) {` |
| `core/warfare_node.go` | `69` | `func (w *WarfareNode) Logistics() []LogisticsRecord {` |
| `Tokens/syn2900.go` | `42` | `func NewInsuranceRegistry() *InsuranceRegistry {` |
| `Tokens/syn2900.go` | `47` | `func (r *InsuranceRegistry) IssuePolicy(holder, coverage string, premium, payout, deductible, limit uint64, start, end time.Time) *InsurancePolicy {` |
| `Tokens/syn2900.go` | `69` | `func (r *InsuranceRegistry) FileClaim(policyID, desc string, amount uint64) (*ClaimRecord, error) {` |
| `Tokens/syn2900.go` | `84` | `func (r *InsuranceRegistry) GetPolicy(policyID string) (*InsurancePolicy, bool) {` |
| `Tokens/syn2900.go` | `97` | `func (r *InsuranceRegistry) ListPolicies() []*InsurancePolicy {` |
| `core/dao_access_control.go` | `23` | `func (d *DAO) AddMember(addr, role string) error {` |
| `core/dao_access_control.go` | `40` | `func (d *DAO) UpdateMemberRole(requester, addr, role string) error {` |
| `core/dao_access_control.go` | `57` | `func (d *DAO) RemoveMember(addr string) error {` |
| `core/dao_access_control.go` | `68` | `func (d *DAO) MemberRole(addr string) (string, bool) {` |
| `core/dao_access_control.go` | `76` | `func (d *DAO) IsMember(addr string) bool {` |
| `core/dao_access_control.go` | `84` | `func (d *DAO) IsAdmin(addr string) bool {` |
| `core/dao_access_control.go` | `91` | `func (d *DAO) MembersList() map[string]string {` |
| `core/audit_node_test.go` | `7` | `func (m *mockBootstrap) Start() error {` |
| `core/audit_node_test.go` | `12` | `func TestAuditNode_StartAndLog(t *testing.T) {` |
| `core/gateway_node.go` | `22` | `func NewGatewayNode(id nodes.Address, cfg GatewayConfig) *GatewayNode {` |
| `core/gateway_node.go` | `31` | `func (g *GatewayNode) RegisterEndpoint(name string, fn func([]byte) error) {` |
| `core/gateway_node.go` | `36` | `func (g *GatewayNode) Handle(name string, data []byte) error {` |
| `core/gateway_node.go` | `44` | `func (g *GatewayNode) RemoveEndpoint(name string) {` |
| `core/gateway_node.go` | `49` | `func (g *GatewayNode) Endpoints() []string {` |
| `core/loanpool_test.go` | `5` | `func TestLoanPoolProposalLifecycle(t *testing.T) {` |
| `data_operations.go` | `19` | `func NewDataFeed(id string) *DataFeed {` |
| `data_operations.go` | `24` | `func (f *DataFeed) Update(key, value string) {` |
| `data_operations.go` | `32` | `func (f *DataFeed) Get(key string) (string, bool) {` |
| `data_operations.go` | `40` | `func (f *DataFeed) Snapshot() map[string]string {` |
| `data_operations.go` | `51` | `func (f *DataFeed) LastUpdated() time.Time {` |
| `core/connection_pool_test.go` | `5` | `func TestConnectionPool(t *testing.T) {` |
| `core/kademlia_test.go` | `8` | `func TestKademliaStoreFind(t *testing.T) {` |
| `core/kademlia_test.go` | `17` | `func TestDistance(t *testing.T) {` |
| `data_resource_management.go` | `15` | `func NewDataResourceManager() *DataResourceManager {` |
| `data_resource_management.go` | `20` | `func (m *DataResourceManager) Put(key string, data []byte) {` |
| `data_resource_management.go` | `32` | `func (m *DataResourceManager) Get(key string) ([]byte, bool) {` |
| `data_resource_management.go` | `45` | `func (m *DataResourceManager) Delete(key string) {` |
| `data_resource_management.go` | `55` | `func (m *DataResourceManager) Usage() int64 {` |
| `core/ledger_test.go` | `5` | `func TestLedgerApplyTransaction(t *testing.T) {` |
| `core/ledger.go` | `22` | `func NewLedger(path ...string) *Ledger {` |
| `core/ledger.go` | `33` | `func (l *Ledger) replayWAL() {` |
| `core/ledger.go` | `53` | `func (l *Ledger) appendWAL(b *Block) {` |
| `core/ledger.go` | `66` | `func (l *Ledger) Head() (int, string) {` |
| `core/ledger.go` | `77` | `func (l *Ledger) GetBlock(height int) (*Block, bool) {` |
| `core/ledger.go` | `87` | `func (l *Ledger) AddBlock(b *Block) {` |
| `core/ledger.go` | `95` | `func (l *Ledger) GetBalance(addr string) uint64 {` |
| `core/ledger.go` | `102` | `func (l *Ledger) Credit(addr string, amount uint64) {` |
| `core/ledger.go` | `110` | `func (l *Ledger) Mint(addr string, amount uint64) {` |
| `core/ledger.go` | `116` | `func (l *Ledger) Transfer(from, to string, amount, fee uint64) error {` |
| `core/ledger.go` | `124` | `func (l *Ledger) ApplyTransaction(tx *Transaction) error {` |
| `core/opcode.go` | `71` | `func Catalogue() []OpcodeInfo {` |
| `core/opcode.go` | `83` | `func Opcodes() map[Opcode]string {` |
| `core/opcode.go` | `95` | `func Register(op Opcode, fn OpcodeFunc) {` |
| `core/opcode.go` | `105` | `func Dispatch(ctx OpContext, op Opcode) error {` |
| `core/opcode.go` | `121` | `func wrap(name string) OpcodeFunc {` |
| `core/opcode.go` | `1611` | `func init() {` |
| `core/opcode.go` | `1639` | `func (op Opcode) Hex() string { return fmt.Sprintf("0x%06X", uint32(op)) }` |
| `core/opcode.go` | `1642` | `func (op Opcode) Bytes() []byte {` |
| `core/opcode.go` | `1651` | `func (op Opcode) String() string { return op.Hex() }` |
| `core/opcode.go` | `1654` | `func ParseOpcode(b []byte) (Opcode, error) {` |
| `core/opcode.go` | `1662` | `func MustParseOpcode(b []byte) Opcode {` |
| `core/opcode.go` | `1672` | `func DebugDump() []string {` |
| `core/opcode.go` | `1691` | `func ToBytecode(fn string) ([]byte, error) {` |
| `core/opcode.go` | `1700` | `func HexDump(fn string) (string, error) {` |
| `core/authority_node_index.go` | `9` | `func NewAuthorityNodeIndex() *AuthorityNodeIndex {` |
| `core/authority_node_index.go` | `14` | `func (idx *AuthorityNodeIndex) Add(node *AuthorityNode) {` |
| `core/authority_node_index.go` | `22` | `func (idx *AuthorityNodeIndex) Get(addr string) (*AuthorityNode, bool) {` |
| `core/authority_node_index.go` | `28` | `func (idx *AuthorityNodeIndex) Remove(addr string) {` |
| `core/authority_node_index.go` | `33` | `func (idx *AuthorityNodeIndex) List() []*AuthorityNode {` |
| `core/historical_node_test.go` | `9` | `func TestHistoricalNode_ArchiveAndRetrieve(t *testing.T) {` |
| `core/historical_node_test.go` | `26` | `func TestHistoricalNode_Duplicate(t *testing.T) {` |
| `core/identity_verification.go` | `30` | `func NewIdentityService() *IdentityService {` |
| `core/identity_verification.go` | `38` | `func (s *IdentityService) Register(addr, name, dob, nationality string) error {` |
| `core/identity_verification.go` | `49` | `func (s *IdentityService) Verify(addr, method string) error {` |
| `core/identity_verification.go` | `61` | `func (s *IdentityService) Info(addr string) (IdentityInfo, bool) {` |
| `core/identity_verification.go` | `69` | `func (s *IdentityService) Logs(addr string) []VerificationLog {` |
| `core/authority_apply_test.go` | `8` | `func TestAuthorityApplication(t *testing.T) {` |
| `core/syn700.go` | `42` | `func NewIPRegistry() *IPRegistry {` |
| `core/syn700.go` | `47` | `func (r *IPRegistry) Register(tokenID, title, desc, creator, owner string) (*IPTokens, error) {` |
| `core/syn700.go` | `57` | `func (r *IPRegistry) CreateLicense(tokenID, licID, licType, licensee string, royalty uint64) error {` |
| `core/syn700.go` | `70` | `func (r *IPRegistry) RecordRoyalty(tokenID, licID, licensee string, amount uint64) error {` |
| `core/syn700.go` | `83` | `func (r *IPRegistry) Get(tokenID string) (*IPTokens, bool) {` |
| `core/address_zero.go` | `9` | `func IsZeroAddress(addr string) bool {` |
| `core/consensus_test.go` | `8` | `func TestThreshold(t *testing.T) {` |
| `core/consensus_test.go` | `15` | `func TestAdjustWeightsAndAvailability(t *testing.T) {` |
| `core/consensus_test.go` | `28` | `func TestTransitionThreshold(t *testing.T) {` |
| `core/consensus_test.go` | `37` | `func TestDifficultyAdjust(t *testing.T) {` |
| `core/consensus_test.go` | `44` | `func TestSelectValidator(t *testing.T) {` |
| `core/cross_chain_transactions.go` | `39` | `func NewCrossChainTxManager(l *Ledger) *CrossChainTxManager {` |
| `core/cross_chain_transactions.go` | `44` | `func (m *CrossChainTxManager) LockMint(bridgeID int, from, to, assetID string, amount uint64, proof string) (int, error) {` |
| `core/cross_chain_transactions.go` | `62` | `func (m *CrossChainTxManager) BurnRelease(bridgeID int, from, to, assetID string, amount uint64) (int, error) {` |
| `core/cross_chain_transactions.go` | `80` | `func (m *CrossChainTxManager) ListTransfers() []*CrossChainTransfer {` |
| `core/cross_chain_transactions.go` | `91` | `func (m *CrossChainTxManager) GetTransfer(id int) (*CrossChainTransfer, error) {` |
| `core/swarm_test.go` | `6` | `func TestSwarmBroadcast(t *testing.T) {` |
| `core/address.go` | `16` | `func StringToAddress(s string) (Address, error) {` |
| `core/address.go` | `28` | `func (a Address) Hex() string { return string(a) }` |
| `core/address.go` | `32` | `func (a Address) Bytes() []byte {` |
| `core/address.go` | `42` | `func (a Address) Short() string {` |
| `core/address.go` | `51` | `func (a Address) String() string { return a.Hex() }` |
| `core/cross_chain_connection_test.go` | `5` | `func TestChainConnectionManager(t *testing.T) {` |
| `core/virtual_machine_test.go` | `5` | `func TestSimpleVM(t *testing.T) {` |
| `core/audit_management.go` | `24` | `func NewAuditManager() *AuditManager {` |
| `core/audit_management.go` | `29` | `func (m *AuditManager) Log(address, event string, metadata map[string]string) error {` |
| `core/audit_management.go` | `43` | `func (m *AuditManager) List(address string) []AuditEntry {` |
| `core/syn1401.go` | `25` | `func NewInvestmentRegistry() *InvestmentRegistry {` |
| `core/syn1401.go` | `30` | `func (r *InvestmentRegistry) Issue(id, owner string, principal uint64, rate float64, maturity time.Time) (*InvestmentRecord, error) {` |
| `core/syn1401.go` | `40` | `func (r *InvestmentRegistry) Accrue(id string, now time.Time) (uint64, error) {` |
| `core/syn1401.go` | `56` | `func (r *InvestmentRegistry) Get(id string) (*InvestmentRecord, bool) {` |
| `core/biometric_test.go` | `5` | `func TestBiometricService(t *testing.T) {` |
| `core/syn3500_token.go` | `19` | `func NewSYN3500Token(name, symbol, issuer string, rate float64) *SYN3500Token {` |
| `core/syn3500_token.go` | `30` | `func (t *SYN3500Token) SetRate(rate float64) {` |
| `core/syn3500_token.go` | `37` | `func (t *SYN3500Token) Info() (string, string, float64) {` |
| `core/syn3500_token.go` | `44` | `func (t *SYN3500Token) Mint(to string, amt uint64) {` |
| `core/syn3500_token.go` | `51` | `func (t *SYN3500Token) Redeem(from string, amt uint64) error {` |
| `core/syn3500_token.go` | `63` | `func (t *SYN3500Token) BalanceOf(addr string) uint64 {` |
| `core/gas.go` | `13` | `func GasCost(op Opcode) uint64 {` |
| `core/gas.go` | `19` | `func initGasTable() {` |
| `core/consensus_adaptive_management.go` | `12` | `func NewAdaptiveManager(engine *SynnergyConsensus) *AdaptiveManager {` |
| `core/consensus_adaptive_management.go` | `18` | `func (am *AdaptiveManager) Adjust(demand, stake float64) ConsensusWeights {` |
| `core/consensus_adaptive_management.go` | `30` | `func (am *AdaptiveManager) Threshold(demand, stake float64) float64 {` |
| `core/consensus_adaptive_management.go` | `40` | `func (am *AdaptiveManager) Weights() ConsensusWeights {` |
| `core/contract_management_test.go` | `5` | `func TestContractManager(t *testing.T) {` |
| `core/base_node_test.go` | `9` | `func TestBaseNodeLifecycle(t *testing.T) {` |
| `core/loanpool_proposal.go` | `20` | `func NewLoanProposal(id uint64, creator, recipient, typ string, amount uint64, desc string, duration time.Duration) *LoanProposal {` |
| `core/loanpool_proposal.go` | `34` | `func (p *LoanProposal) Vote(voter string) {` |
| `core/loanpool_proposal.go` | `39` | `func (p *LoanProposal) VoteCount() int {` |
| `core/loanpool_proposal.go` | `44` | `func (p *LoanProposal) IsExpired(now time.Time) bool {` |
| `core/regulatory_node_test.go` | `5` | `func TestRegulatoryNode(t *testing.T) {` |
| `core/liquidity_pools_test.go` | `5` | `func TestLiquidityPoolLifecycle(t *testing.T) {` |
| `core/coin.go` | `22` | `func BlockReward(height uint64) uint64 {` |
| `core/coin.go` | `30` | `func CirculatingSupply(height uint64) uint64 {` |
| `core/coin.go` | `43` | `func RemainingSupply(height uint64) uint64 {` |
| `core/coin.go` | `53` | `func InitialPrice(C, R, M, V, T, E float64) float64 {` |
| `core/coin.go` | `59` | `func AlphaFactor(volatility, participation, economicStability, normalization float64) float64 {` |
| `core/coin.go` | `66` | `func MinimumStake(totalTx, currentReward, circulatingSupply, alpha float64) float64 {` |
| `core/coin.go` | `75` | `func LockupDuration(base, V, threshold, sigma float64) float64 {` |
| `core/coin.go` | `84` | `func PriceToSupplyRatio(price float64, height uint64) float64 {` |
| `core/block_test.go` | `10` | `func TestSubBlockCreationAndVerification(t *testing.T) {` |
| `core/block_test.go` | `21` | `func TestBlockHeaderHash(t *testing.T) {` |
| `core/zero_trust_data_channels_test.go` | `5` | `func TestZeroTrustEngine(t *testing.T) {` |
| `core/biometrics_auth.go` | `15` | `func NewBiometricsAuth() *BiometricsAuth {` |
| `core/biometrics_auth.go` | `35` | `func (b *BiometricsAuth) Enroll(addr string, biometric []byte, pub ed25519.PublicKey) error {` |
| `core/biometrics_auth.go` | `55` | `func (b *BiometricsAuth) Verify(addr string, biometric []byte, sig []byte) bool {` |
| `core/biometrics_auth.go` | `38` | `func (b *BiometricsAuth) Remove(addr string) {` |
| `core/mobile_mining_node.go` | `13` | `func NewMobileMiningNode(hashRate, powerLimit uint64) *MobileMiningNode {` |
| `core/mobile_mining_node.go` | `18` | `func (mm *MobileMiningNode) Start() {` |
| `core/mobile_mining_node.go` | `25` | `func (mm *MobileMiningNode) Stop() { mm.base.Stop() }` |
| `core/mobile_mining_node.go` | `28` | `func (mm *MobileMiningNode) IsMining() bool { return mm.base.IsMining() }` |
| `core/mobile_mining_node.go` | `31` | `func (mm *MobileMiningNode) Mine(data []byte) (string, error) {` |
| `core/mobile_mining_node.go` | `41` | `func (mm *MobileMiningNode) SetPowerLimit(limit uint64) {` |
| `core/mobile_mining_node.go` | `48` | `func (mm *MobileMiningNode) PowerLimit() uint64 {` |
| `core/syn2900_test.go` | `8` | `func TestTokenInsurancePolicy(t *testing.T) {` |
| `core/kademlia.go` | `18` | `func NewKademlia() *Kademlia {` |
| `core/kademlia.go` | `23` | `func (k *Kademlia) Store(key string, value []byte) {` |
| `core/kademlia.go` | `30` | `func (k *Kademlia) FindValue(key string) ([]byte, bool) {` |
| `core/kademlia.go` | `41` | `func Distance(a, b string) *big.Int {` |
| `core/central_banking_node.go` | `17` | `func NewCentralBankingNode(id, addr string, ledger *Ledger, policy string, token *tokens.SYN10Token) *CentralBankingNode {` |
| `core/central_banking_node.go` | `27` | `func (n *CentralBankingNode) UpdatePolicy(policy string) {` |
| `core/central_banking_node.go` | `30` | `func (n *CentralBankingNode) MintCBDC(to string, amount uint64) error {` |
| `core/central_banking_node.go` | `39` | `func (n *CentralBankingNode) Mint(to string, amount uint64) error {` |
| `core/ai_enhanced_contract.go` | `19` | `func NewAIContractRegistry(base *ContractRegistry) *AIContractRegistry {` |
| `core/ai_enhanced_contract.go` | `28` | `func (r *AIContractRegistry) DeployAIContract(wasm []byte, modelHash, manifest string, gasLimit uint64, owner string) (string, error) {` |
| `core/ai_enhanced_contract.go` | `41` | `func (r *AIContractRegistry) InvokeAIContract(addr string, input []byte, gasLimit uint64) ([]byte, uint64, error) {` |
| `core/ai_enhanced_contract.go` | `49` | `func (r *AIContractRegistry) ModelHash(addr string) (string, bool) {` |
| `core/replication.go` | `15` | `func NewReplicator(l *Ledger) *Replicator {` |
| `core/replication.go` | `20` | `func (r *Replicator) Start() {` |
| `core/replication.go` | `27` | `func (r *Replicator) Stop() {` |
| `core/replication.go` | `34` | `func (r *Replicator) Status() bool {` |
| `core/replication.go` | `42` | `func (r *Replicator) ReplicateBlock(hash string) bool {` |
| `core/quorum_tracker_test.go` | `5` | `func TestQuorumTracker(t *testing.T) {` |
| `core/network.go` | `17` | `func NewNetwork(auth *BiometricService) *Network {` |
| `core/network.go` | `29` | `func (n *Network) AddNode(node *Node) { n.nodes[node.ID] = node }` |
| `core/network.go` | `32` | `func (n *Network) AddRelay(node *Node) { n.relays[node.ID] = node }` |
| `core/network.go` | `35` | `func (n *Network) EnqueueTransaction(tx *Transaction) { n.queue <- tx }` |
| `core/network.go` | `40` | `func (n *Network) Broadcast(tx *Transaction, userID string, biometric []byte) error {` |
| `core/network.go` | `51` | `func (n *Network) processQueue() {` |
| `core/network.go` | `58` | `func (n *Network) broadcast(tx *Transaction) {` |
| `core/bank_institutional_node.go` | `19` | `func NewBankInstitutionalNode(id, addr string, ledger *Ledger) *BankInstitutionalNode {` |
| `core/bank_institutional_node.go` | `28` | `func (n *BankInstitutionalNode) RegisterInstitution(addr, name string, sig []byte, pubKey ed25519.PublicKey) error {` |
| `core/bank_institutional_node.go` | `47` | `func (n *BankInstitutionalNode) RemoveInstitution(addr, name string, sig []byte, pubKey ed25519.PublicKey) error {` |
| `core/bank_institutional_node.go` | `62` | `func (n *BankInstitutionalNode) ListInstitutions() []string {` |
| `core/bank_institutional_node.go` | `73` | `func (n *BankInstitutionalNode) IsRegistered(name string) bool {` |
| `core/bank_institutional_node.go` | `80` | `func (n *BankInstitutionalNode) MarshalJSON() ([]byte, error) {` |
| `core/light_node.go` | `12` | `func NewLightNode(id nodes.Address) *LightNode {` |
| `core/light_node.go` | `17` | `func (n *LightNode) AddHeader(h nodes.BlockHeader) { n.headers = append(n.headers, h) }` |
| `core/light_node.go` | `20` | `func (n *LightNode) LatestHeader() (nodes.BlockHeader, bool) {` |
| `core/light_node.go` | `28` | `func (n *LightNode) Headers() []nodes.BlockHeader {` |
| `core/government_authority_node_test.go` | `5` | `func TestGovernmentAuthorityNode(t *testing.T) {` |
| `core/token_syn130.go` | `40` | `func NewTangibleAssetRegistry() *TangibleAssetRegistry {` |
| `core/token_syn130.go` | `45` | `func (r *TangibleAssetRegistry) Register(id, owner, meta string, value uint64) (*TangibleAsset, error) {` |
| `core/token_syn130.go` | `55` | `func (r *TangibleAssetRegistry) UpdateValuation(id string, val uint64) error {` |
| `core/token_syn130.go` | `65` | `func (r *TangibleAssetRegistry) RecordSale(id, buyer string, price uint64) error {` |
| `core/token_syn130.go` | `76` | `func (r *TangibleAssetRegistry) StartLease(id, lessee string, payment uint64, start, end time.Time) error {` |
| `core/token_syn130.go` | `86` | `func (r *TangibleAssetRegistry) EndLease(id string) error {` |
| `core/token_syn130.go` | `98` | `func (r *TangibleAssetRegistry) Get(id string) (*TangibleAsset, bool) {` |
| `core/base_node.go` | `19` | `func NewBaseNode(id nodes.Address) *BaseNode {` |
| `core/base_node.go` | `27` | `func (n *BaseNode) ID() nodes.Address { return n.id }` |
| `core/base_node.go` | `30` | `func (n *BaseNode) Start() error {` |
| `core/base_node.go` | `41` | `func (n *BaseNode) Stop() error {` |
| `core/base_node.go` | `52` | `func (n *BaseNode) IsRunning() bool {` |
| `core/base_node.go` | `59` | `func (n *BaseNode) Peers() []nodes.Address {` |
| `core/base_node.go` | `72` | `func (n *BaseNode) DialSeed(addr nodes.Address) error {` |
| `core/base_node.go` | `84` | `func (n *BaseNode) DialSeedSigned(addr nodes.Address, sig []byte, pub ed25519.PublicKey) error {` |
| `core/cross_chain_test.go` | `5` | `func TestBridgeRegistry(t *testing.T) {` |
| `core/consensus_difficulty_test.go` | `5` | `func TestDifficultyManager(t *testing.T) {` |
| `core/biometric_security_node.go` | `14` | `func NewBiometricSecurityNode(base *Node, auth *BiometricsAuth) *BiometricSecurityNode {` |
| `core/biometric_security_node.go` | `22` | `func (b *BiometricSecurityNode) GetID() string {` |
| `core/biometric_security_node.go` | `30` | `func (b *BiometricSecurityNode) Enroll(addr string, biometric []byte) {` |
| `core/biometric_security_node.go` | `35` | `func (b *BiometricSecurityNode) Authenticate(addr string, biometric []byte) bool {` |
| `core/biometric_security_node.go` | `41` | `func (b *BiometricSecurityNode) SecureAddTransaction(addr string, biometric []byte, tx *Transaction) error {` |
| `core/rpc_webrtc_test.go` | `5` | `func TestWebRTCRPC(t *testing.T) {` |
| `core/consensus_difficulty.go` | `17` | `func NewDifficultyManager(engine *SynnergyConsensus, window int, initial, target float64) *DifficultyManager {` |
| `core/consensus_difficulty.go` | `26` | `func (dm *DifficultyManager) AddSample(duration float64) float64 {` |
| `core/consensus_difficulty.go` | `45` | `func (dm *DifficultyManager) Difficulty() float64 {` |
| `core/coin_test.go` | `8` | `func TestBlockRewardHalving(t *testing.T) {` |
| `core/coin_test.go` | `17` | `func TestCirculatingAndRemainingSupply(t *testing.T) {` |
| `core/coin_test.go` | `28` | `func TestEconomicHelpers(t *testing.T) {` |
| `core/transaction_control_test.go` | `8` | `func TestScheduleAndCancel(t *testing.T) {` |
| `core/transaction_control_test.go` | `20` | `func TestReverseTransaction(t *testing.T) {` |
| `core/transaction_control_test.go` | `35` | `func TestConvertToPrivate(t *testing.T) {` |
| `core/transaction_control_test.go` | `51` | `func TestReceiptStore(t *testing.T) {` |
| `core/bank_nodes_test.go` | `5` | `func TestBankNodes(t *testing.T) {` |
| `core/syn3900.go` | `21` | `func NewBenefitRegistry() *BenefitRegistry {` |
| `core/syn3900.go` | `26` | `func (r *BenefitRegistry) RegisterBenefit(recipient, program string, amount uint64) uint64 {` |
| `core/syn3900.go` | `34` | `func (r *BenefitRegistry) Claim(id uint64) error {` |
| `core/syn3900.go` | `47` | `func (r *BenefitRegistry) GetBenefit(id uint64) (*BenefitRecord, bool) {` |
| `core/compliance.go` | `46` | `func NewComplianceService() *ComplianceService {` |
| `core/compliance.go` | `56` | `func (s *ComplianceService) ValidateKYC(address string, kycData []byte) (string, error) {` |
| `core/compliance.go` | `73` | `func (s *ComplianceService) EraseKYC(address string) error {` |
| `core/compliance.go` | `92` | `func (s *ComplianceService) RecordFraud(address string, severity int) error {` |
| `core/compliance.go` | `113` | `func (s *ComplianceService) RiskScore(address string) int {` |
| `core/compliance.go` | `121` | `func (s *ComplianceService) AuditTrail(address string) []AuditEntry {` |
| `core/compliance.go` | `107` | `func (s *ComplianceService) MonitorTransaction(tx ComplianceTransaction, threshold float64) bool {` |
| `core/compliance.go` | `119` | `func (s *ComplianceService) VerifyZKP(blob []byte, commitmentHex, proofHex string) bool {` |
| `core/compliance.go` | `124` | `func (s *ComplianceService) appendAudit(addr, event string, metadata map[string]string) {` |
| `core/cross_chain_agnostic_protocols.go` | `19` | `func NewProtocolRegistry() *ProtocolRegistry {` |
| `core/cross_chain_agnostic_protocols.go` | `24` | `func (r *ProtocolRegistry) Register(name string) int {` |
| `core/cross_chain_agnostic_protocols.go` | `34` | `func (r *ProtocolRegistry) List() []ProtocolDefinition {` |
| `core/cross_chain_agnostic_protocols.go` | `45` | `func (r *ProtocolRegistry) Get(id int) (ProtocolDefinition, bool) {` |
| `core/cross_chain_contracts_test.go` | `5` | `func TestCrossChainRegistry(t *testing.T) {` |
| `core/audit_management_test.go` | `5` | `func TestAuditManager_LogAndList(t *testing.T) {` |
| `core/dao.go` | `24` | `func NewDAOManager() *DAOManager {` |
| `core/dao.go` | `29` | `func (m *DAOManager) Create(name, creator string) *DAO {` |
| `core/dao.go` | `38` | `func (m *DAOManager) Join(id, addr string) error {` |
| `core/dao.go` | `51` | `func (m *DAOManager) Leave(id, addr string) error {` |
| `core/dao.go` | `61` | `func (m *DAOManager) Info(id string) (*DAO, error) {` |
| `core/dao.go` | `70` | `func (m *DAOManager) List() []*DAO {` |
| `core/dao_staking.go` | `21` | `func NewDAOStaking(mgr *DAOManager) *DAOStaking {` |
| `core/dao_staking.go` | `26` | `func (s *DAOStaking) Stake(daoID, addr string, amount uint64) error {` |
| `core/dao_staking.go` | `44` | `func (s *DAOStaking) Unstake(daoID, addr string, amount uint64) error {` |
| `core/dao_staking.go` | `63` | `func (s *DAOStaking) Balance(daoID, addr string) uint64 {` |
| `core/dao_staking.go` | `70` | `func (s *DAOStaking) TotalStaked(daoID string) uint64 {` |
| `core/dao_staking_test.go` | `7` | `func TestDAOStaking(t *testing.T) {` |
| `core/dao_quadratic_voting.go` | `9` | `func QuadraticWeight(tokens uint64) uint64 {` |
| `core/dao_quadratic_voting.go` | `15` | `func (pm *ProposalManager) CastQuadraticVote(dao *DAO, id, voter string, tokens uint64, support bool) error {` |
| `core/wallet.go` | `19` | `func NewWallet() (*Wallet, error) {` |
| `core/wallet.go` | `29` | `func (w *Wallet) Sign(tx *Transaction) ([]byte, error) {` |
| `core/wallet.go` | `45` | `func VerifySignature(tx *Transaction, sig []byte, pub *ecdsa.PublicKey) bool {` |
| `core/node_adapter.go` | `12` | `func NewNodeAdapter(n *Node) *NodeAdapter {` |
| `core/syn223_token.go` | `20` | `func NewSYN223Token(name, symbol, owner string, supply uint64) *SYN223Token {` |
| `core/syn223_token.go` | `34` | `func (t *SYN223Token) AddToWhitelist(addr string) {` |
| `core/syn223_token.go` | `41` | `func (t *SYN223Token) RemoveFromWhitelist(addr string) {` |
| `core/syn223_token.go` | `48` | `func (t *SYN223Token) AddToBlacklist(addr string) {` |
| `core/syn223_token.go` | `55` | `func (t *SYN223Token) RemoveFromBlacklist(addr string) {` |
| `core/syn223_token.go` | `62` | `func (t *SYN223Token) Transfer(from, to string, amount uint64) error {` |
| `core/syn223_token.go` | `82` | `func (t *SYN223Token) BalanceOf(addr string) uint64 {` |
| `core/snvm.go` | `13` | `func NewSNVM() *SNVM { return &SNVM{} }` |
| `core/snvm.go` | `19` | `func (vm *SNVM) Execute(tx *Transaction) (int64, error) {` |
| `core/cross_chain_bridge_test.go` | `5` | `func TestBridgeManager(t *testing.T) {` |
| `core/stake_penalty_test.go` | `5` | `func TestStakePenaltyManager(t *testing.T) {` |
| `core/mobile_mining_node_test.go` | `5` | `func TestMobileMiningNode(t *testing.T) {` |
| `core/consensus_start.go` | `16` | `func NewConsensusService(n *Node) *ConsensusService {` |
| `core/consensus_start.go` | `22` | `func (s *ConsensusService) Start(interval time.Duration) {` |
| `core/consensus_start.go` | `41` | `func (s *ConsensusService) Stop() {` |
| `core/consensus_start.go` | `49` | `func (s *ConsensusService) Info() (height int, running bool) {` |
| `core/transaction_test.go` | `5` | `func TestNewTransactionAndHash(t *testing.T) {` |
| `core/transaction_test.go` | `20` | `func TestAttachBiometric(t *testing.T) {` |
| `core/consensus_specific_test.go` | `5` | `func TestConsensusSwitcher(t *testing.T) {` |
| `core/access_control_test.go` | `5` | `func TestAccessController(t *testing.T) {` |
| `core/syn131_token.go` | `20` | `func NewSYN131Registry() *SYN131Registry {` |
| `core/syn131_token.go` | `25` | `func (r *SYN131Registry) Create(id, name, symbol, owner string, valuation uint64) (*SYN131Token, error) {` |
| `core/syn131_token.go` | `35` | `func (r *SYN131Registry) UpdateValuation(id string, val uint64) error {` |
| `core/syn131_token.go` | `45` | `func (r *SYN131Registry) Get(id string) (*SYN131Token, bool) {` |
| `core/peer_management.go` | `13` | `func NewPeerManager() *PeerManager {` |
| `core/peer_management.go` | `18` | `func (pm *PeerManager) AddPeer(id, addr string) {` |
| `core/peer_management.go` | `25` | `func (pm *PeerManager) RemovePeer(id string) {` |
| `core/peer_management.go` | `32` | `func (pm *PeerManager) GetPeer(id string) (string, bool) {` |
| `core/peer_management.go` | `40` | `func (pm *PeerManager) ListPeers() []string {` |
| `core/consensus.go` | `38` | `func NewSynnergyConsensus() *SynnergyConsensus {` |
| `core/consensus.go` | `55` | `func (sc *SynnergyConsensus) Threshold(D, S float64) float64 {` |
| `core/consensus.go` | `63` | `func (sc *SynnergyConsensus) AdjustWeights(D, S float64) {` |
| `core/consensus.go` | `89` | `func (sc *SynnergyConsensus) Tload(D float64) float64 {` |
| `core/consensus.go` | `97` | `func (sc *SynnergyConsensus) Tsecurity(threat float64) float64 {` |
| `core/consensus.go` | `105` | `func (sc *SynnergyConsensus) Tstake(S float64) float64 {` |
| `core/consensus.go` | `114` | `func (sc *SynnergyConsensus) TransitionThreshold(D, threat, S float64) float64 {` |
| `core/consensus.go` | `120` | `func (sc *SynnergyConsensus) DifficultyAdjust(oldDifficulty, actualTime, expectedTime float64) float64 {` |
| `core/consensus.go` | `128` | `func (sc *SynnergyConsensus) SetAvailability(pow, pos, poh bool) {` |
| `core/consensus.go` | `135` | `func (sc *SynnergyConsensus) SetPoWRewards(enabled bool) {` |
| `core/consensus.go` | `141` | `func (sc *SynnergyConsensus) SelectValidator(stakes map[string]uint64) string {` |
| `core/consensus.go` | `162` | `func (sc *SynnergyConsensus) ValidateSubBlock(sb *SubBlock) bool {` |
| `core/consensus.go` | `169` | `func (sc *SynnergyConsensus) MineBlock(b *Block, difficulty uint8) {` |
| `core/consensus.go` | `184` | `func clamp(v, min, max float64) float64 {` |
| `core/contracts_test.go` | `5` | `func TestContractRegistry(t *testing.T) {` |
| `core/warfare_node_test.go` | `8` | `func TestWarfareNode(t *testing.T) {` |
| `core/gas_table.go` | `20` | `func DefaultGasTable() GasTable {` |
| `core/gas_table.go` | `36` | `func parseGasGuide() map[string]uint64 {` |
| `core/gas_table.go` | `65` | `func SetGasCost(op Opcode, cost uint64) {` |
| `core/gas_table.go` | `71` | `func GasTableSnapshot() GasTable {` |
| `core/syn1600_test.go` | `5` | `func TestMusicToken(t *testing.T) {` |
| `core/transaction_control.go` | `24` | `func ScheduleTransaction(tx *Transaction, exec time.Time) *ScheduledTransaction {` |
| `core/transaction_control.go` | `30` | `func CancelTransaction(st *ScheduledTransaction) bool {` |
| `core/transaction_control.go` | `40` | `func ReverseTransaction(l *Ledger, tx *Transaction) error {` |
| `core/transaction_control.go` | `53` | `func ConvertToPrivate(tx *Transaction, key []byte) (*PrivateTransaction, error) {` |
| `core/transaction_control.go` | `74` | `func (pt *PrivateTransaction) Decrypt(key []byte) (*Transaction, error) {` |
| `core/transaction_control.go` | `98` | `func GenerateReceipt(txID, status, details string) Receipt {` |
| `core/transaction_control.go` | `109` | `func NewReceiptStore() *ReceiptStore {` |
| `core/transaction_control.go` | `114` | `func (rs *ReceiptStore) Store(r Receipt) {` |
| `core/transaction_control.go` | `121` | `func (rs *ReceiptStore) Get(id string) (Receipt, bool) {` |
| `core/transaction_control.go` | `129` | `func (rs *ReceiptStore) Search(keyword string) []Receipt {` |
| `core/cross_chain_agnostic_protocols_test.go` | `5` | `func TestProtocolRegistry(t *testing.T) {` |
| `core/validator_node.go` | `12` | `func NewValidatorNode(id, addr string, ledger *Ledger, minStake uint64, quorum int) *ValidatorNode {` |
| `core/validator_node.go` | `22` | `func (vn *ValidatorNode) AddValidator(addr string, stake uint64) error {` |
| `core/validator_node.go` | `31` | `func (vn *ValidatorNode) RemoveValidator(addr string) {` |
| `core/validator_node.go` | `37` | `func (vn *ValidatorNode) SlashValidator(addr string) {` |
| `core/validator_node.go` | `44` | `func (vn *ValidatorNode) HasQuorum() bool {` |
| `core/high_availability.go` | `20` | `func NewFailoverManager(primary string, timeout time.Duration) *FailoverManager {` |
| `core/high_availability.go` | `29` | `func (m *FailoverManager) RegisterBackup(id string) {` |
| `core/high_availability.go` | `36` | `func (m *FailoverManager) Heartbeat(id string) {` |
| `core/high_availability.go` | `45` | `func (m *FailoverManager) Active() string {` |
| `core/syn4700.go` | `44` | `func NewLegalToken(id, name, symbol, docType, hash, owner string, expiry time.Time, supply uint64, parties []string) *LegalToken {` |
| `core/syn4700.go` | `63` | `func (t *LegalToken) Sign(party, sig string) error {` |
| `core/syn4700.go` | `74` | `func (t *LegalToken) RevokeSignature(party string) {` |
| `core/syn4700.go` | `81` | `func (t *LegalToken) UpdateStatus(status LegalTokenStatus) {` |
| `core/syn4700.go` | `88` | `func (t *LegalToken) Dispute(action, result string) {` |
| `core/syn4700.go` | `95` | `func (t *LegalToken) partyExists(party string) bool {` |
| `core/syn4700.go` | `111` | `func NewLegalTokenRegistry() *LegalTokenRegistry {` |
| `core/syn4700.go` | `116` | `func (r *LegalTokenRegistry) Add(t *LegalToken) {` |
| `core/syn4700.go` | `123` | `func (r *LegalTokenRegistry) Get(id string) (*LegalToken, bool) {` |
| `core/syn4700.go` | `131` | `func (r *LegalTokenRegistry) Remove(id string) {` |
| `core/syn4700.go` | `138` | `func (r *LegalTokenRegistry) List() []*LegalToken {` |
| `core/system_health_logging.go` | `19` | `func NewSystemHealthLogger() *SystemHealthLogger {` |
| `core/system_health_logging.go` | `25` | `func (l *SystemHealthLogger) Collect(peerCount int, height uint64) watchtower.Metrics {` |
| `core/system_health_logging.go` | `44` | `func (l *SystemHealthLogger) Snapshot() watchtower.Metrics {` |
| `core/consensus_specific_node.go` | `10` | `func NewConsensusSpecificNode(mode ConsensusMode, id, addr string, ledger *Ledger) *ConsensusSpecificNode {` |
| `core/consensus_specific_node.go` | `19` | `func (n *ConsensusSpecificNode) configure() {` |
| `core/syn1700_token_test.go` | `5` | `func TestEventTickets(t *testing.T) {` |
| `core/syn800_token.go` | `25` | `func NewAssetRegistry() *AssetRegistry {` |
| `core/syn800_token.go` | `30` | `func (r *AssetRegistry) Register(id, desc string, valuation uint64, loc, typ, cert string) (*AssetMetadata, error) {` |
| `core/syn800_token.go` | `40` | `func (r *AssetRegistry) UpdateValuation(id string, valuation uint64) error {` |
| `core/syn800_token.go` | `51` | `func (r *AssetRegistry) Get(id string) (*AssetMetadata, bool) {` |
| `core/compliance_management.go` | `17` | `func NewComplianceManager() *ComplianceManager {` |
| `core/compliance_management.go` | `27` | `func (m *ComplianceManager) Suspend(addr string) error {` |
| `core/compliance_management.go` | `43` | `func (m *ComplianceManager) Resume(addr string) error {` |
| `core/compliance_management.go` | `56` | `func (m *ComplianceManager) Whitelist(addr string) error {` |
| `core/compliance_management.go` | `72` | `func (m *ComplianceManager) Unwhitelist(addr string) error {` |
| `core/compliance_management.go` | `84` | `func (m *ComplianceManager) Status(addr string) (suspended, whitelisted bool) {` |
| `core/compliance_management.go` | `94` | `func (m *ComplianceManager) ReviewTransaction(tx Transaction) error {` |
| `core/cross_consensus_scaling_networks_test.go` | `5` | `func TestConsensusNetworkManager(t *testing.T) {` |
| `core/historical_node.go` | `21` | `func NewHistoricalNode() *HistoricalNode {` |
| `core/historical_node.go` | `29` | `func (h *HistoricalNode) ArchiveBlock(summary nodes.BlockSummary) error {` |
| `core/historical_node.go` | `44` | `func (h *HistoricalNode) GetBlockByHeight(height uint64) (nodes.BlockSummary, bool) {` |
| `core/historical_node.go` | `52` | `func (h *HistoricalNode) GetBlockByHash(hash string) (nodes.BlockSummary, bool) {` |
| `core/historical_node.go` | `60` | `func (h *HistoricalNode) TotalBlocks() int {` |
| `core/private_transactions.go` | `13` | `func Encrypt(key, plaintext []byte) ([]byte, error) {` |
| `core/private_transactions.go` | `31` | `func Decrypt(key, data []byte) ([]byte, error) {` |
| `core/private_transactions.go` | `62` | `func NewPrivateTxManager() *PrivateTxManager {` |
| `core/private_transactions.go` | `67` | `func (m *PrivateTxManager) Send(tx PrivateTransaction) {` |
| `core/private_transactions.go` | `74` | `func (m *PrivateTxManager) List() []PrivateTransaction {` |
| `core/liquidity_pools.go` | `21` | `func NewLiquidityPool(id, tokenA, tokenB string, feeBps uint16) *LiquidityPool {` |
| `core/liquidity_pools.go` | `32` | `func (p *LiquidityPool) AddLiquidity(provider string, amtA, amtB uint64) (uint64, error) {` |
| `core/liquidity_pools.go` | `57` | `func (p *LiquidityPool) RemoveLiquidity(provider string, lpTokens uint64) (uint64, uint64, error) {` |
| `core/liquidity_pools.go` | `72` | `func (p *LiquidityPool) Swap(tokenIn string, amtIn, minOut uint64) (uint64, error) {` |
| `core/liquidity_pools.go` | `102` | `func sqrt(n uint64) uint64 {` |
| `core/liquidity_pools.go` | `117` | `func NewLiquidityPoolRegistry() *LiquidityPoolRegistry {` |
| `core/liquidity_pools.go` | `122` | `func (r *LiquidityPoolRegistry) Create(id, tokenA, tokenB string, feeBps uint16) (*LiquidityPool, error) {` |
| `core/liquidity_pools.go` | `132` | `func (r *LiquidityPoolRegistry) Get(id string) (*LiquidityPool, bool) {` |
| `core/liquidity_pools.go` | `138` | `func (r *LiquidityPoolRegistry) List() []*LiquidityPool {` |
| `core/blockchain_compression_test.go` | `5` | `func TestLedgerCompressionRoundTrip(t *testing.T) {` |
| `core/immutability_enforcement.go` | `12` | `func NewImmutabilityEnforcer(genesis *Block) *ImmutabilityEnforcer {` |
| `core/immutability_enforcement.go` | `22` | `func (i *ImmutabilityEnforcer) CheckLedger(l *Ledger) error {` |
| `core/faucet.go` | `19` | `func NewFaucet(balance, amount uint64, cooldown time.Duration) *Faucet {` |
| `core/faucet.go` | `29` | `func (f *Faucet) Request(addr string) (uint64, error) {` |
| `core/faucet.go` | `45` | `func (f *Faucet) Balance() uint64 {` |
| `core/faucet.go` | `52` | `func (f *Faucet) UpdateConfig(amount uint64, cooldown time.Duration) {` |
| `core/custodial_node.go` | `10` | `func NewCustodialNode(id, addr string, ledger *Ledger) *CustodialNode {` |
| `core/custodial_node.go` | `18` | `func (n *CustodialNode) Custody(user string, amount uint64) {` |
| `core/custodial_node.go` | `23` | `func (n *CustodialNode) Release(user string, amount uint64) bool {` |
| `core/syn3600_test.go` | `8` | `func TestFuturesContract(t *testing.T) {` |
| `core/blockchain_synchronization_test.go` | `5` | `func TestSyncManagerLifecycle(t *testing.T) {` |
| `core/network_test.go` | `10` | `func TestNetworkBroadcast(t *testing.T) {` |
| `core/nat_traversal_test.go` | `5` | `func TestNATManager(t *testing.T) {` |
| `core/initialization_replication_test.go` | `5` | `func TestInitServiceStartStop(t *testing.T) {` |
| `core/wallet_test.go` | `5` | `func TestWalletSignAndVerify(t *testing.T) {` |
| `core/syn1600.go` | `16` | `func NewMusicToken(title, artist, album string) *MusicToken {` |
| `core/syn1600.go` | `26` | `func (m *MusicToken) Info() (string, string, string) {` |
| `core/syn1600.go` | `31` | `func (m *MusicToken) Update(title, artist, album string) {` |
| `core/syn1600.go` | `44` | `func (m *MusicToken) SetRoyaltyShare(addr string, share uint64) {` |
| `core/syn1600.go` | `51` | `func (m *MusicToken) Distribute(amount uint64) (map[string]uint64, error) {` |
| `core/mining_node_test.go` | `5` | `func TestMiningNode(t *testing.T) {` |
| `core/snvm_test.go` | `5` | `func TestSNVMArithmetic(t *testing.T) {` |
| `core/snvm_test.go` | `21` | `func TestSNVMDivideByZero(t *testing.T) {` |
| `core/validator_node_test.go` | `5` | `func TestValidatorNode(t *testing.T) {` |
| `core/identity_verification_test.go` | `5` | `func TestIdentityService(t *testing.T) {` |
| `core/authority_nodes.go` | `23` | `func NewAuthorityNodeRegistry() *AuthorityNodeRegistry {` |
| `core/authority_nodes.go` | `28` | `func (r *AuthorityNodeRegistry) Register(addr, role string) (*AuthorityNode, error) {` |
| `core/authority_nodes.go` | `38` | `func (r *AuthorityNodeRegistry) Vote(voterAddr, candidateAddr string) error {` |
| `core/authority_nodes.go` | `48` | `func (r *AuthorityNodeRegistry) Electorate(size int) []string {` |
| `core/authority_nodes.go` | `66` | `func (r *AuthorityNodeRegistry) IsAuthorityNode(addr string) bool {` |
| `core/authority_nodes.go` | `72` | `func (r *AuthorityNodeRegistry) Info(addr string) (*AuthorityNode, error) {` |
| `core/authority_nodes.go` | `81` | `func (r *AuthorityNodeRegistry) List() []*AuthorityNode {` |
| `core/authority_nodes.go` | `86` | `func (r *AuthorityNodeRegistry) Deregister(addr string) {` |
| `core/consensus_validator_management.go` | `17` | `func NewValidatorManager(minStake uint64) *ValidatorManager {` |
| `core/consensus_validator_management.go` | `26` | `func (vm *ValidatorManager) Add(addr string, stake uint64) error {` |
| `core/consensus_validator_management.go` | `38` | `func (vm *ValidatorManager) Remove(addr string) {` |
| `core/consensus_validator_management.go` | `46` | `func (vm *ValidatorManager) Slash(addr string) {` |
| `core/consensus_validator_management.go` | `56` | `func (vm *ValidatorManager) Eligible() map[string]uint64 {` |
| `core/consensus_validator_management.go` | `69` | `func (vm *ValidatorManager) Stake(addr string) uint64 {` |
| `core/security_test.go` | `5` | `func TestSetStakeEnforcesMinimum(t *testing.T) {` |
| `core/security_test.go` | `16` | `func TestSlashingAndRehabilitation(t *testing.T) {` |
| `core/security_test.go` | `33` | `func TestEligibleStakesExcludesSlashed(t *testing.T) {` |
| `core/security_test.go` | `48` | `func TestSubBlockSignature(t *testing.T) {` |
| `core/mining_node.go` | `19` | `func NewMiningNode(hashRate uint64) *MiningNode {` |
| `core/mining_node.go` | `24` | `func (mn *MiningNode) Start() {` |
| `core/mining_node.go` | `31` | `func (mn *MiningNode) Stop() {` |
| `core/mining_node.go` | `38` | `func (mn *MiningNode) IsMining() bool {` |
| `core/mining_node.go` | `45` | `func (mn *MiningNode) Mine(data []byte) (string, error) {` |
| `core/mining_node.go` | `60` | `func (mn *MiningNode) HashRateHint() uint64 {` |
| `core/syn3800.go` | `22` | `func NewGrantRegistry() *GrantRegistry {` |
| `core/syn3800.go` | `27` | `func (r *GrantRegistry) CreateGrant(beneficiary, name string, amount uint64) uint64 {` |
| `core/syn3800.go` | `35` | `func (r *GrantRegistry) Disburse(id uint64, amount uint64, note string) error {` |
| `core/syn3800.go` | `51` | `func (r *GrantRegistry) GetGrant(id uint64) (*GrantRecord, bool) {` |
| `core/syn3800.go` | `57` | `func (r *GrantRegistry) ListGrants() []*GrantRecord {` |
| `core/peer_management_test.go` | `5` | `func TestPeerManager(t *testing.T) {` |
| `core/consensus_specific.go` | `21` | `func NewConsensusSwitcher(mode ConsensusMode) *ConsensusSwitcher {` |
| `core/consensus_specific.go` | `27` | `func (cs *ConsensusSwitcher) Evaluate(sc *SynnergyConsensus) ConsensusMode {` |
| `core/consensus_specific.go` | `49` | `func (cs *ConsensusSwitcher) Mode() ConsensusMode {` |
| `core/contracts_opcodes.go` | `22` | `func opcodeByName(name string) Opcode {` |
| `core/watchtower_node.go` | `27` | `func NewWatchtowerNode(id string, logger *log.Logger) *Watchtower {` |
| `core/watchtower_node.go` | `37` | `func (w *Watchtower) ID() string { return w.id }` |
| `core/watchtower_node.go` | `40` | `func (w *Watchtower) Start(ctx context.Context) error {` |
| `core/watchtower_node.go` | `54` | `func (w *Watchtower) monitorLoop(ctx context.Context) {` |
| `core/watchtower_node.go` | `71` | `func (w *Watchtower) Stop() error {` |
| `core/watchtower_node.go` | `83` | `func (w *Watchtower) ReportFork(height uint64, hash string) {` |
| `core/watchtower_node.go` | `90` | `func (w *Watchtower) Metrics() watchtower.Metrics {` |
| `core/watchtower_node.go` | `95` | `func (w *Watchtower) Firewall() *Firewall { return w.firewall }` |
| `core/staking_node.go` | `12` | `func NewStakingNode() *StakingNode {` |
| `core/staking_node.go` | `17` | `func (s *StakingNode) Stake(addr string, amt uint64) {` |
| `core/staking_node.go` | `25` | `func (s *StakingNode) Unstake(addr string, amt uint64) {` |
| `core/staking_node.go` | `37` | `func (s *StakingNode) Balance(addr string) uint64 {` |
| `core/staking_node.go` | `44` | `func (s *StakingNode) TotalStaked() uint64 {` |
| `core/regulatory_management.go` | `23` | `func NewRegulatoryManager() *RegulatoryManager {` |
| `core/regulatory_management.go` | `28` | `func (m *RegulatoryManager) AddRegulation(reg Regulation) error {` |
| `core/regulatory_management.go` | `39` | `func (m *RegulatoryManager) RemoveRegulation(id string) {` |
| `core/regulatory_management.go` | `46` | `func (m *RegulatoryManager) GetRegulation(id string) (Regulation, bool) {` |
| `core/regulatory_management.go` | `54` | `func (m *RegulatoryManager) ListRegulations() []Regulation {` |
| `core/regulatory_management.go` | `65` | `func (m *RegulatoryManager) EvaluateTransaction(tx Transaction) []string {` |
| `core/vm_sandbox_management.go` | `29` | `func NewSandboxManager() *SandboxManager {` |
| `core/vm_sandbox_management.go` | `34` | `func (m *SandboxManager) StartSandbox(id, contractAddr string, gasLimit, memoryLimit uint64) (*SandboxInfo, error) {` |
| `core/vm_sandbox_management.go` | `54` | `func (m *SandboxManager) StopSandbox(id string) error {` |
| `core/vm_sandbox_management.go` | `66` | `func (m *SandboxManager) ResetSandbox(id string) error {` |
| `core/vm_sandbox_management.go` | `78` | `func (m *SandboxManager) SandboxStatus(id string) (*SandboxInfo, bool) {` |
| `core/vm_sandbox_management.go` | `86` | `func (m *SandboxManager) ListSandboxes() []*SandboxInfo {` |
| `core/access_control.go` | `12` | `func NewAccessController() *AccessController {` |
| `core/access_control.go` | `17` | `func (a *AccessController) Grant(role, addr string) {` |
| `core/access_control.go` | `27` | `func (a *AccessController) Revoke(role, addr string) {` |
| `core/access_control.go` | `39` | `func (a *AccessController) HasRole(role, addr string) bool {` |
| `core/access_control.go` | `51` | `func (a *AccessController) List(addr string) []string {` |
| `core/dao_token_test.go` | `5` | `func TestDAOTokenLedger(t *testing.T) {` |
| `core/dao_quadratic_voting_test.go` | `5` | `func TestQuadraticWeight(t *testing.T) {` |
| `core/dao_quadratic_voting_test.go` | `11` | `func TestCastQuadraticVoteZeroTokens(t *testing.T) {` |
| `core/dao_quadratic_voting_test.go` | `28` | `func TestCastQuadraticVoteRequiresMembership(t *testing.T) {` |
| `core/dao_quadratic_voting_test.go` | `45` | `func TestCastQuadraticVoteSuccess(t *testing.T) {` |
| `core/loanpool.go` | `18` | `func NewLoanPool(treasury uint64) *LoanPool {` |
| `core/loanpool.go` | `27` | `func (lp *LoanPool) SubmitProposal(creator, recipient, typ string, amount uint64, desc string) (uint64, error) {` |
| `core/loanpool.go` | `38` | `func (lp *LoanPool) VoteProposal(voter string, id uint64) error {` |
| `core/loanpool.go` | `51` | `func (lp *LoanPool) Tick() {` |
| `core/loanpool.go` | `61` | `func (lp *LoanPool) Disburse(id uint64) error {` |
| `core/loanpool.go` | `78` | `func (lp *LoanPool) GetProposal(id uint64) (*LoanProposal, bool) {` |
| `core/loanpool.go` | `84` | `func (lp *LoanPool) ListProposals() []*LoanProposal {` |
| `core/loanpool.go` | `94` | `func (lp *LoanPool) CancelProposal(creator string, id uint64) error {` |
| `core/loanpool.go` | `107` | `func (lp *LoanPool) ExtendProposal(creator string, id uint64, hrs int) error {` |
| `core/rpc_webrtc.go` | `13` | `func NewWebRTCRPC() *WebRTCRPC {` |
| `core/rpc_webrtc.go` | `19` | `func (r *WebRTCRPC) Connect(id string) <-chan []byte {` |
| `core/rpc_webrtc.go` | `28` | `func (r *WebRTCRPC) Send(id string, msg []byte) bool {` |
| `core/rpc_webrtc.go` | `44` | `func (r *WebRTCRPC) Disconnect(id string) {` |
| `core/audit_node.go` | `20` | `func NewAuditNode(b BootstrapNode, m *AuditManager) *AuditNode {` |
| `core/audit_node.go` | `25` | `func (n *AuditNode) Start() error {` |
| `core/audit_node.go` | `33` | `func (n *AuditNode) LogEvent(address, event string, metadata map[string]string) error {` |
| `core/audit_node.go` | `41` | `func (n *AuditNode) ListEvents(address string) []AuditEntry {` |
| `core/nat_traversal.go` | `13` | `func NewNATManager() *NATManager {` |
| `core/nat_traversal.go` | `18` | `func (n *NATManager) MapPort(id string, port int) {` |
| `core/nat_traversal.go` | `25` | `func (n *NATManager) GetPort(id string) (int, bool) {` |
| `core/nat_traversal.go` | `33` | `func (n *NATManager) RemoveMapping(id string) {` |
| `core/initialization_replication.go` | `14` | `func NewInitService(r *Replicator) *InitService {` |
| `core/initialization_replication.go` | `19` | `func (i *InitService) Start() {` |
| `core/initialization_replication.go` | `30` | `func (i *InitService) Stop() {` |
| `core/elected_authority_node.go` | `19` | `func NewElectedAuthorityNode(addr, role string, term time.Duration) *ElectedAuthorityNode {` |
| `core/elected_authority_node.go` | `25` | `func (n *ElectedAuthorityNode) IsActive(now time.Time) bool {` |
| `core/elected_authority_node.go` | `34` | `func (n *ElectedAuthorityNode) RenewTerm(requester string, dao *DAO, add time.Duration) error {` |
| `core/government_authority_node.go` | `10` | `func NewGovernmentAuthorityNode(addr, role, department string) *GovernmentAuthorityNode {` |
| `core/cross_chain_connection.go` | `24` | `func NewChainConnectionManager() *ChainConnectionManager {` |
| `core/cross_chain_connection.go` | `29` | `func (m *ChainConnectionManager) Open(local, remote string) int {` |
| `core/cross_chain_connection.go` | `39` | `func (m *ChainConnectionManager) Close(id int) error {` |
| `core/cross_chain_connection.go` | `51` | `func (m *ChainConnectionManager) Get(id int) (*ChainConnection, error) {` |
| `core/cross_chain_connection.go` | `62` | `func (m *ChainConnectionManager) List() []*ChainConnection {` |
| `core/dao_test.go` | `5` | `func TestDAOManager(t *testing.T) {` |
| `core/dao_access_control_test.go` | `5` | `func TestDAOAccessControl(t *testing.T) {` |
| `core/compliance_management_test.go` | `5` | `func TestComplianceManager(t *testing.T) {` |
| `core/node.go` | `23` | `func NewNode(id, addr string, ledger *Ledger) *Node {` |
| `core/node.go` | `38` | `func (n *Node) AddTransaction(tx *Transaction) error {` |
| `core/node.go` | `48` | `func (n *Node) ValidateTransaction(tx *Transaction) error {` |
| `core/node.go` | `59` | `func (n *Node) MineBlock() *Block {` |
| `core/node.go` | `98` | `func (n *Node) SetStake(addr string, amount uint64) error {` |
| `core/node.go` | `117` | `func (n *Node) ReportDoubleSign(addr string) {` |
| `core/node.go` | `122` | `func (n *Node) ReportDowntime(addr string) {` |
| `core/node.go` | `127` | `func (n *Node) Rehabilitate(addr string) {` |
| `core/syn500.go` | `23` | `func NewSYN500Token(name, symbol, owner string, decimals uint8, supply uint64) *SYN500Token {` |
| `core/syn500.go` | `28` | `func (t *SYN500Token) Grant(addr string, tier int, max uint64) {` |
| `core/syn500.go` | `33` | `func (t *SYN500Token) Use(addr string) error {` |
| `core/ai_enhanced_contract_test.go` | `5` | `func TestAIContractRegistry(t *testing.T) {` |
| `core/syn2100_test.go` | `8` | `func TestTradeFinanceToken(t *testing.T) {` |
| `core/syn3700_token.go` | `23` | `func NewSYN3700Token(name, symbol string) *SYN3700Token {` |
| `core/syn3700_token.go` | `28` | `func (t *SYN3700Token) AddComponent(token string, weight float64) {` |
| `core/syn3700_token.go` | `35` | `func (t *SYN3700Token) RemoveComponent(token string) error {` |
| `core/syn3700_token.go` | `48` | `func (t *SYN3700Token) ListComponents() []IndexComponent {` |
| `core/syn3700_token.go` | `57` | `func (t *SYN3700Token) Value(prices map[string]float64) float64 {` |
| `core/cross_chain_bridge.go` | `39` | `func NewBridgeManager(l *Ledger) *BridgeManager {` |
| `core/cross_chain_bridge.go` | `48` | `func (m *BridgeManager) RegisterBridge(source, target, relayer string) int {` |
| `core/cross_chain_bridge.go` | `62` | `func (m *BridgeManager) ListBridges() []*AssetBridge {` |
| `core/cross_chain_bridge.go` | `73` | `func (m *BridgeManager) GetBridge(id int) (*AssetBridge, error) {` |
| `core/cross_chain_bridge.go` | `84` | `func (m *BridgeManager) AuthorizeRelayer(id int, addr string) error {` |
| `core/cross_chain_bridge.go` | `96` | `func (m *BridgeManager) RevokeRelayer(id int, addr string) error {` |
| `core/cross_chain_bridge.go` | `108` | `func (m *BridgeManager) Deposit(bridgeID int, from, to string, amount uint64, tokenID string) (int, error) {` |
| `core/cross_chain_bridge.go` | `128` | `func (m *BridgeManager) Claim(transferID int, proof string) error {` |
| `core/cross_chain_bridge.go` | `147` | `func (m *BridgeManager) GetTransfer(id int) (*BridgeTransferRecord, error) {` |
| `core/cross_chain_bridge.go` | `158` | `func (m *BridgeManager) ListTransfers() []*BridgeTransferRecord {` |
| `core/cross_chain_bridge.go` | `187` | `func NewBridgeTransferManager() *BridgeTransferManager {` |
| `core/cross_chain_bridge.go` | `192` | `func (m *BridgeTransferManager) Deposit(bridgeID, from, to string, amount uint64, tokenID string) (*BridgeTransfer, error) {` |
| `core/cross_chain_bridge.go` | `211` | `func (m *BridgeTransferManager) Claim(id, proof string) error {` |
| `core/cross_chain_bridge.go` | `226` | `func (m *BridgeTransferManager) GetTransfer(id string) (*BridgeTransfer, bool) {` |
| `core/cross_chain_bridge.go` | `234` | `func (m *BridgeTransferManager) ListTransfers() []*BridgeTransfer {` |
| `core/zero_trust_data_channels.go` | `22` | `func NewZeroTrustEngine() *ZeroTrustEngine {` |
| `core/zero_trust_data_channels.go` | `27` | `func (e *ZeroTrustEngine) OpenChannel(id string, key []byte) error {` |
| `core/zero_trust_data_channels.go` | `38` | `func (e *ZeroTrustEngine) Send(id string, payload []byte) ([]byte, error) {` |
| `core/zero_trust_data_channels.go` | `56` | `func (e *ZeroTrustEngine) Messages(id string) [][]byte {` |
| `core/zero_trust_data_channels.go` | `73` | `func (e *ZeroTrustEngine) CloseChannel(id string) error {` |
| `core/genesis_wallets_test.go` | `5` | `func TestDefaultGenesisWalletsDeterministic(t *testing.T) {` |
| `core/genesis_wallets_test.go` | `16` | `func TestAllocateToGenesisWallets(t *testing.T) {` |
| `core/loanpool_management.go` | `17` | `func NewLoanPoolManager(p *LoanPool) *LoanPoolManager {` |
| `core/loanpool_management.go` | `22` | `func (m *LoanPoolManager) Pause() {` |
| `core/loanpool_management.go` | `27` | `func (m *LoanPoolManager) Resume() {` |
| `core/loanpool_management.go` | `32` | `func (m *LoanPoolManager) Stats() LoanPoolStats {` |
| `core/syn2700_test.go` | `8` | `func TestVestingSchedule(t *testing.T) {` |
| `core/stake_penalty.go` | `7` | `func NewStakePenaltyManager() *StakePenaltyManager { return &StakePenaltyManager{} }` |
| `core/stake_penalty.go` | `10` | `func (spm *StakePenaltyManager) Slash(sn *StakingNode, addr string, penalty uint64) {` |
| `core/stake_penalty.go` | `15` | `func (spm *StakePenaltyManager) Reward(sn *StakingNode, addr string, reward uint64) {` |
| `core/authority_nodes_test.go` | `5` | `func TestAuthorityNodeRegistry(t *testing.T) {` |
| `core/quorum_tracker.go` | `14` | `func NewQuorumTracker(required int) *QuorumTracker {` |
| `core/quorum_tracker.go` | `22` | `func (qt *QuorumTracker) Join(id string) {` |
| `core/quorum_tracker.go` | `29` | `func (qt *QuorumTracker) Leave(id string) {` |
| `core/quorum_tracker.go` | `36` | `func (qt *QuorumTracker) Count() int {` |
| `core/quorum_tracker.go` | `43` | `func (qt *QuorumTracker) Reached() bool {` |
| `core/contract_management.go` | `11` | `func NewContractManager(reg *ContractRegistry) *ContractManager {` |
| `core/contract_management.go` | `16` | `func (m *ContractManager) Transfer(addr, newOwner string) error {` |
| `core/contract_management.go` | `28` | `func (m *ContractManager) Pause(addr string) error {` |
| `core/contract_management.go` | `40` | `func (m *ContractManager) Resume(addr string) error {` |
| `core/contract_management.go` | `52` | `func (m *ContractManager) Upgrade(addr string, wasm []byte, gasLimit uint64) error {` |
| `core/contract_management.go` | `70` | `func (m *ContractManager) Info(addr string) (*Contract, error) {` |
| `core/liquidity_views.go` | `14` | `func NewLiquidityPoolView(p *LiquidityPool) LiquidityPoolView {` |
| `core/liquidity_views.go` | `26` | `func (r *LiquidityPoolRegistry) PoolInfo(id string) (LiquidityPoolView, bool) {` |
| `core/liquidity_views.go` | `35` | `func (r *LiquidityPoolRegistry) PoolViews() []LiquidityPoolView {` |
| `core/immutability_enforcement_test.go` | `5` | `func TestImmutabilityEnforcer(t *testing.T) {` |
| `core/dao_proposal.go` | `29` | `func NewProposalManager() *ProposalManager {` |
| `core/dao_proposal.go` | `41` | `func (pm *ProposalManager) CreateProposal(dao *DAO, creator, desc string) (*DAOProposal, error) {` |
| `core/dao_proposal.go` | `59` | `func (pm *ProposalManager) Vote(dao *DAO, id, voter string, weight uint64, support bool) error {` |
| `core/dao_proposal.go` | `85` | `func (pm *ProposalManager) Results(id string) (yes, no uint64, err error) {` |
| `core/dao_proposal.go` | `104` | `func (pm *ProposalManager) Execute(dao *DAO, id, requester string) error {` |
| `core/dao_proposal.go` | `124` | `func (pm *ProposalManager) Get(id string) (*DAOProposal, error) {` |
| `core/dao_proposal.go` | `135` | `func (pm *ProposalManager) List() []*DAOProposal {` |
| `core/virtual_machine.go` | `18` | `func NewSimpleVM() *SimpleVM { return &SimpleVM{} }` |
| `core/virtual_machine.go` | `21` | `func (vm *SimpleVM) Start() error {` |
| `core/virtual_machine.go` | `32` | `func (vm *SimpleVM) Stop() error {` |
| `core/virtual_machine.go` | `43` | `func (vm *SimpleVM) Status() bool {` |
| `core/virtual_machine.go` | `52` | `func (vm *SimpleVM) Execute(wasm []byte, method string, args []byte, gasLimit uint64) ([]byte, uint64, error) {` |
| `core/loanpool_apply.go` | `25` | `func NewLoanPoolApply(pool *LoanPool) *LoanPoolApply {` |
| `core/loanpool_apply.go` | `34` | `func (l *LoanPoolApply) Submit(applicant string, amount uint64, termMonths uint32, purpose string) uint64 {` |
| `core/loanpool_apply.go` | `49` | `func (l *LoanPoolApply) Vote(voter string, id uint64) error {` |
| `core/loanpool_apply.go` | `59` | `func (l *LoanPoolApply) Process() {` |
| `core/loanpool_apply.go` | `68` | `func (l *LoanPoolApply) Disburse(id uint64) error {` |
| `core/loanpool_apply.go` | `85` | `func (l *LoanPoolApply) Get(id uint64) (*LoanApplication, bool) {` |
| `core/loanpool_apply.go` | `91` | `func (l *LoanPoolApply) List() []*LoanApplication {` |
| `core/fees_test.go` | `5` | `func TestFeeForTransfer(t *testing.T) {` |
| `core/fees_test.go` | `12` | `func TestDistributeFees(t *testing.T) {` |
| `core/fees_test.go` | `23` | `func TestApplyFeeCapFloor(t *testing.T) {` |
| `core/fees_test.go` | `32` | `func TestFeePolicyEnforce(t *testing.T) {` |
| `core/fees_test.go` | `45` | `func TestAdjustFeeRates(t *testing.T) {` |
| `core/fees_test.go` | `52` | `func TestEstimateFee(t *testing.T) {` |
| `core/fees_test.go` | `59` | `func TestShareProportional(t *testing.T) {` |
| `core/fees_test.go` | `67` | `func TestAdjustForBlockUtilization(t *testing.T) {` |
| `core/biometrics_auth_test.go` | `5` | `func TestBiometricsAuth(t *testing.T) {` |
| `core/vm_sandbox_management_test.go` | `5` | `func TestSandboxManager(t *testing.T) {` |
| `core/syn2100.go` | `28` | `func NewTradeFinanceToken() *TradeFinanceToken {` |
| `core/syn2100.go` | `36` | `func (t *TradeFinanceToken) RegisterDocument(docID, issuer, recipient string, amount uint64, issue, due time.Time, desc string) {` |
| `core/syn2100.go` | `49` | `func (t *TradeFinanceToken) FinanceDocument(docID, financier string) error {` |
| `core/syn2100.go` | `63` | `func (t *TradeFinanceToken) GetDocument(docID string) (*FinancialDocument, bool) {` |
| `core/syn2100.go` | `69` | `func (t *TradeFinanceToken) ListDocuments() []*FinancialDocument {` |
| `core/syn2100.go` | `78` | `func (t *TradeFinanceToken) AddLiquidity(addr string, amt uint64) {` |
| `core/syn2100.go` | `83` | `func (t *TradeFinanceToken) RemoveLiquidity(addr string, amt uint64) error {` |
| `core/regulatory_node.go` | `14` | `func NewRegulatoryNode(id string, mgr *RegulatoryManager) *RegulatoryNode {` |
| `core/regulatory_node.go` | `23` | `func (n *RegulatoryNode) ApproveTransaction(tx Transaction) bool {` |
| `core/regulatory_node.go` | `29` | `func (n *RegulatoryNode) FlagEntity(addr, reason string) {` |
| `core/regulatory_node.go` | `36` | `func (n *RegulatoryNode) Logs(addr string) []string {` |
| `core/genesis_wallets.go` | `22` | `func hashAddress(label string) string {` |
| `core/genesis_wallets.go` | `28` | `func DefaultGenesisWallets() GenesisWallets {` |
| `core/genesis_wallets.go` | `44` | `func AllocateToGenesisWallets(total uint64, wallets GenesisWallets) map[string]uint64 {` |
| `core/blockchain_synchronization.go` | `16` | `func NewSyncManager(l *Ledger) *SyncManager {` |
| `core/blockchain_synchronization.go` | `21` | `func (s *SyncManager) Start() {` |
| `core/blockchain_synchronization.go` | `28` | `func (s *SyncManager) Stop() {` |
| `core/blockchain_synchronization.go` | `35` | `func (s *SyncManager) Status() (bool, int) {` |
| `core/blockchain_synchronization.go` | `42` | `func (s *SyncManager) Once() error {` |
| `internal/security/secrets_manager.go` | `11` | `func NewSecretsManager() *SecretsManager {` |
| `internal/security/secrets_manager.go` | `17` | `func (s *SecretsManager) Store(key, value string) error {` |
| `internal/security/secrets_manager.go` | `30` | `func (s *SecretsManager) Retrieve(key string) (string, error) {` |
