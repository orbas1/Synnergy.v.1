use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct DidDocument {
    pub controller: String,
    pub metadata: String,
    pub revoked: bool,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct DIDRegistry {
    documents: HashMap<String, DidDocument>,
}

impl DIDRegistry {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn register(&mut self, did: String, controller: String, metadata: String) -> Result<(), &'static str> {
        if self.documents.contains_key(&did) {
            return Err("exists");
        }
        self.documents.insert(
            did,
            DidDocument {
                controller,
                metadata,
                revoked: false,
            },
        );
        Ok(())
    }

    pub fn update_metadata(&mut self, did: &str, caller: &str, metadata: String) -> Result<(), &'static str> {
        let doc = self.documents.get_mut(did).ok_or("missing")?;
        if doc.controller != caller {
            return Err("forbidden");
        }
        if doc.revoked {
            return Err("revoked");
        }
        doc.metadata = metadata;
        Ok(())
    }

    pub fn transfer_control(&mut self, did: &str, caller: &str, new_controller: String) -> Result<(), &'static str> {
        let doc = self.documents.get_mut(did).ok_or("missing")?;
        if doc.controller != caller {
            return Err("forbidden");
        }
        doc.controller = new_controller;
        Ok(())
    }

    pub fn revoke(&mut self, did: &str, caller: &str) -> Result<(), &'static str> {
        let doc = self.documents.get_mut(did).ok_or("missing")?;
        if doc.controller != caller {
            return Err("forbidden");
        }
        doc.revoked = true;
        Ok(())
    }

    pub fn get(&self, did: &str) -> Option<&DidDocument> {
        self.documents.get(did)
    }
}
