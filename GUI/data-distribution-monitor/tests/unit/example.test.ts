import { renderDashboard } from '../../src/pages/dashboard';

it('renders dashboard text', () => {
  expect(renderDashboard()).toMatch(/dashboard/);
});
