import { main } from './main';

test('main returns greeting', () => {
  expect(main()).toContain('ai-marketplace');
});
