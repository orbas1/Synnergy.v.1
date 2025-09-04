import { start } from './main';

test('start reports configured port', () => {
  const message = start();
  expect(message).toContain('port 3000');
});
