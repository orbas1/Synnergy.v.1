use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct Reinsurance {}

impl Reinsurance {
    pub fn new() -> Self {
        Self {}
    }
}
