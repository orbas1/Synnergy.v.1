import express, { Request, Response } from 'express';

export interface Contract {
  id: string;
  name: string;
}

/**
 * createServer initialises the Smart Contract Marketplace API.
 * The server currently exposes a single read‑only endpoint that
 * returns a list of available contracts. This stub enables future
 * integration with the Synnergy virtual machine and wallet flows
 * while providing a predictable interface for the GUI and CLI.
 */
export function createServer() {
  const app = express();

  const contracts: Contract[] = [{ id: '1', name: 'Sample Contract' }];

  app.get('/contracts', (_: Request, res: Response) => {
    res.json({ contracts });
  });

  return app;
}

// If executed directly start the HTTP server.  The port is configurable
// via the PORT environment variable so that container deployments can
// bind to any required interface.
if (require.main === module) {
  const port = process.env.PORT || 80;
  createServer().listen(port, () => {
    // eslint-disable-next-line no-console
    console.log(`Smart-contract marketplace listening on port ${port}`);
  });
}
