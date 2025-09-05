# Wallet Admin Interface

Enterprise-grade GUI for Wallet Admin Interface. The service exposes a secure
HTTP API for managing wallets and verifying signatures. It is designed to be
driven from the Synnergy CLI and integrates with the virtual machine through
JSON responses.

## Scripts

- `npm run build` – compile TypeScript to `dist/`
- `npm start` – launch the compiled server
- `npm test` – execute unit tests via Jest and Supertest
- `npm run lint` – run ESLint
- `npm run format` – format sources with Prettier

## Deployment

The repository includes Docker and Kubernetes manifests for containerized
deployments as well as a minimal CI pipeline for automated builds and tests.
