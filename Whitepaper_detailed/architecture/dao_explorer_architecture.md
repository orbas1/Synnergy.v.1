# DAO Explorer Architecture

The DAO explorer provides a thin coordination layer between the Synnergy CLI and
its core DAO governance runtime. The GUI invokes `synnergy dao` commands to
create organisations, manage membership and inspect governance details.

* **Virtual machine integration** – DAO creation and membership updates execute
  through the shared VM and sandbox mechanisms used by other contracts,
  guaranteeing deterministic behaviour.
* **Consensus and wallet integration** – operations consume gas via
  `CreateDAO`, `JoinDAO`, `LeaveDAO` and `DAOInfo` opcodes so wallets and
  consensus nodes can account for costs.
* **Fault tolerance** – the explorer is stateless, delegating persistence and
  concurrency control to the CLI and core modules. Higher level services can
  wrap it for replication or sharding in enterprise deployments.

The GUI inherits existing authentication, node and authority policies from the
CLI, allowing it to slot into dashboards without bespoke plumbing.
