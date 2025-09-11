# Identity and Access Architecture

## Overview
Identity services provide secure onboarding and fine‑grained access control. They issue identity tokens, manage credential proofs and enforce permissions across nodes and smart contracts.

## Key Modules
- `identity_verification.go` – validates user documents and links them to on-chain identities.
- `idwallet_registration.go` – registers identity wallets and stores public keys.
- `access_control.go` – checks permissions for CLI commands and contract calls.
- `biometrics_auth.go` – optional biometric login for high-assurance nodes.
- `biometric_security_node.go` – dedicated node for processing biometric templates.

## Workflow
1. **Registration** – users submit credentials processed by `identity_verification`.
2. **Wallet issuance** – `idwallet_registration` creates identity wallets and issues SYN900 tokens.
3. **Permission assignment** – roles are defined and enforced by `access_control`.
4. **Authentication** – nodes may request biometric factors via `biometrics_auth` for sensitive operations.

## Security Considerations
- Personal data is hashed or encrypted before storage.
- Biometric templates are processed on specialized nodes and never leave secure memory.
- Access control checks are logged to detect misuse or privilege escalation.

## CLI Integration
- `synnergy idwallet` – register identity wallets.
- `synnergy access` – manage roles and permissions.
