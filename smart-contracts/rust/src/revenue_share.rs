use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug, Default)]
pub struct Participant {
    pub share_bps: u32,
    pub accrued: u128,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct RevenueShare {
    pub total_bps: u32,
    participants: HashMap<String, Participant>,
}

impl RevenueShare {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn set_participant(&mut self, account: String, share_bps: u32) {
        self.total_bps = self.total_bps.saturating_sub(self.participants.get(&account).map(|p| p.share_bps).unwrap_or(0));
        if share_bps == 0 {
            self.participants.remove(&account);
        } else {
            self.participants.insert(
                account,
                Participant {
                    share_bps,
                    accrued: 0,
                },
            );
            self.total_bps = self.total_bps.saturating_add(share_bps);
        }
    }

    pub fn distribute(&mut self, amount: u128) {
        if amount == 0 || self.total_bps == 0 {
            return;
        }
        for participant in self.participants.values_mut() {
            let share = amount.saturating_mul(participant.share_bps as u128) / self.total_bps as u128;
            participant.accrued = participant.accrued.saturating_add(share);
        }
    }

    pub fn claim(&mut self, account: &str) -> Result<u128, &'static str> {
        let participant = self.participants.get_mut(account).ok_or("unknown")?;
        let amount = participant.accrued;
        participant.accrued = 0;
        Ok(amount)
    }
}
