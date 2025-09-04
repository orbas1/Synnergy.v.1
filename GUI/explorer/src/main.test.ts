import axios from 'axios';
import { getStatus, main } from './main';

jest.mock('axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

test('getStatus retrieves data from endpoint', async () => {
  mockedAxios.get.mockResolvedValueOnce({ data: 'ok' });
  await expect(getStatus('http://host')).resolves.toBe('ok');
});

test('main prints status', async () => {
  mockedAxios.get.mockResolvedValueOnce({ data: 'ok' });
  const log = jest.spyOn(console, 'log').mockImplementation(() => {});
  await main(['node', 'main.ts']);
  expect(log).toHaveBeenCalledWith('ok');
  log.mockRestore();
});
