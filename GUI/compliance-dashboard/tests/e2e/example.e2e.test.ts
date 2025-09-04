import request from 'supertest';
import { startServer } from '../../src/main';

describe('compliance dashboard e2e', () => {
  let server: any;

  beforeAll(() => {
    server = startServer(0); // let OS choose port
  });

  afterAll(() => {
    server.close();
  });

  it('responds to /health', async () => {
    const address = server.address();
    const res = await request(`http://127.0.0.1:${address.port}`).get('/health');
    expect(res.status).toBe(200);
    expect(res.body.status).toBe('ok');
  });
});
