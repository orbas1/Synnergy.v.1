import { execSync } from 'child_process';

export interface HeadInfo {
  height: number;
  hash: string;
}

export function parseHead(out: string): HeadInfo {
  const [heightStr, hash] = out.trim().split(/\s+/);
  return { height: parseInt(heightStr, 10), hash };
}

export function fetchHead(): HeadInfo {
  const out = execSync('synnergy ledger head', { encoding: 'utf8' });
  return parseHead(out);
}

async function main() {
  try {
    const head = fetchHead();
    console.log(`height: ${head.height} hash: ${head.hash}`);
  } catch (err) {
    console.error('explorer error', err);
  }
}

if (require.main === module) {
  main();
}
