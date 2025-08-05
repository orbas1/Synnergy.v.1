package banknodes

import (
	"testing"

	"synnergy/core"
	"synnergy/nodes"
)

type bankNodeAdapter struct{ *core.BankInstitutionalNode }

func (a *bankNodeAdapter) ID() nodes.Address                 { return nodes.Address(a.Node.ID) }
func (a *bankNodeAdapter) Start() error                      { return nil }
func (a *bankNodeAdapter) Stop() error                       { return nil }
func (a *bankNodeAdapter) IsRunning() bool                   { return true }
func (a *bankNodeAdapter) Peers() []nodes.Address            { return nil }
func (a *bankNodeAdapter) DialSeed(addr nodes.Address) error { return nil }

type centralBankNodeAdapter struct{ *core.CentralBankingNode }

func (a *centralBankNodeAdapter) ID() nodes.Address                 { return nodes.Address(a.Node.ID) }
func (a *centralBankNodeAdapter) Start() error                      { return nil }
func (a *centralBankNodeAdapter) Stop() error                       { return nil }
func (a *centralBankNodeAdapter) IsRunning() bool                   { return true }
func (a *centralBankNodeAdapter) Peers() []nodes.Address            { return nil }
func (a *centralBankNodeAdapter) DialSeed(addr nodes.Address) error { return nil }

type custodialNodeAdapter struct{ *core.CustodialNode }

func (a *custodialNodeAdapter) ID() nodes.Address                 { return nodes.Address(a.Node.ID) }
func (a *custodialNodeAdapter) Start() error                      { return nil }
func (a *custodialNodeAdapter) Stop() error                       { return nil }
func (a *custodialNodeAdapter) IsRunning() bool                   { return true }
func (a *custodialNodeAdapter) Peers() []nodes.Address            { return nil }
func (a *custodialNodeAdapter) DialSeed(addr nodes.Address) error { return nil }

// TestInterfaceCompliance ensures core implementations satisfy the banking node interfaces via adapters.
func TestInterfaceCompliance(t *testing.T) {
	ledger := core.NewLedger()

	var _ BankInstitutionalNode = &bankNodeAdapter{core.NewBankInstitutionalNode("id1", "addr1", ledger)}
	var _ CentralBankingNode = &centralBankNodeAdapter{core.NewCentralBankingNode("id2", "addr2", ledger, "policy")}
	var _ CustodialNode = &custodialNodeAdapter{core.NewCustodialNode("id3", "addr3", ledger)}
}
