export default {
  apiUrl: process.env.API_URL ?? 'http://localhost:8080',
  logLevel: process.env.LOG_LEVEL ?? 'info'
} as const;
