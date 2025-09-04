import config from '../../config/production';

test('has default api url', () => {
  expect(config.apiUrl).toBe('');
});
