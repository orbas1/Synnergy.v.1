use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Invoice {
    pub seller: String,
    pub amount: u128,
    pub due_date: u64,
    pub purchased: bool,
    pub paid: bool,
    pub investor: Option<String>,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct InvoiceFactoring {
    invoices: HashMap<u64, Invoice>,
    pub next_id: u64,
}

impl InvoiceFactoring {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn submit_invoice(&mut self, seller: String, amount: u128, due_date: u64) -> u64 {
        self.next_id += 1;
        self.invoices.insert(
            self.next_id,
            Invoice {
                seller,
                amount,
                due_date,
                purchased: false,
                paid: false,
                investor: None,
            },
        );
        self.next_id
    }

    pub fn purchase(&mut self, id: u64, investor: String, advance: u128) -> Result<u128, &'static str> {
        let invoice = self.invoices.get_mut(&id).ok_or("not found")?;
        if invoice.purchased {
            return Err("already purchased");
        }
        if advance > invoice.amount {
            return Err("advance too high");
        }
        invoice.purchased = true;
        invoice.investor = Some(investor);
        Ok(advance)
    }

    pub fn mark_paid(&mut self, id: u64) -> Result<(), &'static str> {
        let invoice = self.invoices.get_mut(&id).ok_or("not found")?;
        invoice.paid = true;
        Ok(())
    }

    pub fn release_payment(&mut self, id: u64) -> Result<u128, &'static str> {
        let invoice = self.invoices.get_mut(&id).ok_or("not found")?;
        if !invoice.paid || !invoice.purchased {
            return Err("not ready");
        }
        let amount = invoice.amount;
        invoice.purchased = false;
        invoice.paid = false;
        invoice.investor = None;
        Ok(amount)
    }
}
