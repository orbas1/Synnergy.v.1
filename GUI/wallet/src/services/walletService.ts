import { z } from 'zod';
import { walletConfig, WalletGuiConfig, resolveConfig } from '../config';
import { normalizeAddress } from '../utils/address';
import { HttpClient, HttpClientOptions } from './httpClient';

export interface HttpLike {
  get<T>(path: string): Promise<T>;
  post<T>(path: string, body?: unknown): Promise<T>;
}

const BalanceSchema = z.object({
  address: z.string().min(1),
  balance: z.number(),
  pending: z.number().min(0).default(0),
  currency: z.string().min(1).default(walletConfig.defaultCurrency),
  updatedAt: z.preprocess((value) => {
    if (value instanceof Date) {
      return value;
    }

    const date = new Date(value as string | number | Date);
    return Number.isNaN(date.getTime()) ? undefined : date;
  }, z.date())
});

type BalancePayload = z.infer<typeof BalanceSchema>;

const BalanceResponseSchema = z.union([
  BalanceSchema,
  z.object({ data: BalanceSchema })
]);

const TransactionSchema = z.object({
  id: z.string().min(1),
  from: z.string().min(1),
  to: z.string().min(1),
  amount: z.number(),
  timestamp: z.preprocess((value) => {
    const date = new Date(value as string | number | Date);
    return Number.isNaN(date.getTime()) ? undefined : date;
  }, z.date())
});

export type Transaction = z.infer<typeof TransactionSchema>;

const TransactionListSchema = z.array(TransactionSchema);

const TransferReceiptSchema = z.object({
  txId: z.string().min(1),
  accepted: z.boolean(),
  broadcastAt: z.preprocess((value) => {
    const date = new Date(value as string | number | Date);
    return Number.isNaN(date.getTime()) ? undefined : date;
  }, z.date())
});

export type TransferReceipt = z.infer<typeof TransferReceiptSchema>;

export interface TransferRequest {
  to: string;
  amount: number;
  memo?: string;
}

export interface BalanceSummary {
  address: string;
  balance: number;
  pending: number;
  currency: string;
  updatedAt: Date;
  available: number;
}

function createHttpClient(config: WalletGuiConfig): HttpClient {
  const options: HttpClientOptions = {
    baseUrl: config.apiUrl,
    timeoutMs: config.requestTimeoutMs,
    maxRetries: config.maxRetries
  };

  return new HttpClient(options);
}

export class WalletService {
  private readonly client: HttpLike;
  private readonly config: WalletGuiConfig;

  constructor(client?: HttpLike, overrides?: Partial<WalletGuiConfig>) {
    this.config = resolveConfig(overrides);
    this.client = client ?? createHttpClient(this.config);
  }

  async getBalance(address: string): Promise<BalanceSummary> {
    const normalized = normalizeAddress(address);
    const payload = await this.client.get<unknown>(`/wallets/${normalized}/balance`);
    const parsed = BalanceResponseSchema.parse(payload);
    const balance: BalancePayload = 'data' in parsed ? parsed.data : parsed;

    return {
      ...balance,
      currency: balance.currency || this.config.defaultCurrency,
      address: balance.address.toLowerCase(),
      available: Number(balance.balance - balance.pending),
      updatedAt: balance.updatedAt
    };
  }

  async getTransactions(address: string): Promise<Transaction[]> {
    const normalized = normalizeAddress(address);
    const payload = await this.client.get<unknown>(`/wallets/${normalized}/transactions`);
    return TransactionListSchema.parse(payload);
  }

  async submitTransfer(from: string, request: TransferRequest): Promise<TransferReceipt> {
    const normalizedFrom = normalizeAddress(from);
    const normalizedTo = normalizeAddress(request.to);

    const payload = await this.client.post<unknown>(`/wallets/${normalizedFrom}/transfer`, {
      ...request,
      to: normalizedTo
    });

    return TransferReceiptSchema.parse(payload);
  }
}
