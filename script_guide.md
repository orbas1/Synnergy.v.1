
# Synnergy Script Guide

This guide documents the utility scripts located in `cmd/scripts`. Each script wraps common `synnergy` CLI commands to showcase typical workflows and to help automate repetitive tasks during development. The collection now includes helpers for the faucet service, DAO voting, marketplace listings and storage deals. A working Go toolchain and the compiled `synnergy` binary are required.

## Prerequisites

1. Run `./setup_synn.sh` to install Go and fetch project dependencies.
2. Optionally execute `./Synnergy.env.sh` to configure additional environment variables and tools.
3. Build the CLI with `./build_cli.sh` or by invoking `go build -o synnergy ./cmd/synnergy` from the `synnergy-network` directory.

## Running the Demo Network

The `start_synnergy_network.sh` script launches a small local network and demonstrates basic commands:

```bash
./start_synnergy_network.sh
```

This script performs the following steps:

1. Compiles the CLI from `cmd/synnergy/main.go`.
2. Starts the networking, consensus, replication and VM services in the background.
3. Executes a sample `~sec merkle` command against demo data to show the security module.
4. Waits until you press `Ctrl+C`, then terminates all services.

Ensure the `synnergy` binary is present in the working directory and that `go` is on your `PATH` before running the demo.

## Script Reference

The files below illustrate how to invoke individual command groups. Arguments in square brackets are optional and default values are taken from each script.

### build_cli.sh
Compile the `synnergy` CLI with trimmed debug paths.

```bash
./build_cli.sh
```

### network_start.sh
Start the networking daemon.

```bash
./network_start.sh
```

### network_peers.sh
List currently connected peers.

```bash
./network_peers.sh
```

### consensus_start.sh
Launch the consensus service.

```bash
./consensus_start.sh
```

### replication_status.sh
Query replication daemon status.

```bash
./replication_status.sh
```

### vm_start.sh
Run the WebAssembly virtual machine daemon.

```bash
./vm_start.sh
```

### coin_mint.sh
Mint SYNN coins to a target address. Usage:

```bash
./coin_mint.sh <address> [amount]
```

### token_transfer.sh
Transfer tokens between two addresses. Usage:

```bash
./token_transfer.sh <token> <from> <to> [amount]
```

### contracts_deploy.sh
Deploy a smart contract from a WASM file.

```bash
./contracts_deploy.sh <file.wasm>
```

### wallet_create.sh
Create a new HD wallet file. Usage:

```bash
./wallet_create.sh [output.json] [password]
```

### transactions_submit.sh
Submit a signed transaction JSON blob.

```bash
./transactions_submit.sh <tx.json>
```

### security_merkle.sh
Compute a Merkle root for auditing. Example:

```bash
./security_merkle.sh "deadbeef,baadf00d"
```

### governance_propose.sh
Create a governance proposal.

```bash
./governance_propose.sh [title] [body.md]
```

### cross_chain_register.sh
Register a cross-chain bridge relayer.

```bash
./cross_chain_register.sh <srcChain> <dstChain> <relayerAddress>
```

### rollup_submit_batch.sh
Submit an optimistic roll-up batch.

```bash
./rollup_submit_batch.sh <batch.json>
```

### sharding_leader.sh
Query the current shard leader.

```bash
./sharding_leader.sh
```

### sidechain_sync.sh
List registered side-chains.

```bash
./sidechain_sync.sh
```

### fault_check.sh
Capture a fault-tolerance snapshot.

```bash
./fault_check.sh
```

### state_channel_open.sh
Open a payment channel. Usage:

```bash
./state_channel_open.sh <from> <to> [amount]
```

### storage_pin.sh
Pin a file in the storage subsystem.

```bash
./storage_pin.sh <file>
```

### ../../scripts/devnet_start.sh
Spin up a multi-node developer network. Example:

```bash
../../scripts/devnet_start.sh 3
```

### ../../scripts/testnet_start.sh
Launch a testnet from a YAML config:

```bash
../../scripts/testnet_start.sh path/to/testnet.yaml
### faucet_fund.sh
Request coins from the local faucet service.

```bash
./faucet_fund.sh <address>
```

### dao_vote.sh
Cast a vote on a governance proposal.

```bash
./dao_vote.sh <proposal-id> [approve]
```

### marketplace_list.sh
Create an AI marketplace listing for a model.

```bash
./marketplace_list.sh <price> <cid>
```

### storage_marketplace_pin.sh
Pin data and publish a storage listing.

```bash
./storage_marketplace_pin.sh <file> [provider] [price] [capacity]
```

### loanpool_apply.sh
Submit a loan proposal transaction.

```bash
./loanpool_apply.sh <creator> <recipient> [type] [amount] [desc]
```

### authority_apply.sh
Register an authority-node candidate.

```bash
./authority_apply.sh <address> [role]
```

These scripts are intentionally minimal to keep the focus on demonstrating CLI usage. Feel free to modify them or combine commands to suit your workflow.
