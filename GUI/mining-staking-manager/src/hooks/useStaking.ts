import { fetchStakeInfo, StakeInfo } from '../services/stakingService';

/**
 * Simple hook-like helper that retrieves stake info and returns the total
 * amount staked. Although this isn't a React hook, the function mirrors the
 * behaviour expected from one and can be adapted for a real frontend later.
 */
export async function useStakingTotal(): Promise<number> {
  const stakes: StakeInfo[] = await fetchStakeInfo();
  return stakes.reduce((sum, s) => sum + s.amount, 0);
}
