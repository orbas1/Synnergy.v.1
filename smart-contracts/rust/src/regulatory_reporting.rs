use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Report {
    pub reporter: String,
    pub category: String,
    pub payload: String,
    pub timestamp: u64,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct RegulatoryReporting {
    pub reports: Vec<Report>,
}

impl RegulatoryReporting {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn submit(&mut self, reporter: String, category: String, payload: String, timestamp: u64) {
        self.reports.push(Report {
            reporter,
            category,
            payload,
            timestamp,
        });
    }

    pub fn reports_by(&self, reporter: &str) -> Vec<&Report> {
        self.reports
            .iter()
            .filter(|report| report.reporter == reporter)
            .collect()
    }

    pub fn latest(&self, category: &str) -> Option<&Report> {
        self.reports
            .iter()
            .rev()
            .find(|report| report.category == category)
    }
}
