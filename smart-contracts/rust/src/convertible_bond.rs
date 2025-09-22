use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct BondPosition {
    pub principal: u128,
    pub accrued_interest: u128,
    pub converted: bool,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct ConvertibleBond {
    pub issuer: String,
    pub conversion_rate: u128,
    pub maturity: u64,
    positions: HashMap<String, BondPosition>,
}

impl ConvertibleBond {
    pub fn new(issuer: String, conversion_rate: u128, maturity: u64) -> Self {
        Self {
            issuer,
            conversion_rate,
            maturity,
            positions: HashMap::new(),
        }
    }

    pub fn issue(&mut self, caller: &str, holder: String, principal: u128) -> Result<(), &'static str> {
        if caller != self.issuer {
            return Err("only issuer");
        }
        let entry = self.positions.entry(holder).or_default();
        entry.principal = entry.principal.saturating_add(principal);
        Ok(())
    }

    pub fn accrue_interest(&mut self, holder: &str, interest: u128) -> Result<(), &'static str> {
        let pos = self.positions.get_mut(holder).ok_or("no position")?;
        if pos.converted {
            return Err("already converted");
        }
        pos.accrued_interest = pos.accrued_interest.saturating_add(interest);
        Ok(())
    }

    pub fn convert(&mut self, holder: &str, now: u64) -> Result<u128, &'static str> {
        let pos = self.positions.get_mut(holder).ok_or("no position")?;
        if now < self.maturity {
            return Err("not matured");
        }
        if pos.converted {
            return Err("already converted");
        }
        pos.converted = true;
        let total_value = pos.principal.saturating_add(pos.accrued_interest);
        Ok(total_value.saturating_mul(self.conversion_rate))
    }

    pub fn balance_of(&self, holder: &str) -> Option<&BondPosition> {
        self.positions.get(holder)
    }
}
