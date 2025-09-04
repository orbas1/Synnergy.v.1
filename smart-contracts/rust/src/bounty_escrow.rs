use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct BountyEscrow {}

impl BountyEscrow {
    pub fn new() -> Self {
        Self {}
    }
}
