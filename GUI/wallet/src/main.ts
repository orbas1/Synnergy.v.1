export function main(): string {
  return 'Hello from wallet';
}

if (require.main === module) {
  console.log(main());
}
