import { renderBalance } from '../components/WalletBalance';

export function renderHome(balance: number): string {
  return `Wallet Home\n${renderBalance(balance)}`;
}
