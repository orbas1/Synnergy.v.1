import { main } from '../../src/main';

test('cli execution outputs total stake', async () => {
  const output = await main();
  expect(output).toMatch(/Total Stake/);
});
