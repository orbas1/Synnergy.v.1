import axios from 'axios';
import { main } from '../../src/main';

jest.mock('axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

test('CLI exits with status information', async () => {
  mockedAxios.get.mockResolvedValueOnce({ data: 'ok' });
  const log = jest.spyOn(console, 'log').mockImplementation(() => {});
  await main(['node', 'main.ts']);
  expect(log).toHaveBeenCalledWith('ok');
  log.mockRestore();
});
