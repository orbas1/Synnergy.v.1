export interface DistributionStatus {
  online: boolean;
}

export async function fetchStatus(
  fetcher: typeof fetch = fetch
): Promise<DistributionStatus> {
  const res = await fetcher('/api/status');
  return res.json();
}
