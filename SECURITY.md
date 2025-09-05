# Security

This project uses several static analysis tools to harden the codebase:

- `staticcheck` for general code issues.
- `gosec` for security-oriented scanning.
- `govulncheck` for dependency vulnerability scanning.

Run `make security` to execute all security checks locally.

Sensitive model artifacts are encrypted using AES‑GCM with 32‑byte keys and
all transactions are signed using standard ECDSA digital signatures. Review
`scripts/pki_setup.sh` for managing certificates and keys.
