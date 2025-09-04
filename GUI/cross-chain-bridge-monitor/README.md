# Cross-Chain Bridge Monitor

This module provides a lightweight monitoring interface for bridge activity across networks. It is built with TypeScript and ships with formatting, linting and unit test support for enterprise workflows.

## Usage

1. Install dependencies:
   ```bash
   npm ci
   ```
2. Run tests:
   ```bash
   npm test
   ```
3. Build the project:
   ```bash
   npm run build
   ```
4. Start the monitor locally:
   ```bash
   npm start
   ```

The service reads the `API_URL` environment variable to determine the bridge service to query. If unset, it defaults to `http://localhost:8080`.

## Scripts

- `npm run build` – compile TypeScript sources into `dist/`
- `npm start` – run the compiled monitor
- `npm test` – execute the Jest test suite and generate coverage
- `npm run lint` – check code formatting and style
- `npm run format` – auto-format sources using Prettier
