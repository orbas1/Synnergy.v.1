use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq)]
pub enum OptionType {
    Call,
    Put,
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct OptionContract {
    pub writer: String,
    pub holder: Option<String>,
    pub strike_price: u128,
    pub premium: u128,
    pub expiry: u64,
    pub option_type: OptionType,
    pub exercised: bool,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct Options {
    pub listings: HashMap<u64, OptionContract>,
    pub next_id: u64,
}

impl Options {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn list_option(
        &mut self,
        writer: String,
        strike_price: u128,
        premium: u128,
        expiry: u64,
        option_type: OptionType,
    ) -> u64 {
        self.next_id += 1;
        self.listings.insert(
            self.next_id,
            OptionContract {
                writer,
                holder: None,
                strike_price,
                premium,
                expiry,
                option_type,
                exercised: false,
            },
        );
        self.next_id
    }

    pub fn purchase(&mut self, id: u64, buyer: String, payment: u128) -> Result<(), &'static str> {
        let listing = self.listings.get_mut(&id).ok_or("not found")?;
        if listing.holder.is_some() {
            return Err("sold");
        }
        if payment < listing.premium {
            return Err("insufficient premium");
        }
        listing.holder = Some(buyer);
        Ok(())
    }

    pub fn exercise(&mut self, id: u64, caller: &str, price: u128, now: u64) -> Result<u128, &'static str> {
        let listing = self.listings.get_mut(&id).ok_or("not found")?;
        let holder = listing.holder.as_ref().ok_or("not sold")?;
        if holder != caller {
            return Err("forbidden");
        }
        if now > listing.expiry {
            return Err("expired");
        }
        if listing.exercised {
            return Err("exercised");
        }
        let payoff = match listing.option_type {
            OptionType::Call => price.saturating_sub(listing.strike_price),
            OptionType::Put => listing.strike_price.saturating_sub(price),
        };
        listing.exercised = true;
        Ok(payoff)
    }
}
