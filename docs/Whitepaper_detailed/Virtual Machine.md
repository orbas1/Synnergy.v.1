# Virtual Machine

## Overview
The Synnergy Virtual Machine (SVM) is the execution backbone of Blackridge Group Ltd.'s blockchain ecosystem. It provides a deterministic environment for running WebAssembly-based smart contracts, ensuring that every node evaluates transactions with identical results. The SVM's design emphasises modularity, resource efficiency and security so that the platform can support applications ranging from high-throughput financial services to lightweight mobile interactions.

## Execution Model
At its core, the SVM interprets contracts encoded as sequences of 24-bit opcodes. Each opcode maps to a handler function responsible for transforming input bytes into new state. Unknown instructions default to a no-operation echo handler, keeping execution deterministic even when the opcode set evolves. Execution is driven through `ExecuteContext(ctx, wasm, method, args, gasLimit)` so callers can provide a cancellation context, raw opcode stream, optional method hint, initial arguments and a maximum gas allowance.

The virtual machine lifecycle is straightforward:

1. **Start** – initialise internal resources and mark the instance as running.
2. **Execute** – interpret bytecode, dispatching each opcode in sequence while accounting for gas consumption and respecting context cancellation.
3. **Stop** – release resources when the instance is no longer required.
4. **Status** – expose whether the engine is currently active.

This lifecycle conforms to the `VirtualMachine` interface used by the contract registry, enabling different VM implementations to be swapped without affecting higher-level modules or external tooling.

## Resource Profiles
To serve diverse deployment scenarios, the SVM supports three resource profiles identified in code as `VMHeavy`, `VMLight` and `VMSuperLight`:

- **Heavy** – allows up to ten concurrent executions for authority nodes or data centres that handle significant throughput.
- **Light** – the default profile permitting five parallel executions, suitable for general-purpose validator nodes.
- **Super Light** – serialises execution to a single request at a time, optimised for constrained environments such as mobile clients.

Profiles are enforced through a buffered channel that acts as a concurrency limiter. Calls exceeding capacity immediately return "vm busy", giving operators predictable back-pressure and preventing resource exhaustion.

## Instruction Encoding and Opcode Registry
Opcode words are big-endian 24-bit values, yielding over sixteen million possible instructions. Opcode `0x000000` is reserved for a default NOP/echo handler used whenever a specific code is unassigned. The registry `SNVMOpcodes` ships with handlers for all core modules—covering artificial intelligence tooling, cross-chain workflows, token economics and node operations. Developers can bind additional behaviour at runtime through `RegisterOpcode`, allowing the ecosystem to evolve without rebuilding the VM.

Key properties of the opcode system include:

- **Deterministic Defaults** – unrecognised codes route to an echo handler, ensuring repeatable execution across versions.
- **Modular Coverage** – every network service registers its operations, enabling a single bytecode stream to orchestrate cross-cutting features.
- **Extensibility** – runtime registration permits domain‑specific opcodes and experimental features without destabilising existing contracts.

## Opcode Catalogue and Inter-module Coverage
The opcode catalogue acts as a connective tissue across the Synnergy Network. Thousands of opcodes are pre-registered, giving bytecode the ability to drive:

- **AI and Analytics** – model publication, fraud prediction and baseline drift monitoring.
- **Cross-chain Bridges** – lock–mint and burn–release flows, sidechain registration and Plasma exits.
- **Token and Identity Management** – minting, staking, compliance checks and wallet operations.
- **Node Infrastructure** – provisioning mining, staking, light and watchtower nodes, alongside system health logging and firewall rules.
- **Governance and Compliance** – KYC validation, risk scoring and regulatory audit trails.
- **Security Services** – firewall configuration, biometric enrolment and zero‑trust data channels.

Because each opcode resolves to a handler within the registry, contracts orchestrate these disparate services without embedding network-specific logic. This modularity also allows Blackridge Group Ltd. to introduce new capabilities simply by registering additional opcodes.

## Lifecycle and Thread Safety
The SVM is engineered for concurrent workloads. Internally, a read/write mutex guards lifecycle transitions so `Start` and `Stop` can be invoked repeatedly without race conditions. Each execution acquires a token from a buffered channel sized according to the chosen profile, ensuring that no more than the allowed number of contracts run simultaneously. This pattern delivers predictable throughput while preserving determinism across validator implementations.

## Detailed Execution Flow
When `ExecuteContext` is called, the VM:

1. Validates that the instance is running and that bytecode is present.
2. Reserves a slot in the concurrency limiter, respecting any context cancellation during the wait.
3. Calculates the gas requirement based on the number of 24‑bit opcodes and checks it against the supplied limit.
4. Iterates through the bytecode three bytes at a time, resolving each opcode to a registered handler or the default echo function.
5. Aggregates the output of each handler as the new state.
6. Introduces a deterministic one‑millisecond delay before returning to align timing across nodes.

This deterministic flow means that any validator executing the same bytecode with the same inputs will arrive at the same output and gas usage.

## Gas Accounting and Cost Management
The SVM measures gas as the number of opcodes executed, with a minimum charge of one unit to cover empty payloads. For fine-grained pricing, the `GasTable` maps mnemonic names to explicit costs, loaded lazily from `docs/reference/gas_table_list.md` via a thread-safe cache. Runtime functions allow:

- **Lookup** – `GasCost` retrieves the fee for a given opcode name, returning a `DefaultGasCost` when unknown.
- **Validation** – `HasOpcode` verifies that a cost is defined.
- **Dynamic Pricing** – `RegisterGasCost` injects new fees and `ResetGasTable` clears cached values so tests and tooling can model evolving economics.

All operations are guarded by mutexes and instrumented with OpenTelemetry spans and structured JSON logs, keeping fee calculation transparent while preserving deterministic execution across nodes.

## Sandbox Management
For additional isolation, the Synnergy platform employs a sandbox manager that provisions dedicated execution environments for smart contracts. Each sandbox tracks:

- **Identifier and Contract Address** – unique references for managing instances.
- **Gas and Memory Limits** – execution boundaries that prevent contracts from monopolising resources.
- **Lifecycle Timestamps** – creation and last-reset times for auditability.

The manager exposes thread-safe operations to start, stop, reset and delete sandboxes, alongside queries for individual status or a complete inventory. By recording creation and reset timestamps, operators can audit lifecycle events and periodically recycle environments to enforce policy or reclaim resources.

Administrative tooling can query `SandboxStatus` for targeted inspection or `ListSandboxes` for a network-wide view, enabling fleet oversight and automated maintenance scripts.

## Error Handling and Recovery
Clear error semantics aid debugging and resilience. Typical responses include:

- **"vm not running"** when execution is attempted before `Start`.
- **"vm busy"** if the concurrency limiter is saturated.
- **"bytecode required"** for empty payloads.
- **"gas limit exceeded"** when estimated cost surpasses the caller's allowance.
- Propagated context cancellations and handler-specific failures.

These messages provide actionable insight, allowing operators to distinguish between transient contention and genuine contract faults.

## Observability and Auditability
Enterprise deployments require deep insight into runtime behaviour. The SVM and its supporting utilities emit structured JSON logs through an internal logger and expose OpenTelemetry tracing hooks. Gas table loading and opcode execution can therefore be correlated across nodes, while sandbox timestamps and reset events provide an auditable trail for compliance reviews.

## Integration with Synnergy Network
The contract registry relies on the SVM for every deployment and invocation. By abstracting execution behind the `VirtualMachine` interface, other modules—such as cross-chain bridges, data services and token registries—invoke contracts without needing to understand low-level execution details. The opcode table also binds many system functions, allowing the VM to orchestrate actions across the wider Synnergy Network in a uniform manner.

## Security and Determinism
Blackridge Group Ltd. designs the SVM with several safeguards:

- **Context-aware Execution** – operations honour cancellation signals, ensuring stalled calls do not linger during network disruptions.
- **Concurrency Limits** – resource profiles shield nodes from overload and deliberate denial-of-service attempts.
- **Deterministic Delay** – a fixed millisecond pause at the end of each run aligns timing across nodes and reduces timing-based side channels.
- **Graceful Error Handling** – empty payloads, gas overflows or saturated concurrency queues return descriptive errors while unknown opcodes fall back to the default handler.

## Enterprise Integration and Deployment
Because contracts compile to WebAssembly, developers can author logic in any language with a WASM toolchain while maintaining deterministic behaviour on the network. Dynamic opcode and gas registration allow phased rollouts of new modules, and profile-driven concurrency controls let operators tune resource usage for diverse hardware tiers. Combined with structured logging, telemetry and sandbox auditing, these capabilities equip Blackridge Group Ltd. with a production-ready execution layer suitable for mission-critical applications.

## Use Cases and Deployment Scenarios
- **Authority Nodes** leverage the heavy profile to process large batches of transactions and complex contract interactions.
- **Validator Nodes** typically run the light profile, balancing responsiveness with resource usage.
- **Edge and Mobile Clients** utilise the super light profile to verify or execute small workloads without draining local resources.

In all cases, the unified VM architecture ensures that applications behave consistently regardless of where they run within the Synnergy Network.

## Conclusion
The Synnergy Virtual Machine delivers a robust, extensible and resource-aware execution layer for the platform's smart contracts. Through its opcode-driven design, sandbox management, observability stack and multiple performance profiles, the SVM embodies Blackridge Group Ltd.'s commitment to secure and scalable decentralised computing.
