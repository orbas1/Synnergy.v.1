import { main } from './main';

test('main greets provided name', () => {
  const result = main(['node', 'test', '--name', 'user']);
  expect(result).toBe('Hello from user');
});

test('registers user via CLI', () => {
  const result = main(['node', 'test', '--register', 'alice', '--key', 'pub']);
  expect(result).toBe('Registered alice');
});
