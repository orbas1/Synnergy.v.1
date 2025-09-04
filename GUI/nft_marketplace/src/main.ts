import { showHomePage } from './pages/Home';

export async function main(): Promise<string> {
  return showHomePage();
}

if (require.main === module) {
  main()
    .then((output) => console.log(output))
    .catch((err) => {
      console.error('Failed to start NFT marketplace', err);
      process.exit(1);
    });
}
