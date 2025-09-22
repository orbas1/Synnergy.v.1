import { normalizeAddress } from './address';

describe('normalizeAddress', () => {
  it('lowercases and trims addresses', () => {
    expect(normalizeAddress('  SYN1ABC ')).toBe('syn1abc');
  });

  it('throws for empty strings', () => {
    expect(() => normalizeAddress('   ')).toThrow('Wallet address cannot be empty');
  });
});
