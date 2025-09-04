import { main } from './main';

test('main greets provided name', () => {
  const result = main(['node', 'test', '--name', 'user']);
  expect(result).toBe('Hello from user');
});
