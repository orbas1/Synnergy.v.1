use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct BudgetLine {
    pub allocated: u128,
    pub spent: u128,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct GovTreasuryBudget {
    pub authority: String,
    categories: HashMap<String, BudgetLine>,
}

impl GovTreasuryBudget {
    pub fn new(authority: String) -> Self {
        Self {
            authority,
            categories: HashMap::new(),
        }
    }

    pub fn allocate(&mut self, caller: &str, category: String, amount: u128) -> Result<(), &'static str> {
        if caller != self.authority {
            return Err("only authority");
        }
        let entry = self.categories.entry(category).or_default();
        entry.allocated = entry.allocated.saturating_add(amount);
        Ok(())
    }

    pub fn spend(&mut self, category: &str, amount: u128) -> Result<(), &'static str> {
        let entry = self.categories.get_mut(category).ok_or("unknown")?;
        if entry.remaining() < amount {
            return Err("insufficient allocation");
        }
        entry.spent = entry.spent.saturating_add(amount);
        Ok(())
    }

    pub fn remaining(&self, category: &str) -> Option<u128> {
        self.categories.get(category).map(|line| line.remaining())
    }
}

impl BudgetLine {
    pub fn remaining(&self) -> u128 {
        self.allocated.saturating_sub(self.spent)
    }
}
