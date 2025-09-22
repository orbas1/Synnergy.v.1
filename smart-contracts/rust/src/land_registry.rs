use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct ParcelRecord {
    pub owner: String,
    pub metadata: String,
    pub encumbrances: Vec<String>,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct LandRegistry {
    parcels: HashMap<String, ParcelRecord>,
}

impl LandRegistry {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn register_parcel(&mut self, parcel_id: String, owner: String, metadata: String) -> Result<(), &'static str> {
        if self.parcels.contains_key(&parcel_id) {
            return Err("exists");
        }
        self.parcels.insert(
            parcel_id,
            ParcelRecord {
                owner,
                metadata,
                encumbrances: Vec::new(),
            },
        );
        Ok(())
    }

    pub fn transfer(&mut self, parcel_id: &str, caller: &str, new_owner: String) -> Result<(), &'static str> {
        let parcel = self.parcels.get_mut(parcel_id).ok_or("unknown")?;
        if parcel.owner != caller {
            return Err("forbidden");
        }
        parcel.owner = new_owner;
        Ok(())
    }

    pub fn add_encumbrance(&mut self, parcel_id: &str, note: String) -> Result<(), &'static str> {
        let parcel = self.parcels.get_mut(parcel_id).ok_or("unknown")?;
        parcel.encumbrances.push(note);
        Ok(())
    }

    pub fn info(&self, parcel_id: &str) -> Option<&ParcelRecord> {
        self.parcels.get(parcel_id)
    }
}
