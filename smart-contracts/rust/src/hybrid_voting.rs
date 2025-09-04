use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct HybridVoting {}

impl HybridVoting {
    pub fn new() -> Self {
        Self {}
    }
}
