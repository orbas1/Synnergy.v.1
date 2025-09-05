import { getBalance } from '../services/api';

export async function useBalance(address: string): Promise<number> {
  try {
    return await getBalance(address);
  } catch {
    return 0;
  }
}
