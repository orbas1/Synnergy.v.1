package nodes

// dummyConsensusNode implements ConsensusNodeInterface for testing purposes.
type dummyConsensusNode struct{}

func (d dummyConsensusNode) ID() Address                 { return "dummy" }
func (d dummyConsensusNode) Start() error                { return nil }
func (d dummyConsensusNode) Stop() error                 { return nil }
func (d dummyConsensusNode) IsRunning() bool             { return false }
func (d dummyConsensusNode) Peers() []Address            { return nil }
func (d dummyConsensusNode) DialSeed(addr Address) error { return nil }
func (d dummyConsensusNode) ConsensusType() string       { return "pow" }

// Ensure dummyConsensusNode satisfies the interface at compile time.
var _ ConsensusNodeInterface = (*dummyConsensusNode)(nil)
