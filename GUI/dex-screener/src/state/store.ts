import { PoolInfo } from '../services/api';

export class Store {
  private lastQuery?: PoolInfo;

  setLastQuery(info: PoolInfo): void {
    this.lastQuery = info;
  }

  getLastQuery(): PoolInfo | undefined {
    return this.lastQuery;
  }
}

export const store = new Store();
