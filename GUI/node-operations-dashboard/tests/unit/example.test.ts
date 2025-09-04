import { fetchNodeStatus } from '../../src/services/status';

beforeEach(() => {
  // Mock global fetch for unit tests
  ;(global as any).fetch = jest.fn().mockResolvedValue({
    ok: true,
    json: async () => ({ status: 'OK' }),
  });
});

afterEach(() => {
  jest.restoreAllMocks();
});

test('fetchNodeStatus returns OK', async () => {
  await expect(fetchNodeStatus()).resolves.toBe('OK');
});
