use std::{fmt, str::FromStr};

use rusqlite::{
    types::{FromSql, FromSqlError, FromSqlResult, ToSqlOutput, ValueRef},
    ToSql,
};
use tauri::utils::acl::ParseIdentifierError;
use time::Date;

pub struct Posting {
    pub account: String,
    pub value: i64,
    pub currency: String,
    pub comment: String,
}

pub struct Transaction {
    pub transaction_date: Date,
    pub description: String,
    pub postings: Vec<Posting>,
}

impl Transaction {
    pub fn check(&self) -> bool {
        // Has the potential to overflow, but unlikely
        self.postings.iter().fold(0, |x, p| x + p.value) == 0_i64
    }

    pub fn add_posting(&mut self, p: Posting) {
        self.postings.push(p);
    }
}

pub struct Subscription {
    pub description: String,
    pub last_update_date: Date,
    pub frequency: SubscriptionFrequency,
    pub postings: Vec<Posting>,
}

impl Subscription {
    pub fn check(&self) -> bool {
        // Has the potential to overflow, but unlikely
        self.postings.iter().fold(0, |x, p| x + p.value) == 0_i64
    }

    pub fn add_posting(&mut self, p: Posting) {
        self.postings.push(p);
    }
}

#[derive(Debug)]
pub enum SubscriptionFrequency {
    Daily,
    Weekly,
    Biweekly,
    Monthly,
    Yearly,
}

impl fmt::Display for SubscriptionFrequency {
    fn fmt(&self, f: &mut fmt::Formatter) -> Result<(), fmt::Error> {
        match self {
            SubscriptionFrequency::Daily => write!(f, "Daily"),
            SubscriptionFrequency::Weekly => write!(f, "Weekly"),
            SubscriptionFrequency::Biweekly => write!(f, "Biweekly"),
            SubscriptionFrequency::Monthly => write!(f, "Monthly"),
            SubscriptionFrequency::Yearly => write!(f, "Yearly"),
        }
    }
}

impl FromStr for SubscriptionFrequency {
    type Err = ParseIdentifierError;
    fn from_str(input: &str) -> Result<Self, Self::Err> {
        match input {
            "Daily" => Ok(SubscriptionFrequency::Daily),
            "Weekly" => Ok(SubscriptionFrequency::Weekly),
            "Biweekly" => Ok(SubscriptionFrequency::Biweekly),
            "Monthly" => Ok(SubscriptionFrequency::Monthly),
            "Yearly" => Ok(SubscriptionFrequency::Yearly),
            _ => Err(ParseIdentifierError::InvalidFormat),
        }
    }
}

impl ToSql for SubscriptionFrequency {
    fn to_sql(&self) -> rusqlite::Result<ToSqlOutput<'_>> {
        Ok(self.to_string().into())
    }
}

impl FromSql for SubscriptionFrequency {
    fn column_result(value: ValueRef<'_>) -> FromSqlResult<Self> {
        value
            .as_str()?
            .parse()
            .map_err(|e| FromSqlError::Other(Box::new(e)))
    }
}
