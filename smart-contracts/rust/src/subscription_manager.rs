use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct SubscriptionManager {}

impl SubscriptionManager {
    pub fn new() -> Self {
        Self {}
    }
}
