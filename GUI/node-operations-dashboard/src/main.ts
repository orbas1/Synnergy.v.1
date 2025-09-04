import { fetchNodeStatus } from './services/status';

export async function main(): Promise<string> {
  try {
    const status = await fetchNodeStatus();
    return `Node status: ${status}`;
  } catch {
    // Ensure callers always receive a string even on unexpected errors
    return 'Node status: ERROR';
  }
}

if (require.main === module) {
  main()
    .then((output) => console.log(output))
    .catch((err) => {
      console.error('Failed to start node-operations-dashboard', err);
      process.exit(1);
    });
}
