use std::collections::HashMap;
use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Model {
    pub seller: String,
    pub price: u128,
    pub uri: String,
    pub sold: bool,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct AiModelMarket {
    pub models: HashMap<u64, Model>,
    pub next_id: u64,
}

impl AiModelMarket {
    pub fn list_model(&mut self, seller: String, price: u128, uri: String) -> u64 {
        self.next_id += 1;
        let m = Model { seller, price, uri, sold: false };
        self.models.insert(self.next_id, m);
        self.next_id
    }

    pub fn purchase(&mut self, id: u64, paid: u128) -> Result<(), &'static str> {
        let m = self.models.get_mut(&id).ok_or("not found")?;
        if m.sold { return Err("sold"); }
        if paid < m.price { return Err("price"); }
        m.sold = true;
        Ok(())
    }
}
