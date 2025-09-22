use std::collections::HashSet;
use serde::{Deserialize, Serialize};

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct SanctionsScreen {
    listed: HashSet<String>,
}

impl SanctionsScreen {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn add(&mut self, account: String) {
        self.listed.insert(account);
    }

    pub fn remove(&mut self, account: &str) {
        self.listed.remove(account);
    }

    pub fn is_listed(&self, account: &str) -> bool {
        self.listed.contains(account)
    }
}
