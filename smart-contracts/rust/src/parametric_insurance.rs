use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Policy {
    pub holder: String,
    pub premium: u128,
    pub trigger_value: i64,
    pub payout: u128,
    pub active: bool,
    pub triggered: bool,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct ParametricInsurance {
    pub insurer: String,
    policies: HashMap<u64, Policy>,
    pub next_id: u64,
}

impl ParametricInsurance {
    pub fn new(insurer: String) -> Self {
        Self {
            insurer,
            policies: HashMap::new(),
            next_id: 0,
        }
    }

    pub fn create_policy(
        &mut self,
        holder: String,
        premium: u128,
        trigger_value: i64,
        payout: u128,
    ) -> u64 {
        self.next_id += 1;
        self.policies.insert(
            self.next_id,
            Policy {
                holder,
                premium,
                trigger_value,
                payout,
                active: true,
                triggered: false,
            },
        );
        self.next_id
    }

    pub fn report_measurement(&mut self, id: u64, value: i64) -> Result<bool, &'static str> {
        let policy = self.policies.get_mut(&id).ok_or("not found")?;
        if !policy.active {
            return Err("inactive");
        }
        if value >= policy.trigger_value {
            policy.triggered = true;
            policy.active = false;
            Ok(true)
        } else {
            Ok(false)
        }
    }

    pub fn claim(&mut self, id: u64, caller: &str) -> Result<u128, &'static str> {
        let policy = self.policies.get_mut(&id).ok_or("not found")?;
        if policy.holder != caller {
            return Err("forbidden");
        }
        if !policy.triggered {
            return Err("not triggered");
        }
        policy.triggered = false;
        Ok(policy.payout)
    }
}
