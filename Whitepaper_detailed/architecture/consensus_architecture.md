# Consensus Architecture

This group defines the mechanisms that drive agreement across the network. Core modules handle consensus algorithms, difficulty adjustments, adaptive management, and validator coordination. CLI commands expose administrative interfaces for interacting with consensus features.

**Key Modules**
- consensus.go
- consensus_adaptive_management.go
- consensus_difficulty.go
- consensus_mode.go
- consensus_service.go
- consensus_specific.go
- consensus_specific_node.go
- dynamic_consensus_hopping.go
- validator_management.go

**Related CLI Files**
- cli/consensus.go
- cli/consensus_adaptive_management.go
- cli/consensus_difficulty.go
- cli/consensus_mode.go
- cli/consensus_service.go
- cli/consensus_specific_node.go
- cli/validator_management.go

These components coordinate to ensure blocks are validated and added securely while allowing dynamic algorithm changes and validator oversight.

Stage 7 introduces a coded errors package and OpenTelemetry tracing across consensus modules, providing contextual diagnostics and observability for validator operations.
