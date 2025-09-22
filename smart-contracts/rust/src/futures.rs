use std::collections::HashMap;
use std::convert::TryFrom;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct FutureContract {
    pub buyer: String,
    pub seller: String,
    pub quantity: u128,
    pub price: u128,
    pub margin_buyer: u128,
    pub margin_seller: u128,
    pub settled: bool,
    pub settlement_price: Option<u128>,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct Futures {
    pub contracts: HashMap<u64, FutureContract>,
    pub next_id: u64,
}

impl Futures {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn create_contract(
        &mut self,
        buyer: String,
        seller: String,
        quantity: u128,
        price: u128,
    ) -> u64 {
        self.next_id += 1;
        let contract = FutureContract {
            buyer,
            seller,
            quantity,
            price,
            margin_buyer: 0,
            margin_seller: 0,
            settled: false,
            settlement_price: None,
        };
        self.contracts.insert(self.next_id, contract);
        self.next_id
    }

    pub fn deposit_margin(&mut self, id: u64, caller: &str, amount: u128) -> Result<(), &'static str> {
        let contract = self.contracts.get_mut(&id).ok_or("not found")?;
        if contract.settled {
            return Err("settled");
        }
        if caller == contract.buyer {
            contract.margin_buyer = contract.margin_buyer.saturating_add(amount);
            Ok(())
        } else if caller == contract.seller {
            contract.margin_seller = contract.margin_seller.saturating_add(amount);
            Ok(())
        } else {
            Err("unauthorised")
        }
    }

    pub fn settle(&mut self, id: u64, caller: &str, settlement_price: u128) -> Result<i128, &'static str> {
        let contract = self.contracts.get_mut(&id).ok_or("not found")?;
        if caller != contract.buyer && caller != contract.seller {
            return Err("unauthorised");
        }
        if contract.settled {
            return Err("settled");
        }
        let settlement_price_i = i128::try_from(settlement_price).map_err(|_| "overflow")?;
        let price_i = i128::try_from(contract.price).map_err(|_| "overflow")?;
        let quantity_i = i128::try_from(contract.quantity).map_err(|_| "overflow")?;
        let pnl = (settlement_price_i - price_i) * quantity_i;
        contract.settled = true;
        contract.settlement_price = Some(settlement_price);
        Ok(pnl)
    }
}
