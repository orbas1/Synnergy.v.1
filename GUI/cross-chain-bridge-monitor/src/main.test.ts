import { main } from './main';

describe('main', () => {
  test('reflects API_URL environment variable', () => {
    process.env.API_URL = 'http://example.com';
    expect(main()).toContain('http://example.com');
  });
});
