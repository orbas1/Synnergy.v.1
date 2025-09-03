use std::collections::HashMap;
use serde::{Serialize, Deserialize};

/// Minimal NFT minting model storing ownership mappings.
#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct NftMinting {
    pub name: String,
    pub symbol: String,
    pub token_count: u64,
    pub owner_of: HashMap<u64, String>,
    pub balance_of: HashMap<String, u64>,
}

impl NftMinting {
    pub fn new(name: String, symbol: String) -> Self {
        Self { name, symbol, ..Default::default() }
    }

    pub fn mint(&mut self, to: String) -> u64 {
        self.token_count += 1;
        self.owner_of.insert(self.token_count, to.clone());
        *self.balance_of.entry(to).or_insert(0) += 1;
        self.token_count
    }
}
