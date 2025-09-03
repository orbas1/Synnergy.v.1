import { spawn } from 'child_process';

export function createListing(hash: string, price: number, owner: string): Promise<string> {
  return new Promise((resolve, reject) => {
    const child = spawn('synnergy', ['storage_marketplace', 'list', hash, String(price), owner]);
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

export function listListings(): Promise<any[]> {
  return new Promise((resolve, reject) => {
    const child = spawn('synnergy', ['storage_marketplace', 'listings']);
    let out = '';
    let err = '';
    child.stdout.on('data', d => out += d);
    child.stderr.on('data', d => err += d);
    child.on('close', code => {
      if (code === 0) {
        try {
          resolve(JSON.parse(out));
        } catch (e) {
          reject(e);
        }
      } else {
        reject(new Error(err.trim() || `exit ${code}`));
      }
    });
  });
}

export function openDeal(listingID: string, buyer: string): Promise<string> {
  return new Promise((resolve, reject) => {
    const child = spawn('synnergy', ['storage_marketplace', 'deal', listingID, buyer]);
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
