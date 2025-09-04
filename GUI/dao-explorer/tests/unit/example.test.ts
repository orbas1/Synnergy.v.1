import { main } from '../../src/main';

describe('main unit', () => {
  it('greets dao-explorer by default', () => {
    expect(main()).toBe('Hello from dao-explorer');
  });
});
