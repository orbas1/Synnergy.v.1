### Stage 1
# synnergy/core
core/cross_chain_contracts.go:13:6: ContractRegistry redeclared in this block
	core/contracts.go:32:6: other declaration of ContractRegistry
core/cross_chain_contracts.go:19:6: NewContractRegistry redeclared in this block
	core/contracts.go:39:6: other declaration of NewContractRegistry
core/cross_chain_contracts.go:20:27: unknown field mappings in struct literal of type ContractRegistry
core/cross_chain_contracts.go:27:4: r.mappings undefined (type *ContractRegistry has no field or method mappings)

### Stage 2
core/cross_chain_contracts.go:34:13: r.mappings undefined (type *ContractRegistry has no field or method mappings)
core/cross_chain_contracts.go:42:43: r.mappings undefined (type *ContractRegistry has no field or method mappings)
core/cross_chain_contracts.go:43:22: r.mappings undefined (type *ContractRegistry has no field or method mappings)
core/cross_chain_contracts.go:53:11: r.mappings undefined (type *ContractRegistry has no field or method mappings)
core/gas_table.go:4:5: gasTable redeclared in this block
	core/gas.go:14:5: other declaration of gasTable
core/gas_table.go:9:6: initGasTable redeclared in this block

### Stage 3
	core/gas.go:22:6: other declaration of initGasTable
core/gas_table.go:22:5: gasTable redeclared in this block
	core/gas.go:14:5: other declaration of gasTable
core/gas_table.go:26:6: initGasTable redeclared in this block
	core/gas.go:22:6: other declaration of initGasTable
core/gas_table.go:32:6: GasCost redeclared in this block

### Stage 4
	core/gas.go:17:6: other declaration of GasCost
core/opcodes_basic.go:5:2: OpNoop redeclared in this block
	core/opcodes_base.go:7:2: other declaration of OpNoop
core/opcodes_basic.go:6:2: OpPush redeclared in this block
	core/opcodes_base.go:8:2: other declaration of OpPush
core/opcodes_basic.go:7:2: OpAdd redeclared in this block

### Stage 5
	core/opcodes_base.go:9:2: other declaration of OpAdd
core/opcodes_basic.go:8:2: OpSub redeclared in this block
	core/opcodes_base.go:10:2: other declaration of OpSub
core/opcodes_basic.go:9:2: OpMul redeclared in this block
	core/opcodes_base.go:11:2: other declaration of OpMul
core/opcodes_basic.go:10:2: OpDiv redeclared in this block

### Stage 7
core/transaction.go:43:57: unknown field Fee in struct literal of type Transaction
core/transaction.go:43:67: unknown field Nonce in struct literal of type Transaction
core/transaction.go:43:81: unknown field Timestamp in struct literal of type Transaction
core/transaction.go:43:111: unknown field Type in struct literal of type Transaction
core/transaction.go:51:30: t.BiometricHash undefined (type *Transaction has no field or method BiometricHash)
core/transaction.go:52:84: t.Fee undefined (type *Transaction has no field or method Fee)

### Stage 8
core/transaction.go:52:91: t.Nonce undefined (type *Transaction has no field or method Nonce)
core/transaction.go:52:100: t.Timestamp undefined (type *Transaction has no field or method Timestamp)
core/transaction.go:58:30: t.Signature undefined (type *Transaction has no field or method Signature)
core/transaction.go:74:4: t.BiometricHash undefined (type *Transaction has no field or method BiometricHash)
core/transaction_control.go:43:26: tx.Fee undefined (type *Transaction has no field or method Fee)
core/transaction_control.go:44:25: invalid operation: l.balances[tx.To] < tx.Amount (mismatched types uint64 and float64)

### Stage 9
core/transaction_control.go:47:2: invalid operation: l.balances[tx.To] -= tx.Amount (mismatched types uint64 and float64)
core/transaction_control.go:53:6: PrivateTransaction redeclared in this block
	core/private_transactions.go:49:6: other declaration of PrivateTransaction
core/transaction_control.go:75:47: unknown field Nonce in struct literal of type PrivateTransaction
core/transaction_control.go:85:36: pt.Nonce undefined (type *PrivateTransaction has no field or method Nonce)
core/base_node.go:12:16: undefined: nodes.Address

### Stage 10
core/base_node.go:13:20: undefined: nodes.Address
core/base_node.go:19:27: undefined: nodes.Address
core/base_node.go:22:25: undefined: nodes.Address
core/base_node.go:27:31: undefined: nodes.Address
core/base_node.go:52:36: undefined: nodes.Address
core/base_node.go:55:22: undefined: nodes.Address

### Stage 11
core/base_node.go:63:40: undefined: nodes.Address
core/charity.go:79:25: undefined: Address
core/charity.go:80:25: undefined: Address
core/charity.go:86:28: undefined: StringToAddress
core/charity.go:91:32: undefined: StringToAddress
core/charity.go:102:23: undefined: Address

### Stage 12
core/charity.go:109:44: undefined: StateRW
core/charity.go:109:88: undefined: CharityPool
core/charity.go:110:10: undefined: CharityPool
core/charity.go:117:11: undefined: CharityPool
core/charity.go:117:37: undefined: Address
core/charity.go:125:11: undefined: CharityPool

### Stage 13
core/charity.go:125:38: undefined: Address
core/charity.go:145:7: undefined: CharityRegistration
core/charity.go:146:23: undefined: mustJSON
core/charity.go:155:11: undefined: CharityPool
core/charity.go:155:44: undefined: Address
core/charity.go:170:10: undefined: voteKey

### Stage 14
core/charity.go:175:10: undefined: CharityRegistration
core/charity.go:178:42: undefined: mustJSON
core/charity.go:183:18: undefined: Address
core/charity.go:186:30: undefined: Hash
core/charity.go:187:9: undefined: Hash
core/charity.go:199:11: undefined: CharityPool

### Stage 15
core/charity.go:219:11: undefined: CharityPool
core/charity.go:239:11: undefined: CharityPool
core/charity.go:241:37: undefined: CharityRegistration
core/charity.go:244:9: undefined: CharityRegistration
core/charity.go:248:16: undefined: Address
core/charity.go:261:33: undefined: mustJSON

### Stage 16
core/charity.go:264:11: undefined: CharityPool
core/charity.go:264:51: undefined: Address
core/charity.go:269:12: undefined: Address
core/charity.go:278:11: undefined: CharityPool
core/charity.go:278:49: undefined: Address
core/charity.go:286:16: undefined: Address

### Stage 17
core/charity.go:295:11: undefined: CharityPool
core/charity.go:295:59: undefined: Address
core/charity.go:295:69: undefined: CharityRegistration
core/charity.go:299:10: undefined: CharityRegistration
core/charity.go:317:11: undefined: CharityPool
core/charity.go:324:11: undefined: CharityPool

### Stage 18
core/charity.go:327:11: undefined: CharityPool
core/charity.go:335:32: undefined: Address
core/charity.go:341:11: undefined: CharityPool
core/charity.go:345:9: undefined: CharityRegistration
core/light_node.go:8:18: undefined: nodes.BlockHeader
core/light_node.go:12:28: undefined: nodes.Address

### Stage 19
core/light_node.go:17:39: undefined: nodes.BlockHeader
core/light_node.go:20:43: undefined: nodes.BlockHeader
core/light_node.go:22:16: undefined: nodes.BlockHeader
core/full_node.go:22:27: undefined: nodes.Address
core/gateway_node.go:22:30: undefined: nodes.Address
core/ledger.go:127:26: tx.Fee undefined (type *Transaction has no field or method Fee)

### Stage 20
core/ledger.go:132:2: invalid operation: l.balances[tx.To] += tx.Amount (mismatched types uint64 and float64)
core/node.go:15:31: undefined: nodes.Address
core/node_engine.go:49:49: tx.Fee undefined (type *Transaction has no field or method Fee)
core/node_engine.go:77:19: tx.Fee undefined (type *Transaction has no field or method Fee)
core/snvm.go:21:25: tx.Program undefined (type *Transaction has no field or method Program)
core/wallet.go:39:5: tx.Signature undefined (type *Transaction has no field or method Signature)

