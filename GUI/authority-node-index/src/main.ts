import config from '../config/production';

/**
 * start boots the authority node index service and returns a status message.
 * The function centralises configuration usage so both the CLI and GUI
 * surface the same runtime details.
 */
export function start(): string {
  const message = `Authority Node Index running on port ${config.port}`;
  return message;
}

if (require.main === module) {
  // Log startup status when executed directly from the CLI.
  // Wrapping in try/catch keeps the process from crashing on init failures.
  try {
    // eslint-disable-next-line no-console
    console.log(start());
  } catch (err) {
    // eslint-disable-next-line no-console
    console.error('failed to start Authority Node Index', err);
    process.exit(1);
  }
}
