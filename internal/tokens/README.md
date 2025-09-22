# Tokens

The tokens package implements the enterprise token stack used by the Synnergy
CLI, consensus services and the function web. Major capabilities include:

* **Concurrent-safe primitives** – the `BaseToken` exposes hooks, deterministic
  account snapshots and hardened validation (overflow detection, max supply
  enforcement and zero-amount guards) to keep ledger operations fault tolerant.
* **Registry observability** – the token registry and SYN1000 index both emit
  lifecycle events so the CLI and web UI can render live dashboards without
  polling internal state.
* **Specialised token models** – SYN10 CBDC tokens track issuer metadata and an
  auditable exchange-rate history, while SYN1000 stablecoins expose reserve
  breakdowns, collateralisation ratios and strict reserve validation.
* **CLI alignment** – reserve-management helpers now surface actionable errors
  for command line operators, avoiding silent failures during onboarding or
  compliance workflows.

Refer to the associated tests for examples that cover concurrency, governance
flows and reserve valuation scenarios.
