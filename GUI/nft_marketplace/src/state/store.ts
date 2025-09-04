import { NFT } from '../services/nftService';

class NFTStore {
  private items: NFT[] = [];

  set(nfts: NFT[]): void {
    this.items = nfts;
  }

  get(): NFT[] {
    return this.items;
  }
}

export const nftStore = new NFTStore();
