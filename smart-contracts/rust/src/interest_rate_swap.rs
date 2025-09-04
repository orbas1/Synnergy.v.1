use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct InterestRateSwap {}

impl InterestRateSwap {
    pub fn new() -> Self {
        Self {}
    }
}
