# Mining Staking Manager

Enterprise-grade GUI for Mining Staking Manager. This scaffold includes TypeScript, linting, formatting, and testing configuration.

## Setup
Copy `.env.example` to `.env` and adjust values:

```
API_URL=http://localhost:3000
PORT=4000
LOG_LEVEL=info
```

Install dependencies and run tests:

```
make install
make test
```

## Scripts
- `npm run build` – compile TypeScript
- `npm test` – run unit and e2e tests
- `npm run lint` – lint sources with ESLint
- `npm run format` – apply Prettier formatting

## Docker
```
docker build -t mining-manager .
docker run -p 4000:4000 --env-file .env mining-manager
```
