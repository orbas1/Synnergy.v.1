import { main } from '../../src/main';

it('renders NFTs end-to-end', async () => {
  const output = await main();
  expect(output.split('\n')).toContain('#1 - Genesis Token');
});
