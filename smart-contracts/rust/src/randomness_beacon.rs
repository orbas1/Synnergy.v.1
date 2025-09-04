use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct RandomnessBeacon {}

impl RandomnessBeacon {
    pub fn new() -> Self {
        Self {}
    }
}
