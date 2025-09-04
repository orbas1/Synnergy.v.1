import { main } from './main';

describe('main', () => {
  it('returns default greeting', () => {
    expect(main()).toBe('Hello from dao-explorer');
  });

  it('supports custom name', () => {
    expect(main('tester')).toBe('Hello from tester');
  });
});
