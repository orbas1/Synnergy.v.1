use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct CliffGrant {
    pub total_amount: u128,
    pub start: u64,
    pub cliff_seconds: u64,
    pub duration_seconds: u64,
    pub claimed: u128,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct EquityCliff {
    grants: HashMap<String, CliffGrant>,
}

impl EquityCliff {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn add_grant(
        &mut self,
        employee: String,
        total_amount: u128,
        start: u64,
        cliff_seconds: u64,
        duration_seconds: u64,
    ) {
        self.grants.insert(
            employee,
            CliffGrant {
                total_amount,
                start,
                cliff_seconds,
                duration_seconds,
                claimed: 0,
            },
        );
    }

    pub fn releasable_amount(&self, employee: &str, now: u64) -> Option<u128> {
        self.grants.get(employee).map(|grant| {
            if now < grant.start + grant.cliff_seconds {
                return 0;
            }
            let elapsed = now.saturating_sub(grant.start);
            let vested = if elapsed >= grant.duration_seconds {
                grant.total_amount
            } else {
                grant
                    .total_amount
                    .saturating_mul(elapsed as u128)
                    / grant.duration_seconds.max(1) as u128
            };
            vested.saturating_sub(grant.claimed)
        })
    }

    pub fn claim(&mut self, employee: &str, now: u64) -> Result<u128, &'static str> {
        let releasable = self.releasable_amount(employee, now).unwrap_or(0);
        if releasable == 0 {
            return Err("nothing vested");
        }
        let grant = self.grants.get_mut(employee).ok_or("no grant")?;
        grant.claimed = grant.claimed.saturating_add(releasable);
        Ok(releasable)
    }
}
