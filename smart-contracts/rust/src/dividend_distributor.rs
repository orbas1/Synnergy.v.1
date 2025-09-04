use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct DividendDistributor {}

impl DividendDistributor {
    pub fn new() -> Self {
        Self {}
    }
}
