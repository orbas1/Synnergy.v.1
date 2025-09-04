import { main } from './main';

beforeEach(() => {
  ;(global as any).fetch = jest.fn().mockResolvedValue({
    ok: true,
    json: async () => ({ status: 'OK' }),
  });
});

afterEach(() => {
  jest.restoreAllMocks();
});

it('returns node status', async () => {
  await expect(main()).resolves.toContain('Node status: OK');
});
