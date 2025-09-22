use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct Project {
    pub name: String,
    pub contributions: HashMap<String, u128>,
    pub matching_awarded: u128,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct QuadraticFunding {
    pub projects: HashMap<String, Project>,
    pub matching_pool: u128,
}

impl QuadraticFunding {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn add_project(&mut self, id: String, name: String) {
        self.projects.insert(
            id,
            Project {
                name,
                contributions: HashMap::new(),
                matching_awarded: 0,
            },
        );
    }

    pub fn contribute(&mut self, id: &str, contributor: String, amount: u128) -> Result<(), &'static str> {
        let project = self.projects.get_mut(id).ok_or("unknown")?;
        let entry = project.contributions.entry(contributor).or_default();
        *entry = entry.saturating_add(amount);
        Ok(())
    }

    pub fn add_matching_funds(&mut self, amount: u128) {
        self.matching_pool = self.matching_pool.saturating_add(amount);
    }

    pub fn finalize(&mut self, id: &str) -> Result<u128, &'static str> {
        let project = self.projects.get_mut(id).ok_or("unknown")?;
        let mut sum_sqrt = 0u128;
        let mut total = 0u128;
        for amount in project.contributions.values() {
            sum_sqrt = sum_sqrt.saturating_add(integer_sqrt(*amount));
            total = total.saturating_add(*amount);
        }
        let potential = sum_sqrt.saturating_mul(sum_sqrt);
        let match_amount = potential.saturating_sub(total);
        let awarded = match_amount.min(self.matching_pool);
        self.matching_pool = self.matching_pool.saturating_sub(awarded);
        project.matching_awarded = awarded;
        Ok(awarded)
    }
}

fn integer_sqrt(value: u128) -> u128 {
    if value == 0 {
        return 0;
    }
    let mut x0 = value;
    let mut x1 = (x0 + value / x0) / 2;
    while x1 < x0 {
        x0 = x1;
        x1 = (x0 + value / x0) / 2;
    }
    x0
}
