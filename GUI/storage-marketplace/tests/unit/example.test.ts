import { main } from '../../src/main';

test('main exposes endpoint', async () => {
  process.env.STORAGE_MARKETPLACE_ENDPOINT = 'http://unit-endpoint';
  const output = await main();
  expect(output).toContain('http://unit-endpoint');
});
