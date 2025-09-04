import { fetchNFTs } from '../../src/services/nftService';

it('fetchNFTs returns static NFT data', async () => {
  const nfts = await fetchNFTs();
  expect(nfts).toHaveLength(3);
  expect(nfts[0].title).toBe('Genesis Token');
});
