import { renderDashboard } from './components/MiningDashboard';

/**
 * Entrypoint for the mining-staking-manager CLI. It renders a textual dashboard
 * showing aggregate staking information.
 */
export async function main(): Promise<string> {
  return renderDashboard();
}

if (require.main === module) {
  main()
    .then(console.log)
    .catch(err => {
      console.error(err);
      process.exit(1);
    });
}
