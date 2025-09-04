import { main } from '../../src/main';

describe('smart contract marketplace e2e', () => {
  it('greets from main entry', () => {
    expect(main()).toContain('smart-contract-marketplace');
  });
});
