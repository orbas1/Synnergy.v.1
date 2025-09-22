use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct ProviderSla {
    pub target_uptime: u8,
    pub penalties: u128,
    pub history: Vec<(u64, u8)>,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct StorageSLA {
    providers: HashMap<String, ProviderSla>,
    pub penalty_rate: u128,
}

impl StorageSLA {
    pub fn new(penalty_rate: u128) -> Self {
        Self {
            providers: HashMap::new(),
            penalty_rate,
        }
    }

    pub fn register_provider(&mut self, provider: String, target_uptime: u8) {
        self.providers.insert(
            provider,
            ProviderSla {
                target_uptime,
                penalties: 0,
                history: Vec::new(),
            },
        );
    }

    pub fn record_uptime(&mut self, provider: &str, timestamp: u64, uptime: u8) -> Result<(), &'static str> {
        let entry = self.providers.get_mut(provider).ok_or("unknown")?;
        entry.history.push((timestamp, uptime));
        if uptime < entry.target_uptime {
            entry.penalties = entry.penalties.saturating_add(self.penalty_rate);
        }
        Ok(())
    }

    pub fn penalties(&self, provider: &str) -> Option<u128> {
        self.providers.get(provider).map(|p| p.penalties)
    }

    pub fn average_uptime(&self, provider: &str) -> Option<f64> {
        self.providers.get(provider).and_then(|p| {
            if p.history.is_empty() {
                None
            } else {
                let sum: u64 = p.history.iter().map(|(_, uptime)| *uptime as u64).sum();
                Some(sum as f64 / p.history.len() as f64)
            }
        })
    }
}
