export function main(): string {
  return 'Hello from authority-node-index';
}

if (require.main === module) {
  console.log(main());
}
