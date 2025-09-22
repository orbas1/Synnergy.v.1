use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Channel {
    pub participants: [String; 2],
    pub balances: HashMap<String, u128>,
    pub version: u64,
    pub closing_time: Option<u64>,
    pub open: bool,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct RollupStateChannel {
    pub channels: HashMap<u64, Channel>,
    pub challenge_period: u64,
    pub next_id: u64,
}

impl RollupStateChannel {
    pub fn new(challenge_period: u64) -> Self {
        Self {
            channels: HashMap::new(),
            challenge_period,
            next_id: 0,
        }
    }

    pub fn open_channel(
        &mut self,
        party_a: String,
        party_b: String,
        deposit_a: u128,
        deposit_b: u128,
    ) -> u64 {
        self.next_id += 1;
        let mut balances = HashMap::new();
        balances.insert(party_a.clone(), deposit_a);
        balances.insert(party_b.clone(), deposit_b);
        self.channels.insert(
            self.next_id,
            Channel {
                participants: [party_a, party_b],
                balances,
                version: 0,
                closing_time: None,
                open: true,
            },
        );
        self.next_id
    }

    pub fn update_state(
        &mut self,
        id: u64,
        caller: &str,
        version: u64,
        balances: HashMap<String, u128>,
    ) -> Result<(), &'static str> {
        let channel = self.channels.get_mut(&id).ok_or("not found")?;
        if !channel.open {
            return Err("closed");
        }
        if !channel.participants.iter().any(|p| p == caller) {
            return Err("unauthorised");
        }
        if version <= channel.version {
            return Err("stale");
        }
        for participant in &channel.participants {
            if !balances.contains_key(participant) {
                return Err("invalid state");
            }
        }
        channel.version = version;
        channel.balances = balances;
        channel.closing_time = None;
        Ok(())
    }

    pub fn initiate_close(&mut self, id: u64, caller: &str, now: u64) -> Result<(), &'static str> {
        let channel = self.channels.get_mut(&id).ok_or("not found")?;
        if !channel.participants.iter().any(|p| p == caller) {
            return Err("unauthorised");
        }
        channel.closing_time = Some(now + self.challenge_period);
        Ok(())
    }

    pub fn finalize(&mut self, id: u64, now: u64) -> Result<HashMap<String, u128>, &'static str> {
        let channel = self.channels.get_mut(&id).ok_or("not found")?;
        if let Some(close_time) = channel.closing_time {
            if now < close_time {
                return Err("challenge active");
            }
        } else {
            return Err("not closing");
        }
        channel.open = false;
        Ok(channel.balances.clone())
    }
}
