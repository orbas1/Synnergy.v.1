use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Employee {
    pub rate_per_second: u128,
    pub last_update: u64,
    pub active: bool,
    pub accrued: u128,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct DynamicPayroll {
    pub employer: String,
    employees: HashMap<String, Employee>,
}

impl DynamicPayroll {
    pub fn new(employer: String) -> Self {
        Self {
            employer,
            employees: HashMap::new(),
        }
    }

    fn accrue(&mut self, employee: &str, now: u64) {
        if let Some(info) = self.employees.get_mut(employee) {
            if info.active && now > info.last_update {
                let elapsed = now - info.last_update;
                info.accrued = info
                    .accrued
                    .saturating_add(info.rate_per_second.saturating_mul(elapsed as u128));
                info.last_update = now;
            } else if now > info.last_update {
                info.last_update = now;
            }
        }
    }

    pub fn add_employee(
        &mut self,
        caller: &str,
        employee: String,
        rate_per_second: u128,
        start_time: u64,
    ) -> Result<(), &'static str> {
        if caller != self.employer {
            return Err("only employer");
        }
        self.employees.insert(
            employee,
            Employee {
                rate_per_second,
                last_update: start_time,
                active: true,
                accrued: 0,
            },
        );
        Ok(())
    }

    pub fn update_rate(
        &mut self,
        caller: &str,
        employee: &str,
        new_rate: u128,
        now: u64,
    ) -> Result<(), &'static str> {
        if caller != self.employer {
            return Err("only employer");
        }
        self.accrue(employee, now);
        let info = self.employees.get_mut(employee).ok_or("unknown")?;
        info.rate_per_second = new_rate;
        Ok(())
    }

    pub fn set_active(
        &mut self,
        caller: &str,
        employee: &str,
        active: bool,
        now: u64,
    ) -> Result<(), &'static str> {
        if caller != self.employer {
            return Err("only employer");
        }
        self.accrue(employee, now);
        let info = self.employees.get_mut(employee).ok_or("unknown")?;
        info.active = active;
        Ok(())
    }

    pub fn payout(&mut self, employee: &str, now: u64) -> Result<u128, &'static str> {
        self.accrue(employee, now);
        let info = self.employees.get_mut(employee).ok_or("unknown")?;
        let amount = info.accrued;
        info.accrued = 0;
        Ok(amount)
    }
}
