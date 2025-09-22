use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug, PartialEq, Eq)]
pub enum BountyStatus {
    Open,
    Assigned,
    Submitted,
    Released,
    Refunded,
}

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Bounty {
    pub creator: String,
    pub deposit: u128,
    pub hunter: Option<String>,
    pub proof: Option<String>,
    pub status: BountyStatus,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct BountyEscrow {
    pub bounties: HashMap<u64, Bounty>,
    pub next_id: u64,
}

impl BountyEscrow {
    pub fn new() -> Self {
        Self::default()
    }

    pub fn create_bounty(&mut self, creator: String, deposit: u128) -> u64 {
        self.next_id += 1;
        let bounty = Bounty {
            creator,
            deposit,
            hunter: None,
            proof: None,
            status: BountyStatus::Open,
        };
        self.bounties.insert(self.next_id, bounty);
        self.next_id
    }

    pub fn assign_hunter(
        &mut self,
        id: u64,
        caller: &str,
        hunter: String,
    ) -> Result<(), &'static str> {
        let bounty = self.bounties.get_mut(&id).ok_or("not found")?;
        if bounty.creator != caller {
            return Err("unauthorised");
        }
        if bounty.status != BountyStatus::Open {
            return Err("not open");
        }
        bounty.hunter = Some(hunter);
        bounty.status = BountyStatus::Assigned;
        Ok(())
    }

    pub fn submit_proof(
        &mut self,
        id: u64,
        caller: &str,
        proof: String,
    ) -> Result<(), &'static str> {
        let bounty = self.bounties.get_mut(&id).ok_or("not found")?;
        match (&bounty.hunter, bounty.status.clone()) {
            (Some(h), BountyStatus::Assigned) if h == caller => {
                bounty.proof = Some(proof);
                bounty.status = BountyStatus::Submitted;
                Ok(())
            }
            _ => Err("unauthorised"),
        }
    }

    pub fn release(&mut self, id: u64, caller: &str) -> Result<u128, &'static str> {
        let bounty = self.bounties.get_mut(&id).ok_or("not found")?;
        if bounty.creator != caller {
            return Err("unauthorised");
        }
        if bounty.status != BountyStatus::Submitted {
            return Err("not submitted");
        }
        bounty.status = BountyStatus::Released;
        Ok(bounty.deposit)
    }

    pub fn refund(&mut self, id: u64, caller: &str) -> Result<u128, &'static str> {
        let bounty = self.bounties.get_mut(&id).ok_or("not found")?;
        if bounty.creator != caller {
            return Err("unauthorised");
        }
        if bounty.status == BountyStatus::Released {
            return Err("already released");
        }
        bounty.status = BountyStatus::Refunded;
        Ok(bounty.deposit)
    }
}
