# Network First Initiation Scripts

The first initiation scripts bootstrap a Synnergy network by creating the Synthron (SYNN) system coin, genesis wallet, required treasuries and initial authority node keys. Coin parameters such as name, symbol, max supply and genesis reward are loaded from `configs/coin.json` to ensure consistent network settings. Each script emits a CSV file containing the generated addresses and private keys; keep this file secure.

## Mainnet
Run the mainnet bootstrap:

```bash
./scripts/Mainnet_first_initiation.sh
```

* Keys stored at `data/mainnet/init_keys.csv`
* Sample nodes listen on `127.0.0.1:30301-30303`
* Treasuries registered with loanpool and charity modules when the `synnergy` CLI is available

## Testnet
Initialize a public testing environment:

```bash
./scripts/Testnet_first_initiation.sh
```

* Keys stored at `data/testnet/init_keys.csv`
* Nodes listen on `127.0.0.1:40301-40303`
* Treasuries registered with loanpool and charity modules when the `synnergy` CLI is available

## Devnet
For local development, use the devnet script:

```bash
./scripts/Devnet_first_initiation.sh
```

* Keys stored at `data/devnet/init_keys.csv`
* Nodes listen on `127.0.0.1:50301-50303`
* Treasuries registered with loanpool and charity modules when the `synnergy` CLI is available

## Output
All scripts generate a CSV with columns `entity,address,private_key`. The file is **not** distributed; store it safely and restrict access.

To connect to the nodes, use:

```bash
synnergy node dial <address>:<port>
```

Replace `<address>` with the server's IP and `<port>` with one of the listener ports above.
