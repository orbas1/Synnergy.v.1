# Validator Governance Portal

Enterprise-grade GUI for Validator Governance Portal. This scaffold includes TypeScript, linting, formatting, and testing configuration.

## Scripts
- `npm install` – install dependencies
- `npm run build` – compile TypeScript
- `npm test` – run unit tests
- `npm run lint` – run ESLint
- `npm run format` – apply Prettier formatting

## Development

Copy `.env.example` to `.env` and adjust as necessary. Run `npm start` to execute the compiled CLI after building.

## Docker

Build and run using Docker:

```
docker build -t validator-governance-portal .
docker run --rm -p 3000:3000 validator-governance-portal
```
