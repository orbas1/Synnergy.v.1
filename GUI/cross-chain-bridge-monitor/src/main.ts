export function main(): string {
  const apiUrl = process.env.API_URL || 'http://localhost:8080';
  return `Bridge monitor API: ${apiUrl}`;
}

if (require.main === module) {
  // eslint-disable-next-line no-console
  console.log(main());
}
