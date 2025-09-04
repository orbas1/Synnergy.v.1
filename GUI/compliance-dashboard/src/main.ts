import express from 'express';
import dotenv from 'dotenv';
import winston from 'winston';

dotenv.config();

const logger = winston.createLogger({
  level: 'info',
  transports: [new winston.transports.Console()],
});

export function createApp() {
  const app = express();
  app.get('/health', (_req, res) => {
    res.json({ status: 'ok' });
  });
  return app;
}

export function startServer(port: number) {
  const app = createApp();
  return app.listen(port, () => {
    logger.info(`Compliance dashboard listening on port ${port}`);
  });
}

if (require.main === module) {
  const port = Number(process.env.PORT) || 3000;
  startServer(port);
}
