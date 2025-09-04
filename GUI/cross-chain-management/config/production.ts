export interface ProductionConfig {
  apiUrl: string;
  logLevel: string;
}

const config: ProductionConfig = {
  apiUrl: process.env.API_URL || 'http://localhost:3000',
  logLevel: process.env.LOG_LEVEL || 'info',
};

export default config;
