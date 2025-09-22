use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct ProvenanceEvent {
    pub timestamp: u64,
    pub actor: String,
    pub action: String,
    pub reference: Option<String>,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct AssetProvenance {
    records: HashMap<String, Vec<ProvenanceEvent>>,
}

impl AssetProvenance {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn register_asset(
        &mut self,
        asset_id: String,
        timestamp: u64,
        actor: String,
        reference: Option<String>,
    ) -> Result<(), &'static str> {
        if self.records.contains_key(&asset_id) {
            return Err("asset exists");
        }
        let event = ProvenanceEvent {
            timestamp,
            actor,
            action: "registered".to_string(),
            reference,
        };
        self.records.insert(asset_id, vec![event]);
        Ok(())
    }

    pub fn add_event(
        &mut self,
        asset_id: &str,
        action: String,
        actor: String,
        timestamp: u64,
        reference: Option<String>,
    ) -> Result<(), &'static str> {
        let events = self.records.get_mut(asset_id).ok_or("missing asset")?;
        events.push(ProvenanceEvent {
            timestamp,
            actor,
            action,
            reference,
        });
        Ok(())
    }

    pub fn history(&self, asset_id: &str) -> Option<&[ProvenanceEvent]> {
        self.records.get(asset_id).map(|events| events.as_slice())
    }
}
