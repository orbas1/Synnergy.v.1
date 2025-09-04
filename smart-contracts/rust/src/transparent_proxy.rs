use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct TransparentProxy {}

impl TransparentProxy {
    pub fn new() -> Self {
        Self {}
    }
}
