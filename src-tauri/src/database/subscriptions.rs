use rusqlite::{named_params, Connection, Error, Transaction};

use crate::types::SubscriptionFrequency;

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
    comment: &str,
) -> Result<(), Error> {
    tx.execute(
        "INSERT INTO subscription_templates (subscription_id, account, value, currency, comment)
        VALUES (:sub_id, :acct, :val, :currency, :comment)",
        named_params! {
            ":sub_id": id,
            ":acct": acct,
            ":val": val,
            ":currency": currency,
            ":comment": comment,
        },
    )?;

    Ok(())
}

/// Deletes a subscription template
pub fn delete_subscription(id: i64) -> Result<(), Error> {
    let mut conn = Connection::open(super::get_db_file())?;

    let tx = conn.transaction()?;

    tx.execute(
        "DELETE FROM subscriptions WHERE id=:id",
        named_params! {
            ":id": id
        },
    )?;

    tx.commit()?;

    Ok(())
}
