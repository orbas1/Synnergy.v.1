use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct ConvertibleBond {}

impl ConvertibleBond {
    pub fn new() -> Self {
        Self {}
    }
}
