import { Command } from 'commander';
import axios from 'axios';

export async function getStatus(url: string): Promise<string> {
  const res = await axios.get(`${url}/status`);
  return res.data;
}

export async function main(argv: string[] = process.argv): Promise<void> {
  const program = new Command();
  program
    .name('synnergy-explorer')
    .description('CLI explorer for the Synnergy network')
    .option('-u, --url <url>', 'base URL of a Synnergy node', 'http://localhost:8080')
    .action(async (opts) => {
      try {
        const status = await getStatus(opts.url);
        console.log(status);
      } catch (err) {
        console.error('failed to fetch status', err);
        process.exitCode = 1;
      }
    });

  await program.parseAsync(argv);
}

if (require.main === module) {
  main();
}
