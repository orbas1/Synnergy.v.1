# Synnergy Web Control Panel

This Next.js application provides a web interface for the `synnergy` CLI. It dynamically enumerates all CLI commands and exposes input fields for each flag so the entire CLI can be exercised from a browser. The project is deployable to [Vercel](https://vercel.com/).

## Development

```bash
npm install
npm run dev
```

## Usage

1. Open the running app in your browser.
2. Select a command from the drop-down list.
3. Provide values for any listed flags and optional additional arguments.
4. Click **Run** to execute the command and view the output.

The backend invokes `go run ../cmd/synnergy/main.go`, so Go must be available in the deployment environment.

## Deployment

Deploy the `web` directory to Vercel. API routes used by the interface include:

- `GET /api/commands` – list available CLI commands
- `GET /api/help?cmd=<command>` – list flags for a command
- `POST /api/run` – execute a command with arguments

Restrict access appropriately when exposing this interface in production environments.
