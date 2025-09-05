# Wallet Admin Interface Documentation

This interface exposes administrative wallet functionality through a hardened
HTTP API. It is designed for integration with the Synnergy CLI and virtual
machine by providing JSON‑based endpoints that can be driven from automated
workflows.

## Endpoints

- `GET /health` – lightweight health probe used by orchestration tools and
  Kubernetes readiness checks.
- `POST /verify` – verifies a message signature against a provided public key.
  Requests are validated and rejected when malformed.

## Security

The server applies [helmet](https://github.com/helmetjs/helmet) for sane
defaults, validates request bodies and runs behind TLS in production. All
operations are intended to be signed client‑side to preserve privacy and
integrity.

## Integration

Deployments can be automated using the accompanying Docker and Kubernetes
manifests and exercised through Jest and Supertest based unit tests.
