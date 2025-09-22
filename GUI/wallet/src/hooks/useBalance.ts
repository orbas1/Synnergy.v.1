import { walletStore, BalanceProvider } from '../state/walletStore';
import { BalanceSummary } from '../services/walletService';
import { normalizeAddress } from '../utils/address';

export interface UseBalanceOptions {
  refresh?: boolean;
  store?: BalanceProvider;
}

export async function useBalance(address: string, options: UseBalanceOptions = {}): Promise<BalanceSummary> {
  const provider = options.store ?? walletStore;
  const normalized = normalizeAddress(address);

  if (!options.refresh) {
    const cached = provider.getBalance(normalized);
    if (cached) {
      return cached;
    }
  }

  return provider.refreshBalance(normalized);
}
