/**
 * Entry point for the storage‑marketplace module.
 * Reads the endpoint from the environment so the CLI and GUI
 * can be configured without recompilation.
 */
export async function main(): Promise<string> {
  const endpoint =
    process.env.STORAGE_MARKETPLACE_ENDPOINT ?? 'http://localhost:8080';
  return `Storage Marketplace service available at ${endpoint}`;
}

// Allow the file to be executed directly via `node dist/main.js`
// or `ts-node src/main.ts`.
if (require.main === module) {
  main()
    .then((msg) => console.log(msg))
    .catch((err) => {
      console.error('Startup failed', err);
      process.exit(1);
    });
}
