export function main(): string {
  return 'Hello from data-distribution-monitor';
}

if (require.main === module) {
  console.log(main());
}
