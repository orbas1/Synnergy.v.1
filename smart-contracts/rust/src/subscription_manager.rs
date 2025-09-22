use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize, Clone, Debug)]
pub struct Subscription {
    pub plan: String,
    pub expiry: u64,
    pub active: bool,
}

#[derive(Default, Serialize, Deserialize, Clone, Debug)]
pub struct SubscriptionManager {
    subscriptions: HashMap<String, Subscription>,
    pub default_duration: u64,
}

impl SubscriptionManager {
    pub fn new(default_duration: u64) -> Self {
        Self {
            subscriptions: HashMap::new(),
            default_duration,
        }
    }

    pub fn subscribe(&mut self, user: String, plan: String, start: u64) {
        self.subscriptions.insert(
            user,
            Subscription {
                plan,
                expiry: start + self.default_duration,
                active: true,
            },
        );
    }

    pub fn cancel(&mut self, user: &str) -> Result<(), &'static str> {
        let sub = self.subscriptions.get_mut(user).ok_or("no subscription")?;
        sub.active = false;
        Ok(())
    }

    pub fn renew(&mut self, user: &str, now: u64) -> Result<(), &'static str> {
        let sub = self.subscriptions.get_mut(user).ok_or("no subscription")?;
        if now > sub.expiry {
            sub.expiry = now + self.default_duration;
        } else {
            sub.expiry += self.default_duration;
        }
        sub.active = true;
        Ok(())
    }

    pub fn is_active(&self, user: &str, now: u64) -> bool {
        self.subscriptions
            .get(user)
            .map(|sub| sub.active && now <= sub.expiry)
            .unwrap_or(false)
    }
}
