# Security Operations Center Architecture

The security operations center (SOC) is a lightweight TypeScript service.  It loads configuration from `config/production.ts`, starts via `src/main.ts` and exposes hooks for UI components under `src/components`.

```
config/production.ts   -> runtime configuration
src/main.ts            -> entry point with error handling
src/components/        -> React/Vue components (future)
src/services/          -> API wrappers
```

The SOC communicates with other Synnergy services through REST endpoints defined by `API_URL`.  All network calls should be resilient and surface failures to the calling CLI or GUI.
