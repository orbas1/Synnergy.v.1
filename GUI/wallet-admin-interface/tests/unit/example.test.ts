import request from 'supertest';
import { createServer } from '../../src/main';

describe('signature verification', () => {
  it('rejects invalid payload', async () => {
    const app = createServer();
    const res = await request(app).post('/verify').send({});
    expect(res.status).toBe(400);
  });
});
