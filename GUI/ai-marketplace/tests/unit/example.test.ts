import { store } from '../../src/state/store';

test('reset clears state', () => {
  store.set('count', 1);
  store.reset();
  expect(store.get('count')).toBeUndefined();
});
