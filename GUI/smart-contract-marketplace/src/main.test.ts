import request from 'supertest';
import { createServer } from './main';

describe('createServer', () => {
  it('exposes contract listings', async () => {
    const app = createServer();
    const res = await request(app).get('/contracts');
    expect(res.status).toBe(200);
    expect(res.body.contracts).toEqual(
      expect.arrayContaining([{ id: '1', name: 'Sample Contract' }])
    );
  });
});

