import yargs from 'yargs';
import { hideBin } from 'yargs/helpers';

export function main(name: string = 'dao-explorer'): string {
  return `Hello from ${name}`;
}

if (require.main === module) {
  const argv = yargs(hideBin(process.argv))
    .option('name', {
      type: 'string',
      description: 'Name to greet',
      default: 'dao-explorer',
    })
    .strict()
    .help()
    .parseSync();

  console.log(main(argv.name));
}
