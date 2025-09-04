import { fetchPool } from '../../src/services/api';

test('fetchPool uses provided fetcher', async () => {
  const fake = jest.fn().mockResolvedValue({
    json: () => Promise.resolve({ pair: 'ABC/DEF', reserve: 42 })
  });
  const res = await fetchPool('ABC/DEF', fake as any);
  expect(res.reserve).toBe(42);
  expect(fake).toHaveBeenCalled();
});
