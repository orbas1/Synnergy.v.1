import { spawnSync } from 'child_process';
import path from 'path';

test('CLI prints configured endpoint', () => {
  const result = spawnSync(
    'node',
    ['-r', 'ts-node/register', path.join(__dirname, '../../src/main.ts')],
    {
      env: { ...process.env, STORAGE_MARKETPLACE_ENDPOINT: 'http://e2e-endpoint' },
    }
  );
  const output = result.stdout.toString();
  expect(output).toContain('http://e2e-endpoint');
});
