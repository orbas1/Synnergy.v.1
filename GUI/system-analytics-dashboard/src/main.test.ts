import { main } from './main';

beforeEach(() => {
  (global as any).fetch = jest.fn().mockResolvedValue({
    ok: true,
    json: async () => ({ status: 'OK' }),
  });
});

afterEach(() => {
  jest.restoreAllMocks();
});

test('main returns analytics status', async () => {
  await expect(main()).resolves.toContain('System analytics: OK');
});
