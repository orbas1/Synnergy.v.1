import { main } from '../../src/main';

describe('main', () => {
  it('uses default API url when none provided', () => {
    delete process.env.API_URL;
    expect(main()).toBe(
      'Security Operations Center started with API at http://localhost:8080'
    );
  });
});
