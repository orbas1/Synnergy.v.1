/**
 * Production configuration for the mining-staking-manager module. The API URL
 * defaults to a placeholder endpoint but can be overridden via the `API_URL`
 * environment variable.  Keeping the configuration strongly typed helps catch
 * misconfigurations at build time and eases future extension.
 */
export interface ProductionConfig {
  /** Base URL for backend API calls */
  apiUrl: string;
}

const config: ProductionConfig = {
  apiUrl: process.env.API_URL || 'https://api.example.com',
};

export default config;
