import { EventEmitter } from 'events';
import { WalletService, BalanceSummary } from '../services/walletService';

export type BalanceListener = (summary: BalanceSummary) => void;

export interface BalanceProvider {
  getBalance(address: string): BalanceSummary | undefined;
  refreshBalance(address: string): Promise<BalanceSummary>;
}

function eventName(address: string): string {
  return `balance:${address}`;
}

type WalletBalanceService = {
  getBalance(address: string): Promise<BalanceSummary>;
};

export class WalletStore extends EventEmitter implements BalanceProvider {
  private readonly balances = new Map<string, BalanceSummary>();

  constructor(private readonly service: WalletBalanceService = new WalletService()) {
    super();
  }

  getBalance(address: string): BalanceSummary | undefined {
    return this.balances.get(address);
  }

  async refreshBalance(address: string): Promise<BalanceSummary> {
    const summary = await this.service.getBalance(address);
    this.balances.set(address, summary);
    this.emit('balance', summary);
    this.emit(eventName(address), summary);
    return summary;
  }

  subscribe(address: string, listener: BalanceListener): () => void {
    const key = eventName(address);
    this.on(key, listener);
    return () => this.off(key, listener);
  }
}

export const walletStore = new WalletStore();
