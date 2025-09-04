export function main(): string {
  return 'Hello from explorer';
}

if (require.main === module) {
  console.log(main());
}
