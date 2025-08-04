package core

// Bank node type identifiers.
const (
	BankInstitutionalNodeType = "BANK_INSTITUTIONAL"
	CentralBankingNodeType    = "CENTRAL_BANKING"
	CustodialNodeType         = "CUSTODIAL"
)

// BankNodeTypes lists all bank-related node categories supported by the network.
var BankNodeTypes = []string{
	BankInstitutionalNodeType,
	CentralBankingNodeType,
	CustodialNodeType,
}
