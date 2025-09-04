# Network

The Synnergy network is composed of heterogeneous nodes coordinated through a
gas-aware consensus layer.  Stage 46 introduces an end-to-end network harness
used during development to automatically spin up wallet services and in-memory
nodes, broadcast transactions and assert propagation.  The harness provides a
repeatable way to validate fault tolerance, address handling and gas
accounting across the virtual machine and CLI interfaces.
