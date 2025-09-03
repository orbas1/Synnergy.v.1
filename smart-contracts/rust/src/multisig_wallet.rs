use std::collections::{HashSet};
use serde::{Serialize, Deserialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Transaction {
    pub to: String,
    pub value: u128,
    pub data: Vec<u8>,
    pub executed: bool,
    pub confirmations: HashSet<String>,
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct MultisigWallet {
    pub owners: Vec<String>,
    pub required: usize,
    pub transactions: Vec<Transaction>,
}

impl MultisigWallet {
    pub fn new(owners: Vec<String>, required: usize) -> Self {
        Self { owners, required, transactions: Vec::new() }
    }

    fn is_owner(&self, addr: &str) -> bool { self.owners.iter().any(|o| o == addr) }

    pub fn submit(&mut self, caller: &str, to: String, value: u128, data: Vec<u8>) -> Result<usize, &'static str> {
        if !self.is_owner(caller) { return Err("not owner"); }
        let tx = Transaction { to, value, data, executed: false, confirmations: HashSet::new() };
        self.transactions.push(tx);
        Ok(self.transactions.len() - 1)
    }

    pub fn approve(&mut self, caller: &str, id: usize) -> Result<(), &'static str> {
        if !self.is_owner(caller) { return Err("not owner"); }
        let tx = self.transactions.get_mut(id).ok_or("not found")?;
        if tx.executed { return Err("executed"); }
        tx.confirmations.insert(caller.to_string());
        if tx.confirmations.len() >= self.required {
            tx.executed = true; // in model we mark executed automatically
        }
        Ok(())
    }
}
