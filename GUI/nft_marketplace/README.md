# Nft_marketplace

Enterprise-grade GUI for the Synnergy NFT marketplace. The project is built
with TypeScript and comes preconfigured with linting, formatting and unit test
tooling.

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
docker build -t nft-marketplace .
docker run -p 3000:3000 nft-marketplace
```
