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

export interface TokenOptions {
  name: string;
  symbol: string;
  owner: string;
  decimals: number;
  supply: number;
}

export function createToken(kind: string, opts: TokenOptions): Promise<string> {
  const args = [
    kind,
    'create',
    '--name', opts.name,
    '--symbol', opts.symbol,
    '--owner', opts.owner,
    '--dec', String(opts.decimals),
    '--supply', String(opts.supply)
  ];
  return run(args);
}

export function listTokens(kind: string): Promise<string> {
  return run([kind, 'list']);
}

if (require.main === module) {
  // Example usage: node src/main.ts syn500 MyToken MTK owner 2 1000
  const [kind, name, symbol, owner, dec, supply] = process.argv.slice(2);
  createToken(kind, {
    name,
    symbol,
    owner,
    decimals: Number(dec),
    supply: Number(supply)
  }).then(res => {
    console.log(res);
  }).catch(err => {
    console.error('token creation error', err);
    process.exit(1);
  });
}
