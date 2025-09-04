use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct LandRegistry {}

impl LandRegistry {
    pub fn new() -> Self {
        Self {}
    }
}
