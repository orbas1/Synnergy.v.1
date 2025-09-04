use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct QuadraticFunding {}

impl QuadraticFunding {
    pub fn new() -> Self {
        Self {}
    }
}
