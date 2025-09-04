import baseConfig from '../config/production';

export function main(): string {
  const cfg = { ...baseConfig, apiUrl: process.env.API_URL ?? baseConfig.apiUrl };
  return `Security Operations Center started with API at ${cfg.apiUrl}`;
}

if (require.main === module) {
  try {
    console.log(main());
  } catch (err) {
    console.error('Failed to start Security Operations Center', err);
    process.exit(1);
  }
}
