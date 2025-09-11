# Storage and Data Architecture

## Overview
Storage components handle data persistence, distribution and secure retrieval for contracts and nodes. They support off-chain resources, encrypted storage and marketplace mechanisms for sharing datasets.

## Key Modules
- `data_distribution.go` – replicates data across nodes for redundancy.
- `data_resource_management.go` – tracks dataset metadata and ownership.
- `data_operations.go` – provides read/write primitives with access checks.
- `ai_secure_storage.go` – encrypts sensitive model artifacts.
- Storage marketplace GUI and contracts for trading dataset access.

## Workflow
1. **Registration** – datasets are registered via `data_resource_management` with hashes and permissions.
2. **Distribution** – `data_distribution` propagates chunks to approved nodes.
3. **Access** – contracts call `data_operations` to read or mutate data; checks ensure caller permissions.
4. **Marketplace** – optional marketplace contracts allow leasing or selling access to datasets.

## Security Considerations
- All data operations require permission checks and are logged.
- Encryption keys are stored separately and rotated regularly.
- Distribution routines verify integrity using content hashes.

## CLI Integration
- `synnergy data` – manage resources, distribution and access policies.
