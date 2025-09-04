use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct TokenVesting {}

impl TokenVesting {
    pub fn new() -> Self {
        Self {}
    }
}
