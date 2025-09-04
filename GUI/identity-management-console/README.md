# Identity Management Console

Enterprise-grade GUI for Identity Management Console. This scaffold includes TypeScript, linting, formatting, and testing configuration.

## Getting Started

```bash
npm ci
npm run build
npm start -- --name yourname
```

Register a user:

```bash
npm start -- --register alice --key pubkey
```

## Development

- `npm test` - run unit tests
- `npm run lint` - lint sources
- `npm run format` - format sources with Prettier

Docker and Kubernetes manifests are provided for containerized deployments.
