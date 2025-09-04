use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct StorageSLA {}

impl StorageSLA {
    pub fn new() -> Self {
        Self {}
    }
}
