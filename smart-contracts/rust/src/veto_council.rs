use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct VetoCouncil {}

impl VetoCouncil {
    pub fn new() -> Self {
        Self {}
    }
}
