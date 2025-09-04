use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct GDPRCompliantStorage {}

impl GDPRCompliantStorage {
    pub fn new() -> Self {
        Self {}
    }
}
