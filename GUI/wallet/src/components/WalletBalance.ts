import { BalanceSummary } from '../services/walletService';

const numberFormatters = new Map<string, Intl.NumberFormat>();

function formatCurrency(value: number, currency: string): string {
  const key = currency.toUpperCase();
  if (!numberFormatters.has(key)) {
    numberFormatters.set(
      key,
      new Intl.NumberFormat('en-US', {
        style: 'currency',
        currency: key,
        minimumFractionDigits: 2,
        maximumFractionDigits: 8
      })
    );
  }

  return numberFormatters.get(key)!.format(value);
}

export function renderBalance(summary: BalanceSummary): string {
  const lines = [
    `Address        : ${summary.address}`,
    `Available      : ${formatCurrency(summary.available, summary.currency)}`,
    `Total Balance  : ${formatCurrency(summary.balance, summary.currency)}`,
    `Pending Amount : ${formatCurrency(summary.pending, summary.currency)}`,
    `Last Updated   : ${summary.updatedAt.toISOString()}`
  ];

  return lines.join('\n');
}
