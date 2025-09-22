import { renderBalance } from '../components/WalletBalance';
import { BalanceSummary } from '../services/walletService';

export function renderHome(summary: BalanceSummary): string {
  const sections = [
    '=== Synnergy Wallet Overview ===',
    renderBalance(summary),
    'Use `wallet --help` for available commands.'
  ];

  return sections.join('\n\n');
}
