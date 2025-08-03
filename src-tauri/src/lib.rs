use std::collections::HashMap;

use database::{get_db_file, transactions::{delete, delete_posting, get_transaction_count, get_transaction_details, get_transactions_between, get_transactions_by_account, get_transactions_with_description, update}};
use rusqlite::Connection;
use rust_decimal::Decimal;
use time::{macros::date, Date, Month};
use types::{Posting, Transaction};

mod database;
mod error;
mod types;

/// Gets a list of the most recent transactions
#[tauri::command]
fn get_all_transaction_list() -> Result<HashMap<i64,Transaction>, String> {
    let conn = Connection::open(get_db_file()).map_err(|e| e.to_string())?;
    conn.pragma_update(None, "foreign_keys", "ON").map_err(|e| e.to_string())?;
    get_transactions_between(&conn, date!{1970 - 1 - 1}, Date::MAX).map_err(|e| e.to_string())
}

/// Get a complete transaction
#[tauri::command]
fn get_transaction_by_id(id: i64) -> Result<Transaction, String> {
    let conn = Connection::open(get_db_file()).map_err(|e| e.to_string())?;
    conn.pragma_update(None, "foreign_keys", "ON").map_err(|e| e.to_string())?;
    get_transaction_details(&conn, id).map_err(|e| e.to_string())
}

/// Get list of transactions between given dates
#[tauri::command]
fn get_all_transactions_between(since: Option<Date>, until: Option<Date>) -> Result<HashMap<i64, Transaction>, String> {
    let conn = Connection::open(get_db_file()).map_err(|e| e.to_string())?;
    conn.pragma_update(None, "foreign_keys", "ON").map_err(|e| e.to_string())?;
    get_transactions_between(&conn, since.unwrap_or(date!{1970 - 1 - 1}), until.unwrap_or(Date::MAX)).map_err(|e| e.to_string())
}

// Get list of transactions with associated account
#[tauri::command]
fn get_transactions_involving_account(acct: String) -> Result<HashMap<i64, Transaction>, String> {
    let conn = Connection::open(get_db_file()).map_err(|e| e.to_string())?;
    conn.pragma_update(None, "foreign_keys", "ON").map_err(|e| e.to_string())?;
    get_transactions_by_account(&conn, &acct).map_err(|e| e.to_string())
}

// Get list of transactions with given string in description or posting comments
#[tauri::command]
fn get_transactions_by_text(text: String) -> Result<HashMap<i64, Transaction>, String> {
    let conn = Connection::open(get_db_file()).map_err(|e| e.to_string())?;
    conn.pragma_update(None, "foreign_keys", "ON").map_err(|e| e.to_string())?;
    get_transactions_with_description(&conn, &text).map_err(|e| e.to_string())
}

#[tauri::command]
fn add_transaction(
    year: i32,
    month: usize,
    day: u8,
    postings: Vec<Posting>,
    desc: String,
) -> Result<i64, String> {
    let t = Transaction {
        id: None,
        transaction_date: Date::from_calendar_date(year, get_month(month), day).unwrap(),
        description: desc,
        postings,
    };
    if t.check() {
        println!("Balanced, adding to db");
        let mut conn = Connection::open(get_db_file()).map_err(|e| e.to_string())?;
        conn.pragma_update(None, "foreign_keys", "ON").map_err(|e| e.to_string())?;
        let tx = conn.transaction().map_err(|e| e.to_string())?;
        let result = match database::transactions::insert(&tx, t) {
            Ok(i) => Ok(i),
            Err(e) => Err(e.to_string()),
        };
        tx.commit().map_err(|e| e.to_string())?;
        result
    } else {
        println!("Unbalanced, returning");
        Err(t.balance().to_string())
    }
}

#[tauri::command]
fn update_transaction(
    id: i64,
    year: i32,
    month: usize,
    day: u8,
    postings: Vec<Posting>,
    desc: String,
    delete_list: Vec<i64>
) -> Result<(), String> {
    let t = Transaction {
        id: Some(id),
        transaction_date: Date::from_calendar_date(year, get_month(month), day).unwrap(),
        description: desc,
        postings,
    };
    let mut conn = Connection::open(get_db_file()).map_err(|e| e.to_string())?;
    conn.pragma_update(None, "foreign_keys", "ON").map_err(|e| e.to_string())?;
    let tx = conn.transaction().map_err(|e| e.to_string())?;
    update(&tx, t, delete_list).map_err(|e| e.to_string())?;
    tx.commit().map_err(|e| e.to_string())?;
    Ok(())
}

#[tauri::command]
fn delete_transaction(id: i64) -> Result<bool, String> {
    let mut conn = Connection::open(get_db_file()).map_err(|e| e.to_string())?;
    conn.pragma_update(None, "foreign_keys", "ON").map_err(|e| e.to_string())?;
    let tx = conn.transaction().map_err(|e| e.to_string())?;
    let output = match delete(&tx, id) {
        Ok(_) => Ok(true),
        Err(e) => Err(e.to_string()),
    };
    tx.commit().map_err(|e| e.to_string())?;
    output
}

#[tauri::command]
fn get_row_count() -> Result<i32, String> {
    match get_transaction_count() {
        Ok(i) => Ok(i),
        Err(e) => Err(e.to_string()),
    }
}

#[tauri::command]
fn validate(postings: Vec<Posting>) -> bool {
    println!("hello");
    println!("{}", postings.len());
    if postings.len() == 0 {
        return false;
    }
    let mut sum = Decimal::ZERO;
    for p in postings.iter() {
        println!("{} {} {}", p.account, p.value, p.currency);
        if p.account.len() == 0 {
            return false;
        }
        if p.currency.len() == 0 {
            return false;
        }
        sum += p.value;
    }
    println!("sum: {}", sum);
    return sum == Decimal::ZERO;
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .invoke_handler(tauri::generate_handler![
            get_all_transaction_list,
            get_transaction_by_id,
            get_all_transactions_between,
            get_transactions_involving_account,
            add_transaction,
            update_transaction,
            delete_transaction,
        ])
        .setup(|_| {
            database::init();
            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}

fn get_month(month: usize) -> Month {
    let months: [Month; 12] = [
        Month::January,
        Month::February,
        Month::March,
        Month::April,
        Month::May,
        Month::June,
        Month::July,
        Month::August,
        Month::September,
        Month::October,
        Month::November,
        Month::December,
    ];
    months[month-1]
}
