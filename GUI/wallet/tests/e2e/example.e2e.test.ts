import http from 'http';
import type { AddressInfo } from 'net';
import { main } from '../../src/main';

async function withTestServer(handler: (server: http.Server) => Promise<void>) {
  const server = http.createServer((req, res) => {
    if (!req.url) {
      res.statusCode = 400;
      res.end('missing url');
      return;
    }

    if (req.url.endsWith('/balance')) {
      res.setHeader('Content-Type', 'application/json');
      res.end(
        JSON.stringify({
          address: 'syn1e2e',
          balance: 200,
          pending: 20,
          currency: 'SYN',
          updatedAt: new Date('2024-01-01T00:00:00Z').toISOString()
        })
      );
      return;
    }

    res.statusCode = 404;
    res.end('not found');
  });

  await new Promise<void>((resolve) => server.listen(0, resolve));

  try {
    await handler(server);
  } finally {
    await new Promise<void>((resolve) => server.close(() => resolve()));
  }
}

describe('wallet CLI e2e', () => {
  it('renders formatted balance output from a live server', async () => {
    await withTestServer(async (server) => {
      const { port } = server.address() as AddressInfo;
      process.env.WALLET_API_URL = `http://127.0.0.1:${port}`;

      const output = await main('syn1e2e');

      expect(output).toContain('syn1e2e');
      expect(output).toContain('Available');
      expect(output).toContain('Synnergy Wallet Overview');
    });
  });
});
