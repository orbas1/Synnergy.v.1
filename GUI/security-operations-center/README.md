# Security Operations Center

Enterprise-grade GUI for monitoring and responding to network security events. This scaffold ships with TypeScript, linting,
formatting and containerization to ease production deployments.

## Scripts
- `npm run build` – compile TypeScript to JavaScript
- `npm test` – run unit tests (placeholder)
- `npm run lint` – lint sources with ESLint
- `npm run format` – format sources with Prettier

## Development
1. Copy `.env.example` to `.env` and adjust values.
2. `make install` to install dependencies.
3. `make build` or `make test` as needed.

## Docker
A multi-stage `Dockerfile` is provided. Build with:
```sh
docker build -t soc .
```
