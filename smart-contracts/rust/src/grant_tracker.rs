use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct GrantTracker {}

impl GrantTracker {
    pub fn new() -> Self {
        Self {}
    }
}
