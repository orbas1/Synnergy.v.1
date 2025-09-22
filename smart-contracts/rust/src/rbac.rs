use std::collections::{HashMap, HashSet};
use serde::{Deserialize, Serialize};

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct RBAC {
    roles: HashMap<String, HashSet<String>>,
    pub admin_role: String,
}

impl RBAC {
    pub fn new(admin_role: String) -> Self {
        Self {
            roles: HashMap::new(),
            admin_role,
        }
    }

    pub fn grant_role(&mut self, role: String, account: String) {
        self.roles.entry(role).or_default().insert(account);
    }

    pub fn revoke_role(&mut self, role: &str, account: &str) {
        if let Some(members) = self.roles.get_mut(role) {
            members.remove(account);
        }
    }

    pub fn has_role(&self, role: &str, account: &str) -> bool {
        self.roles
            .get(role)
            .map(|set| set.contains(account))
            .unwrap_or(false)
    }

    pub fn list_members(&self, role: &str) -> Vec<String> {
        self.roles
            .get(role)
            .map(|set| set.iter().cloned().collect())
            .unwrap_or_default()
    }
}
