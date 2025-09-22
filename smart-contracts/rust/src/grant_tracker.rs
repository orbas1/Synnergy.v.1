use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq)]
pub enum GrantStatus {
    Submitted,
    Approved,
    Rejected,
    Completed,
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Milestone {
    pub description: String,
    pub amount: u128,
    pub completed: bool,
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Grant {
    pub applicant: String,
    pub requested_amount: u128,
    pub approved_amount: u128,
    pub status: GrantStatus,
    pub milestones: Vec<Milestone>,
    pub disbursed: u128,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct GrantTracker {
    pub authority: String,
    pub grants: HashMap<u64, Grant>,
    pub next_id: u64,
}

impl GrantTracker {
    pub fn new(authority: String) -> Self {
        Self {
            authority,
            grants: HashMap::new(),
            next_id: 0,
        }
    }

    pub fn submit_grant(
        &mut self,
        applicant: String,
        requested_amount: u128,
        milestones: Vec<(String, u128)>,
    ) -> u64 {
        self.next_id += 1;
        let items = milestones
            .into_iter()
            .map(|(description, amount)| Milestone {
                description,
                amount,
                completed: false,
            })
            .collect();
        let grant = Grant {
            applicant,
            requested_amount,
            approved_amount: 0,
            status: GrantStatus::Submitted,
            milestones: items,
            disbursed: 0,
        };
        self.grants.insert(self.next_id, grant);
        self.next_id
    }

    pub fn approve(&mut self, caller: &str, id: u64, amount: u128) -> Result<(), &'static str> {
        if caller != self.authority {
            return Err("only authority");
        }
        let grant = self.grants.get_mut(&id).ok_or("not found")?;
        if grant.status != GrantStatus::Submitted {
            return Err("invalid state");
        }
        grant.status = GrantStatus::Approved;
        grant.approved_amount = amount;
        Ok(())
    }

    pub fn reject(&mut self, caller: &str, id: u64) -> Result<(), &'static str> {
        if caller != self.authority {
            return Err("only authority");
        }
        let grant = self.grants.get_mut(&id).ok_or("not found")?;
        grant.status = GrantStatus::Rejected;
        Ok(())
    }

    pub fn mark_milestone(&mut self, caller: &str, id: u64, index: usize) -> Result<(), &'static str> {
        if caller != self.authority {
            return Err("only authority");
        }
        let grant = self.grants.get_mut(&id).ok_or("not found")?;
        if let Some(milestone) = grant.milestones.get_mut(index) {
            milestone.completed = true;
            Ok(())
        } else {
            Err("invalid milestone")
        }
    }

    pub fn disburse(&mut self, caller: &str, id: u64, amount: u128) -> Result<(), &'static str> {
        if caller != self.authority {
            return Err("only authority");
        }
        let grant = self.grants.get_mut(&id).ok_or("not found")?;
        if grant.status != GrantStatus::Approved {
            return Err("not approved");
        }
        if grant.disbursed.saturating_add(amount) > grant.approved_amount {
            return Err("exceeds approval");
        }
        grant.disbursed = grant.disbursed.saturating_add(amount);
        if grant.disbursed == grant.approved_amount
            && grant.milestones.iter().all(|m| m.completed)
        {
            grant.status = GrantStatus::Completed;
        }
        Ok(())
    }
}
