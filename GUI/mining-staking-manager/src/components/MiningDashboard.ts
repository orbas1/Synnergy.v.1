import { fetchStakeInfo } from '../services/stakingService';
import { StakingStore } from '../state/store';

/**
 * Render a basic textual dashboard with total stake information.
 */
export async function renderDashboard(): Promise<string> {
  const store = new StakingStore();
  const data = await fetchStakeInfo();
  store.setStakes(data);
  return `Total Stake: ${store.totalStake}`;
}
