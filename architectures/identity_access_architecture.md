# Identity and Access Architecture

Modules in this group manage user identities, authentication, and permission control for network operations.

**Key Modules**
- access_control.go
- identity_verification.go
- idwallet_registration.go
- biometric_security_node.go
- biometrics_auth.go
- address_zero.go

**Related CLI Files**
- cli/identity.go
- cli/idwallet.go
- cli/access.go
- cli/address.go
- cli/address_zero.go
- cli/biometric_security_node.go
- cli/biometrics_auth.go

These components provide secure onboarding and identity enforcement across the platform. The access controller exposes
thread‑safe functions (`GrantRole`, `RevokeRole`, `HasRole`, `ListRoles`) used by the CLI and virtual machine to enforce
permissions at runtime. Biometric modules hash templates and bind them to
ECDSA public keys so that all enrollments and authentications are
cryptographically signed and tamper evident.
