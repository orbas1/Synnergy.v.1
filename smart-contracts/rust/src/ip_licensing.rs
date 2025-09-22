use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct License {
    pub licensee: String,
    pub terms: String,
    pub active: bool,
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct IPAsset {
    pub owner: String,
    pub metadata: String,
    pub licenses: Vec<License>,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct IPLicensing {
    assets: HashMap<String, IPAsset>,
}

impl IPLicensing {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn register_asset(&mut self, asset_id: String, owner: String, metadata: String) {
        self.assets.insert(
            asset_id,
            IPAsset {
                owner,
                metadata,
                licenses: Vec::new(),
            },
        );
    }

    pub fn transfer(&mut self, asset_id: &str, caller: &str, new_owner: String) -> Result<(), &'static str> {
        let asset = self.assets.get_mut(asset_id).ok_or("unknown")?;
        if asset.owner != caller {
            return Err("forbidden");
        }
        asset.owner = new_owner;
        Ok(())
    }

    pub fn grant_license(
        &mut self,
        asset_id: &str,
        caller: &str,
        licensee: String,
        terms: String,
    ) -> Result<(), &'static str> {
        let asset = self.assets.get_mut(asset_id).ok_or("unknown")?;
        if asset.owner != caller {
            return Err("forbidden");
        }
        asset.licenses.push(License {
            licensee,
            terms,
            active: true,
        });
        Ok(())
    }

    pub fn revoke_license(&mut self, asset_id: &str, caller: &str, licensee: &str) -> Result<(), &'static str> {
        let asset = self.assets.get_mut(asset_id).ok_or("unknown")?;
        if asset.owner != caller {
            return Err("forbidden");
        }
        if let Some(license) = asset
            .licenses
            .iter_mut()
            .find(|lic| lic.licensee == licensee)
        {
            license.active = false;
            Ok(())
        } else {
            Err("not found")
        }
    }

    pub fn active_licenses(&self, asset_id: &str) -> Option<Vec<License>> {
        self.assets.get(asset_id).map(|asset| {
            asset
                .licenses
                .iter()
                .filter(|lic| lic.active)
                .cloned()
                .collect()
        })
    }
}
