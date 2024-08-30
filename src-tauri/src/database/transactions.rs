use rusqlite::{named_params, Connection, Error, Transaction};
use time::{macros::date, Date};

/// Adds a single transaction with no associated postings
pub fn insert_transactions(
    tx: &Transaction,
    transaction_date: Date,
    description: &str,
) -> Result<i64, Error> {
    tx.execute(
        "INSERT INTO transactions (transaction_date, description) VALUES (:tx_date, :desc)",
        named_params! {
            ":tx_date": transaction_date,
            ":desc": description,
        },
    )?;

    Ok(tx.last_insert_rowid())
}

/// Adds a single posting
pub fn insert_posting(
    tx: &Transaction,
    id: i64,
    acct: &str,
    val: i64,
    currency: &str,
    comment: &str,
) -> Result<(), Error> {
    tx.execute(
        "INSERT INTO postings (transaction_id, account, value, currency, comment) VALUES (:tx_id, :acct, :val, :currency, :comment)",
        named_params! {
            ":tx_id": id,
            ":acct": acct,
            ":val": val,
            ":currency": currency,
            ":comment": comment,
        }
    )?;

    Ok(())
}

/// Deletes a transaction by id and all associated postings
pub fn delete_transaction(id: i64) -> Result<(), Error> {
    let mut conn = Connection::open(super::get_db_file())?;

    let tx = conn.transaction()?;

    tx.execute(
        "DELETE FROM transactions WHERE id=:id",
        named_params! {
            ":id": id
        },
    )?;

    tx.commit()?;

    Ok(())
}

pub fn update_transaction(
    id: i64,
    transaction_date: Option<Date>,
    description: Option<&str>,
) -> Result<(), Error> {
    let mut statement = String::from("UPDATE transactions SET ");

    statement.push_str(match (transaction_date, description) {
        (Some(_), Some(_)) => "transaction_date = :t_date, description = :desc",
        (Some(_), None) => "transaction_date = :t_date",
        (None, Some(_)) => "description = :desc",
        (None, None) => return Ok(()),
    });

    statement.push_str(" WHERE id = :id;");

    let mut conn = Connection::open(super::get_db_file())?;

    let tx = conn.transaction()?;

    tx.execute(
        &statement,
        named_params! {
            ":t_date": transaction_date.unwrap_or(date!(1970 - 1 - 1)),
            ":desc": description.unwrap_or(""),
            ":id": id,
        },
    )?;

    tx.commit()?;

    Ok(())
}

pub fn update_posting(
    t_id: i64,
    id: i64,
    acct: Option<&str>,
    value: Option<i64>,
    currency: Option<&str>,
    comment: Option<&str>,
) -> Result<(), Error> {
    let mut statement = String::from("UPDATE postings SET ");

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

    statement.push_str(" WHERE t_id = :t_id AND id = :id;");

    let mut conn = Connection::open(super::get_db_file())?;

    let tx = conn.transaction()?;

    tx.execute(
        &statement,
        named_params! {
            ":acct": acct.unwrap_or(""),
            ":val": value.unwrap_or(0),
            ":currency": currency.unwrap_or(""),
            ":comment": comment.unwrap_or(""),
            ":t_id": t_id,
            ":id": id
        },
    )?;

    tx.commit()?;

    Ok(())
}
