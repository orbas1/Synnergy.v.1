use std::collections::HashMap;
use std::convert::TryFrom;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Swap {
    pub fixed_payer: String,
    pub floating_payer: String,
    pub notional: u128,
    pub fixed_rate_bps: u32,
    pub payment_interval: u64,
    pub last_payment_time: u64,
    pub floating_rates: Vec<(u64, u32)>,
    pub cashflows: Vec<i128>,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct InterestRateSwap {
    pub swaps: HashMap<u64, Swap>,
    pub next_id: u64,
}

impl InterestRateSwap {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn open_swap(
        &mut self,
        fixed_payer: String,
        floating_payer: String,
        notional: u128,
        fixed_rate_bps: u32,
        start_time: u64,
        payment_interval: u64,
    ) -> u64 {
        self.next_id += 1;
        self.swaps.insert(
            self.next_id,
            Swap {
                fixed_payer,
                floating_payer,
                notional,
                fixed_rate_bps,
                payment_interval,
                last_payment_time: start_time,
                floating_rates: Vec::new(),
                cashflows: Vec::new(),
            },
        );
        self.next_id
    }

    pub fn record_floating_rate(
        &mut self,
        id: u64,
        timestamp: u64,
        rate_bps: u32,
    ) -> Result<(), &'static str> {
        let swap = self.swaps.get_mut(&id).ok_or("not found")?;
        swap.floating_rates.push((timestamp, rate_bps));
        Ok(())
    }

    pub fn settle(&mut self, id: u64, now: u64) -> Result<i128, &'static str> {
        let swap = self.swaps.get_mut(&id).ok_or("not found")?;
        if now < swap.last_payment_time + swap.payment_interval {
            return Err("too early");
        }
        let floating = swap
            .floating_rates
            .iter()
            .rev()
            .find(|(ts, _)| *ts <= now)
            .map(|(_, rate)| *rate)
            .ok_or("missing rate")?;
        let floating_i = i128::try_from(floating).map_err(|_| "overflow")?;
        let fixed_i = i128::from(swap.fixed_rate_bps);
        let notional_i = i128::try_from(swap.notional).map_err(|_| "overflow")?;
        let delta_rate = floating_i - fixed_i;
        let cashflow = delta_rate * notional_i / 10_000;
        swap.last_payment_time = now;
        swap.cashflows.push(cashflow);
        Ok(cashflow)
    }
}
