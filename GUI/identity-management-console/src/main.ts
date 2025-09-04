import { IdentityService } from './services/identityService';

function parseArgs(argv: string[]): Record<string, string | boolean> {
  const args: Record<string, string | boolean> = {};
  for (let i = 2; i < argv.length; i++) {
    const arg = argv[i];
    if (arg.startsWith('--')) {
      const key = arg.slice(2);
      const next = argv[i + 1];
      if (!next || next.startsWith('-')) {
        args[key] = true;
      } else {
        args[key] = next;
        i++;
      }
    } else if (arg.startsWith('-')) {
      const key = arg.slice(1);
      const next = argv[i + 1];
      if (!next || next.startsWith('-')) {
        args[key] = true;
      } else {
        args[key] = next;
        i++;
      }
    }
  }
  return args;
}

export function main(argv: string[] = process.argv): string {
  const service = new IdentityService();
  const opts = parseArgs(argv) as { name?: string; register?: string; key?: string };
  if (opts.register) {
    if (!opts.key) throw new Error('public key required');
    service.register(opts.register, opts.key);
    return `Registered ${opts.register}`;
  }
  const name = opts.name || 'identity-management-console';
  return `Hello from ${name}`;
}

if (require.main === module) {
  console.log(main());
}
