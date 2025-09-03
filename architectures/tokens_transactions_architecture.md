# Tokens and Transactions Architecture

This group governs token standards, contract execution, and transaction processing.

**Key Modules**
- contracts.go
- contracts_opcodes.go
- virtual_machine.go
- gas_table.go
- private_transactions.go
- transaction.go
- vm_sandbox_management.go
- coin.go
- blockchain_compression.go
- blockchain_synchronization.go
- charity.go

**Related CLI Files**
- cli/contracts.go
- cli/contracts_opcodes.go
- cli/virtual_machine.go
- cli/gas_table.go
- cli/private_transactions.go
- cli/transaction.go
- cli/coin.go
- cli/compression.go
- cli/synchronization.go
- cli/charity.go

These modules define how smart contracts run and how tokens and transactions are validated and executed.

Stage 16 introduces a concurrency‑safe token registry and base token with
benchmarks to measure transfer performance across the network.
