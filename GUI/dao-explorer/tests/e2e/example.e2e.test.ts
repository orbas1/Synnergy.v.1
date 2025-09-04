import { main } from '../../src/main';

test('CLI main executes with custom name', () => {
  expect(main('e2e')).toContain('e2e');
});
