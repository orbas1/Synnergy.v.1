# Synnergy AI Marketplace

A minimal TypeScript interface for deploying AI-enhanced contracts through the Synnergy CLI.

## Usage

```
npm install
npm start <wasm_file> <model_hash> <manifest> <gas_limit> <owner>
```

The marketplace shells out to `synnergy ai_contract deploy` and prints the resulting contract address. Ensure the `synnergy` binary is on your `PATH` and a node is running.

## Testing

```
npm test
```
