use serde::{Serialize, Deserialize};

/// Lock and release model representing a simplistic cross-chain bridge.
#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct CrossChainBridge {
    pub admin: String,
    pub balances: std::collections::HashMap<String, u128>,
}

impl CrossChainBridge {
    pub fn new(admin: String) -> Self {
        Self { admin, balances: Default::default() }
    }

    pub fn lock(&mut self, from: String, amount: u128) {
        let bal = self.balances.entry(from).or_default();
        *bal += amount;
    }

    pub fn release(&mut self, caller: &str, to: String, amount: u128) -> Result<(), &'static str> {
        if caller != self.admin { return Err("admin"); }
        let bal = self.balances.entry(to).or_default();
        if *bal < amount { return Err("insufficient"); }
        *bal -= amount;
        Ok(())
    }
}
