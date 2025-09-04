export interface PoolInfo {
  pair: string;
  reserve: number;
}

export async function fetchPool(
  pair: string,
  fetcher: typeof fetch = fetch
): Promise<PoolInfo> {
  const res = await fetcher(`/api/pools/${pair}`);
  return res.json();
}
