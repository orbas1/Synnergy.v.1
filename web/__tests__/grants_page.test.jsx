import { describe, expect, it, vi, beforeEach, afterEach } from 'vitest';
import { render, screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import GrantsConsole from '../pages/grants.jsx';

describe('GrantsConsole', () => {
  const originalConsoleError = console.error;

  beforeEach(() => {
    vi.spyOn(console, 'error').mockImplementation((message, ...args) => {
      if (typeof message === 'string' && message.includes('act(')) {
        return;
      }
      originalConsoleError(message, ...args);
    });
    vi.spyOn(global, 'fetch').mockResolvedValue({
      ok: true,
      json: async () => ({ grants: [], status: {} }),
    });
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('submits a wallet-signed create request', async () => {
    const user = userEvent.setup();
    render(<GrantsConsole />);
    await waitFor(() => expect(fetch).toHaveBeenCalledTimes(1));
    await waitFor(() => expect(screen.queryByText('Processing request…')).not.toBeInTheDocument());

    fetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({
        message: '1',
        grants: [{ id: 1, beneficiary: 'alice', status: 'ACTIVE' }],
        status: { total: 1, active: 1 },
      }),
    });

    const createSection = screen.getByRole('heading', { name: /Create grant/i }).closest('section');
    expect(createSection).toBeTruthy();
    const scoped = within(createSection);

    await user.type(scoped.getByLabelText(/Beneficiary/i), 'alice');
    await user.type(scoped.getByLabelText(/Program name/i), 'research');
    await user.type(scoped.getByLabelText(/^Amount$/i), '100');
    await user.type(scoped.getByLabelText(/Authorizer wallets/i), 'wallet1:pw1');
    await user.type(scoped.getByLabelText(/Creator wallet path/i), '/tmp/wallet');
    await user.type(scoped.getByLabelText(/Creator wallet password/i), 'secret');
    await user.click(scoped.getByRole('button', { name: /Create grant/i }));
    await waitFor(() => expect(fetch).toHaveBeenCalledTimes(2));
    await waitFor(() => expect(screen.queryByText('Processing request…')).not.toBeInTheDocument());
    const [, request] = fetch.mock.calls[1];
    const payload = JSON.parse(request.body);
    expect(payload.wallet).toBe('/tmp/wallet');
    expect(payload.password).toBe('secret');
    await waitFor(() => expect(screen.getByRole('status')).toHaveTextContent('1'));
    expect(scoped.getByLabelText(/Creator wallet password/i)).toHaveValue('');
    expect(scoped.getByLabelText(/Creator wallet path/i)).toHaveValue('/tmp/wallet');
  });
});
