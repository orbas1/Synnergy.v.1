import http from 'http';
import { AddressInfo } from 'net';

export function main(port: number = Number(process.env.PORT) || 3000): http.Server {
  const server = http.createServer((_req, res) => {
    res.writeHead(200, { 'Content-Type': 'text/plain' });
    res.end('ai-marketplace alive');
  });
  server.listen(port);
  return server;
}

if (require.main === module) {
  const portArg = parseInt(process.argv[2] || '', 10);
  const port = Number.isNaN(portArg) ? undefined : portArg;
  const server = main(port);
  const address = server.address() as AddressInfo;
  console.log(`AI marketplace listening on port ${address.port}`);
}
