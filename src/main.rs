use std::fs;
use std::path::PathBuf;

use crispy::import_hledger;
use rusqlite::Connection;
use rusqlite::Result;

fn main() -> Result<()> {
    let db_file = dirs::data_dir().unwrap().join("crispy/data.db");
    let conn = if !db_file.try_exists().unwrap() {
        init(db_file)?
    } else {
        println!("found existing file");
        Connection::open(db_file)?
    };

    import_hledger(&"/home/seanaoki/repos/ledger/2023.journal".to_string());

    Ok(())
}

fn init(db_file: PathBuf) -> Result<Connection> {
    fs::create_dir_all(db_file.parent().unwrap()).expect("File could not be created");
    let conn = Connection::open(db_file)?;
    println!("creating new file");

    conn.execute(
        "CREATE TABLE IF NOT EXISTS transactions (
            id                          INTEGER PRIMARY KEY,
            transaction_date            TEXT,
            last_updated                TEXT DEFAULT CURRENT_DATE,
            description                 TEXT
        );",
        (),
    )?;

    conn.execute(
        "CREATE TABLE IF NOT EXISTS posting (
            id                          INTEGER PRIMARY KEY,
            transaction_id              INTEGER,
            account                     TEXT NOT NULL,
            value                       REAL NOT NULL,
            currency                    TEXT NOT NULL DEFAULT 'USD',
            comment                     TEXT,
            FOREIGN KEY(transaction_id) REFERENCES transactions(id)
        );",
        (),
    )?;

    conn.execute(
        "CREATE TABLE IF NOT EXISTS subscription (
            id                          INTEGER PRIMARY KEY,
            description                 TEXT,
            last_updated                TEXT DEFAULT CURRENT_DATE,
            frequency                   TEXT
        );",
        (),
    )?;

    conn.execute(
        "CREATE TABLE IF NOT EXISTS subscription_template (
            subscription_id             INTEGER,
            account                     TEXT NOT NULL,
            currency                    TEXT NOT NULL,
            comment                     TEXT,
            FOREIGN KEY(subscription_id) REFERENCES subscription(id)
        );",
        (),
    )?;

    Ok(conn)
}
