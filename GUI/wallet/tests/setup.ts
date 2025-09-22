afterEach(() => {
  jest.clearAllMocks();
  jest.restoreAllMocks();
  delete process.env.WALLET_API_URL;
  delete process.env.API_URL;
  delete process.env.WALLET_REQUEST_TIMEOUT_MS;
  delete process.env.WALLET_REQUEST_RETRIES;
  delete process.env.WALLET_CURRENCY;
});

export {};
