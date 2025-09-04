import { createApp } from '../../src/main';

test('app exposes health endpoint', () => {
  const app = createApp();
  const hasHealth = app._router.stack.some((r: any) => r.route && r.route.path === '/health');
  expect(hasHealth).toBe(true);
});
