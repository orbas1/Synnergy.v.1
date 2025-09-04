import { fetchNFTs, NFT } from '../services/nftService';

export async function useNFTs(): Promise<NFT[]> {
  return fetchNFTs();
}
