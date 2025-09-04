use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct DynamicPayroll {}

impl DynamicPayroll {
    pub fn new() -> Self {
        Self {}
    }
}
