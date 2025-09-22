# How to Deploy a Contract

Neto Solaris provides a deterministic smart‑contract environment in the
Synnergy Network. Contracts are compiled to WebAssembly (WASM) and deployed via
the `synnergy` command‑line interface (CLI), yielding an immutable address
derived from the bytecode hash. This guide describes how to publish and manage
contracts on the network.

## Prerequisites

- **Synnergy CLI** – build the `cmd/synnergy` binary or obtain a release build.
  Building from source ensures the toolchain matches the runtime used by the
  network:

  ```bash
  go build ./cmd/synnergy
  ./synnergy --help
  ```

  The binary exposes the `contracts`, `contract-mgr`, `xcontract` and
  `ai_contract` subcommands. Ensure the executable is on your `PATH` or export
  `BIN_PATH` for helper scripts.
- **Compiled contract** – a WASM file or template provided by the project.
- **Owner address** – the account that will control the contract.
- **Gas limit** – maximum execution gas per call. Defaults to `100000` but can
  be adjusted during deployment.
- **(Optional) Ricardian manifest** – JSON manifest providing human‑readable
  context for legal or compliance purposes.
- **Network access** – a node endpoint configured with sufficient permissions
  to accept deployment transactions.

### Stage 79 Bootstrap Support
Before deploying, operators can run `synnergy orchestrator bootstrap` to guarantee the heavy VM profile is started, consensus relayers are registered and ledger replication is active. The command invokes `core.EnterpriseOrchestrator.BootstrapNetwork`, sealing the bootstrap with the orchestrator wallet and wiring governance roles so contract calls immediately inherit enterprise-grade security guarantees.【F:cli/orchestrator.go†L58-L117】【F:core/enterprise_orchestrator.go†L71-L209】 Startup ensures the Stage 79 gas schedule is loaded beside Stage 78 costs, keeping deployment budgets predictable across automation, while the control panel exposes the same workflow for GUI-based operations.【F:cmd/synnergy/main.go†L63-L106】【F:web/pages/index.js†L1-L214】【F:web/pages/api/bootstrap.js†L1-L45】 Tests exercise unit, situational, stress and real-world bootstrap scenarios to confirm the environment remains deterministic before and after contract deployment.【F:core/enterprise_orchestrator_test.go†L73-L178】

## Step‑by‑Step Deployment

### 1. Compile the contract

Use the CLI to compile source WAT/WASM into deterministic bytecode. Under the
hood `CompileWASM` validates the source and computes a SHA‑256 digest that
serves as the future contract address:

```bash
synnergy contracts compile path/to/src.wat
```

### 2. Deploy the contract

The `deploy` subcommand registers the compiled bytecode with the on‑chain
registry. Provide the WASM path, optional manifest, gas limit and owner:

```bash
synnergy contracts deploy \
  --wasm path/to/contract.wasm \
  --ric path/to/manifest.json \
  --gas 150000 \
  --owner <owner_address>
```

The command reads the specified files, stores the manifest if supplied and
returns the contract address for subsequent interactions.

#### Ricardian manifest structure

A manifest is a small JSON document that links the bytecode to a human‑readable
description and legal terms. Typical fields include a title, descriptive text
and optional jurisdictional clauses:

```json
{
  "title": "StorageMarket",
  "description": "Decentralised file‑rental contract",
  "terms": "See https://blackridge.example/terms/storage"
}
```

Manifests are immutable once deployed and should be stored in version control
alongside the WASM to guarantee reproducibility.

### 3. Verify deployment

List deployed contracts or inspect a specific one:

```bash
synnergy contracts list
synnergy contracts info <address>
```

### 4. Invoke the contract

Execute exported methods with `contracts invoke`. Gas defaults to the limit set
at deployment unless overridden:

```bash
synnergy contracts invoke <address> --method transfer --args "<hex_bytes>" --gas 50000
```

Outputs include the return bytes and gas consumed, enabling deterministic
testing of business logic.

## Internal Architecture

Understanding how Synnergy processes contracts clarifies why addresses and gas
usage are deterministic.

### Contract registry

The registry persists contract metadata—owner, manifest, gas limit and paused
state—and derives addresses from the SHA‑256 hash of the WASM bytecode.
Concurrent access is guarded by read/write mutexes to ensure thread safety
during deployments and invocations【F:contracts.go†L10-L18】【F:contracts.go†L57-L78】.

### Virtual machine

Execution is handled by the `SimpleVM`, a lightweight interpreter that supports
resource‑tuned modes. Profiles range from `VMHeavy` for high throughput to
`VMSuperLight` for constrained devices. Each mode controls a concurrency limiter
and enforces gas caps before dispatching opcodes【F:core/virtual_machine.go†L11-L72】【F:core/virtual_machine.go†L96-L108】.

## Deploying Contract Templates

Synnergy ships precompiled templates for common use cases. Enumerate available
artifacts and deploy them with a single command:

```bash
synnergy contracts list-templates
synnergy contracts deploy-template \
  --name token_faucet \
  --owner <owner_address> \
  --gas 100000
```

Templates provide predictable gas costs and vetted bytecode, expediting
development for standard token faucets, storage markets, DAO governance modules,
NFT minting, AI model markets and more. Template names correspond to `.wasm`
files stored under `smart-contracts/`.

## Scripted Deployment

For automation, use the provided helper script. It verifies prerequisites and
invokes the CLI on your behalf:

```bash
scripts/deploy_contract.sh path/to/contract.wasm
```

The script ensures the binary exists and aborts on missing or invalid input. It
respects `BIN_PATH` to locate a custom CLI build and exits non‑zero if the
deployment fails, making it suitable for continuous‑integration pipelines.
Override the binary location by exporting `BIN_PATH` prior to invocation:

```bash
export BIN_PATH=/opt/synnergy/bin/synnergy
scripts/deploy_contract.sh contract.wasm
```

## Managing Deployed Contracts

After deployment, the registry and manager modules expose maintenance
operations:

- **Transfer ownership** – reassign a contract to a new account.
- **Pause or resume** – temporarily disable execution.
- **Upgrade bytecode** – replace the WASM and optionally adjust the gas limit.

These actions are available programmatically through the `ContractManager`
abstraction and enable controlled life‑cycle management without redeploying from
scratch. The CLI exposes these capabilities via the `contract-mgr` command:

```bash
synnergy contract-mgr transfer <addr> <new_owner>
synnergy contract-mgr pause <addr>
synnergy contract-mgr resume <addr>
synnergy contract-mgr upgrade <addr> <wasm_hex> <gas_limit>
synnergy contract-mgr info <addr>
```

Under the hood these commands call `ContractManager` methods that update the
registry while emitting OpenTelemetry spans. Transfer, pause, resume, upgrade
and info operations validate addresses and mutate ownership or bytecode in a
thread‑safe manner【F:core/contract_management.go†L20-L84】【F:core/contract_management.go†L86-L96】.

Each operation emits structured errors with telemetry spans, facilitating
enterprise monitoring and audit trails.

### Advanced Registries

Synnergy extends the base registry with specialised modules for complex
deployments:

- **AI contracts** – `ai_contract deploy` requires both WASM and a model hash;
  gas is validated against minimum thresholds before storing the mapping【F:core/ai_enhanced_contract.go†L38-L55】.
  Invoke inference with:

  ```bash
  synnergy ai_contract invoke <addr> <input_hex> <gas_limit>
  ```

- **Cross‑chain mappings** – `xcontract register` links a local contract to a
  remote chain and exposes `list`, `get` and `remove` operations. Each command
  optionally emits JSON for machine consumption【F:cli/cross_chain_contracts.go†L25-L98】.

- **Contract marketplace** – `marketplace deploy` uploads bytecode while
  enforcing gas pricing, and `marketplace trade` reassigns ownership via the
  embedded manager【F:core/smart_contract_marketplace.go†L11-L53】.

These registries extend the core deployment flow for specialised enterprise
scenarios without altering the deterministic base address.

## Enterprise Deployment Considerations

- **Environment separation** – maintain distinct testnet, staging and
  production nodes. Deterministic addresses make migrations predictable across
  environments.
- **Secrets management** – store owner keys in hardware modules or vaults and
  inject them at runtime rather than hard‑coding.
- **Continuous integration** – run `go test` and `mkdocs build` before
  promotion; the helper script's non‑zero exits integrate with CI failure
  thresholds.
- **Observability** – telemetry spans emitted by `ContractManager` operations
  feed into OpenTelemetry collectors for audit trails and performance analysis.
- **Disaster recovery** – export manifests and bytecode to version control so
  identical addresses can be recreated after catastrophic loss.
- **Throughput planning** – select a `SimpleVM` mode (`VMHeavy`, `VMLight`,
  `VMSuperLight`) that aligns with node capacity. Concurrency limits prevent a
  single deployment from exhausting resources【F:core/virtual_machine.go†L11-L72】.
- **Gas governance** – integrate `GasCost` lookups in CI to ensure supplied gas
  meets protocol‑defined minima before submitting transactions.

## Operational Monitoring and Error Codes

`contract-mgr` and related CLIs surface errors with machine‑parsable codes using
the internal error package, printing messages like `error (NOT_FOUND): contract
not found` for integration with alerting pipelines【F:cli/contract_management.go†L87-L93】.
Telemetry spans emitted from manager and marketplace operations provide latency
metrics and success/failure tags that feed directly into enterprise dashboards.

## Best Practices

- Audit contracts and manifests before deployment.
- Allocate sufficient gas to prevent out‑of‑gas failures.
- Record the returned contract address securely; it is derived from the SHA‑256
  hash of the bytecode and acts as the canonical identifier.
- Use version control and test suites to validate upgrades prior to production
  deployment.
- Tag releases with the contract address to maintain a clear provenance trail.
- Run `contracts info` regularly to verify manifests have not been tampered with
  and `contract-mgr info` to confirm owner and paused status.
- Consider using the AI and cross‑chain registries only after security reviews,
  as they expand the attack surface.

## Troubleshooting

- **`--wasm required`** – ensure the `--wasm` flag points to a valid file.
- **`contract already deployed`** – bytecode hash matches an existing contract.
  Modify the code or use `ContractManager` to upgrade.
- **`contract not found`** – verify the address when invoking or managing a
  contract.
- **`contract paused`** – resume via `contract-mgr resume` or review pause
  policy.
- **`wasm bytecode required`** – indicates empty or corrupted file during an
  upgrade. Re‑compile and retry.

## Conclusion

By following these steps, developers can confidently deploy and administer
smart contracts on the Synnergy Network under the Neto Solaris
umbrella. The deterministic tooling and accompanying templates streamline
contract lifecycles while preserving security and auditability. Whether rolling
out a single template or orchestrating a fleet of cross‑chain and AI contracts,
the platform provides a reproducible foundation for enterprise workflows.

