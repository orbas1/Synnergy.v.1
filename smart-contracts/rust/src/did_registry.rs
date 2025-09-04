use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct DIDRegistry {}

impl DIDRegistry {
    pub fn new() -> Self {
        Self {}
    }
}
