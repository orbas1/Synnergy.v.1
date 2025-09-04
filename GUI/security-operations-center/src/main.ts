export function main(): string {
  return 'Hello from security-operations-center';
}

if (require.main === module) {
  console.log(main());
}
