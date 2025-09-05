export default {
  apiUrl: process.env.API_URL || 'https://api.synnergy.local',
  port: Number(process.env.PORT) || 3000,
  logLevel: process.env.LOG_LEVEL || 'info'
};
