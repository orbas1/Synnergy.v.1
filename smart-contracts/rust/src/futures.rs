use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct Futures {}

impl Futures {
    pub fn new() -> Self {
        Self {}
    }
}
