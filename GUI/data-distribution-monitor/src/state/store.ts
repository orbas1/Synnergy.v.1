import { DistributionStatus } from '../services/api';

export class Store {
  private status: DistributionStatus = { online: false };

  setStatus(s: DistributionStatus): void {
    this.status = s;
  }

  getStatus(): DistributionStatus {
    return this.status;
  }
}

export const store = new Store();
