# GUI Quickstart

Synnergy ships several graphical interfaces for interacting with the network, including wallets, explorers and marketplaces. All GUIs are built with Node.js and React.

## Installing Dependencies

- Node.js 18+
- yarn or npm

Install dependencies in a GUI module:

```bash
cd GUI/wallet
yarn install
```

## Running in Development

```bash
yarn start
```

The development server reloads on file changes and proxies API requests to the configured backend.

## Building for Production

```bash
yarn build
```

Outputs a static bundle under `dist/` ready to be served by the wallet server or any HTTP server.

## Environment Variables

Each GUI reads configuration from `.env` files. Common settings include API endpoints and feature flags. Copy `.env.example` to `.env` and adjust as needed.

## Docker

Most GUIs include a `Dockerfile` for containerised deployment:

```bash
docker build -t syn-wallet GUI/wallet
```

## Available Modules

- `wallet` – manage accounts and send transactions
- `explorer` – query ledger state
- `ai-marketplace` – browse and purchase AI models
- `nft_marketplace` – mint and trade NFTs

Other interfaces follow the same commands shown above.

