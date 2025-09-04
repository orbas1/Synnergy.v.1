import { IdentityService } from '../../src/services/identityService';

test('register and retrieve user', () => {
  const svc = new IdentityService();
  svc.register('bob', 'pubkey');
  expect(svc.getUser('bob')).toBe('pubkey');
});

