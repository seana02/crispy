use std::collections::HashMap;

use rusqlite::{named_params, Connection, Error};
use time::{
    macros::format_description,
    Date,
};

use crate::{
    error::TransactionError,
    types::{Posting, Transaction},
};

use super::get_db_file;

/// Gets all transactions only (no postings) within specified dates
/// since is inclusive, until is exclusive
pub fn get_transactions_between(
    tx: &rusqlite::Transaction,
    since: Date,
    until: Date,
) -> Result<HashMap<i64, Transaction>, Error> {
    let str = String::from(
        "SELECT id, transaction_date, description
        FROM transactions
        WHERE transaction_date BETWEEN :since AND :until
        ORDER BY transaction_date DESC, id DESC",
    );
    let format = format_description!("[year]-[month]-[day]");

    let mut stmt = tx.prepare(&str)?;

    let mut rows = stmt.query(named_params! {
        ":since": &since.format(format).unwrap(),
        ":until": &until.format(format).unwrap(),
    })?;
    let mut ts = HashMap::<i64, Transaction>::new();
    while let Some(row) = rows.next()? {
        ts.insert(row.get(0).unwrap(), Transaction {
            id: row.get(0).unwrap(),
            transaction_date: row.get(1).unwrap(),
            description: row.get(2).unwrap(),
            postings: Vec::new(),
        });
    }
    Ok(ts)
}

/// Gets a single complete transaction by id
pub fn get_transaction_details(tx: &rusqlite::Transaction, id: i64) -> Result<Transaction, Error> {
    let postings = get_postings_by_transaction_id(tx, id)?;
    let result = tx.query_row(
        "SELECT transaction_date, description
        FROM transactions
        WHERE id = :id",
        named_params! {
            ":id": id
        },
        |row| {
            Ok(Transaction {
                id: Some(id),
                transaction_date: row.get(0).unwrap(),
                description: row.get(1).unwrap(),
                postings,
            })
        },
    )?;
    Ok(result)
}

/// Gets all transactions associated with the specified account
pub fn get_transactions_by_account(
    tx: &rusqlite::Transaction,
    acct: &str
) -> Result<HashMap<i64, Transaction>, Error> {
    let str = String::from(
        "SELECT transactions.id, transaction_date, description
        FROM transacations
        LEFT JOIN postings
        ON postings.transaction_id = transactions.id
        WHERE account LIKE \"%:acct%\"",
    );

    let mut stmt = tx.prepare(&str)?;
    
    let mut rows = stmt.query(named_params! {
        ":acct": acct,
    })?;

    let mut ts = HashMap::new();
    while let Some(row) = rows.next()? {
        let t_id = row.get::<usize, i64>(1).unwrap();
        ts.insert(t_id, Transaction {
            id: row.get(0).unwrap(),
            transaction_date:  row.get(1).unwrap(),
            description: row.get(2).unwrap(),
            postings: Vec::new(),
        });
    }
    Ok(ts)
}

/// Get all transactions by description/posting comment text
pub fn get_transactions_with_description(tx: &rusqlite::Transaction, text: &str) -> Result<HashMap<i64, Transaction>, Error> {
    let str = String::from(
        "SELECT DISTINCT transactions.id, transaction_date, description FROM transactions
        LEFT JOIN postings
        ON postings.transaction_id = transactions.id
        WHERE postings.comment LIKE \"%:text%\"
        OR transactions.description LIKE \":text\""
    );
    let mut stmt = tx.prepare(&str)?;

    let mut rows = stmt.query(named_params! {
        ":text": text,
    })?;

    let mut ts = HashMap::new();
    while let Some(row) = rows.next()? {
        let t_id = row.get::<usize, i64>(1).unwrap();
        ts.insert(t_id, Transaction {
            id:  row.get(0).unwrap(),
            transaction_date: row.get(1).unwrap(),
            description: row.get(2).unwrap(),
            postings: Vec::new(),
        });
    }
    Ok(ts)
}

/// Adds a complete transaction to the database
pub fn insert(t: Transaction) -> Result<(), TransactionError> {
    if !t.check() {
        return Err(TransactionError::UnbalancedPostingError);
    }
    let mut conn = Connection::open(get_db_file())?;
    conn.pragma_update(None, "foreign_keys", "ON")?;
    let tx = conn.transaction()?;
    let rowid = insert_transaction(&tx, t.transaction_date, &t.description)?;
    for p in t.postings.iter() {
        insert_posting(
            &tx,
            rowid,
            &p.account,
            &p.value.to_string(),
            &p.currency,
            &p.comment,
        )?;
    }
    tx.commit()?;
    Ok(())
}

/// Updates a complete transaction in the database
pub fn update(t: Transaction) -> Result<(), TransactionError> {
    if !t.check() {
        return Err(TransactionError::UnbalancedPostingError);
    }
    let mut conn = Connection::open(get_db_file())?;
    conn.pragma_update(None, "foreign_keys", "ON")?;
    let tx = conn.transaction()?;
    update_transaction(
        &tx,
        t.id.unwrap(),
        Some(t.transaction_date),
        Some(&t.description),
    )?;
    for p in t.postings.iter() {
        update_posting(
            &tx,
            t.id.unwrap(),
            p.id.unwrap(),
            Some(&p.account),
            Some(&p.value.to_string()),
            Some(&p.currency),
            Some(&p.comment),
        )?;
    }
    tx.commit()?;
    Ok(())
}

/// Adds a single transaction with no associated postings
fn insert_transaction(
    tx: &rusqlite::Transaction,
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
fn insert_posting(
    tx: &rusqlite::Transaction,
    id: i64,
    acct: &str,
    val: &str,
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

pub fn delete(id: i64) -> Result<(), Error> {
    let mut conn = Connection::open(get_db_file())?;
    conn.pragma_update(None, "foreign_keys", "ON")?;
    let tx = conn.transaction()?;
    delete_transaction(&tx, id)?;
    tx.commit()?;
    Ok(())
}

/// Deletes a transaction by id and all associated postings
fn delete_transaction(tx: &rusqlite::Transaction, id: i64) -> Result<(), Error> {
    // SQLite enforces posting deletions
    tx.execute(
        "DELETE FROM transactions WHERE id=:id;",
        named_params! {
            ":id": id
        },
    )?;

    Ok(())
}

/// Deletes a posting by transaction id and posting id
pub fn delete_posting(tx: &rusqlite::Transaction, t_id: i64, id: i64) -> Result<(), Error> {
    tx.execute(
        "DELETE FROM postings WHERE transaction_id=:t_id AND id=:id;",
        named_params! {
            ":t_id": t_id,
            ":id": id
        },
    )?;

    Ok(())
}

/// Updates a transaction date and/or description by id
fn update_transaction(
    tx: &rusqlite::Transaction,
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

    tx.execute(
        &statement,
        named_params! {
            ":t_date": transaction_date.unwrap(),
            ":desc": description.unwrap(),
            ":id": id,
        },
    )?;

    Ok(())
}

// Updates a posting by id
fn update_posting(
    tx: &rusqlite::Transaction,
    t_id: i64,
    id: i64,
    acct: Option<&str>,
    value: Option<&str>,
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

    statement.push_str(" WHERE transaction_id = :t_id AND id = :id;");

    tx.execute(
        &statement,
        named_params! {
            ":acct": acct.unwrap_or(""),
            ":val": value.unwrap_or("0"),
            ":currency": currency.unwrap_or(""),
            ":comment": comment.unwrap_or(""),
            ":t_id": t_id,
            ":id": id
        },
    )?;

    Ok(())
}

/// Gets the count of transactions
pub fn get_transaction_count() ->  Result<i32, Error> {
    let conn = Connection::open(get_db_file())?;
    conn.pragma_update(None, "foreign_keys", "ON")?;
    conn.query_row("SELECT COUNT() FROM transactions", [], |row| row.get(0))
}

/// Gets a vector of postings based on transaction id
fn get_postings_by_transaction_id(
    tx: &rusqlite::Transaction,
    id: i64,
) -> Result<Vec<Posting>, Error> {
    let mut stmt = tx.prepare(
        "SELECT id, account, value, currency, comment FROM postings WHERE transaction_id = :id",
    )?;
    let mut rows = stmt.query(named_params! { ":id": id })?;
    let mut ps = Vec::new();
    while let Some(row) = rows.next()? {
        ps.push(Posting::new(
            row.get(0).unwrap(),
            row.get(1).unwrap(),
            row.get(2).unwrap(),
            row.get(3).unwrap(),
            row.get(4).unwrap_or(String::new()),
        ))
    }
    Ok(ps)
}
