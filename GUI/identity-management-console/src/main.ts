export function main(): string {
  return 'Hello from identity-management-console';
}

if (require.main === module) {
  console.log(main());
}
