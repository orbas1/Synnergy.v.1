import { renderDashboard } from '../../src/components/MiningDashboard';

describe('MiningDashboard', () => {
  test('renderDashboard aggregates stake amounts', async () => {
    const output = await renderDashboard();
    expect(output).toBe('Total Stake: 150');
  });
});
