interface Config {
  apiUrl: string;
}

const config: Config = {
  apiUrl: process.env.API_URL || 'https://api.synnergy.network'
};

export default config;
