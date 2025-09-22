use std::collections::{HashMap, HashSet};
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Proposal {
    pub description: String,
    pub deadline: u64,
    pub token_yes: u128,
    pub token_no: u128,
    pub identity_yes: u64,
    pub identity_no: u64,
    voters: HashSet<String>,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct HybridVoting {
    token_balances: HashMap<String, u128>,
    identity_registry: HashSet<String>,
    proposals: HashMap<u64, Proposal>,
    pub next_id: u64,
}

impl HybridVoting {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn set_token_balance(&mut self, account: String, amount: u128) {
        if amount == 0 {
            self.token_balances.remove(&account);
        } else {
            self.token_balances.insert(account, amount);
        }
    }

    pub fn register_identity(&mut self, account: String) {
        self.identity_registry.insert(account);
    }

    pub fn create_proposal(&mut self, description: String, deadline: u64) -> u64 {
        self.next_id += 1;
        let proposal = Proposal {
            description,
            deadline,
            token_yes: 0,
            token_no: 0,
            identity_yes: 0,
            identity_no: 0,
            voters: HashSet::new(),
        };
        self.proposals.insert(self.next_id, proposal);
        self.next_id
    }

    pub fn vote(&mut self, caller: &str, id: u64, support: bool, now: u64) -> Result<(), &'static str> {
        let proposal = self.proposals.get_mut(&id).ok_or("not found")?;
        if now > proposal.deadline {
            return Err("closed");
        }
        if !proposal.voters.insert(caller.to_string()) {
            return Err("already voted");
        }
        let token_weight = self.token_balances.get(caller).copied().unwrap_or(0);
        let identity_weight = if self.identity_registry.contains(caller) { 1 } else { 0 };
        if support {
            proposal.token_yes = proposal.token_yes.saturating_add(token_weight);
            proposal.identity_yes = proposal.identity_yes.saturating_add(identity_weight);
        } else {
            proposal.token_no = proposal.token_no.saturating_add(token_weight);
            proposal.identity_no = proposal.identity_no.saturating_add(identity_weight);
        }
        Ok(())
    }

    pub fn hybrid_result(&self, id: u64) -> Option<(u128, u128)> {
        self.proposals.get(&id).map(|p| {
            let yes_score = p.token_yes.saturating_add(u128::from(p.identity_yes));
            let no_score = p.token_no.saturating_add(u128::from(p.identity_no));
            (yes_score, no_score)
        })
    }
}
