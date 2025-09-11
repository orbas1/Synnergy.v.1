# Desktop Shell

A lightweight Electron-based launcher that loads Synnergy GUI modules on demand.

## Features
- Centralized wallet/key authentication shared across all spawned modules.
- Dynamic menu generated from `modules.json` allowing multiple GUI windows to
  run concurrently.
- Each module window receives a session token via query parameter and remains
  stateless, relying on on-chain data.

## Configuration
1. Set a private key used for session token derivation:
   ```
   export DESKTOP_SHELL_PRIVATE_KEY="<hex-private-key>"
   ```
2. Adjust `modules.json` to match the URLs of running GUI modules.

## Development
```
npm install
npm start
```
Ensure dependent modules are running on the URLs defined in `modules.json`.
