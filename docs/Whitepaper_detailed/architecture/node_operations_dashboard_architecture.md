# Node Operations Dashboard Architecture

The Node Operations Dashboard provides operators with a web-based and CLI interface to monitor node health across the Synnergy network. It consumes REST endpoints exposed by authority nodes and displays aggregated status information.

Key components:

- **Client** – TypeScript application offering both CLI and browser interfaces.
- **Service layer** – fetches node metrics via `fetchNodeStatus` and related APIs.
- **Integration** – connects with the virtual machine and consensus modules to surface opcode and gas data.
