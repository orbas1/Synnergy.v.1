use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct IPLicensing {}

impl IPLicensing {
    pub fn new() -> Self {
        Self {}
    }
}
