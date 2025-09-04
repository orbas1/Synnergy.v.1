export interface NFT {
  id: number;
  title: string;
}

export async function fetchNFTs(): Promise<NFT[]> {
  // In a production system this would call an external API or smart contract.
  // Static data keeps the example deterministic for tests.
  return [
    { id: 1, title: 'Genesis Token' },
    { id: 2, title: 'Founders Badge' },
    { id: 3, title: 'Community Award' }
  ];
}
