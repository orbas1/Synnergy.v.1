use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct RBAC {}

impl RBAC {
    pub fn new() -> Self {
        Self {}
    }
}
