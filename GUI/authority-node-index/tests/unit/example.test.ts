import config from '../../config/production';

test('uses default configuration values', () => {
  expect(config.port).toBe(3000);
  expect(config.logLevel).toBe('info');
});
