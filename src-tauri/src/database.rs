mod subscriptions;
mod transactions;

use std::{
    fs::{self, File},
    path::PathBuf,
};

use rusqlite::Connection;

fn get_db_file() -> PathBuf {
    dirs::data_dir()
        .expect("Could not resolve data directory")
        .join("crispy/data.db")
}

pub fn init() {
    let file = get_db_file();
    fs::create_dir_all(file.parent().unwrap()).expect("Could not create data directory");
    if !file.exists() {
        create_db_file(&file).expect("Could not initialize SQLite database");
    }
}

fn create_db_file(file: &PathBuf) -> Result<(), rusqlite::Error> {
    File::create_new(file).expect("Could not create database file");

    let conn = Connection::open(file)?;

    conn.execute_batch(
        "CREATE TABLE IF NOT EXISTS transactions (
            id                          INTEGER PRIMARY KEY,
            transaction_date            TEXT,
            description                 TEXT
        );
        CREATE TABLE IF NOT EXISTS postings (
            id                          INTEGER PRIMARY KEY,
            transaction_id              INTEGER,
            account                     TEXT NOT NULL,
            value                       INTEGER NOT NULL,
            currency                    TEXT NOT NULL DEFAULT 'USD',
            comment                     TEXT,
            FOREIGN KEY(transaction_id) REFERENCES transactions(id) ON DELETE CASCADE
        );
        CREATE TABLE IF NOT EXISTS subscriptions (
            id                          INTEGER PRIMARY KEY,
            description                 TEXT,
            last_updated                TEXT DEFAULT CURRENT_DATE,
            frequency                   TEXT
        );
        CREATE TABLE IF NOT EXISTS subscription_templates (
            id                          INTEGER PRIMARY KEY,
            subscription_id             INTEGER,
            account                     TEXT NOT NULL,
            value                       INTEGER NOT NULL,
            currency                    TEXT NOT NULL,
            comment                     TEXT,
            FOREIGN KEY(subscription_id) REFERENCES subscription(id) ON DELETE CASCADE
        );",
    )?;

    Ok(())
}
