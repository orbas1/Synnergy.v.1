use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct NotarizedDocument {
    pub owner: String,
    pub hash: String,
    pub timestamp: u64,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct DocumentNotary {
    records: HashMap<String, NotarizedDocument>,
}

impl DocumentNotary {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn notarize(&mut self, doc_id: String, owner: String, hash: String, timestamp: u64) -> Result<(), &'static str> {
        if self.records.contains_key(&doc_id) {
            return Err("already notarized");
        }
        self.records.insert(
            doc_id,
            NotarizedDocument {
                owner,
                hash,
                timestamp,
            },
        );
        Ok(())
    }

    pub fn verify(&self, doc_id: &str, hash: &str) -> bool {
        self.records
            .get(doc_id)
            .map(|doc| doc.hash == hash)
            .unwrap_or(false)
    }

    pub fn owner_of(&self, doc_id: &str) -> Option<&String> {
        self.records.get(doc_id).map(|doc| &doc.owner)
    }
}
