use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct ETFToken {
    pub manager: String,
    pub share_supply: u128,
    holdings_value: HashMap<String, u128>,
    cash_reserve: u128,
    pub balances: HashMap<String, u128>,
}

impl ETFToken {
    pub fn new(manager: String) -> Self {
        Self {
            manager,
            share_supply: 0,
            holdings_value: HashMap::new(),
            cash_reserve: 0,
            balances: HashMap::new(),
        }
    }

    pub fn deposit_asset(&mut self, caller: &str, asset: String, value: u128) -> Result<(), &'static str> {
        if caller != self.manager {
            return Err("only manager");
        }
        let entry = self.holdings_value.entry(asset).or_default();
        *entry = entry.saturating_add(value);
        self.cash_reserve = self.cash_reserve.saturating_add(value);
        Ok(())
    }

    pub fn update_asset_value(&mut self, caller: &str, asset: String, value: u128) -> Result<(), &'static str> {
        if caller != self.manager {
            return Err("only manager");
        }
        self.holdings_value.insert(asset, value);
        Ok(())
    }

    pub fn mint_shares(
        &mut self,
        caller: &str,
        to: String,
        shares: u128,
        collateral_value: u128,
    ) -> Result<(), &'static str> {
        if caller != self.manager {
            return Err("only manager");
        }
        if shares == 0 {
            return Err("invalid shares");
        }
        if collateral_value > self.cash_reserve {
            return Err("insufficient collateral");
        }
        self.cash_reserve -= collateral_value;
        self.share_supply = self.share_supply.saturating_add(shares);
        let entry = self.balances.entry(to).or_default();
        *entry = entry.saturating_add(shares);
        Ok(())
    }

    pub fn redeem_shares(&mut self, holder: &str, shares: u128) -> Result<u128, &'static str> {
        if shares == 0 {
            return Err("invalid shares");
        }
        if self.share_supply == 0 {
            return Err("no supply");
        }
        let nav = self.total_nav();
        let balance = self.balances.get_mut(holder).ok_or("no balance")?;
        if *balance < shares {
            return Err("insufficient");
        }
        let redemption_value = nav.saturating_mul(shares) / self.share_supply.max(1);
        *balance -= shares;
        self.share_supply -= shares;
        if redemption_value > self.cash_reserve {
            return Err("insufficient liquidity");
        }
        self.cash_reserve -= redemption_value;
        Ok(redemption_value)
    }

    pub fn total_nav(&self) -> u128 {
        self.cash_reserve
            .saturating_add(self.holdings_value.values().copied().sum::<u128>())
    }
}
