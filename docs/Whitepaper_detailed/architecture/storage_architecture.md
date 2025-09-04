# Storage and Data Architecture

Storage-related modules handle data operations, secure storage, and sandboxed environments for smart contract execution.

**Key Modules**
- data.go
- data_distribution.go
- data_operations.go
- data_resource_management.go
- ai_secure_storage.go
- vm_sandbox_management.go
- zero_trust_data_channels.go
- storage_marketplace.go

**Related CLI Files**
- cli/state_rw.go
- cli/zero_trust_data_channels.go
- cli/storage_marketplace.go

**Scripts**
- scripts/backup_ledger.sh

These modules provide robust data handling, ensuring integrity and isolation for contract and node storage requirements.
Stage 35 adds a storage marketplace spanning the core, CLI and a reference GUI under `GUI/storage-marketplace` where users can list and lease capacity through gas-priced opcodes.
