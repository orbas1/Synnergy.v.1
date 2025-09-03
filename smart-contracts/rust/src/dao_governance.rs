use std::collections::{HashMap, HashSet};
use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Proposal {
    pub description: String,
    pub votes_for: u64,
    pub votes_against: u64,
    pub deadline: u64,
    pub executed: bool,
    voters: HashSet<String>,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct DaoGovernance {
    pub proposals: HashMap<u64, Proposal>,
    pub next_id: u64,
}

impl DaoGovernance {
    pub fn create_proposal(&mut self, description: String, now: u64, duration: u64) -> u64 {
        self.next_id += 1;
        let p = Proposal {
            description,
            votes_for: 0,
            votes_against: 0,
            deadline: now + duration,
            executed: false,
            voters: HashSet::new(),
        };
        self.proposals.insert(self.next_id, p);
        self.next_id
    }

    pub fn vote(&mut self, id: u64, voter: String, support: bool, now: u64) -> Result<(), &'static str> {
        let p = self.proposals.get_mut(&id).ok_or("not found")?;
        if now >= p.deadline { return Err("ended"); }
        if !p.voters.insert(voter) { return Err("voted"); }
        if support { p.votes_for += 1; } else { p.votes_against += 1; }
        Ok(())
    }

    pub fn execute(&mut self, id: u64, now: u64) -> Result<bool, &'static str> {
        let p = self.proposals.get_mut(&id).ok_or("not found")?;
        if now < p.deadline { return Err("active"); }
        if p.executed { return Err("done"); }
        p.executed = true;
        Ok(p.votes_for > p.votes_against)
    }
}
