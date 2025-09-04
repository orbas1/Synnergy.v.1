interface ProdConfig {
  apiUrl: string;
  port: number;
}

const config: ProdConfig = {
  apiUrl: process.env.API_URL || 'https://api.synnergy.local',
  port: Number(process.env.PORT) || 3000
};

export default config;
