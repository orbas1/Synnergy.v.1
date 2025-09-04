use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct Options {}

impl Options {
    pub fn new() -> Self {
        Self {}
    }
}
