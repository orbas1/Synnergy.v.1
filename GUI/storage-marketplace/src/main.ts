export function main(): string {
  return 'Hello from storage-marketplace';
}

if (require.main === module) {
  console.log(main());
}
