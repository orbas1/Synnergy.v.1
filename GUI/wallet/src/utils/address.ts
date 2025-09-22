export function normalizeAddress(address: string): string {
  const trimmed = address.trim();

  if (!trimmed) {
    throw new Error('Wallet address cannot be empty.');
  }

  return trimmed.toLowerCase();
}
