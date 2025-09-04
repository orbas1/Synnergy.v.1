use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct ZKTransaction {}

impl ZKTransaction {
    pub fn new() -> Self {
        Self {}
    }
}
