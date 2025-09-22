import { main } from './main';
import { useBalance } from './hooks/useBalance';
import type { BalanceSummary } from './services/walletService';

jest.mock('./hooks/useBalance');

describe('main', () => {
  const summary: BalanceSummary = {
    address: 'syn1example',
    balance: 125.5,
    pending: 10,
    currency: 'SYN',
    updatedAt: new Date('2024-01-01T00:00:00Z'),
    available: 115.5
  };

  beforeEach(() => {
    (useBalance as jest.Mock).mockResolvedValue(summary);
  });

  afterEach(() => {
    jest.clearAllMocks();
  });

  it('renders the home view with resolved balance data', async () => {
    const output = await main('syn1example');
    expect(output).toContain('Synnergy Wallet Overview');
    expect(output).toContain('syn1example');
    expect(useBalance).toHaveBeenCalledWith(
      'syn1example',
      expect.objectContaining({ refresh: true, store: expect.anything() })
    );
  });

  it('throws when no wallet address is provided', async () => {
    await expect(main('')).rejects.toThrow('wallet address');
  });
});
