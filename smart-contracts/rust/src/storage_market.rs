use std::collections::HashMap;
use serde::{Serialize, Deserialize};

/// In-memory key/value storage market model requiring a fee per write.
#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct StorageMarket {
    pub owner: String,
    pub price: u128,
    pub data: HashMap<String, String>,
}

impl StorageMarket {
    pub fn new(owner: String, price: u128) -> Self {
        Self { owner, price, data: HashMap::new() }
    }

    pub fn save(&mut self, key: String, value: String, paid: u128) -> Result<(), &'static str> {
        if paid < self.price { return Err("fee"); }
        self.data.insert(key, value);
        Ok(())
    }

    pub fn retrieve(&self, key: &str) -> Option<&String> {
        self.data.get(key)
    }
}
