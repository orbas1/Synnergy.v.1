interface Config {
  apiUrl: string;
}

const config: Config = {
  apiUrl: process.env.API_URL || 'https://api.synnergy.example.com'
};

export default config;
