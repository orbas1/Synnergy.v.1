use serde::{Deserialize, Serialize};

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct RandomnessBeacon {
    pub seed: u64,
    pub round: u64,
    pub history: Vec<u64>,
}

impl RandomnessBeacon {
    pub fn new(seed: u64) -> Self {
        Self {
            seed,
            round: 0,
            history: vec![seed],
        }
    }

    pub fn contribute_entropy(&mut self, entropy: u64) {
        self.seed ^= entropy.rotate_left(13);
        self.seed = self.seed.wrapping_mul(0x9e3779b97f4a7c15);
    }

    pub fn next_random(&mut self) -> u64 {
        self.round += 1;
        let mut x = self.seed ^ (self.seed << 7);
        x ^= x >> 9;
        x ^= x << 8;
        self.seed = x;
        self.history.push(x);
        x
    }
}
