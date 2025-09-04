import { main } from './main';

describe('main', () => {
  it('includes configured API url', () => {
    process.env.API_URL = 'http://example';
    expect(main()).toContain('http://example');
  });
});
