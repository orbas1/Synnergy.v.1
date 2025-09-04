test('overrides api url via environment variable', () => {
  process.env.API_URL = 'https://api.internal';
  jest.resetModules();
  // eslint-disable-next-line global-require
  const cfg = require('../../config/production').default;
  expect(cfg.apiUrl).toBe('https://api.internal');
  delete process.env.API_URL;
});
