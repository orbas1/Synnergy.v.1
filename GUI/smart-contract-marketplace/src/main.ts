export function main(): string {
  return 'Hello from smart-contract-marketplace';
}

if (require.main === module) {
  console.log(main());
}
