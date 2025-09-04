export function main(): string {
  return 'Hello from cross-chain-management';
}

if (require.main === module) {
  console.log(main());
}
