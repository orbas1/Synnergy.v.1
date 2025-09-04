use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct EquityCliff {}

impl EquityCliff {
    pub fn new() -> Self {
        Self {}
    }
}
