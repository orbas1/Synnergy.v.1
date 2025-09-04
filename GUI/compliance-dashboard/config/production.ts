interface Config {
  apiUrl: string;
  port: number;
  logLevel: string;
}

const config: Config = {
  apiUrl: process.env.API_URL || 'http://localhost:3000',
  port: Number(process.env.PORT) || 3000,
  logLevel: process.env.LOG_LEVEL || 'info'
};

export default config;
