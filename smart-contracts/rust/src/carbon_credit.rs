use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct CarbonCredit {
    pub authority: String,
    balances: HashMap<String, u128>,
    pub total_issued: u128,
    pub total_retired: u128,
}

impl CarbonCredit {
    pub fn new(authority: String) -> Self {
        Self {
            authority,
            balances: HashMap::new(),
            total_issued: 0,
            total_retired: 0,
        }
    }

    pub fn balance_of(&self, account: &str) -> u128 {
        *self.balances.get(account).unwrap_or(&0)
    }

    pub fn issue(&mut self, caller: &str, to: String, amount: u128) -> Result<(), &'static str> {
        if caller != self.authority {
            return Err("only authority");
        }
        let entry = self.balances.entry(to).or_default();
        *entry = entry.saturating_add(amount);
        self.total_issued = self.total_issued.saturating_add(amount);
        Ok(())
    }

    pub fn transfer(&mut self, from: &str, to: String, amount: u128) -> Result<(), &'static str> {
        let from_balance = self.balances.get_mut(from).ok_or("insufficient")?;
        if *from_balance < amount {
            return Err("insufficient");
        }
        *from_balance -= amount;
        let entry = self.balances.entry(to).or_default();
        *entry = entry.saturating_add(amount);
        Ok(())
    }

    pub fn retire(&mut self, caller: &str, amount: u128) -> Result<(), &'static str> {
        let balance = self.balances.get_mut(caller).ok_or("insufficient")?;
        if *balance < amount {
            return Err("insufficient");
        }
        *balance -= amount;
        self.total_retired = self.total_retired.saturating_add(amount);
        Ok(())
    }
}
