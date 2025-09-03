use std::collections::HashMap;
use serde::{Serialize, Deserialize};

/// Simple token faucet model that tracks last request time for each account.
#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct TokenFaucet {
    pub owner: String,
    pub drip_amount: u64,
    pub wait_seconds: u64,
    last_request: HashMap<String, u64>,
}

impl TokenFaucet {
    pub fn new(owner: String, drip_amount: u64, wait_seconds: u64) -> Self {
        Self { owner, drip_amount, wait_seconds, last_request: HashMap::new() }
    }

    /// Request tokens from the faucet. Returns amount dispensed on success.
    pub fn request(&mut self, caller: &str, now: u64) -> Result<u64, &'static str> {
        let last = self.last_request.get(caller).copied().unwrap_or(0);
        if now < last + self.wait_seconds { return Err("wait"); }
        self.last_request.insert(caller.to_string(), now);
        Ok(self.drip_amount)
    }
}
