use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct PerpetualSwap {}

impl PerpetualSwap {
    pub fn new() -> Self {
        Self {}
    }
}
