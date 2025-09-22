use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct DataRecord {
    pub owner: String,
    pub value: String,
    pub consent: bool,
    pub deleted: bool,
    pub last_updated: u64,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct GDPRCompliantStorage {
    records: HashMap<String, DataRecord>,
}

impl GDPRCompliantStorage {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn store(
        &mut self,
        key: String,
        owner: String,
        value: String,
        consent: bool,
        timestamp: u64,
    ) {
        self.records.insert(
            key,
            DataRecord {
                owner,
                value,
                consent,
                deleted: false,
                last_updated: timestamp,
            },
        );
    }

    pub fn update_value(
        &mut self,
        key: &str,
        caller: &str,
        value: String,
        consent: bool,
        timestamp: u64,
    ) -> Result<(), &'static str> {
        let record = self.records.get_mut(key).ok_or("missing")?;
        if record.owner != caller || record.deleted {
            return Err("forbidden");
        }
        record.value = value;
        record.consent = consent;
        record.last_updated = timestamp;
        Ok(())
    }

    pub fn revoke_consent(&mut self, key: &str, caller: &str) -> Result<(), &'static str> {
        let record = self.records.get_mut(key).ok_or("missing")?;
        if record.owner != caller {
            return Err("forbidden");
        }
        record.consent = false;
        Ok(())
    }

    pub fn erase(&mut self, key: &str, caller: &str) -> Result<(), &'static str> {
        let record = self.records.get_mut(key).ok_or("missing")?;
        if record.owner != caller {
            return Err("forbidden");
        }
        record.deleted = true;
        record.value.clear();
        Ok(())
    }

    pub fn read(&self, key: &str) -> Option<&DataRecord> {
        self.records
            .get(key)
            .filter(|record| record.consent && !record.deleted)
    }
}
