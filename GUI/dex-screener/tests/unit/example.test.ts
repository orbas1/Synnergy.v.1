import { renderHome } from '../../src/pages/home';

describe('home page', () => {
  it('returns the exact home page title', () => {
    expect(renderHome()).toBe('Dex Screener Home');
  });
});
