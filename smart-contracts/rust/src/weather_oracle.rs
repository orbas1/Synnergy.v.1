use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct WeatherOracle {}

impl WeatherOracle {
    pub fn new() -> Self {
        Self {}
    }
}
