import { main } from './main';

test('main returns greeting', () => {
  expect(main()).toContain('authority-node-index');
});
