import { main } from './main';

describe('main', () => {
  it('includes configured endpoint', async () => {
    process.env.STORAGE_MARKETPLACE_ENDPOINT = 'http://test-endpoint';
    const msg = await main();
    expect(msg).toContain('http://test-endpoint');
  });
});
