export default {
  apiUrl: process.env.API_URL || '',
  port: Number(process.env.PORT) || 3000,
  logLevel: process.env.LOG_LEVEL || 'info'
};
