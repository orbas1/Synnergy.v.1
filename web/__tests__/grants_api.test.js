import { afterEach, describe, expect, it, vi } from 'vitest';
import handler from '../pages/api/grants.js';

vi.mock('child_process', () => {
  const execFile = vi.fn();
  return {
    execFile,
    default: { execFile },
  };
});

const { execFile } = await import('child_process');

const createResponse = () => {
  const res = {};
  res.status = vi.fn(() => res);
  res.json = vi.fn();
  return res;
};

afterEach(() => {
  execFile.mockReset();
});

describe('grants API', () => {
  it('returns grants and status summaries', async () => {
    execFile.mockImplementation((cmd, args, options, callback) => {
      const action = args[3];
      if (action === 'list') {
        callback(null, '[{"id":1}]', '');
        return;
      }
      if (action === 'status') {
        callback(null, '{"total":1}', '');
        return;
      }
      callback(new Error(`unexpected action ${action}`));
    });
    const res = createResponse();
    await handler({ method: 'GET', query: {} }, res);
    expect(res.status).toHaveBeenCalledWith(200);
    expect(res.json).toHaveBeenCalledWith({ grants: [{ id: 1 }], status: { total: 1 } });
  });

  it('creates a grant and returns refreshed state', async () => {
    execFile.mockImplementation((cmd, args, options, callback) => {
      const action = args[3];
      if (action === 'create') {
        callback(null, '1\n', '');
        return;
      }
      if (action === 'list') {
        callback(null, '[{"id":1}]', '');
        return;
      }
      if (action === 'status') {
        callback(null, '{"total":1}', '');
        return;
      }
      callback(new Error(`unexpected action ${action}`));
    });
    const body = {
      action: 'create',
      beneficiary: 'alice',
      name: 'program',
      amount: '100',
      wallet: '/tmp/wallet',
      password: 'pw',
      authorizers: ['auth:path'],
    };
    const res = createResponse();
    await handler({ method: 'POST', body }, res);
    expect(res.status).toHaveBeenCalledWith(200);
    expect(res.json).toHaveBeenCalledWith({
      message: '1',
      grants: [{ id: 1 }],
      status: { total: 1 },
    });
    const createCall = execFile.mock.calls.find(([, args]) => args[3] === 'create');
    expect(createCall[1]).toContain('--wallet');
    expect(createCall[1]).toContain('--password');
    expect(createCall[1]).toContain('auth:path');
  });

  it('rejects create requests without authentication', async () => {
    const res = createResponse();
    await handler({ method: 'POST', body: { action: 'create', beneficiary: 'a', name: 'b', amount: '1' } }, res);
    expect(res.status).toHaveBeenCalledWith(500);
    expect(res.json).toHaveBeenCalledWith({ error: 'wallet and password are required' });
  });
});
