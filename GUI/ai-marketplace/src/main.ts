import { execSync } from 'child_process';

export function parseDeploy(out: string): string {
  return out.trim().split(/\s+/).pop() || '';
}

export function deploy(wasm: string, modelHash: string, manifest: string, gas: string, owner: string): string {
  const cmd = `synnergy ai_contract deploy ${wasm} ${modelHash} ${manifest} ${gas} ${owner}`;
  const out = execSync(cmd, { encoding: 'utf8' });
  return parseDeploy(out);
}

async function main() {
  const [, , wasm, modelHash, manifest, gas, owner] = process.argv;
  if (!wasm || !modelHash || !manifest || !gas || !owner) {
    console.error('usage: ts-node main.ts <wasm> <model_hash> <manifest> <gas_limit> <owner>');
    process.exit(1);
  }
  try {
    const addr = deploy(wasm, modelHash, manifest, gas, owner);
    console.log(`contract address: ${addr}`);
  } catch (err) {
    console.error('deployment failed', err);
  }
}

if (require.main === module) {
  main();
}
