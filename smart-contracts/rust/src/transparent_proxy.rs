use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct TransparentProxy {
    pub admin: String,
    pub implementation: String,
    pub pending_implementation: Option<String>,
    pub version: u32,
}

impl TransparentProxy {
    pub fn new(admin: String, implementation: String) -> Self {
        Self {
            admin,
            implementation,
            pending_implementation: None,
            version: 1,
        }
    }

    pub fn upgrade_to(&mut self, caller: &str, implementation: String) -> Result<(), &'static str> {
        if caller != self.admin {
            return Err("only admin");
        }
        self.pending_implementation = Some(implementation);
        Ok(())
    }

    pub fn finalize_upgrade(&mut self, caller: &str) -> Result<(), &'static str> {
        if caller != self.admin {
            return Err("only admin");
        }
        if let Some(new_impl) = self.pending_implementation.take() {
            self.implementation = new_impl;
            self.version += 1;
            Ok(())
        } else {
            Err("no upgrade")
        }
    }

    pub fn change_admin(&mut self, caller: &str, new_admin: String) -> Result<(), &'static str> {
        if caller != self.admin {
            return Err("only admin");
        }
        self.admin = new_admin;
        Ok(())
    }
}
