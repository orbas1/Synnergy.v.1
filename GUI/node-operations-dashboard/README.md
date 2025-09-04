# Node Operations Dashboard

Enterprise-grade GUI for monitoring Synnergy network nodes. Built with TypeScript and configured with linting, formatting, and testing tools.

## Development

```bash
npm ci
npm start
```

## Scripts

- `npm run build` – compile TypeScript sources into `dist/`
- `npm test` – run the Jest test suite
- `npm run lint` – run ESLint over the codebase

## Docker

Build and run the container locally:

```bash
docker build -t node-operations-dashboard .
docker run -p 3000:3000 node-operations-dashboard
```
