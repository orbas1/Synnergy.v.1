import { NFT } from '../services/nftService';

export function renderNFTs(nfts: NFT[]): string {
  return nfts.map((nft) => `#${nft.id} - ${nft.title}`).join('\n');
}
