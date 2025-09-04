import { main } from '../../src/main';

describe('security operations center e2e', () => {
  it('boots with provided API url', () => {
    process.env.API_URL = 'http://e2e.local';
    expect(main()).toBe(
      'Security Operations Center started with API at http://e2e.local'
    );
  });
});
