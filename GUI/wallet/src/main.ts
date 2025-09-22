import { renderHome } from './pages/Home';
import { useBalance } from './hooks/useBalance';
import { WalletStore, BalanceProvider } from './state/walletStore';
import { WalletService } from './services/walletService';

export interface MainOptions {
  store?: BalanceProvider;
}

export async function main(addressInput?: string, options: MainOptions = {}): Promise<string> {
  const candidate = addressInput ?? process.env.WALLET_ADDRESS ?? '';
  const address = candidate.trim();

  if (!address) {
    throw new Error('A wallet address must be provided via an argument or the WALLET_ADDRESS environment variable.');
  }

  const provider = options.store ?? new WalletStore(new WalletService());
  const summary = await useBalance(address, { refresh: true, store: provider });
  return renderHome(summary);
}

if (require.main === module) {
  const argumentAddress = process.argv[2];
  main(argumentAddress)
    .then((output) => {
      console.log(output);
    })
    .catch((error: unknown) => {
      const message = error instanceof Error ? error.message : 'Unknown error';
      console.error(`Wallet CLI failed: ${message}`);
      process.exitCode = 1;
    });
}
