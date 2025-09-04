import { main } from '../../src/main';

test('CLI outputs successful analytics status', async () => {
  (global as any).fetch = jest.fn().mockResolvedValue({
    ok: true,
    json: async () => ({ status: 'OK' }),
  });

  const output = await main();
  expect(output).toContain('System analytics: OK');
});
