use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct GovernedMintBurnToken {}

impl GovernedMintBurnToken {
    pub fn new() -> Self {
        Self {}
    }
}
