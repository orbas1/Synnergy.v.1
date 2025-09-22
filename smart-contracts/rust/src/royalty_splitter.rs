use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct Payee {
    pub share: u32,
    pub accrued: u128,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct RoyaltySplitter {
    payees: HashMap<String, Payee>,
    pub total_share: u32,
}

impl RoyaltySplitter {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn set_payee(&mut self, account: String, share: u32) {
        self.total_share = self
            .total_share
            .saturating_sub(self.payees.get(&account).map(|p| p.share).unwrap_or(0));
        if share == 0 {
            self.payees.remove(&account);
        } else {
            self.payees.insert(
                account,
                Payee {
                    share,
                    accrued: 0,
                },
            );
            self.total_share = self.total_share.saturating_add(share);
        }
    }

    pub fn distribute(&mut self, amount: u128) {
        if amount == 0 || self.total_share == 0 {
            return;
        }
        for payee in self.payees.values_mut() {
            let payout = amount.saturating_mul(payee.share as u128) / self.total_share as u128;
            payee.accrued = payee.accrued.saturating_add(payout);
        }
    }

    pub fn claim(&mut self, account: &str) -> Result<u128, &'static str> {
        let payee = self.payees.get_mut(account).ok_or("unknown")?;
        let amount = payee.accrued;
        payee.accrued = 0;
        Ok(amount)
    }
}
