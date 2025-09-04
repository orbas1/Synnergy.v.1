use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct InvoiceFactoring {}

impl InvoiceFactoring {
    pub fn new() -> Self {
        Self {}
    }
}
