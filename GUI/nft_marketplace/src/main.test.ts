import { main } from './main';

it('outputs a rendered NFT list', async () => {
  const output = await main();
  expect(output).toContain('Genesis Token');
});
