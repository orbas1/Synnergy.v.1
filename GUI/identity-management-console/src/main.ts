import { Command } from 'commander';

export function main(argv: string[] = process.argv): string {
  const program = new Command();
  program
    .name('identity-console')
    .description('Identity management console CLI')
    .option('-n, --name <name>', 'name to greet', 'identity-management-console');

  program.parse(argv);
  const opts = program.opts<{ name: string }>();
  const greeting = `Hello from ${opts.name}`;
  return greeting;
}

if (require.main === module) {
  console.log(main());
}
