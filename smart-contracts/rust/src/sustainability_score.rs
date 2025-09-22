use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct Metrics {
    pub emissions_intensity: u8,
    pub renewable_ratio: u8,
    pub waste_recycled: u8,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct SustainabilityScore {
    metrics: HashMap<String, Metrics>,
}

impl SustainabilityScore {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn update_metrics(
        &mut self,
        entity: String,
        emissions_intensity: u8,
        renewable_ratio: u8,
        waste_recycled: u8,
    ) {
        self.metrics.insert(
            entity,
            Metrics {
                emissions_intensity,
                renewable_ratio,
                waste_recycled,
            },
        );
    }

    pub fn score(&self, entity: &str) -> Option<u8> {
        self.metrics.get(entity).map(|m| {
            let emissions_component = 100u32.saturating_sub(m.emissions_intensity as u32);
            let renewable_component = m.renewable_ratio as u32;
            let waste_component = m.waste_recycled as u32;
            let weighted = emissions_component * 4 + renewable_component * 4 + waste_component * 2;
            (weighted / 10).min(100) as u8
        })
    }
}
