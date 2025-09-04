use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct RollupStateChannel {}

impl RollupStateChannel {
    pub fn new() -> Self {
        Self {}
    }
}
