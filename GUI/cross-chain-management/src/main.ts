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

export interface BridgeInfo {
  chainId: string;
  endpoint: string;
}

export function listBridges(): Promise<string> {
  return run(['cross-chain-management', 'bridges']);
}

export function connectChain(info: BridgeInfo): Promise<string> {
  return run([
    'cross-chain-management',
    'connect',
    '--chain', info.chainId,
    '--endpoint', info.endpoint
  ]);
}

if (require.main === module) {
  const [cmd, ...rest] = process.argv.slice(2);
  if (cmd === 'bridges') {
    listBridges().then(console.log).catch(err => {
      console.error('bridge list error', err);
      process.exit(1);
    });
  } else if (cmd === 'connect') {
    const [chainId, endpoint] = rest;
    connectChain({ chainId, endpoint }).then(console.log).catch(err => {
      console.error('connect chain error', err);
      process.exit(1);
    });
  } else {
    console.error('usage: node main.ts [bridges|connect <chainId> <endpoint>]');
    process.exit(1);
  }
}
