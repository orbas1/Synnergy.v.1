import request from 'supertest';
import { createApp } from './main';

describe('compliance dashboard server', () => {
  it('returns ok on /health', async () => {
    const app = createApp();
    const res = await request(app).get('/health');
    expect(res.status).toBe(200);
    expect(res.body.status).toBe('ok');
  });
});
