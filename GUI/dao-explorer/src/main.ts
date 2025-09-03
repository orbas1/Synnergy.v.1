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

export function createDAO(name: string, creator: string): Promise<string> {
  return run(['dao', 'create', name, creator]);
}

export function joinDAO(id: string, addr: string): Promise<void> {
  return run(['dao', 'join', id, addr]).then(() => {});
}

export function leaveDAO(id: string, addr: string): Promise<void> {
  return run(['dao', 'leave', id, addr]).then(() => {});
}

export function getDAOInfo(id: string): Promise<string> {
  return run(['dao', 'info', id]);
}

export function listDAOs(): Promise<string> {
  return run(['dao', 'list']);
}

if (require.main === module) {
  listDAOs().then(console.log).catch(err => {
    console.error('dao explorer error', err);
    process.exit(1);
  });
}
