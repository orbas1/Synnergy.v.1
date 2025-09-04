import request from 'supertest';
import { createServer } from '../../src/main';

describe('smart contract marketplace e2e', () => {
  it('serves contract listings over HTTP', async () => {
    const app = createServer();
    const res = await request(app).get('/contracts');
    expect(res.status).toBe(200);
    expect(res.body.contracts[0].name).toBe('Sample Contract');
  });
});
