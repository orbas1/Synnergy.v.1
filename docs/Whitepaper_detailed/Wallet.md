# Wallet

Neto Solaris’s Synnergy Network wallet is the primary interface for creating, securing and authorizing value on-chain. It combines hardened cryptography with flexible tooling so developers, institutions and end users can transact confidently across the ecosystem.

## Cryptographic Foundation
- **Key Generation:** Each wallet is built on an ECDSA key pair derived from the P-256 curve. The public key is hashed with SHA-256 and truncated to 20 bytes, yielding a 40-character hexadecimal address that uniquely identifies the account on-chain【F:core/wallet.go†L26-L41】.
- **Signature Workflow:** Transactions are hashed and signed using the wallet’s private key. Counterparties verify the signature against the signer’s public key, ensuring non-repudiation and tamper resistance【F:core/wallet.go†L44-L68】.

## Secure Key Storage
- **Encrypted Wallet Files:** Private keys can be written to disk using AES‑256 GCM encryption with an `scrypt`-derived key. The process stores the address, salt, nonce and ciphertext, safeguarding the wallet against offline attacks【F:core/wallet.go†L71-L108】.
- **Recovery and Portability:** Encrypted files are reopened with the correct passphrase, reconstituting the P-256 key pair and address so wallets can move between environments without exposing secrets【F:core/wallet.go†L111-L150】.

## Genesis Allocation and Auditing
Synnergy’s genesis event pre-allocates funds to well-known roles using deterministic addresses derived from hashed labels. These genesis wallets cover internal development, charity initiatives, loan pools, validator rewards and more, enabling transparent auditing from day one【F:core/genesis_wallets.go†L8-L41】. Fee inflows can then be programmatically split across those addresses by `AllocateToGenesisWallets`, which maps the network’s distribution policy to concrete wallet targets【F:core/genesis_wallets.go†L44-L58】【F:core/fees.go†L101-L108】.

Operational scripts such as `mainnet_setup.sh` exercise this model by generating distribution, reserve and treasury wallets and computing initial allocations before nodes are launched【F:scripts/mainnet_setup.sh†L25-L40】.

## Stage 80 Treasury Integration
The Synthron Treasury orchestrator initialises a dedicated wallet during bootstrap, credits it with the genesis allocation and registers it with consensus relayers and authority nodes so monetary policy actions are authenticated end-to-end. The wallet exposes mint, burn, transfer, operator governance and bridge flows through `synnergy coin telemetry`, giving operators a single entry point for treasury issuance with optional JSON diagnostics for governance pipelines【F:treasury/synthron_treasury.go†L41-L215】【F:cli/coin.go†L23-L130】. Telemetry snapshots include wallet address, minted and burned totals, circulating supply, operator roster, ledger height, subsystem health and a signed audit trail, ensuring stakeholders can audit treasury movements across CLI, web dashboards and automated tests while guarding against unauthorised actions and proving provenance for every signature【F:treasury/synthron_treasury.go†L214-L612】【F:docs/Whitepaper_detailed/Synthron Coin_test.go†L1-L108】【F:web/pages/index.js†L1-L260】.

## Interfaces and Tooling
### Command-Line Interface
The `synnergy` CLI includes a `wallet` module for deterministic wallet creation. Operators can output an encrypted file and specify the encryption password at generation time, enabling scripted provisioning and secure backups【F:cli/wallet.go†L8-L34】.

Unit tests assert that the command emits JSON with a 40‑character address and writes the encrypted wallet to disk, providing a safety net for automated pipelines【F:cli/wallet_cli_test.go†L9-L30】.

### Wallet Server
A lightweight HTTP service exposes `/health` and `/wallet/new` endpoints. The server uses the core wallet library to generate addresses on demand and returns them as JSON, allowing GUIs or external tooling to request wallets without direct key management【F:walletserver/handlers.go†L15-L33】【F:walletserver/main.go†L9-L18】.

Handler tests cover both endpoints to ensure consistent responses and proper error handling before the service is deployed【F:walletserver/handlers_test.go†L10-L35】.

### GUI and Test Harness
Integration tests demonstrate a GUI-driven workflow where the wallet server issues addresses, the CLI validates them and transactions are signed end-to-end. These tests ensure every interface remains interoperable across releases【F:tests/gui_wallet_test.go†L39-L72】.

## Funding and Test Networks
New wallets can be seeded without manual token transfers through a rate‑limited faucet. The core faucet tracks balances, dispense amounts, cooldowns and per‑address timestamps to prevent abuse【F:core/faucet.go†L9-L56】. Operators initialise and interact with the faucet via CLI commands for setup, requests, balance checks and configuration updates, enabling rapid developer onboarding【F:cli/faucet.go†L13-L83】. Unit tests validate cooldown enforcement and configuration changes to ensure predictable behaviour during network bootstrapping【F:core/faucet_test.go†L8-L24】.

## Identity and Compliance
Regulated deployments often require wallets to be tied to verified identities. The `IDRegistry` component records metadata for wallet addresses, allowing applications to confirm registration status before granting privileges【F:idwallet_registration.go†L8-L44】. The `idwallet` CLI wraps this registry so operators can register and query identity-tagged wallets at the command line【F:cli/idwallet.go†L13-L44】, with an accompanying automation stub ready for integration into onboarding workflows【F:scripts/idwallet_register.sh†L1-L17】.

Fee schedules further reinforce governance; `FeeForWalletVerification` estimates the cost of higher assurance checks so compliance programs can budget accurately【F:core/fees.go†L84-L88】. For sanctioned or compromised accounts, the network-wide firewall can block wallet addresses, tokens or peer IPs in real time, and an accompanying CLI allows operators to list, allow or block IPs on demand【F:firewall.go†L5-L44】【F:cli/firewall.go†L13-L58】.

## Advanced Custody Models
For high-value accounts, Synnergy supports multisignature authorisation through a Solidity `MultisigWallet` contract where owners submit, approve and execute transactions once a quorum is met【F:smart-contracts/solidity/MultisigWallet.sol†L4-L59】. The repository also lays groundwork for session-key and smart contract-based wallets, with templates that will expand into programmable spending policies【F:smart-contracts/solidity/SessionKeyWallet.sol†L1-L6】【F:smart-contracts/solidity/SmartWallet.sol†L1-L6】. Automation stubs are provided to bootstrap multisignature setups in deployment scripts, allowing enterprises to codify approval flows【F:scripts/wallet_multisig_setup.sh†L1-L17】.

## Enterprise Automation
The repository includes a suite of shell scripts that automate key lifecycle tasks such as initialization, rotation, multisignature configuration, offline signing and hardware wallet integration. These utilities provide a foundation for institutional policy enforcement and cold‑storage workflows【F:scripts/wallet_init.sh†L1-L17】【F:scripts/wallet_key_rotation.sh†L1-L17】【F:scripts/wallet_multisig_setup.sh†L1-L17】【F:scripts/wallet_offline_sign.sh†L1-L17】【F:scripts/wallet_hardware_integration.sh†L1-L17】. Additional stubs prepare automated wallet server deployments and end‑to‑end mainnet bootstraps where CLI tools create distribution wallets and compute genesis allocations【F:scripts/wallet_server_setup.sh†L1-L17】【F:scripts/mainnet_setup.sh†L19-L40】.

## Deployment and Operations
Production environments can run the wallet server as a replicated Kubernetes deployment with resource quotas and health probes. The `wallet.yaml` manifest provisions two pods from the `synnergy/walletserver` image, exposes port 8081, mounts configuration and defines readiness and liveness endpoints for self-healing clusters【F:deploy/k8s/wallet.yaml†L1-L67】. This operational blueprint enables horizontal scaling and consistent lifecycle management across cloud providers.

## Verification and Test Coverage
Comprehensive tests validate wallet functionality at every layer. Core tests confirm address uniqueness, signature integrity and encrypted file round‑trips【F:core/wallet_test.go†L8-L56】. CLI tests verify that `wallet new` outputs valid JSON and persists encrypted keys【F:cli/wallet_cli_test.go†L9-L30】, while server tests exercise the HTTP endpoints for health checks and wallet generation【F:walletserver/handlers_test.go†L10-L35】. Integration tests drive the server and CLI together to mirror a production GUI workflow【F:tests/gui_wallet_test.go†L39-L72】.

## Best Practices
- Use strong, unique passwords for encrypted wallet files and store them separately from key backups.
- Rotate keys periodically and leverage multisignature arrangements for high-value accounts.
- Employ offline signing or hardware wallets when operating in adversarial environments.
- Register wallets with the identity registry when operating under compliance regimes.
- Apply firewall rules to block sanctioned or suspicious addresses.

By combining robust cryptography, layered security controls and a spectrum of automation tools, the Synnergy wallet empowers organizations to interact with the Neto Solaris blockchain confidently while meeting stringent operational and regulatory requirements.
