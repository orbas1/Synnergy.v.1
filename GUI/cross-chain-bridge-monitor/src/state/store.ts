export interface BridgeState {
  status: string;
}

export const defaultState: BridgeState = {
  status: 'idle',
};
