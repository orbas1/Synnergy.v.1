import { WalletService, BalanceSummary, Transaction, TransferReceipt, TransferRequest } from './walletService';

export const walletService = new WalletService();

export async function getBalance(address: string): Promise<BalanceSummary> {
  return walletService.getBalance(address);
}

export async function getTransactions(address: string): Promise<Transaction[]> {
  return walletService.getTransactions(address);
}

export async function transfer(from: string, request: TransferRequest): Promise<TransferReceipt> {
  return walletService.submitTransfer(from, request);
}

export type { BalanceSummary, Transaction, TransferReceipt, TransferRequest };
