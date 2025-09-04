import http from 'http';
import { AddressInfo } from 'net';
import { main } from './main';

test('server responds with alive message', async () => {
  const server = main(0);
  const { port } = server.address() as AddressInfo;
  const data = await new Promise<string>((resolve) => {
    http.get(`http://127.0.0.1:${port}`, (res) => {
      let buf = '';
      res.on('data', (chunk) => (buf += chunk));
      res.on('end', () => resolve(buf));
    });
  });
  expect(data).toBe('ai-marketplace alive');
  server.close();
});
