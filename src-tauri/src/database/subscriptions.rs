use rusqlite::{named_params, Connection, Error, Transaction};
use rust_decimal::Decimal;

use crate::types::{Posting, Subscription, SubscriptionFrequency};

/// Adds a single subscription with no associated template items
pub fn insert_subscription(
    tx: &Transaction,
    desc: &str,
    frequency: SubscriptionFrequency,
) -> Result<i64, Error> {
    tx.execute(
        "INSERT INTO subscriptions (description, frequency)
        VALUES (:desc, :freq)",
        named_params! {
            ":desc:": desc,
            ":freq": frequency,
        },
    )?;

    Ok(tx.last_insert_rowid())
}

/// Adds a single subscription template item
pub fn insert_sub_template(
    tx: &Transaction,
    id: i64,
    acct: &str,
    val: i64,
    currency: &str,
    comment: Option<&str>,
) -> Result<(), Error> {
    tx.execute(
        "INSERT INTO subscription_templates (subscription_id, account, value, currency, comment)
        VALUES (:sub_id, :acct, :val, :currency, :comment)",
        named_params! {
            ":sub_id": id,
            ":acct": acct,
            ":val": val,
            ":currency": currency,
            ":comment": comment.unwrap_or(""),
        },
    )?;

    Ok(())
}

/// Deletes a subscription and associated subscription templates
pub fn delete_subscription(tx: &Transaction, id: i64) -> Result<(), Error> {
    tx.execute(
        "DELETE FROM subscriptions WHERE id=:id",
        named_params! {
            ":id": id
        },
    )?;

    Ok(())
}

/// Delete subscription template item by sub id and template id
pub fn delete_sub_template(tx: &Transaction, s_id: i64, id: i64) -> Result<(), Error> {
    tx.execute(
        "DELETE FROM subscription_templates
        WHERE subscription_id=:s_id AND id=:id",
        named_params! {
            ":s_id": s_id,
            ":id": id
        },
    )?;

    Ok(())
}

/// Updates subscription by id
pub fn update_subscription(
    tx: &Transaction,
    id: i64,
    desc: Option<&str>,
    frequency: Option<SubscriptionFrequency>,
) -> Result<(), Error> {
    let mut statement = String::from("UPDATE subscriptions SET ");

    statement.push_str(match (desc, &frequency) {
        (Some(_), Some(_)) => "frequency = :frequency, description = :desc",
        (Some(_), None) => "frequency = :frequency",
        (None, Some(_)) => "description = :desc",
        (None, None) => return Ok(()),
    });

    statement.push_str(" WHERE id = :id;");

    tx.execute(
        &statement,
        named_params! {
            ":frequency": frequency.unwrap(),
            ":desc": desc.unwrap_or(""),
            ":id": id,
        },
    )?;

    Ok(())
}

/// Updates subscription template by id
pub fn update_sub_template(
    tx: &rusqlite::Transaction,
    s_id: i64,
    id: i64,
    acct: Option<&str>,
    value: Option<i64>,
    currency: Option<&str>,
    comment: Option<&str>,
) -> Result<(), Error> {
    let mut statement = String::from("UPDATE subscription_templates SET ");

    let mut first = true;

    if let Some(_) = acct {
        statement.push_str("account = :acct");
        first = false;
    }

    if let Some(_) = value {
        if !first {
            statement.push_str(", ");
        }
        statement.push_str("value = :val");
        first = false;
    }

    if let Some(_) = currency {
        if !first {
            statement.push_str(", ");
        }
        statement.push_str("currency = :currency");
        first = false;
    }

    if let Some(_) = comment {
        if !first {
            statement.push_str(", ");
        }
        statement.push_str("comment = :comment");
    }

    if first {
        return Ok(());
    };

    statement.push_str(" WHERE s_id = :s_id AND id = :id;");

    tx.execute(
        &statement,
        named_params! {
            ":acct": acct.unwrap_or(""),
            ":val": value.unwrap_or(0),
            ":currency": currency.unwrap_or(""),
            ":comment": comment.unwrap_or(""),
            ":s_id": s_id,
            ":id": id
        },
    )?;

    Ok(())
}

/// Get single subscription by id without postings
pub fn get_subscription_by_id(id: i64) -> Result<Subscription, Error> {
    let conn = Connection::open(super::get_db_file())?;

    let result = conn.query_row(
        "SELECT description, last_updated, frequency
        FROM subscriptions
        WHERE id = :id",
        named_params! {
            ":id": id
        },
        |row| {
            Ok(Subscription {
                id,
                description: row.get(0).unwrap(),
                last_update_date: row.get(1).unwrap(),
                frequency: row.get(2).unwrap(),
                postings: Vec::new(),
            })
        },
    )?;
    Ok(result)
}

/// Get all subscription template postings by subscription id
pub fn get_sub_templates_by_sub_id(
    tx: &rusqlite::Transaction,
    s_id: i64,
) -> Result<Vec<Posting>, Error> {
    let mut stmt = tx.prepare(
        "SELECT id, account, value, currency, comment FROM postings WHERE subscription_id = :s_id",
    )?;
    let mut rows = stmt.query(named_params! { ":s_id": s_id })?;
    let mut ps = Vec::new();
    while let Some(row) = rows.next()? {
        ps.push(Posting {
            id: row.get(0).unwrap(),
            account: row.get(1).unwrap(),
            value: Decimal::from_str_exact(&row.get::<usize, String>(2).unwrap()).unwrap(),
            currency: row.get(3).unwrap(),
            comment: row.get(4).unwrap(),
        });
    }
    Ok(ps)
}
