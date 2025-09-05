import request from 'supertest';
import { createServer } from '../../src/main';

describe('wallet admin interface e2e', () => {
  it('responds to health check', async () => {
    const app = createServer();
    const res = await request(app).get('/health');
    expect(res.status).toBe(200);
    expect(res.body).toEqual({ status: 'ok' });
  });
});
