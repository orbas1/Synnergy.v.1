import { main } from './main';

test('main renders dashboard with total stake', async () => {
  const output = await main();
  expect(output).toContain('Total Stake');
});
