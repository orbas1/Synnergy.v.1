# Cross Chain Management

Enterprise-grade CLI for managing cross-chain bridges. This scaffold uses
TypeScript with linting, formatting and testing tools.

## Development
1. Copy `.env.example` to `.env` and adjust values.
2. Install dependencies with `npm ci`.
3. Run `npm test` to execute unit and e2e tests.
4. Use `npm run build` to compile TypeScript.
5. Start the CLI with `npm start <command>` after building.

## Usage

After building, run `npm start status` to view bridge health. Use
`npm start connect <chain>` to simulate connecting to a new chain.

## Environment Variables
- `API_URL` – base URL for backend services
- `LOG_LEVEL` – application log verbosity
- `DB_URL` – connection string for the metadata database

## Scripts
- `npm run build` – compile TypeScript
- `npm start` – run the CLI
- `npm test` – run unit and e2e tests
- `npm run lint` – run ESLint
- `npm run format` – format code with Prettier
