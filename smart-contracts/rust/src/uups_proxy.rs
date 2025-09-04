use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct UUPSProxy {}

impl UUPSProxy {
    pub fn new() -> Self {
        Self {}
    }
}
