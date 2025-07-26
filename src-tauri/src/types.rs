use std::{fmt, str::FromStr};

use rusqlite::{
    types::{FromSql, FromSqlError, FromSqlResult, ToSqlOutput, ValueRef},
    ToSql,
};
use rust_decimal::Decimal;
use serde::{ser::SerializeStruct, Deserialize, Serialize, Serializer};
use tauri::utils::acl::ParseIdentifierError;
use time::{format_description, Date};

#[derive(Serialize, Deserialize, Debug)]
pub struct Posting {
    pub id: Option<i64>,
    pub account: String,
    pub value: Decimal,
    pub currency: String,
    pub comment: String,
}

impl Posting {
    pub fn new(
        id: Option<i64>,
        account: String,
        value: String,
        currency: String,
        comment: String,
    ) -> Self {
        Posting {
            id,
            account,
            value: Decimal::from_str_exact(&value).unwrap(),
            currency,
            comment,
        }
    }
}

#[derive(Deserialize, Debug)]
pub struct Transaction {
    pub id: Option<i64>,
    pub transaction_date: Date,
    pub description: String,
    pub postings: Vec<Posting>,
}

impl Transaction {
    pub fn add_posting(&mut self, p: Posting) {
        self.postings.push(p);
    }

    pub fn check(&self) -> bool {
        self.balance() == Decimal::ZERO
    }

    pub fn balance(&self) -> Decimal {
        self.postings
            .iter()
            .fold(Decimal::ZERO, |x, p| x + p.value)
            .normalize()
    }
}

impl Serialize for Transaction {
    fn serialize<S: Serializer>(&self, serializer: S) -> Result<S::Ok, S::Error> {
        //let date = (
        //    self.transaction_date.year(),
        //    self.transaction_date.month(),
        //    self.transaction_date.day(),
        //);
        let mut s = serializer.serialize_struct("Transaction", 4)?;
        s.serialize_field("id", &self.id)?;
        //s.serialize_field("transaction_date", &date)?;
        s.serialize_field("transaction_date", &self.transaction_date.format(&format_description::parse("[year]-[month]-[day]").unwrap()).unwrap())?;
        s.serialize_field("description", &self.description)?;
        s.serialize_field("postings", &self.postings)?;
        s.end()
    }
}

pub struct Subscription {
    pub id: i64,
    pub description: String,
    pub last_update_date: Date,
    pub frequency: SubscriptionFrequency,
    pub postings: Vec<Posting>,
}

impl Subscription {
    pub fn check(&self) -> bool {
        // Has the potential to overflow, but unlikely
        self.postings.iter().fold(Decimal::ZERO, |x, p| x + p.value) == Decimal::ZERO
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
