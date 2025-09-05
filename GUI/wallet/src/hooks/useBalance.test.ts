import { useBalance } from './useBalance';
import * as api from '../services/api';

test('useBalance returns balance from service', async () => {
  jest.spyOn(api, 'getBalance').mockResolvedValue(10);
  await expect(useBalance('addr')).resolves.toBe(10);
});
