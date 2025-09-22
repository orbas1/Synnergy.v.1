use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Sensor {
    pub owner: String,
    pub last_reading: Option<f64>,
    pub last_timestamp: u64,
    pub active: bool,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct IoTOracle {
    sensors: HashMap<String, Sensor>,
}

impl IoTOracle {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn register_sensor(&mut self, sensor_id: String, owner: String) {
        self.sensors.insert(
            sensor_id,
            Sensor {
                owner,
                last_reading: None,
                last_timestamp: 0,
                active: true,
            },
        );
    }

    pub fn deactivate(&mut self, sensor_id: &str, caller: &str) -> Result<(), &'static str> {
        let sensor = self.sensors.get_mut(sensor_id).ok_or("unknown")?;
        if sensor.owner != caller {
            return Err("forbidden");
        }
        sensor.active = false;
        Ok(())
    }

    pub fn submit_reading(
        &mut self,
        sensor_id: &str,
        caller: &str,
        reading: f64,
        timestamp: u64,
    ) -> Result<(), &'static str> {
        let sensor = self.sensors.get_mut(sensor_id).ok_or("unknown")?;
        if sensor.owner != caller {
            return Err("forbidden");
        }
        if !sensor.active {
            return Err("inactive");
        }
        if timestamp < sensor.last_timestamp {
            return Err("stale");
        }
        sensor.last_reading = Some(reading);
        sensor.last_timestamp = timestamp;
        Ok(())
    }

    pub fn latest_reading(&self, sensor_id: &str) -> Option<(f64, u64)> {
        self.sensors.get(sensor_id).and_then(|sensor| {
            sensor
                .last_reading
                .map(|value| (value, sensor.last_timestamp))
        })
    }
}
