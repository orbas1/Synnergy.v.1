use serde::{Serialize, Deserialize};

/// Simple escrow model releasing funds based on arbiter decision.
#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct EscrowPayment {
    pub payer: String,
    pub payee: String,
    pub arbiter: String,
    pub balance: u128,
    pub released: bool,
}

impl EscrowPayment {
    pub fn new(payer: String, payee: String, arbiter: String, deposit: u128) -> Self {
        Self { payer, payee, arbiter, balance: deposit, released: false }
    }

    pub fn release(&mut self, caller: &str) -> Result<u128, &'static str> {
        if caller != self.arbiter { return Err("arbiter"); }
        if self.released { return Err("done"); }
        self.released = true;
        Ok(self.balance)
    }

    pub fn refund(&mut self, caller: &str) -> Result<u128, &'static str> {
        if caller != self.arbiter { return Err("arbiter"); }
        if self.released { return Err("done"); }
        self.released = true;
        Ok(0) // indicates refund to payer
    }
}
