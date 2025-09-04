use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct IoTOracle {}

impl IoTOracle {
    pub fn new() -> Self {
        Self {}
    }
}
