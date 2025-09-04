import { renderHome } from '../../src/pages/home';

test('renders home text', () => {
  expect(renderHome()).toContain('Home');
});
