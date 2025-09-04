# Compliance Dashboard

Enterprise-grade dashboard for monitoring network compliance. The project uses
TypeScript with strict linting, formatting and test tooling to ensure
production reliability.

## Requirements
- Node.js 18+
- npm 9+

## Setup
```bash
npm ci
cp .env.example .env
```

## Scripts
- `npm start` – start the dashboard locally
- `npm run build` – compile TypeScript
- `npm test` – run unit tests
- `npm run lint` – lint sources

## Docker
```bash
docker build -t compliance-dashboard .
docker run -p 3000:3000 compliance-dashboard
```

## Kubernetes
Deployment manifests are provided in `k8s/` for cluster deployment.

