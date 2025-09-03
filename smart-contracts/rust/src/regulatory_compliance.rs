use std::collections::HashSet;
use serde::{Serialize, Deserialize};

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct RegulatoryCompliance {
    pub regulator: String,
    approved: HashSet<String>,
}

impl RegulatoryCompliance {
    pub fn new(regulator: String) -> Self {
        Self { regulator, approved: HashSet::new() }
    }

    pub fn set_approval(&mut self, caller: &str, user: String, status: bool) -> Result<(), &'static str> {
        if caller != self.regulator { return Err("only regulator"); }
        if status { self.approved.insert(user); } else { self.approved.remove(&user); }
        Ok(())
    }

    pub fn is_approved(&self, user: &str) -> bool {
        self.approved.contains(user)
    }
}
