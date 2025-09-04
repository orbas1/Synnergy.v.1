import { main } from './main';

test('main returns greeting', () => {
  expect(main()).toContain('token-creation-tool');
});
