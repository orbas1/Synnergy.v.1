import { fetchNodeStatus } from './services/status';

export async function main(): Promise<string> {
  const status = await fetchNodeStatus();
  return `Node status: ${status}`;
}

if (require.main === module) {
  main()
    .then((output) => console.log(output))
    .catch((err) => {
      console.error('Failed to start node-operations-dashboard', err);
      process.exit(1);
    });
}
