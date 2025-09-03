import { execSync } from 'child_process';

function createWallet(password: string, out: string) {
  const cmd = `synnergy wallet new --out ${out} --password ${password} --json`;
  const res = execSync(cmd, { encoding: 'utf8' });
  return JSON.parse(res);
}

function getBalance(address: string) {
  const cmd = `synnergy ledger balance ${address}`;
  return execSync(cmd, { encoding: 'utf8' }).trim();
}

async function main() {
  const walletFile = 'wallet.json';
  const password = 'changeit';
  const wallet = createWallet(password, walletFile);
  const bal = getBalance(wallet.address);
  console.log(`Wallet ${wallet.address} balance: ${bal}`);
}

main().catch(err => {
  console.error('wallet error', err);
});
