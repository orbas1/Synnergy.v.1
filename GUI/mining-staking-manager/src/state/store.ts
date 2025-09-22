import { StakeInfo } from '../services/stakingService';

export type StakeSubscriber = (state: StakeState) => void;
export type StakeValidator = (stake: StakeInfo) => void;

export interface StakeState {
  readonly totalStake: number;
  readonly averageStake: number;
  readonly validatorCount: number;
  readonly lastUpdated: number;
  readonly validators: ReadonlyArray<StakeInfo>;
  readonly topValidators: ReadonlyArray<StakeInfo>;
}

export interface StakeSnapshot {
  readonly timestamp: number;
  readonly totalStake: number;
  readonly validators: ReadonlyArray<StakeInfo>;
}

export interface StakingStoreOptions {
  readonly historySize?: number;
  readonly validator?: StakeValidator;
  readonly clock?: () => number;
}

interface StakeRecord extends StakeInfo {
  readonly lastUpdated: number;
}

const DEFAULT_HISTORY_SIZE = 32;

const defaultValidator: StakeValidator = (stake) => {
  if (!stake || typeof stake.miner !== 'string' || stake.miner.trim().length === 0) {
    throw new Error('miner identifier is required');
  }
  if (typeof stake.amount !== 'number' || Number.isNaN(stake.amount)) {
    throw new Error('stake amount must be a valid number');
  }
  if (stake.amount < 0) {
    throw new Error('stake amount must be non-negative');
  }
};

const defaultClock = () => Date.now();

export class StakingStore {
  private readonly historySize: number;

  private readonly validator: StakeValidator;

  private readonly clock: () => number;

  private stakes: Map<string, StakeRecord> = new Map();

  private total = 0;

  private lastUpdated = 0;

  private subscribers: Set<StakeSubscriber> = new Set();

  private history: StakeSnapshot[] = [];

  constructor(options: StakingStoreOptions = {}) {
    this.historySize = options.historySize ?? DEFAULT_HISTORY_SIZE;
    this.validator = options.validator ?? defaultValidator;
    this.clock = options.clock ?? defaultClock;
  }

  setStakes(data: StakeInfo[]): void {
    if (!Array.isArray(data)) {
      throw new TypeError('stakes payload must be an array');
    }

    const now = this.clock();
    const next = new Map<string, StakeRecord>();
    let total = 0;

    for (const item of data) {
      this.validator(item);
      const miner = item.miner.trim();
      const current = next.get(miner);
      const amount = item.amount;
      if (current) {
        const mergedAmount = current.amount + amount;
        next.set(miner, { miner, amount: mergedAmount, lastUpdated: now });
        total += amount;
      } else {
        next.set(miner, { miner, amount, lastUpdated: now });
        total += amount;
      }
    }

    let hasChanged = next.size !== this.stakes.size || this.total !== total;
    if (!hasChanged) {
      for (const [miner, record] of next) {
        const existing = this.stakes.get(miner);
        if (!existing || existing.amount !== record.amount) {
          hasChanged = true;
          break;
        }
      }
    }

    if (!hasChanged) {
      return;
    }

    this.stakes = next;
    this.total = total;
    this.lastUpdated = now;
    this.recordHistory();
    this.emit();
  }

  updateStake(update: StakeInfo): void {
    this.validator(update);
    const miner = update.miner.trim();
    const now = this.clock();
    const existing = this.stakes.get(miner);
    const amount = update.amount;

    if (existing) {
      this.total += amount - existing.amount;
    } else {
      this.total += amount;
    }

    this.stakes.set(miner, { miner, amount, lastUpdated: now });
    this.lastUpdated = now;
    this.recordHistory();
    this.emit();
  }

  removeStake(miner: string): boolean {
    const normalized = miner.trim();
    const existing = this.stakes.get(normalized);
    if (!existing) {
      return false;
    }

    this.total -= existing.amount;
    this.stakes.delete(normalized);
    this.lastUpdated = this.clock();
    this.recordHistory();
    this.emit();
    return true;
  }

  clear(): void {
    if (this.stakes.size === 0) {
      return;
    }
    this.stakes = new Map();
    this.total = 0;
    this.lastUpdated = this.clock();
    this.recordHistory();
    this.emit();
  }

  get totalStake(): number {
    return this.total;
  }

  get averageStake(): number {
    return this.stakes.size === 0 ? 0 : this.total / this.stakes.size;
  }

  get validatorCount(): number {
    return this.stakes.size;
  }

  get lastUpdateTimestamp(): number {
    return this.lastUpdated;
  }

  getStake(miner: string): StakeInfo | undefined {
    const record = this.stakes.get(miner.trim());
    return record ? { miner: record.miner, amount: record.amount } : undefined;
  }

  hasStake(miner: string): boolean {
    return this.stakes.has(miner.trim());
  }

  getAll(): StakeInfo[] {
    return Array.from(this.stakes.values()).map(({ miner, amount }) => ({ miner, amount }));
  }

  getTopValidators(count = 5): StakeInfo[] {
    return this.getAll()
      .sort((a, b) => b.amount - a.amount || a.miner.localeCompare(b.miner))
      .slice(0, Math.max(0, count));
  }

  getDistribution(): Record<string, number> {
    const distribution: Record<string, number> = {};
    if (this.total === 0) {
      return distribution;
    }
    for (const { miner, amount } of this.stakes.values()) {
      distribution[miner] = amount / this.total;
    }
    return distribution;
  }

  subscribe(subscriber: StakeSubscriber): () => void {
    if (typeof subscriber !== 'function') {
      throw new TypeError('subscriber must be a function');
    }
    this.subscribers.add(subscriber);
    subscriber(this.snapshot());
    return () => {
      this.subscribers.delete(subscriber);
    };
  }

  getHistory(): StakeSnapshot[] {
    return this.history.map((entry) => ({
      timestamp: entry.timestamp,
      totalStake: entry.totalStake,
      validators: entry.validators.map((stake) => ({ ...stake })),
    }));
  }

  private snapshot(): StakeState {
    const validators = this.getAll();
    return Object.freeze({
      totalStake: this.total,
      averageStake: this.averageStake,
      validatorCount: validators.length,
      lastUpdated: this.lastUpdated,
      validators,
      topValidators: this.getTopValidators(),
    });
  }

  private emit(): void {
    const state = this.snapshot();
    for (const subscriber of this.subscribers) {
      subscriber(state);
    }
  }

  private recordHistory(): void {
    const snapshot: StakeSnapshot = {
      timestamp: this.lastUpdated,
      totalStake: this.total,
      validators: this.getTopValidators(),
    };
    this.history.push(snapshot);
    if (this.history.length > this.historySize) {
      this.history = this.history.slice(-this.historySize);
    }
  }
}
