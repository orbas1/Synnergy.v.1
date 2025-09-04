use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct AssetProvenance {}

impl AssetProvenance {
    pub fn new() -> Self {
        Self {}
    }
}
