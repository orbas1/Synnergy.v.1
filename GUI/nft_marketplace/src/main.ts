import { spawn } from 'child_process';

export function mintNFT(id: string, owner: string, metadata: string, price: string): Promise<void> {
  return new Promise((resolve, reject) => {
    const child = spawn('synnergy', ['nft', 'mint', id, owner, metadata, price]);
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

export function buyNFT(id: string, newOwner: string): Promise<void> {
  return new Promise((resolve, reject) => {
    const child = spawn('synnergy', ['nft', 'buy', id, newOwner]);
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

export function listNFT(id: string): Promise<string> {
  return new Promise((resolve, reject) => {
    const child = spawn('synnergy', ['nft', 'list', id]);
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
