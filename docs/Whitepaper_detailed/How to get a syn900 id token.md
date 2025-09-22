# How to Get a SYN900 ID Token

## Overview
Neto Solaris issues the **SYN900** identity token as the cornerstone of user authentication on the Synnergy Network. Each token binds a wallet address to verified personal details, enabling compliant participation across modules such as the LoanPool, Charity Pool and governance portals. Unlike traditional credentials, the SYN900 token is privacy‑preserving and auditable, allowing services to confirm identity without exposing raw documentation.

## Why the SYN900 Token Matters
- **Access control** – Contracts and services gate sensitive actions to proven identities. For example, the Charity Pool verifies token holders before counting a vote, rejecting any address without a registered ID【F:core/charity.go†L103-L110】.
- **Regulatory compliance** – LoanPool submissions require a verified identity token and up‑to‑date KYC/AML logs, ensuring accountability and legal adherence【F:docs/Whitepaper_detailed/How apply for a grant or loan from loanpool.md†L11-L18】【F:docs/Whitepaper_detailed/How apply for a grant or loan from loanpool.md†L154-L159】.
- **Reusable authentication** – Once issued, the SYN900 token can authenticate the holder across any Synnergy module, removing redundant checks while retaining a complete audit trail.
- **Programmatic registry integration** – The IDRegistry maps wallet addresses to descriptive metadata and exposes constant‑time lookups for services that must confirm enrolment before delivering benefits or fee distributions【F:idwallet_registration.go†L8-L44】【F:docs/Whitepaper_detailed/Transaction fee distribution.md†L106-L110】.
- **Auditable KYC trail** – The `ComplianceService` hashes submitted documents and records fraud signals, preserving a tamper‑evident audit history without revealing sensitive data【F:compliance.go†L61-L73】【F:compliance.go†L83-L90】【F:compliance.go†L102-L109】.

## Prerequisites
Before requesting a SYN900 token, ensure that:
- You possess a Synnergy wallet address registered in the ID wallet registry【F:idwallet_registration.go†L8-L27】.
- Legal identification documents for KYC/AML checks, such as a passport or government‑issued ID.
- CLI or API access to submit commands through the identity and idwallet modules【F:cli/idwallet.go†L12-L44】【F:cli/identity.go†L13-L67】.
- Network connectivity to submit registration data to the IdentityService.
- Your address is not suspended by the ComplianceManager; suspended or flagged addresses cannot receive identity tokens【F:compliance_management.go†L30-L35】【F:compliance_management.go†L67-L77】.

Stage 79 bootstrap tooling ensures the compliance stack is ready before identity enrolment begins. Operators can execute `synnergy orchestrator bootstrap --authority compliance=regnode` to invoke `core.EnterpriseOrchestrator.BootstrapNetwork`, which registers the regulatory node, starts ledger replication and produces a signed bootstrap signature for audit logs so SYN900 requests inherit enterprise-grade controls from the first RPC call.【F:cli/orchestrator.go†L58-L117】【F:core/enterprise_orchestrator.go†L71-L209】 Startup now loads Stage 79 gas costs alongside Stage 78, while the web control panel exposes the same workflow to keep CLI, API and GUI issuance paths aligned.【F:cmd/synnergy/main.go†L63-L106】【F:web/pages/index.js†L1-L214】【F:web/pages/api/bootstrap.js†L1-L45】 Comprehensive bootstrap tests cover unit, situational, stress, functional and real-world cases, preserving privacy guarantees and regulatory compliance across all identity operations.【F:core/enterprise_orchestrator_test.go†L73-L178】

## Step‑by‑Step Process
1. **Register the wallet** – Use the on‑chain `IDRegistry` to associate your wallet with metadata. The registry rejects duplicates by locking a mutex during writes and provides read‑only `Info` and `IsRegistered` helpers for external checks【F:idwallet_registration.go†L8-L44】.
2. **Submit identity details** – Call the `IdentityService.Register` method with your name, date of birth and nationality. Data is stored in a mutex‑protected map to maintain ledger integrity【F:identity_verification.go†L37-L46】.
3. **Complete verification** – Provide required KYC/AML evidence. Each verification method is persisted via `IdentityService.Verify`, appending a timestamped log entry for auditors and regulators【F:identity_verification.go†L48-L57】.
4. **Record compliance evidence** – Hash documents through `ComplianceService.ValidateKYC`, which stores an immutable commitment and audit entry for later review【F:compliance.go†L61-L73】.
5. **Issuance via IdentityTokenAPI** – After logs are validated, an authority node mints the SYN900 token using dedicated modules such as `tokens_syn900.go` and `tokens_syn900_index.go` which expose the IdentityTokenAPI【F:docs/Whitepaper_detailed/guide/module_guide.md†L460-L461】.
6. **Receive the token** – The minted token is deposited into the registered wallet. Standard operations consume 90 gas units, matching the `StdSYN900` opcode budget【F:docs/Whitepaper_detailed/guide/opcode_and_gas_guide.md†L1458-L1462】.
7. **Use across the network** – Loan proposals, charity votes and governance actions validate the identity token automatically before proceeding.

### Command‑Line Example
```bash
synnergy idwallet register 0xABCD... "KYC reference"
synnergy identity register 0xABCD... "Alice" "1990-01-01" "UK"
synnergy identity verify 0xABCD... passport
synnergy identity info 0xABCD...
synnergy identity logs 0xABCD...
```
These commands wrap the `IDRegistry` and `IdentityService` interfaces, enabling scripted enrolment and audits【F:cli/idwallet.go†L12-L44】【F:cli/identity.go†L13-L67】.

## Smart‑Contract and Opcode Integration
- **Opcode hooks** – Identity functions surface at the VM layer through opcodes like `RegisterIdentity`, `VerifyIdentity`, and `RemoveIdentity`, allowing contracts to enforce enrollment checks directly【F:contracts_opcodes.go†L1052-L1055】【F:snvm._opcodes.go†L141-L145】.
- **Deterministic gas costs** – The opcode schedule assigns fixed budgets (500 gas for registration, 100 for verification, 200 for removal), helping enterprises plan resource consumption【F:docs/Whitepaper_detailed/guide/opcode_and_gas_guide.md†L1410-L1418】.

## Enterprise Integration and Compliance
- **Regulatory oversight** – `RegulatoryManager` catalogues jurisdiction‑specific rules and evaluates transactions, while `RegulatoryNode` flags non‑compliant addresses and records reasons for review【F:regulatory_management.go†L8-L65】【F:regulatory_node.go†L8-L43】. Identity tokens provide the link between flagged addresses and real‑world entities.
- **Sovereign deployments** – Central banking nodes embed the `IdentityService` to vet participants before distributing CBDCs, leveraging the same audit logs for governmental compliance【F:docs/Whitepaper_detailed/Central banks.md†L18-L19】.
- **Smart‑contract support** – Scaffolded DID registries in Solidity and Rust allow enterprises to integrate custom on‑chain identity workflows or extend SYN900 issuance logic【F:smart-contracts/solidity/DIDRegistry.sol†L1-L6】【F:smart-contracts/rust/src/did_registry.rs†L1-L10】.
- **Address suspension and whitelisting** – `ComplianceManager` can suspend, whitelist, or review transactions involving an address, preventing sanctioned entities from registering or transacting【F:compliance_management.go†L30-L77】.
- **Audit trails and risk scoring** – `ComplianceService` stores hashed KYC commitments, logs fraud signals and exposes risk scores and audit trails for investigators【F:compliance.go†L61-L90】【F:compliance.go†L94-L109】.
- **Role-based access** – `AccessController` grants, revokes and checks roles for operator wallets, supporting separation of duties in large organisations【F:access_control.go†L5-L24】【F:access_control.go†L38-L47】.

## Token Lifecycle and Compliance
- **Immutable records** – The `IdentityService` exposes `Info` and `Logs` to audit registered details and verification events without altering history【F:identity_verification.go†L60-L75】.
- **Ongoing verification** – Run `Verify` whenever documents expire or regulations change; updated logs keep the token in good standing.
- **Revocation and recovery** – If a wallet is compromised, re‑register a new address in the `IDRegistry` and repeat the verification flow to obtain a fresh token while preserving regulatory traceability【F:idwallet_registration.go†L20-L27】【F:identity_verification.go†L37-L57】.
- **Data privacy** – Only hashed or encrypted identifiers are stored on‑chain; personal documents remain off‑chain but linked via secure references.

## Security Best Practices
- Enable multi‑factor or hardware‑wallet authentication to guard the token.
- Monitor verification logs for unusual activity and escalate discrepancies to Neto Solaris support.
- Use offline storage for seed phrases and rotate credentials when personnel change.
- Review regulatory flags periodically via `RegulatoryNode.Logs` to ensure ongoing compliance【F:regulatory_node.go†L32-L43】.
- Track `ComplianceService.RiskScore` to detect wallets accumulating fraud signals and respond before regulatory thresholds are breached【F:compliance.go†L83-L90】【F:compliance.go†L94-L100】.

## Additional Resources
- Token architecture and additional standards are documented in the Synnergy Network token catalogue【F:docs/Whitepaper_detailed/Tokens.md†L184-L187】.
- Module guides and gas tables provide deeper insight into token operations and opcode budgeting【F:docs/Whitepaper_detailed/guide/module_guide.md†L460-L461】【F:docs/Whitepaper_detailed/guide/opcode_and_gas_guide.md†L1458-L1462】.
- For economic incentives such as Passive Income distributions, consult the transaction fee distribution guide which shows how the `IDRegistry` underpins targeted payouts【F:docs/Whitepaper_detailed/Transaction fee distribution.md†L106-L110】.

---
© Neto Solaris All rights reserved.
