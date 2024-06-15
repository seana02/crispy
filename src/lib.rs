use std::error::Error;
use std::fs;

use chrono::NaiveDate;

use rusqlite::Result;
use rusqlite::{named_params, Connection};

pub struct Transaction {
    pub date: NaiveDate,
    pub description: String,
    pub postings: Vec<Posting>,
}

pub struct Posting {
    pub account: String,
    pub value: f32, // up to 7 decimal digits
    pub currency: String,
    pub comment: String,
}

pub fn import_hledger(file: &String) -> Result<(), Box<dyn Error>> {
    for group in fs::read_to_string(file)?.split("\n\n") {
        if !group.is_empty() && !group.starts_with("\n") {
            let mut lines = group.lines();
            let transaction = lines.next().expect("unempty transaction");
            println!("{}", transaction);
            for post in lines {
                // post.trim()
                // account:subaccount     $$.$$  ; comment
                println!("..{}", post.trim());
            }
        }
    }

    Ok(())
}

pub fn add_transaction(conn: &Connection, t: &Transaction) -> Result<()> {
    let transaction_id: u32 = next_transaction_id(conn);

    conn.execute(
        "INSERT INTO transactions (id, transaction_date, description)
         VALUES (:id, :t_date, :desc)",
        named_params! {
           ":id": &transaction_id,
           ":t_date": t.date.to_string(),
           ":desc": t.description
        },
    )?;

    for p in &t.postings {
        conn.execute(
            "INSERT INTO posting (transaction_id, account, value, currency, comment)
            VALUES (:t_id, :acct, :val, :curr, :comment)",
            named_params! {
                ":t_id": &transaction_id,
                ":acct": p.account,
                ":val": p.value,
                ":curr": p.currency,
                ":comment": p.comment
            },
        )?;
    }

    Ok(())
}

fn next_transaction_id(conn: &Connection) -> u32 {
    conn.query_row("SELECT MAX(id) FROM transactions", [], |row| row.get(0))
        .unwrap_or(1)
}
