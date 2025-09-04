import { main } from '../../src/main';

beforeEach(() => {
  (global as any).fetch = jest.fn().mockResolvedValue({
    ok: true,
    json: async () => ({ status: 'OK' }),
  });
});

afterEach(() => {
  jest.restoreAllMocks();
});

test('main integrates with API and returns analytics status', async () => {
  await expect(main()).resolves.toBe('System analytics: OK');
});
