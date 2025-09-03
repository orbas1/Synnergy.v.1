# Smart-Contract Marketplace GUI

This reference interface demonstrates how generic WebAssembly contracts can be
deployed and traded through the Synnergy CLI.

## Usage

```bash
# Deploy a contract
node src/main.ts deploy.wasm alice

# Trade ownership
node src/main.ts trade <address> bob
```

The implementation spawns the local `synnergy` binary and therefore inherits all
CLI configuration such as network endpoints and gas table settings.
