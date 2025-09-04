import { execSync } from 'child_process';
import path from 'path';

test('CLI status command reports operational', () => {
  const cli = path.resolve(__dirname, '../../src/main.ts');
  const output = execSync(`node -r ts-node/register ${cli} status`).toString();
  expect(output).toContain('operational');
});

