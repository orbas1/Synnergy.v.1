use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct RoyaltySplitter {}

impl RoyaltySplitter {
    pub fn new() -> Self {
        Self {}
    }
}
