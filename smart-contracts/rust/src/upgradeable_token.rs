use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct UpgradeableToken {}

impl UpgradeableToken {
    pub fn new() -> Self {
        Self {}
    }
}
