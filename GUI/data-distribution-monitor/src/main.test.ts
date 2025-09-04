import { main } from './main';

test('status command returns operational message', () => {
  expect(main(['status'])).toContain('operational');
});

test('unknown command throws error', () => {
  expect(() => main(['bad'])).toThrow('Unknown command');
});
