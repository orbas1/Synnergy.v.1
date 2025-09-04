use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct RevenueShare {}

impl RevenueShare {
    pub fn new() -> Self {
        Self {}
    }
}
