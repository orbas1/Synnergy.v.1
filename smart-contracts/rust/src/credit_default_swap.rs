use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct CreditDefaultSwap {
    pub buyer: String,
    pub seller: String,
    pub notional: u128,
    pub premium: u128,
    pub premiums_paid: u64,
    pub default_triggered: bool,
    pub payout_confirmed: bool,
}

impl CreditDefaultSwap {
    pub fn new(buyer: String, seller: String, notional: u128, premium: u128) -> Self {
        Self {
            buyer,
            seller,
            notional,
            premium,
            premiums_paid: 0,
            default_triggered: false,
            payout_confirmed: false,
        }
    }

    pub fn pay_premium(&mut self, caller: &str) -> Result<u128, &'static str> {
        if caller != self.buyer {
            return Err("only buyer");
        }
        if self.default_triggered {
            return Err("default triggered");
        }
        self.premiums_paid += 1;
        Ok(self.premium)
    }

    pub fn trigger_default(&mut self, caller: &str) -> Result<u128, &'static str> {
        if caller != self.buyer {
            return Err("only buyer");
        }
        if self.default_triggered {
            return Err("already triggered");
        }
        self.default_triggered = true;
        Ok(self.notional)
    }

    pub fn confirm_payout(&mut self, caller: &str) -> Result<(), &'static str> {
        if caller != self.seller {
            return Err("only seller");
        }
        if !self.default_triggered {
            return Err("no default");
        }
        if self.payout_confirmed {
            return Err("already confirmed");
        }
        self.payout_confirmed = true;
        Ok(())
    }
}
