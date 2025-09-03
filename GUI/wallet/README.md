# Synnergy GUI Wallet

Stage 31 introduces a lightweight TypeScript wallet that interacts with the
Synnergy CLI. It demonstrates how graphical interfaces can drive on-chain
operations without embedding private keys in the browser.

## Usage

```bash
npm install
npm start
```

The script will invoke the `synnergy` CLI to create a new encrypted wallet,
store it in `wallet.json`, and display its current balance.

## Integration

Future stages connect this UI to the `walletserver` backend for remote key
management. Commands are executed through the CLI, which exposes JSON output for
front-end consumption.
