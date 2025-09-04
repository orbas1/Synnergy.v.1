use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct ParametricInsurance {}

impl ParametricInsurance {
    pub fn new() -> Self {
        Self {}
    }
}
