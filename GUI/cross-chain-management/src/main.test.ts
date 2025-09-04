import { buildCLI, getStatus } from './main';

test('getStatus returns operational message', () => {
  expect(getStatus()).toContain('operational');
});

test('CLI status command outputs status', () => {
  const log = jest.spyOn(console, 'log').mockImplementation(() => {});
  const program = buildCLI();
  program.exitOverride();
  program.parse(['status'], { from: 'user' });
  expect(log.mock.calls.join('\n')).toContain('operational');
  log.mockRestore();
});

