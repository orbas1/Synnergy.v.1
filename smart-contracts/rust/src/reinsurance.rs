use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct ReinsuranceContract {
    pub cedent: String,
    pub reinsurer: String,
    pub premium: u128,
    pub limit: u128,
    pub claims_paid: u128,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct Reinsurance {
    pub contracts: HashMap<u64, ReinsuranceContract>,
    pub next_id: u64,
}

impl Reinsurance {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn create_contract(
        &mut self,
        cedent: String,
        reinsurer: String,
        premium: u128,
        limit: u128,
    ) -> u64 {
        self.next_id += 1;
        self.contracts.insert(
            self.next_id,
            ReinsuranceContract {
                cedent,
                reinsurer,
                premium,
                limit,
                claims_paid: 0,
            },
        );
        self.next_id
    }

    pub fn record_claim(&mut self, id: u64, amount: u128) -> Result<(), &'static str> {
        let contract = self.contracts.get_mut(&id).ok_or("not found")?;
        if contract.claims_paid.saturating_add(amount) > contract.limit {
            return Err("limit exceeded");
        }
        contract.claims_paid = contract.claims_paid.saturating_add(amount);
        Ok(())
    }
}
