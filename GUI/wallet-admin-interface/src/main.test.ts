import { main } from './main';

test('main returns greeting', () => {
  expect(main()).toContain('wallet-admin-interface');
});
