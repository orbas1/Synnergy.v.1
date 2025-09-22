import { WalletService, TransferRequest } from './walletService';
import type { HttpLike } from './walletService';

class StubHttpClient implements HttpLike {
  getMock = jest.fn<Promise<unknown>, [string]>();
  postMock = jest.fn<Promise<unknown>, [string, unknown]>();

  get<T>(path: string): Promise<T> {
    return this.getMock(path) as Promise<T>;
  }

  post<T>(path: string, body?: unknown): Promise<T> {
    return this.postMock(path, body) as Promise<T>;
  }
}

describe('WalletService', () => {
  const baseResponse = {
    address: 'SYN1ABC',
    balance: 100,
    pending: 5,
    currency: 'SYN',
    updatedAt: '2024-01-01T00:00:00Z'
  };

  it('normalizes addresses and computes available balances', async () => {
    const client = new StubHttpClient();
    client.getMock.mockResolvedValueOnce(baseResponse);
    const service = new WalletService(client);

    const result = await service.getBalance('  SYN1ABC  ');

    expect(client.getMock).toHaveBeenCalledWith('/wallets/syn1abc/balance');
    expect(result.address).toBe('syn1abc');
    expect(result.available).toBe(95);
    expect(result.updatedAt.toISOString()).toBe('2024-01-01T00:00:00.000Z');
  });

  it('fetches transactions for an address', async () => {
    const client = new StubHttpClient();
    client.getMock.mockResolvedValueOnce([
      {
        id: 'tx-1',
        from: 'syn1abc',
        to: 'syn1def',
        amount: 10,
        timestamp: '2024-01-01T00:00:00Z'
      }
    ]);

    const service = new WalletService(client);
    const transactions = await service.getTransactions('syn1abc');

    expect(client.getMock).toHaveBeenCalledWith('/wallets/syn1abc/transactions');
    expect(transactions).toHaveLength(1);
    expect(transactions[0].id).toBe('tx-1');
  });

  it('submits transfers with normalized addresses', async () => {
    const client = new StubHttpClient();
    client.postMock.mockResolvedValueOnce({
      txId: 'tx-2',
      accepted: true,
      broadcastAt: '2024-01-01T01:00:00Z'
    });

    const service = new WalletService(client);
    const request: TransferRequest = { to: 'SYN1DEF', amount: 42 };

    const receipt = await service.submitTransfer(' SYN1ABC ', request);

    expect(client.postMock).toHaveBeenCalledWith('/wallets/syn1abc/transfer', {
      ...request,
      to: 'syn1def'
    });
    expect(receipt.txId).toBe('tx-2');
  });
});
