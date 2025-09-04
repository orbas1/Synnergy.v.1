use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct SanctionsScreen {}

impl SanctionsScreen {
    pub fn new() -> Self {
        Self {}
    }
}
