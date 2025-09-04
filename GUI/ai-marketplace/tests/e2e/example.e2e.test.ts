import { store } from '../../src/state/store';

test('stores and retrieves values', () => {
  store.set('foo', 'bar');
  expect(store.get<string>('foo')).toBe('bar');
  store.reset();
});
