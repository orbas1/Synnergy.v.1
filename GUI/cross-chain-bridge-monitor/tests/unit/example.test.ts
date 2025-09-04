import { getBridgeStatus } from '../../src/services/bridgeService';

test('bridge service returns ok', () => {
  expect(getBridgeStatus()).toBe('ok');
});
