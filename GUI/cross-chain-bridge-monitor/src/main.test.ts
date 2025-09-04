import { main } from './main';

test('main returns greeting', () => {
  expect(main()).toContain('cross-chain-bridge-monitor');
});
