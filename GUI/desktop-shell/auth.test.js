const { test } = require('node:test');
const assert = require('assert');
const { authenticate } = require('./auth');

test('authenticate returns deterministic token for same key', () => {
  const key = 'test-key';
  const token1 = authenticate(key);
  const token2 = authenticate(key);
  assert.strictEqual(token1, token2);
});

test('authenticate throws when key is missing', () => {
  assert.throws(() => authenticate(), /DESKTOP_SHELL_PRIVATE_KEY/);
});

test('authenticate tokens differ for different keys', () => {
  const token1 = authenticate('key-one');
  const token2 = authenticate('key-two');
  assert.notStrictEqual(token1, token2);
});
