# AI Marketplace Architecture

## Overview
The AI Marketplace serves as a lightweight web interface for deploying AI‑enhanced smart contracts on the Synnergy Network. Rather than re‑implementing the full blockchain stack in the browser, the GUI acts as a guided wrapper around the `synnergy ai_contract deploy` command. This design keeps the web surface area small while still giving developers an approachable way to publish models.

## Component Layers
- **User Interface** – A React/TypeScript single page application that collects the WebAssembly module, model hash, manifest and gas limit. The form validates basic requirements such as file type, hash length and numeric ranges before enabling deployment.
- **CLI Adapter** – The browser passes the validated artifacts to a small Node service that shells out to the Synnergy CLI. The adapter streams stdout and stderr back to the UI so users can observe progress in real time.
- **Synnergy CLI** – Handles signing, consensus interaction and virtual machine integration. Because the CLI manages keys locally, the marketplace avoids exposing private keys to the browser.
- **Result Parser** – Once the CLI returns, the adapter extracts the contract address and returns it to the front‑end for display or further interaction.

## Deployment Workflow
1. **Package preparation** – Developers compile their model into WebAssembly and generate a manifest describing inputs, expected outputs and gas requirements.
2. **Upload and validation** – The user selects the module and manifest in the GUI. Client‑side validation checks file size limits and ensures the hash in the manifest matches the uploaded artifact.
3. **CLI execution** – The adapter writes the files to a temporary directory and invokes `synnergy ai_contract deploy --wasm <file> --manifest <file> --gas <limit>`.
4. **Network confirmation** – The CLI broadcasts the transaction and waits for inclusion in a block. Progress is streamed to the UI so the user can cancel or retry if necessary.
5. **Address delivery** – On success, the resulting contract address is displayed and stored for quick access from the marketplace.

## Data and Storage
- Uploaded artifacts are kept only long enough to complete deployment and are deleted after confirmation to avoid lingering sensitive data.
- Future iterations may integrate an artifact registry so popular models can be reused without re‑uploading large binaries.

## Security Considerations
- **No in‑browser keys** – All signing occurs within the CLI, eliminating the risk of key theft through the web layer.
- **Strict validation** – Both client and adapter validate manifest fields, file sizes and command arguments to prevent command injection or resource exhaustion.
- **Rate limiting** – The adapter can throttle requests per user to discourage denial‑of‑service attacks.
- **Audit logging** – Every deployment request and CLI response is logged with a request identifier, providing traceability for governance or billing.

## Future Enhancements
- Marketplace listings that allow developers to publish and monetize their models after deployment.
- Integration with on‑chain inference services so contracts can call models hosted by third parties.
- Richer progress reporting, including estimated gas costs and block confirmations.
- Hooks for community review and governance to vet models before they become discoverable.

This architecture prioritizes simplicity and security while leaving room for the marketplace to grow into a full ecosystem for AI contracts.
