use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct DelegatedRepresentation {}

impl DelegatedRepresentation {
    pub fn new() -> Self {
        Self {}
    }
}
