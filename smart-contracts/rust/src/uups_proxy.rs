use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct UUPSProxy {
    pub implementation: String,
    pub upgrader: String,
    pub version: u32,
}

impl UUPSProxy {
    pub fn new(implementation: String, upgrader: String) -> Self {
        Self {
            implementation,
            upgrader,
            version: 1,
        }
    }

    pub fn upgrade(&mut self, caller: &str, new_implementation: String) -> Result<(), &'static str> {
        if caller != self.upgrader {
            return Err("only upgrader");
        }
        self.implementation = new_implementation;
        self.version += 1;
        Ok(())
    }

    pub fn transfer_upgrader(&mut self, caller: &str, new_upgrader: String) -> Result<(), &'static str> {
        if caller != self.upgrader {
            return Err("only upgrader");
        }
        self.upgrader = new_upgrader;
        Ok(())
    }
}
