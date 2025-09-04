use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct RevenueShareToken {}

impl RevenueShareToken {
    pub fn new() -> Self {
        Self {}
    }
}
