# Security Architecture

Security modules enforce network protection mechanisms including firewalls, anomaly detection, and private transactions.

**Key Modules**
- firewall.go
- anomaly_detection.go
- private_transactions.go
- stake_penalty.go
- zero_trust_data_channels.go

**Related CLI Files**
- cli/firewall.go
- cli/private_transactions.go
- cli/stake_penalty.go
- cli/zero_trust_data_channels.go

Collectively, these components harden the platform against threats and enforce secure communication channels. Stage 13 upgrades the zero trust engine with ed25519 signatures layered over AES-GCM encryption so that every message is both confidential and authenticated.
