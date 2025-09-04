export type State = Record<string, unknown>;

class Store {
  private state: State = {};

  set(key: string, value: unknown): void {
    this.state[key] = value;
  }

  get<T>(key: string): T | undefined {
    return this.state[key] as T | undefined;
  }

  reset(): void {
    this.state = {};
  }
}

export const store = new Store();
