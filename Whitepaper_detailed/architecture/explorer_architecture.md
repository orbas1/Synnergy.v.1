# Explorer Architecture

The Explorer GUI is a thin TypeScript layer that shells out to the `synnergy`
CLI.  It issues read-only commands such as `synnergy ledger head` and parses the
structured output.  The GUI avoids exposing private keys or node ports and can be
deployed alongside the CLI or remote nodes with SSH access.  It serves as a
foundation for richer dashboards that may fetch blocks, transactions and node
telemetry through authenticated CLI calls.
