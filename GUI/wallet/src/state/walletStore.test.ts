import { WalletStore } from './walletStore';
import type { BalanceSummary } from '../services/walletService';

class FakeService {
  constructor(private readonly summary: BalanceSummary) {}
  getBalance = jest.fn(async (_address: string) => this.summary);
}

describe('WalletStore', () => {
  const summary: BalanceSummary = {
    address: 'syn1abc',
    balance: 100,
    pending: 20,
    currency: 'SYN',
    updatedAt: new Date('2024-01-01T00:00:00Z'),
    available: 80
  };

  it('caches balances after refresh', async () => {
    const service = new FakeService(summary);
    const store = new WalletStore(service);

    await store.refreshBalance('syn1abc');
    const cached = store.getBalance('syn1abc');

    expect(service.getBalance).toHaveBeenCalledTimes(1);
    expect(cached).toEqual(summary);
  });

  it('notifies subscribers on refresh', async () => {
    const service = new FakeService(summary);
    const store = new WalletStore(service);

    const listener = jest.fn();
    const unsubscribe = store.subscribe('syn1abc', listener);
    await store.refreshBalance('syn1abc');

    expect(listener).toHaveBeenCalledWith(summary);
    unsubscribe();
  });
});
