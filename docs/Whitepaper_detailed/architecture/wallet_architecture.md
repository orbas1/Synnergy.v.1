# Wallet Architecture

## Overview
Wallet tooling provides key management and transaction capabilities for end users. A thin GUI wraps the `synnergy` CLI while a `walletserver` backend can broker requests for web or mobile clients.

## Key Components
- **CLI wallet commands** – generate keys, sign transactions and export encrypted JSON keystores.
- **walletserver** – lightweight HTTP server that forwards signed transactions to the network.
- **GUI wallet** – Electron/React interface that shells out to the CLI.

## Workflow
1. **Key generation** – users create wallets via the CLI which stores keys locally.
2. **Transaction building** – the GUI or server constructs transactions and requests signatures from the CLI.
3. **Broadcast** – signed transactions are submitted to nodes through the wallet server or directly via the CLI.
4. **Backup** – keys can be exported as encrypted JSON for cold storage.

## Security Considerations
- Keys never leave the user's machine unencrypted.
- Scrypt-derived keys with AES‑GCM encryption protect keystores at rest.
- The wallet server operates without access to private keys, reducing compromise risk.

## CLI Integration
- `synnergy wallet` – generate keys and sign transactions.
- `walletserver` binary – run an API-compatible backend for remote clients.
