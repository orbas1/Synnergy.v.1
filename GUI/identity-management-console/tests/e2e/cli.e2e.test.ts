import { main } from '../../src/main';

test('cli registers user', () => {
  const output = main(['node', 'test', '--register', 'carol', '--key', 'pub']);
  expect(output).toBe('Registered carol');
});
