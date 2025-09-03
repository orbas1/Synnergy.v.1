# Synnergy Explorer

A minimal TypeScript explorer that queries the Synnergy blockchain via the CLI.

## Usage

```
npm install
npm start
```

The explorer calls `synnergy ledger head` to display the chain height and latest
block hash. Ensure the `synnergy` binary is on your `PATH` and a node is
running.

## Testing

```
npm test
```
