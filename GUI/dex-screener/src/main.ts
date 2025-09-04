import { spawn } from 'child_process';

function run(args: string[]): Promise<string> {
  return new Promise((resolve, reject) => {
    const child = spawn('synnergy', args);
    let out = '';
    let err = '';
    child.stdout.on('data', d => out += d);
    child.stderr.on('data', d => err += d);
    child.on('close', code => {
      if (code === 0) {
        resolve(out.trim());
      } else {
        reject(new Error(err.trim() || `exit ${code}`));
      }
    });
  });
}

export interface PoolView {
  id: string;
  token_a: string;
  token_b: string;
  reserve_a: number;
  reserve_b: number;
  fee_bps: number;
}

export function listPools(): Promise<PoolView[]> {
  return run(['liquidity_views', 'list']).then(out => JSON.parse(out));
}

export function poolInfo(id: string): Promise<PoolView> {
  return run(['liquidity_views', 'info', id]).then(out => JSON.parse(out));
}

if (require.main === module) {
  listPools().then(pools => {
    console.log(JSON.stringify(pools, null, 2));
  }).catch(err => {
    console.error('dex screener error', err);
    process.exit(1);
  });
}
