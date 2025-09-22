use std::collections::{HashMap, HashSet};
use serde::{Deserialize, Serialize};

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct VetoCouncil {
    members: HashSet<String>,
    vetos: HashMap<u64, HashSet<String>>,
    pub quorum: usize,
}

impl VetoCouncil {
    pub fn new(quorum: usize) -> Self {
        Self {
            members: HashSet::new(),
            vetos: HashMap::new(),
            quorum,
        }
    }

    pub fn add_member(&mut self, member: String) {
        self.members.insert(member);
    }

    pub fn remove_member(&mut self, member: &str) {
        self.members.remove(member);
    }

    pub fn cast_veto(&mut self, proposal_id: u64, member: &str) -> Result<(), &'static str> {
        if !self.members.contains(member) {
            return Err("not member");
        }
        self.vetos
            .entry(proposal_id)
            .or_insert_with(HashSet::new)
            .insert(member.to_string());
        Ok(())
    }

    pub fn is_vetoed(&self, proposal_id: u64) -> bool {
        self.vetos
            .get(&proposal_id)
            .map(|set| set.len() >= self.quorum)
            .unwrap_or(false)
    }
}
