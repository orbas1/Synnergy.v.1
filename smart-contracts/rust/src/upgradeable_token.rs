use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct UpgradeableToken {
    pub name: String,
    pub symbol: String,
    pub owner: String,
    pub total_supply: u128,
    pub version: u32,
    balances: HashMap<String, u128>,
    pub changelog: Vec<String>,
}

impl UpgradeableToken {
    pub fn new(name: String, symbol: String, owner: String) -> Self {
        Self {
            name,
            symbol,
            owner,
            total_supply: 0,
            version: 1,
            balances: HashMap::new(),
            changelog: vec!["initial deployment".to_string()],
        }
    }

    pub fn balance_of(&self, owner: &str) -> u128 {
        *self.balances.get(owner).unwrap_or(&0)
    }

    pub fn mint(&mut self, caller: &str, to: String, amount: u128) -> Result<(), &'static str> {
        if caller != self.owner {
            return Err("only owner");
        }
        let entry = self.balances.entry(to).or_default();
        *entry = entry.saturating_add(amount);
        self.total_supply = self.total_supply.saturating_add(amount);
        Ok(())
    }

    pub fn transfer(&mut self, from: &str, to: String, amount: u128) -> Result<(), &'static str> {
        let balance = self.balances.get_mut(from).ok_or("insufficient")?;
        if *balance < amount {
            return Err("insufficient");
        }
        *balance -= amount;
        let entry = self.balances.entry(to).or_default();
        *entry = entry.saturating_add(amount);
        Ok(())
    }

    pub fn upgrade(&mut self, caller: &str, notes: String) -> Result<(), &'static str> {
        if caller != self.owner {
            return Err("only owner");
        }
        self.version += 1;
        self.changelog.push(notes);
        Ok(())
    }

    pub fn change_owner(&mut self, caller: &str, new_owner: String) -> Result<(), &'static str> {
        if caller != self.owner {
            return Err("only owner");
        }
        self.owner = new_owner;
        Ok(())
    }
}
