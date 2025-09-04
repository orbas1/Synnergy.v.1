export interface StakeInfo {
  miner: string;
  amount: number;
}

/**
 * Fetch stake information. In production this would call a backend service.
 * Here we return mocked data for demonstration and testing purposes.
 */
export async function fetchStakeInfo(): Promise<StakeInfo[]> {
  return [
    { miner: 'demoMiner', amount: 100 },
    { miner: 'testMiner', amount: 50 }
  ];
}
