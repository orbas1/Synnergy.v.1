import { fetchStatus } from '../../src/services/api';

test('fetchStatus returns online', async () => {
  const fake = jest.fn().mockResolvedValue({
    json: () => Promise.resolve({ online: true })
  });
  const res = await fetchStatus(fake as any);
  expect(res.online).toBe(true);
  expect(fake).toHaveBeenCalled();
});
