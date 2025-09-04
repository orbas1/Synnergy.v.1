interface Config {
  apiUrl: string;
  port: number;
}

const config: Config = {
  apiUrl: process.env.API_URL || 'http://localhost:3000',
  port: parseInt(process.env.PORT || '3000', 10),
};

export default config;
