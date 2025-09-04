import { StakeInfo } from '../services/stakingService';

/**
 * Lightweight in-memory store for stake data. This is a placeholder for a more
 * sophisticated state management solution.
 */
export class StakingStore {
  private stakes: StakeInfo[] = [];

  setStakes(data: StakeInfo[]): void {
    this.stakes = data;
  }

  get totalStake(): number {
    return this.stakes.reduce((sum, s) => sum + s.amount, 0);
  }
}
