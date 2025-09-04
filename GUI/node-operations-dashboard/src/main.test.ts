import { main } from './main';

it('returns node status', async () => {
  await expect(main()).resolves.toContain('Node status: OK');
});
