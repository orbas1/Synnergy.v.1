use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct GovernedMintBurnToken {
    pub name: String,
    pub symbol: String,
    pub governance: String,
    pub total_supply: u128,
    balances: HashMap<String, u128>,
}

impl GovernedMintBurnToken {
    pub fn new(name: String, symbol: String, governance: String) -> Self {
        Self {
            name,
            symbol,
            governance,
            total_supply: 0,
            balances: HashMap::new(),
        }
    }

    pub fn balance_of(&self, owner: &str) -> u128 {
        *self.balances.get(owner).unwrap_or(&0)
    }

    pub fn mint(&mut self, caller: &str, to: String, amount: u128) -> Result<(), &'static str> {
        if caller != self.governance {
            return Err("only governance");
        }
        let entry = self.balances.entry(to).or_default();
        *entry = entry.saturating_add(amount);
        self.total_supply = self.total_supply.saturating_add(amount);
        Ok(())
    }

    pub fn burn(&mut self, caller: &str, from: &str, amount: u128) -> Result<(), &'static str> {
        if caller != self.governance {
            return Err("only governance");
        }
        let balance = self.balances.get_mut(from).ok_or("insufficient")?;
        if *balance < amount {
            return Err("insufficient");
        }
        *balance -= amount;
        self.total_supply = self.total_supply.saturating_sub(amount);
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

    pub fn update_governance(&mut self, caller: &str, new_governance: String) -> Result<(), &'static str> {
        if caller != self.governance {
            return Err("only governance");
        }
        self.governance = new_governance;
        Ok(())
    }
}
