use std::collections::hash_map::DefaultHasher;
use std::collections::HashMap;
use std::hash::{Hash, Hasher};
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Transaction {
    pub sender: String,
    pub receiver: String,
    pub amount: u128,
    pub proof: String,
    pub executed: bool,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct ZKTransaction {
    pub verifying_key: String,
    pub transactions: HashMap<u64, Transaction>,
    pub next_id: u64,
}

impl ZKTransaction {
    pub fn new(verifying_key: String) -> Self {
        Self {
            verifying_key,
            transactions: HashMap::new(),
            next_id: 0,
        }
    }

    pub fn submit_transaction(
        &mut self,
        sender: String,
        receiver: String,
        amount: u128,
        proof: String,
    ) -> Result<u64, &'static str> {
        if !self.verify(&sender, &receiver, amount, &proof) {
            return Err("invalid proof");
        }
        self.next_id += 1;
        self.transactions.insert(
            self.next_id,
            Transaction {
                sender,
                receiver,
                amount,
                proof,
                executed: false,
            },
        );
        Ok(self.next_id)
    }

    pub fn execute(&mut self, id: u64) -> Result<&Transaction, &'static str> {
        let tx = self.transactions.get_mut(&id).ok_or("not found")?;
        if tx.executed {
            return Err("executed");
        }
        tx.executed = true;
        Ok(tx)
    }

    fn verify(&self, sender: &str, receiver: &str, amount: u128, proof: &str) -> bool {
        let mut hasher = DefaultHasher::new();
        sender.hash(&mut hasher);
        receiver.hash(&mut hasher);
        amount.hash(&mut hasher);
        self.verifying_key.hash(&mut hasher);
        let expected = format!("{:x}", hasher.finish());
        expected == proof
    }
}
