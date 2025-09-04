# Dao Explorer

Enterprise-grade web interface for exploring decentralized autonomous organizations on the Synnergy network. The application is built with TypeScript and comes preconfigured with linting, formatting, testing and containerization tooling.

## Getting Started

```bash
npm ci
npm run build
npm test
```

## Environment Variables

See `.env.example` for configurable options such as `API_URL`, `LOG_LEVEL` and `DATABASE_URL`.

## Scripts

- `npm run build` – compile TypeScript to JavaScript
- `npm test` – run unit tests
- `npm run lint` – lint source files

## Docker

To run via Docker:

```bash
docker build -t dao-explorer .
docker run -p 3000:3000 dao-explorer
```
