import request from 'supertest';
import { createServer } from '../../src/main';

test('server provides contracts endpoint', async () => {
  const app = createServer();
  const res = await request(app).get('/contracts');
  expect(res.status).toBe(200);
  expect(res.body.contracts).toHaveLength(1);
});
