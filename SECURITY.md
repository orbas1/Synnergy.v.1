# Security

This project uses several static analysis tools to harden the codebase:

- `staticcheck` for general code issues.
- `gosec` for security-oriented scanning.
- `govulncheck` for dependency vulnerability scanning.

Run `make security` to execute all security checks locally.
