import { StakingStore } from '../../src/state/store';
import type { StakeState } from '../../src/state/store';

describe('StakingStore', () => {
  test('manages stake lifecycle and publishes updates', () => {
    let tick = 0;
    const store = new StakingStore({ historySize: 4, clock: () => ++tick });

    const snapshots: StakeState[] = [];
    const unsubscribe = store.subscribe((state) => {
      snapshots.push(state);
    });

    store.setStakes([
      { miner: 'alpha', amount: 100 },
      { miner: 'beta', amount: 40 },
      { miner: 'alpha', amount: 10 }
    ]);

    expect(store.totalStake).toBe(150);
    expect(store.validatorCount).toBe(2);
    expect(store.getStake('alpha')).toEqual({ miner: 'alpha', amount: 110 });
    expect(store.getTopValidators(1)).toEqual([{ miner: 'alpha', amount: 110 }]);

    store.updateStake({ miner: 'beta', amount: 90 });
    expect(store.totalStake).toBe(200);
    expect(store.averageStake).toBe(100);
    expect(store.getDistribution()).toMatchObject({
      alpha: 110 / 200,
      beta: 90 / 200
    });

    expect(store.getHistory()).toHaveLength(2);
    expect(store.getHistory()[1].validators[0]).toEqual({ miner: 'alpha', amount: 110 });

    unsubscribe();

    expect(snapshots).toHaveLength(3); // initial snapshot + setStakes + updateStake
    expect(snapshots[0].totalStake).toBe(0);
    expect(snapshots[1].totalStake).toBe(150);
    expect(snapshots[2].totalStake).toBe(200);
  });

  test('removes stakes, clears state and bounds history', () => {
    let tick = 100;
    const store = new StakingStore({ historySize: 2, clock: () => ++tick });

    store.setStakes([
      { miner: 'gamma', amount: 60 },
      { miner: 'delta', amount: 40 }
    ]);
    store.updateStake({ miner: 'gamma', amount: 80 });

    expect(store.removeStake('delta')).toBe(true);
    expect(store.totalStake).toBe(80);
    expect(store.hasStake('delta')).toBe(false);

    expect(store.removeStake('missing')).toBe(false);

    store.clear();
    expect(store.totalStake).toBe(0);
    expect(store.validatorCount).toBe(0);

    const history = store.getHistory();
    expect(history).toHaveLength(2);
    expect(history[history.length - 1].totalStake).toBe(0);
  });
});
