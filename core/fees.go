package core

import "sort"

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

// CalculateBaseFee computes the base fee using the median of recent fees and
// an adjustment factor representing current network load.
func CalculateBaseFee(recent []uint64, adjustment float64) uint64 {
	if len(recent) == 0 {
		return 0
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

// FeeDistribution represents the allocation of fees across network
// stakeholders.
type FeeDistribution struct {
	InternalDevelopment uint64
	InternalCharity     uint64
	ExternalCharity     uint64
	LoanPool            uint64
	PassiveIncome       uint64
	ValidatorsMiners    uint64
	NodeHosts           uint64
	CreatorWallet       uint64
}

// DistributeFees splits the total fees according to the network's policy.
func DistributeFees(total uint64) FeeDistribution {
	return FeeDistribution{
		InternalDevelopment: total * 5 / 100,
		InternalCharity:     total * 5 / 100,
		ExternalCharity:     total * 5 / 100,
		LoanPool:            total * 5 / 100,
		PassiveIncome:       total * 5 / 100,
		ValidatorsMiners:    total * 69 / 100,
		NodeHosts:           total * 5 / 100,
		CreatorWallet:       total * 1 / 100,
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
