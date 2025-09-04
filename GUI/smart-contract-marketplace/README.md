# Smart Contract Marketplace

Enterprise-grade GUI for Smart Contract Marketplace. This scaffold includes TypeScript, linting, formatting, and testing configuration.

## Scripts
- `npm run build` - compile TypeScript
- `npm test` - run unit and e2e tests
- `npm run lint` - static analysis
- `npm run format` - format source files

## Environment
Copy `.env.example` to `.env` and adjust values for deployment.

## Usage
After installing dependencies with `npm install`:

1. Build the project with `npm run build`.
2. Start the API server using `npm start`.
3. The service exposes `GET /contracts` which returns available contract
   templates in JSON format. This endpoint underpins future UI and CLI
   integrations with the Synnergy Network.

The Jest test suite (`npm test`) exercises the REST endpoint with Supertest
to provide a baseline for further development.
