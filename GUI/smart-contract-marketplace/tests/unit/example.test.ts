import { main } from '../../src/main';

test('main returns greeting', () => {
  expect(main()).toBe('Hello from smart-contract-marketplace');
});
