use std::collections::HashMap;

use database::transactions::{delete, get_all_transaction_list, get_recent_transactions, get_single_transaction, get_transaction_count, update};
use rust_decimal::Decimal;
use time::{Date, Month};
use types::{Posting, Transaction};

mod database;
mod error;
mod state;
mod types;

// Learn more about Tauri commands at https://tauri.app/v1/guides/features/command
#[tauri::command]
fn get_transactions() -> Result<HashMap<i64,Transaction>, String> {
    match get_all_transaction_list() {
        Ok(t) => Ok(t),
        Err(e) => {
            println!("{}", e.to_string());
            Err(e.to_string())
        }
    }
}

#[tauri::command]
fn get_transaction_details(id: i64) -> Result<Transaction, String> {
    match get_single_transaction(id) {
        Ok(t) => Ok(t),
        Err(e) => {
            println!("{}", e.to_string());
            Err(e.to_string())
        }
    }
}

#[tauri::command]
fn add_transaction(
    year: i32,
    month: usize,
    day: u8,
    postings: Vec<Posting>,
    desc: String,
) -> Result<(), String> {
    let t = Transaction {
        id: None,
        transaction_date: Date::from_calendar_date(year, get_month(month), day).unwrap(),
        description: desc,
        postings,
    };
    if t.check() {
        println!("Balanced, adding to db");
        match database::transactions::insert(t) {
            Ok(_) => Ok(()),
            Err(e) => Err(e.to_string()),
        }
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
) -> Result<(), String> {
    let t = Transaction {
        id: Some(id),
        transaction_date: Date::from_calendar_date(year, get_month(month), day).unwrap(),
        description: desc,
        postings,
    };
    if t.check() {
        println!("Balanced, updating db");
        match update(t) {
            Ok(_) => Ok(()),
            Err(e) => Err(e.to_string()),
        }
    } else {
        println!("Unbalanced, returning");
        Err(t.balance().to_string())
    }
}

#[tauri::command]
fn delete_transaction(id: i64) -> Result<(), String> {
    match delete(id) {
        Ok(_) => Ok(()),
        Err(e) => Err(e.to_string()),
    }
}

#[tauri::command]
fn get_latest_transactions(limit: i64, offset: i64) -> Result<HashMap<i64,Transaction>, String> {
    let result = get_recent_transactions(limit, Some(offset));
    match result {
        Ok(t_list) => Ok(t_list),
        Err(e) => Err(e.to_string()),
    }
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
            get_transactions,
            get_transaction_details,
            add_transaction,
            update_transaction,
            get_latest_transactions,
            get_row_count,
            delete_transaction,
            validate
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
    months[month]
}
