import request from 'supertest';
import { createServer } from './main';

describe('wallet-admin-interface server', () => {
  const app = createServer();

  it('responds to health check', async () => {
    const res = await request(app).get('/health');
    expect(res.status).toBe(200);
    expect(res.body.status).toBe('ok');
  });

  it('rejects invalid signatures', async () => {
    const res = await request(app)
      .post('/verify')
      .send({ message: 'test', signature: '00', publicKey: '00' });
    expect(res.status).toBe(400);
  });
});
