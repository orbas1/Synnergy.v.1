# Wallet

Enterprise-grade GUI for Wallet. This scaffold provides TypeScript tooling,
linting, formatting, tests and production deployment assets.

## Requirements
- Node.js 18+
- npm

## Scripts
- `npm run build` - compile TypeScript
- `npm start` - start the compiled application
- `npm test` - run unit tests
- `npm run lint` - run ESLint
- `npm run format` - run Prettier

## Development
1. Install dependencies with `npm ci`.
2. Copy `.env.example` to `.env` and adjust settings.
3. Run `npm start` to launch the wallet locally.

## Deployment
The provided `Dockerfile` builds a production image using a non-root user.
`docker-compose.yml` can be used to run the service alongside dependencies.
