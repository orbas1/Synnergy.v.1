import { resolveConfig } from '../../src/config';

describe('resolveConfig', () => {
  it('prefers explicit overrides over environment defaults', () => {
    process.env.API_URL = 'http://env-host';
    process.env.WALLET_REQUEST_TIMEOUT_MS = '2000';

    const config = resolveConfig({ apiUrl: 'http://override', maxRetries: 5 });

    expect(config.apiUrl).toBe('http://override');
    expect(config.requestTimeoutMs).toBe(2000);
    expect(config.maxRetries).toBe(5);
  });
});
