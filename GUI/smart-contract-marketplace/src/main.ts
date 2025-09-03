import { spawn } from 'child_process';

export function deployContract(path: string, owner: string): Promise<string> {
  return new Promise((resolve, reject) => {
    const child = spawn('synnergy', ['marketplace', 'deploy', path, owner]);
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

export function tradeContract(addr: string, newOwner: string): Promise<void> {
  return new Promise((resolve, reject) => {
    const child = spawn('synnergy', ['marketplace', 'trade', addr, newOwner]);
    let err = '';
    child.stderr.on('data', d => err += d);
    child.on('close', code => {
      if (code === 0) {
        resolve();
      } else {
        reject(new Error(err.trim() || `exit ${code}`));
      }
    });
  });
}
