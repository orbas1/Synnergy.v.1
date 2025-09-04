import http from 'http';
import { AddressInfo } from 'net';

// main starts an HTTP server that reports the portal is alive. The port defaults
// to 3000 but can be overridden via an argument or the PORT env var.
export function main(port: number = Number(process.env.PORT) || 3000): http.Server {
  const server = http.createServer((_req, res) => {
    res.writeHead(200, { 'Content-Type': 'text/plain' });
    res.end('validator-governance-portal alive');
  });
  server.listen(port);
  return server;
}

if (require.main === module) {
  const portArg = parseInt(process.argv[2] || '', 10);
  const port = Number.isNaN(portArg) ? undefined : portArg;
  const server = main(port);
  const address = server.address() as AddressInfo;
  console.log(`Validator governance portal listening on port ${address.port}`);
}
