import { main } from '../../src/main';

beforeEach(() => {
  ;(global as any).fetch = jest.fn().mockResolvedValue({
    ok: true,
    json: async () => ({ status: 'OK' }),
  });
});

afterEach(() => {
  jest.restoreAllMocks();
});

test('reports node status end-to-end', async () => {
  const output = await main();
  expect(output).toBe('Node status: OK');
});
