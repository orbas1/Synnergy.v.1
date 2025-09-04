use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct GovTreasuryBudget {}

impl GovTreasuryBudget {
    pub fn new() -> Self {
        Self {}
    }
}
