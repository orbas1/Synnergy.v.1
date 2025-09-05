import { main } from './main';

test('main renders home with balance', () => {
  expect(main(5)).toContain('Current balance: 5');
});
