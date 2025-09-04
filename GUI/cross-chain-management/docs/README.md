# Cross Chain Management Documentation

This module provides a CLI for supervising Synnergy cross-chain bridges. It
exposes commands to query bridge status and initiate new connections.

## CLI

Run `npm run build` followed by `npm start status` to display the current bridge
status. Use `npm start connect <chain>` to simulate establishing a connection to
another chain.

## Testing

The project uses Jest with TypeScript. Execute `npm test` to run both unit and
end-to-end tests.
