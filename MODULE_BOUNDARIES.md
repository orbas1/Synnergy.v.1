# Module Boundaries

This document captures the initial module audit and domain grouping for the Synnergy repository. It serves as a guide for the
ongoing refactor that will move code into the new `internal/` and `pkg/` directories.

## Domain Groupings

- **Core**
  - Foundational blockchain logic: consensus, state management, VM and sandbox utilities. Representative files include `core/`,
    `virtual_machine.go`, `dynamic_consensus_hopping.go` and `vm_sandbox_management.go`.
- **Nodes**
  - Specialized network nodes such as `watchtower_node.go`, `regulatory_node.go`, `energy_efficient_node.go`, `biometric_security_node.go` and `geospatial_node.go`.
- **Tokens**
  - Token and asset features including the `Tokens/` directory, `private_transactions.go`, `faucet.go` and related tests.
- **Cross-Chain**
  - Interoperability logic: `cross_chain.go`, `cross_chain_bridge.go`, `cross_chain_transactions.go` and `cross_chain_contracts.go`.
- **Security**
  - Authentication, authorization and monitoring such as `firewall.go`, `identity_verification.go`, `access_control.go` and `zero_trust_data_channels.go`.
- **Utilities**
  - Shared helpers that will become public libraries under `pkg/` including configuration, logging and common data structures.

## Directory Strategy

- `internal/` houses code meant exclusively for Synnergy binaries and services.
- `pkg/` provides stable, reusable libraries for external applications.

Subsequent production stages will relocate existing packages into these directories and update imports accordingly.
