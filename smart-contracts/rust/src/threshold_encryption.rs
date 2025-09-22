use std::collections::{HashMap, HashSet};
use serde::{Deserialize, Serialize};

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct ThresholdEncryption {
    pub threshold: usize,
    participants: HashMap<String, u64>,
    submitted: HashSet<String>,
}

impl ThresholdEncryption {
    pub fn new(threshold: usize) -> Self {
        Self {
            threshold,
            participants: HashMap::new(),
            submitted: HashSet::new(),
        }
    }

    pub fn register_participant(&mut self, id: String, weight: u64) {
        self.participants.insert(id, weight);
    }

    pub fn submit_share(&mut self, id: &str) -> Result<(), &'static str> {
        if !self.participants.contains_key(id) {
            return Err("unknown");
        }
        self.submitted.insert(id.to_string());
        Ok(())
    }

    pub fn can_decrypt(&self) -> bool {
        self.submitted.len() >= self.threshold
    }

    pub fn reset(&mut self) {
        self.submitted.clear();
    }
}
