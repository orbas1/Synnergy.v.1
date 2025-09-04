import { useNFTs } from '../hooks/useNFTs';
import { nftStore } from '../state/store';
import { renderNFTs } from '../components/NFTList';

export async function showHomePage(): Promise<string> {
  const nfts = await useNFTs();
  nftStore.set(nfts);
  return renderNFTs(nftStore.get());
}
