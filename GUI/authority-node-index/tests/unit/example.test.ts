import config from '../../config/production';

test('uses default port', () => {
  expect(config.port).toBe(3000);
});
