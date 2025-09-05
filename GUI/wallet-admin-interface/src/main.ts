import express, { Application } from 'express';
import helmet from 'helmet';
import crypto from 'crypto';

export function createServer(): Application {
  const app = express();
  app.use(helmet());
  app.use(express.json());

  app.get('/health', (_req, res) => {
    res.json({ status: 'ok' });
  });

  app.post('/verify', (req, res) => {
    const { message, signature, publicKey } = req.body || {};
    try {
      const verifier = crypto.createVerify('SHA256');
      verifier.update(message || '');
      const valid = verifier.verify(publicKey, signature, 'hex');
      res.json({ valid });
    } catch {
      res.status(400).json({ error: 'invalid payload' });
    }
  });

  return app;
}

export function main(): void {
  const port = Number(process.env.PORT) || 3000;
  const app = createServer();
  app.listen(port, () => {
    console.log(`Wallet Admin Interface listening on port ${port}`);
  });
}

if (require.main === module) {
  main();
}
