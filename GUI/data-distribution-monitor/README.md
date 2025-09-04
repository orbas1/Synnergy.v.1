# Data Distribution Monitor

Enterprise-grade GUI and CLI for monitoring data distribution across the Synnergy network. The module exposes a lightweight command line interface and web-ready hooks.

## Scripts
- `npm run build` – compile TypeScript sources
- `npm test` – execute unit tests with Jest
- `npm run lint` – lint sources
- `npm run format` – auto-format code

## Usage
Run the monitor locally:

```bash
npm install
npm run build
node dist/main.js status
```

The command above prints the service status, enabling integration with higher level orchestrators and the function web.
