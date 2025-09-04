import { getStatus } from '../../src/main';

describe('cross-chain management status', () => {
  test('reports all bridges operational', () => {
    expect(getStatus()).toBe('All cross-chain bridges operational');
  });
});
