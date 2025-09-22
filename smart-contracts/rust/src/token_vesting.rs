use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct VestingSchedule {
    pub total: u128,
    pub released: u128,
    pub start: u64,
    pub cliff: u64,
    pub duration: u64,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct TokenVesting {
    schedules: HashMap<String, VestingSchedule>,
}

impl TokenVesting {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn create_schedule(
        &mut self,
        beneficiary: String,
        total: u128,
        start: u64,
        cliff: u64,
        duration: u64,
    ) {
        self.schedules.insert(
            beneficiary,
            VestingSchedule {
                total,
                released: 0,
                start,
                cliff,
                duration,
            },
        );
    }

    pub fn releasable(&self, beneficiary: &str, now: u64) -> Option<u128> {
        self.schedules.get(beneficiary).map(|schedule| {
            if now < schedule.start + schedule.cliff {
                return 0;
            }
            let elapsed = now.saturating_sub(schedule.start);
            let vested = if elapsed >= schedule.duration {
                schedule.total
            } else {
                schedule.total.saturating_mul(elapsed as u128) / schedule.duration.max(1) as u128
            };
            vested.saturating_sub(schedule.released)
        })
    }

    pub fn claim(&mut self, beneficiary: &str, now: u64) -> Result<u128, &'static str> {
        let releasable = self.releasable(beneficiary, now).unwrap_or(0);
        if releasable == 0 {
            return Err("nothing vested");
        }
        let schedule = self.schedules.get_mut(beneficiary).ok_or("no schedule")?;
        schedule.released = schedule.released.saturating_add(releasable);
        Ok(releasable)
    }
}
