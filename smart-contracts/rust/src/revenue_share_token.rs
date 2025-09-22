use std::collections::HashMap;
use serde::{Deserialize, Serialize};

const PRECISION: u128 = 1_000_000_000_000;

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
struct HolderInfo {
    balance: u128,
    reward_debt: u128,
    pending: u128,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct RevenueShareToken {
    pub name: String,
    pub symbol: String,
    pub owner: String,
    pub total_supply: u128,
    reward_per_token: u128,
    holders: HashMap<String, HolderInfo>,
}

impl RevenueShareToken {
    pub fn new(name: String, symbol: String, owner: String) -> Self {
        Self {
            name,
            symbol,
            owner,
            total_supply: 0,
            reward_per_token: 0,
            holders: HashMap::new(),
        }
    }

    fn ensure_holder(&mut self, account: &str) {
        self.holders.entry(account.to_string()).or_insert_with(|| HolderInfo {
            balance: 0,
            reward_debt: self.reward_per_token,
            pending: 0,
        });
    }

    fn sync_holder(&mut self, account: &str) {
        if let Some(info) = self.holders.get_mut(account) {
            let accrued = info
                .balance
                .saturating_mul(self.reward_per_token.saturating_sub(info.reward_debt))
                / PRECISION;
            info.pending = info.pending.saturating_add(accrued);
            info.reward_debt = self.reward_per_token;
        }
    }

    pub fn mint(&mut self, caller: &str, to: String, amount: u128) -> Result<(), &'static str> {
        if caller != self.owner {
            return Err("only owner");
        }
        self.total_supply = self.total_supply.saturating_add(amount);
        self.ensure_holder(&to);
        self.sync_holder(&to);
        let entry = self.holders.get_mut(&to).expect("exists");
        entry.balance = entry.balance.saturating_add(amount);
        entry.reward_debt = self.reward_per_token;
        Ok(())
    }

    pub fn burn(&mut self, caller: &str, from: &str, amount: u128) -> Result<(), &'static str> {
        if caller != self.owner {
            return Err("only owner");
        }
        self.sync_holder(from);
        let info = self.holders.get_mut(from).ok_or("no balance")?;
        if info.balance < amount {
            return Err("insufficient");
        }
        info.balance -= amount;
        info.reward_debt = self.reward_per_token;
        self.total_supply = self.total_supply.saturating_sub(amount);
        Ok(())
    }

    pub fn transfer(&mut self, from: &str, to: String, amount: u128) -> Result<(), &'static str> {
        self.sync_holder(from);
        {
            let from_info = self.holders.get_mut(from).ok_or("no balance")?;
            if from_info.balance < amount {
                return Err("insufficient");
            }
            from_info.balance -= amount;
            from_info.reward_debt = self.reward_per_token;
        }

        self.ensure_holder(&to);
        self.sync_holder(&to);
        let to_info = self.holders.get_mut(&to).expect("exists");
        to_info.balance = to_info.balance.saturating_add(amount);
        to_info.reward_debt = self.reward_per_token;
        Ok(())
    }

    pub fn distribute(&mut self, amount: u128) {
        if amount == 0 || self.total_supply == 0 {
            return;
        }
        let increment = amount.saturating_mul(PRECISION) / self.total_supply;
        if increment == 0 {
            return;
        }
        self.reward_per_token = self.reward_per_token.saturating_add(increment);
    }

    pub fn claim(&mut self, account: &str) -> Result<u128, &'static str> {
        self.sync_holder(account);
        let info = self.holders.get_mut(account).ok_or("no balance")?;
        let amount = info.pending;
        info.pending = 0;
        Ok(amount)
    }
}
