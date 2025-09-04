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

export interface AuthorityNodeInfo {
  address: string;
  publicKey: string;
}

export function listNodes(): Promise<string> {
  return run(['authority-node-index', 'list']);
}

export function registerNode(info: AuthorityNodeInfo): Promise<string> {
  return run([
    'authority-node-index',
    'register',
    '--address', info.address,
    '--pubkey', info.publicKey
  ]);
}

if (require.main === module) {
  const [cmd, ...rest] = process.argv.slice(2);
  if (cmd === 'list') {
    listNodes().then(console.log).catch(err => {
      console.error('node list error', err);
      process.exit(1);
    });
  } else if (cmd === 'register') {
    const [address, pubkey] = rest;
    registerNode({ address, publicKey: pubkey }).then(console.log).catch(err => {
      console.error('node register error', err);
      process.exit(1);
    });
  } else {
    console.error('usage: node main.ts [list|register <address> <pubkey>]');
    process.exit(1);
  }
}
