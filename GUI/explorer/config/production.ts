interface ExplorerConfig {
  apiUrl: string;
  logLevel: string;
  port: number;
}

const config: ExplorerConfig = {
  apiUrl: process.env.API_URL || 'http://localhost:3000',
  logLevel: process.env.LOG_LEVEL || 'info',
  port: parseInt(process.env.PORT || '3000', 10),
};

export default config;
