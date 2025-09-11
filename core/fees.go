package core

import (
	"fmt"
	"sort"
	"sync"
)

var (
	creatorDistMu      sync.RWMutex
	creatorDistEnabled = true
)

// SetCreatorDistribution enables or disables the creator wallet fee share.
// When disabled, the creator allocation is redirected to node hosts.
func SetCreatorDistribution(enabled bool) {
	creatorDistMu.Lock()
	creatorDistEnabled = enabled
	creatorDistMu.Unlock()
}

// IsCreatorDistributionEnabled reports whether the creator wallet currently
// receives its fee allocation.
func IsCreatorDistributionEnabled() bool {
	creatorDistMu.RLock()
	defer creatorDistMu.RUnlock()
	return creatorDistEnabled
}

// TransactionType represents high level categories of transactions used for
// fee calculations and policy decisions.
type TransactionType int

const (
	TxTypeTransfer TransactionType = iota
	TxTypePurchase
	TxTypeTokenInteraction
	TxTypeContract
	TxTypeWalletVerification
)

// FeeBreakdown captures the components of a transaction fee.
type FeeBreakdown struct {
	Base     uint64
	Variable uint64
	Priority uint64
	Total    uint64
}

// CalculateBaseFee computes the base fee using the median of the most recent
// 1000 block fees and an adjustment factor representing current network load.
//
// If more than 1000 fee values are supplied only the latest 1000 are used,
// allowing the caller to provide a rolling history without manual trimming.
func CalculateBaseFee(recent []uint64, adjustment float64) uint64 {
	if len(recent) == 0 {
		return 0
	}

	if len(recent) > 1000 {
		recent = recent[len(recent)-1000:]
	}

	data := append([]uint64(nil), recent...)
	sort.Slice(data, func(i, j int) bool { return data[i] < data[j] })
	median := data[len(data)/2]
	return uint64(float64(median) * (1 + adjustment))
}

// CalculateVariableFee multiplies gas units by the gas price per unit.
func CalculateVariableFee(gasUnits, gasPrice uint64) uint64 {
	return gasUnits * gasPrice
}

// CalculatePriorityFee returns the tip specified by the user.
func CalculatePriorityFee(tip uint64) uint64 { return tip }

// FeeForTransfer calculates fees for a simple transfer based on data size.
func FeeForTransfer(dataSize, baseFee, variableRate, tip uint64) FeeBreakdown {
	variable := dataSize * variableRate
	total := baseFee + variable + tip
	return FeeBreakdown{Base: baseFee, Variable: variable, Priority: tip, Total: total}
}

// FeeForPurchase calculates fees for purchase transactions based on contract calls.
func FeeForPurchase(calls, baseFee, variableRate, tip uint64) FeeBreakdown {
	variable := calls * variableRate
	total := baseFee + variable + tip
	return FeeBreakdown{Base: baseFee, Variable: variable, Priority: tip, Total: total}
}

// FeeForTokenUsage calculates fees for interactions with deployed tokens.
func FeeForTokenUsage(computationUnits, baseFee, variableRate, tip uint64) FeeBreakdown {
	variable := computationUnits * variableRate
	total := baseFee + variable + tip
	return FeeBreakdown{Base: baseFee, Variable: variable, Priority: tip, Total: total}
}

// FeeForContract calculates fees for contract creation or modification.
func FeeForContract(complexityFactor, baseFee, variableRate, tip uint64) FeeBreakdown {
	variable := complexityFactor * variableRate
	total := baseFee + variable + tip
	return FeeBreakdown{Base: baseFee, Variable: variable, Priority: tip, Total: total}
}

// FeeForWalletVerification calculates fees for wallet verification steps.
func FeeForWalletVerification(securityLevel, baseFee, variableRate, tip uint64) FeeBreakdown {
	variable := securityLevel * variableRate
	total := baseFee + variable + tip
	return FeeBreakdown{Base: baseFee, Variable: variable, Priority: tip, Total: total}
}

// FeeForValidatedTransfer returns zero fees for transfers that have been
// validated as eligible for fee-less execution. If validated is false it
// falls back to the standard transfer fee calculation.
func FeeForValidatedTransfer(dataSize, baseFee, variableRate, tip uint64, validated bool) FeeBreakdown {
	if validated {
		return FeeBreakdown{}
	}
	return FeeForTransfer(dataSize, baseFee, variableRate, tip)
}

// FeeDistribution represents the allocation of fees across network
// stakeholders.
type FeeDistribution struct {
	InternalDevelopment uint64
	InternalCharity     uint64
	ExternalCharity     uint64
	LoanPool            uint64
	PassiveIncome       uint64
	ValidatorsMiners    uint64
	AuthorityNodes      uint64
	NodeHosts           uint64
	CreatorWallet       uint64
}

type FeeSplitPolicy struct {
	InternalDevelopment uint64
	InternalCharity     uint64
	ExternalCharity     uint64
	LoanPool            uint64
	PassiveIncome       uint64
	ValidatorsMiners    uint64
	AuthorityNodes      uint64
	NodeHosts           uint64
	CreatorWallet       uint64
}

// Validate ensures the fee split policy sums to 100 percent.
func (p FeeSplitPolicy) Validate() error {
	sum := p.InternalDevelopment + p.InternalCharity + p.ExternalCharity + p.LoanPool + p.PassiveIncome + p.ValidatorsMiners + p.AuthorityNodes + p.NodeHosts + p.CreatorWallet
	if sum != 100 {
		return fmt.Errorf("fee split must sum to 100, got %d", sum)
	}
	return nil
}

var DefaultFeeSplitPolicy = FeeSplitPolicy{5, 5, 5, 10, 5, 59, 5, 5, 1}

// DistributeFeesWithPolicy splits total fees according to the provided policy.
func DistributeFeesWithPolicy(total uint64, p FeeSplitPolicy) (FeeDistribution, error) {
	if err := p.Validate(); err != nil {
		return FeeDistribution{}, err
	}
	return FeeDistribution{
		InternalDevelopment: total * p.InternalDevelopment / 100,
		InternalCharity:     total * p.InternalCharity / 100,
		ExternalCharity:     total * p.ExternalCharity / 100,
		LoanPool:            total * p.LoanPool / 100,
		PassiveIncome:       total * p.PassiveIncome / 100,
		ValidatorsMiners:    total * p.ValidatorsMiners / 100,
		AuthorityNodes:      total * p.AuthorityNodes / 100,
		NodeHosts:           total * p.NodeHosts / 100,
		CreatorWallet:       total * p.CreatorWallet / 100,
	}, nil
}

// DistributeFees splits the total fees using the default network policy.
func DistributeFees(total uint64) FeeDistribution {

	dist, _ := DistributeFeesWithPolicy(total, DefaultFeeSplitPolicy)
	return dist

	creatorShare := total * 1 / 100
	nodeHostShare := total * 5 / 100
	if !IsCreatorDistributionEnabled() {
		nodeHostShare += creatorShare
		creatorShare = 0
	}

	return FeeDistribution{
		InternalDevelopment: total * 5 / 100,
		InternalCharity:     total * 5 / 100,
		ExternalCharity:     total * 5 / 100,
		LoanPool:            total * 10 / 100,
		PassiveIncome:       total * 5 / 100,
		ValidatorsMiners:    total * 59 / 100,
		AuthorityNodes:      total * 5 / 100,
		NodeHosts:           nodeHostShare,
		CreatorWallet:       creatorShare,
	}

}

// ApplyFeeCapFloor constrains fees to the provided cap and floor values.
func ApplyFeeCapFloor(fee, cap, floor uint64) uint64 {
	if cap > 0 && fee > cap {
		fee = cap
	}
	if fee < floor {
		fee = floor
	}
	return fee
}

// FeePolicy defines network-wide fee limits.
type FeePolicy struct {
	Cap   uint64
	Floor uint64
}

// Enforce applies the cap and floor and returns the adjusted fee along with a
// message if a threshold was triggered.
func (p FeePolicy) Enforce(fee uint64) (uint64, string) {
	adjusted := ApplyFeeCapFloor(fee, p.Cap, p.Floor)
	var note string
	if p.Cap > 0 && fee > p.Cap {
		note = fmt.Sprintf("fee capped at %d", p.Cap)
	} else if fee < p.Floor {
		note = fmt.Sprintf("fee raised to floor %d", p.Floor)
	}
	return adjusted, note
}

// AdjustFeeRates adjusts the base fee and variable rate according to the
// provided network load factor. A load of 0 leaves the fees unchanged while a
// load of 0.5 increases them by 50%.
func AdjustFeeRates(baseFee, variableRate uint64, load float64) (uint64, uint64) {
	if load < 0 {
		load = 0
	}
	adjustedBase := uint64(float64(baseFee) * (1 + load))
	adjustedVariable := uint64(float64(variableRate) * (1 + load))
	return adjustedBase, adjustedVariable
}

// EstimateFee provides a generic fee estimation for any supported transaction
// type. The units argument represents the metric relevant to the given
// transaction type, such as data size for transfers or contract calls for
// purchases.
func EstimateFee(txType TransactionType, units, baseFee, variableRate, tip uint64) FeeBreakdown {
	switch txType {
	case TxTypeTransfer:
		return FeeForTransfer(units, baseFee, variableRate, tip)
	case TxTypePurchase:
		return FeeForPurchase(units, baseFee, variableRate, tip)
	case TxTypeTokenInteraction:
		return FeeForTokenUsage(units, baseFee, variableRate, tip)
	case TxTypeContract:
		return FeeForContract(units, baseFee, variableRate, tip)
	case TxTypeWalletVerification:
		return FeeForWalletVerification(units, baseFee, variableRate, tip)
	default:
		total := baseFee + tip
		return FeeBreakdown{Base: baseFee, Priority: tip, Total: total}
	}
}

// ShareProportional splits total fees according to provided weights.
// Remaining units from integer division are assigned to the first address.
func ShareProportional(total uint64, weights map[string]uint64) map[string]uint64 {
	shares := make(map[string]uint64)
	var weightTotal uint64
	var firstAddr string
	var maxWeight uint64
	for addr, w := range weights {
		weightTotal += w
		if firstAddr == "" || w > maxWeight || (w == maxWeight && addr < firstAddr) {
			firstAddr = addr
			maxWeight = w
		}
	}
	if weightTotal == 0 {
		return shares
	}
	var distributed uint64
	for addr, w := range weights {
		share := total * w / weightTotal
		shares[addr] = share
		distributed += share
	}
	if distributed < total {
		shares[firstAddr] += total - distributed
	}
	return shares
}

// FeeDistributionContract simulates a smart contract that credits fee shares to the ledger.
type FeeDistributionContract struct {
	Ledger *Ledger
}

// NewFeeDistributionContract creates a distribution contract.
func NewFeeDistributionContract(l *Ledger) *FeeDistributionContract {
	return &FeeDistributionContract{Ledger: l}
}

// Distribute credits each participant with its fee share.
func (f *FeeDistributionContract) Distribute(shares map[string]uint64) {
	for addr, amt := range shares {
		f.Ledger.Credit(addr, amt)
	}
}

// AdjustForBlockUtilization modifies the validators/miners fee pool based on block usage.
// A high utilization (>90%) increases rewards by 10% while low utilization (<50%) reduces by 10%.
func AdjustForBlockUtilization(pool uint64, used, capacity int) uint64 {
	if capacity == 0 {
		return pool
	}
	ratio := float64(used) / float64(capacity)
	switch {
	case ratio > 0.9:
		return uint64(float64(pool) * 1.1)
	case ratio < 0.5:
		return uint64(float64(pool) * 0.9)
	default:
		return pool
	}
}
