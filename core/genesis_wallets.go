package core

import (
	"crypto/sha256"
	"encoding/hex"
)

// GenesisWallets holds addresses for initial ecosystem allocations.
type GenesisWallets struct {
	Genesis             string
	InternalDevelopment string
	InternalCharity     string
	ExternalCharity     string
	LoanPool            string
	PassiveIncome       string
	ValidatorsMiners    string
	NodeHosts           string
	CreatorWallet       string
}

// hashAddress returns a deterministic hex string for a given label.
func hashAddress(label string) string {
	h := sha256.Sum256([]byte(label))
	return hex.EncodeToString(h[:])
}

// DefaultGenesisWallets provides the default set of genesis wallet addresses.
func DefaultGenesisWallets() GenesisWallets {
	return GenesisWallets{
		Genesis:             hashAddress("genesis"),
		InternalDevelopment: hashAddress("internal_development"),
		InternalCharity:     hashAddress("internal_charity"),
		ExternalCharity:     hashAddress("external_charity"),
		LoanPool:            hashAddress("loan_pool"),
		PassiveIncome:       hashAddress("passive_income"),
		ValidatorsMiners:    hashAddress("validators_miners"),
		NodeHosts:           hashAddress("node_hosts"),
		CreatorWallet:       hashAddress("creator_wallet"),
	}
}

// AllocateToGenesisWallets splits total fees across the genesis wallet
// addresses according to the network's fee distribution policy.
func AllocateToGenesisWallets(total uint64, wallets GenesisWallets) map[string]uint64 {
	dist := DistributeFees(total)
	return map[string]uint64{
		wallets.InternalDevelopment: dist.InternalDevelopment,
		wallets.InternalCharity:     dist.InternalCharity,
		wallets.ExternalCharity:     dist.ExternalCharity,
		wallets.LoanPool:            dist.LoanPool,
		wallets.PassiveIncome:       dist.PassiveIncome,
		wallets.ValidatorsMiners:    dist.ValidatorsMiners,
		wallets.NodeHosts:           dist.NodeHosts,
		wallets.CreatorWallet:       dist.CreatorWallet,
	}
}
