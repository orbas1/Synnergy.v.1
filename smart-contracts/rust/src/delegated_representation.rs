use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct DelegatedRepresentation {
    delegations: HashMap<String, String>,
}

impl DelegatedRepresentation {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn delegate(&mut self, delegator: String, delegate: String) {
        self.delegations.insert(delegator, delegate);
    }

    pub fn revoke(&mut self, delegator: &str) {
        self.delegations.remove(delegator);
    }

    pub fn representative_of(&self, delegator: &str) -> Option<&String> {
        self.delegations.get(delegator)
    }

    pub fn delegators_for(&self, delegate: &str) -> Vec<String> {
        self.delegations
            .iter()
            .filter_map(|(delegator, d)| if d == delegate { Some(delegator.clone()) } else { None })
            .collect()
    }
}
