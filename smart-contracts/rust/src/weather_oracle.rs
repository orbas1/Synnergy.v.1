use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct WeatherReading {
    pub temperature_c: f64,
    pub humidity: f64,
    pub timestamp: u64,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct WeatherOracle {
    stations: HashMap<String, WeatherReading>,
}

impl WeatherOracle {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn submit_reading(&mut self, station: String, temperature_c: f64, humidity: f64, timestamp: u64) {
        let existing = self.stations.get(&station);
        if existing.map(|r| r.timestamp <= timestamp).unwrap_or(true) {
            self.stations.insert(
                station,
                WeatherReading {
                    temperature_c,
                    humidity,
                    timestamp,
                },
            );
        }
    }

    pub fn get_latest(&self, station: &str) -> Option<&WeatherReading> {
        self.stations.get(station)
    }
}
