import { main } from '../../src/main';

test('reports node status end-to-end', async () => {
  const output = await main();
  expect(output).toBe('Node status: OK');
});
