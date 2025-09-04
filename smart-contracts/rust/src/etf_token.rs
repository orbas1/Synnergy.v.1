use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct ETFToken {}

impl ETFToken {
    pub fn new() -> Self {
        Self {}
    }
}
