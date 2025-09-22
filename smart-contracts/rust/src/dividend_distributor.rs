use std::collections::HashMap;
use serde::{Deserialize, Serialize};

const PRECISION: u128 = 1_000_000_000_000;

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
struct Shareholder {
    shares: u128,
    accrued: u128,
    checkpoint: u128,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct DividendDistributor {
    shareholders: HashMap<String, Shareholder>,
    pub total_shares: u128,
    cumulative_per_share: u128,
    pub undistributed: u128,
}

impl DividendDistributor {
    pub fn new() -> Self {
        Self::default()
    }

    fn accrue(&mut self, account: &str) {
        if let Some(holder) = self.shareholders.get_mut(account) {
            let delta = self.cumulative_per_share.saturating_sub(holder.checkpoint);
            if delta > 0 && holder.shares > 0 {
                let owed = holder.shares.saturating_mul(delta) / PRECISION;
                holder.accrued = holder.accrued.saturating_add(owed);
                holder.checkpoint = self.cumulative_per_share;
            } else {
                holder.checkpoint = self.cumulative_per_share;
            }
        }
    }

    pub fn set_shares(&mut self, account: String, shares: u128) {
        if !self.shareholders.contains_key(&account) {
            self.shareholders.insert(
                account.clone(),
                Shareholder {
                    shares: 0,
                    accrued: 0,
                    checkpoint: self.cumulative_per_share,
                },
            );
        }
        self.accrue(&account);
        if let Some(holder) = self.shareholders.get_mut(&account) {
            if self.total_shares >= holder.shares {
                self.total_shares = self.total_shares - holder.shares + shares;
            } else {
                self.total_shares = shares;
            }
            holder.shares = shares;
            holder.checkpoint = self.cumulative_per_share;
        }
    }

    pub fn deposit_dividends(&mut self, amount: u128) {
        if amount == 0 {
            return;
        }
        if self.total_shares == 0 {
            self.undistributed = self.undistributed.saturating_add(amount);
            return;
        }
        let scaled = amount.saturating_mul(PRECISION);
        let per_share = scaled / self.total_shares;
        if per_share == 0 {
            self.undistributed = self.undistributed.saturating_add(amount);
            return;
        }
        self.cumulative_per_share = self.cumulative_per_share.saturating_add(per_share);
        let distributed = per_share.saturating_mul(self.total_shares) / PRECISION;
        if amount > distributed {
            self.undistributed = self.undistributed.saturating_add(amount - distributed);
        }
    }

    pub fn recycle_undistributed(&mut self) {
        let pending = self.undistributed;
        self.undistributed = 0;
        self.deposit_dividends(pending);
    }

    pub fn claim(&mut self, account: &str) -> Result<u128, &'static str> {
        let exists = self.shareholders.contains_key(account);
        if !exists {
            return Err("unknown shareholder");
        }
        self.accrue(account);
        let holder = self.shareholders.get_mut(account).expect("checked");
        let amount = holder.accrued;
        holder.accrued = 0;
        Ok(amount)
    }

    pub fn info(&self, account: &str) -> Option<(u128, u128)> {
        self.shareholders
            .get(account)
            .map(|h| (h.shares, h.accrued))
    }
}
