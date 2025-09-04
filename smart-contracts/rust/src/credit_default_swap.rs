use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct CreditDefaultSwap {}

impl CreditDefaultSwap {
    pub fn new() -> Self {
        Self {}
    }
}
