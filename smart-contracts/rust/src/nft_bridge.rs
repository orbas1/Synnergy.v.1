use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct NFTBridge {}

impl NFTBridge {
    pub fn new() -> Self {
        Self {}
    }
}
