use std::collections::{HashMap, HashSet};
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct BridgeRequest {
    pub owner: String,
    pub token_id: u64,
    pub source_chain: String,
    pub target_chain: String,
    pub target_owner: Option<String>,
    pub completed: bool,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct NFTBridge {
    pub relayers: HashSet<String>,
    pub requests: HashMap<u64, BridgeRequest>,
    pub locked_tokens: HashSet<u64>,
    pub next_id: u64,
}

impl NFTBridge {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn add_relayer(&mut self, relayer: String) {
        self.relayers.insert(relayer);
    }

    pub fn lock_nft(
        &mut self,
        owner: String,
        token_id: u64,
        source_chain: String,
        target_chain: String,
    ) -> Result<u64, &'static str> {
        if self.locked_tokens.contains(&token_id) {
            return Err("locked");
        }
        self.locked_tokens.insert(token_id);
        self.next_id += 1;
        self.requests.insert(
            self.next_id,
            BridgeRequest {
                owner,
                token_id,
                source_chain,
                target_chain,
                target_owner: None,
                completed: false,
            },
        );
        Ok(self.next_id)
    }

    pub fn release(
        &mut self,
        caller: &str,
        request_id: u64,
        target_owner: String,
    ) -> Result<u64, &'static str> {
        if !self.relayers.contains(caller) {
            return Err("forbidden");
        }
        let request = self.requests.get_mut(&request_id).ok_or("not found")?;
        if request.completed {
            return Err("completed");
        }
        request.completed = true;
        request.target_owner = Some(target_owner);
        self.locked_tokens.remove(&request.token_id);
        Ok(request.token_id)
    }

    pub fn is_locked(&self, token_id: u64) -> bool {
        self.locked_tokens.contains(&token_id)
    }
}
