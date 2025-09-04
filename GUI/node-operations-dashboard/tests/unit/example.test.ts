import { fetchNodeStatus } from '../../src/services/status';

test('fetchNodeStatus returns OK', async () => {
  await expect(fetchNodeStatus()).resolves.toBe('OK');
});
