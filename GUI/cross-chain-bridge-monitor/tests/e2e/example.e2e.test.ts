import { getBridgeStatus } from '../../src/services/bridgeService';

test('basic bridge workflow returns ok status', () => {
  expect(getBridgeStatus()).toBe('ok');
});
