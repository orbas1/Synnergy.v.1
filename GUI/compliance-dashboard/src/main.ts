export function main(): string {
  return 'Hello from compliance-dashboard';
}

if (require.main === module) {
  console.log(main());
}
