use std::collections::HashMap;
use std::convert::TryFrom;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Position {
    pub size: i128,
    pub margin: u128,
    pub entry_price: u128,
    pub funding_paid: i128,
    pub open: bool,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct PerpetualSwap {
    pub positions: HashMap<String, Position>,
    pub funding_rate_bps: i32,
}

impl PerpetualSwap {
    pub fn new() -> Self {
        Self {
            positions: HashMap::new(),
            funding_rate_bps: 0,
        }
    }

    pub fn open_position(
        &mut self,
        trader: String,
        size: i128,
        entry_price: u128,
        margin: u128,
    ) -> Result<(), &'static str> {
        if size == 0 {
            return Err("invalid size");
        }
        self.positions.insert(
            trader,
            Position {
                size,
                margin,
                entry_price,
                funding_paid: 0,
                open: true,
            },
        );
        Ok(())
    }

    pub fn apply_funding(&mut self, rate_bps: i32) {
        self.funding_rate_bps = rate_bps;
        for position in self.positions.values_mut() {
            if position.open {
                let rate = i128::from(rate_bps);
                position.funding_paid += position.size * rate / 10_000;
            }
        }
    }

    pub fn close_position(&mut self, trader: &str, exit_price: u128) -> Result<i128, &'static str> {
        let position = self.positions.get_mut(trader).ok_or("no position")?;
        if !position.open {
            return Err("closed");
        }
        let exit_price_i = i128::try_from(exit_price).map_err(|_| "overflow")?;
        let entry_price_i = i128::try_from(position.entry_price).map_err(|_| "overflow")?;
        let price_diff = exit_price_i - entry_price_i;
        let pnl = price_diff * position.size - position.funding_paid;
        position.open = false;
        Ok(pnl)
    }
}
