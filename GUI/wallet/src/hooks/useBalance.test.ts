import { useBalance } from './useBalance';
import type { BalanceProvider } from '../state/walletStore';
import type { BalanceSummary } from '../services/walletService';

const summary: BalanceSummary = {
  address: 'syn1cache',
  balance: 100,
  pending: 25,
  currency: 'SYN',
  updatedAt: new Date('2024-01-01T00:00:00Z'),
  available: 75
};

class StubStore implements BalanceProvider {
  constructor(private cached?: BalanceSummary, private refreshed?: BalanceSummary) {}

  getBalance = jest.fn((address: string) => {
    return this.cached && this.cached.address === address ? this.cached : undefined;
  });

  refreshBalance = jest.fn(async () => {
    if (!this.refreshed) {
      throw new Error('No refreshed value provided');
    }

    return this.refreshed;
  });
}

describe('useBalance', () => {
  it('returns cached balance when available', async () => {
    const store = new StubStore(summary);
    await expect(useBalance('syn1cache', { store })).resolves.toEqual(summary);
    expect(store.getBalance).toHaveBeenCalledTimes(1);
    expect(store.refreshBalance).not.toHaveBeenCalled();
  });

  it('refreshes when forced', async () => {
    const refreshed: BalanceSummary = { ...summary, available: 60, balance: 110, pending: 50 };
    const store = new StubStore(summary, refreshed);

    await expect(useBalance('syn1cache', { store, refresh: true })).resolves.toEqual(refreshed);
    expect(store.refreshBalance).toHaveBeenCalledTimes(1);
  });

  it('throws for empty addresses', async () => {
    const store = new StubStore();
    await expect(useBalance('  ', { store })).rejects.toThrow('Wallet address cannot be empty');
  });
});
