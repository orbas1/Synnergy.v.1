# Wallet Server

The wallet server exposes a minimal HTTP interface for generating
wallets used by the Synnergy network. It provides endpoints for
health checks and wallet creation and is intended to be consumed by
CLI tools and future GUI applications.

## Building and Running

```bash
go run ./walletserver
```

By default the server listens on `:8080`.

## API

### `GET /health`
Returns `{ "status": "ok" }` when the service is running.

### `POST /wallet/new`
Generates a new wallet and responds with a JSON object containing the
wallet address.

```json
{ "address": "<hex address>" }
```

