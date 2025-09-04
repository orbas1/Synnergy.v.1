use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct DocumentNotary {}

impl DocumentNotary {
    pub fn new() -> Self {
        Self {}
    }
}
