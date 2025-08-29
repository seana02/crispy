use std::str::FromStr;
use rust_decimal::Decimal;
use serde::{Deserialize, Serialize};

#[derive(Debug, PartialEq, Eq, Serialize, Deserialize)]
pub struct AccountNode {
    label: String,
    total_currency: Vec<Wallet>,
    sub_accounts: Vec<AccountNode>,
}

#[derive(Debug, PartialEq, Eq, Serialize, Deserialize)]
pub struct Wallet {
    ccy: String,
    total_value: Decimal,
}

impl AccountNode {
    fn new(label: &str) -> Self {
        AccountNode {
            label: label.to_string(),
            total_currency: Vec::new(),
            sub_accounts: Vec::new(),
        }
    }

    fn insert(&mut self, path: &[&str], value: Decimal, ccy: &str) {
        let currency = self.total_currency.iter_mut().find(|c| c.ccy == ccy);
        match currency {
            Some(cur) => cur.total_value += value,
            None => {
                self.total_currency.push(
                    Wallet {
                        ccy: ccy.to_string(),
                        total_value: value
                    }
                );
                self.total_currency.sort_by(|a,b| {
                    get_ordering(&a.ccy).cmp(&get_ordering(&b.ccy))
                });
            }
        }
        if path.is_empty() {
            return;
        }

        let first = path[0];
        let child = self.sub_accounts.iter_mut().find(|c| c.label == first);

        match child {
            Some(child_node) => {
                child_node.insert(&path[1..], value, ccy);
            }
            None => {
                let mut new_node = AccountNode::new(first);
                new_node.insert(&path[1..], value, ccy);
                self.sub_accounts.push(new_node);
            }
        }
    }
}

fn get_ordering(ccy: &str) -> i32 {
    match ccy {
        "USD" => 1,
        "JPY" => 2,
        "CAD" => 3,
        "GBP" => 4,
        "EUR" => 5,
        "CNY" => 6,
        "AUD" => 7,
        "MXN" => 8,
        "KRW" => 9,
        "HKD" => 10,
        "IND" => 11,
        "BTC" => 1001,
        "ETH" => 1002,
        _ => 999,
    }
}

pub fn build_tree(data: Vec<(&str, Vec<(&str, Decimal)>)>) -> AccountNode {
    let mut root_node = AccountNode::new("root");

    for (path, wallet) in data {
        let parts: Vec<&str> = path.split(':').collect();

        if parts.is_empty() {
            continue;
        }

        for (currency, value) in wallet {
            root_node.insert(&parts, value, currency);
        }
    }

    // println!("{:#?}", root_node);

    root_node
}

#[cfg(test)]
mod tests {
    use super::*;

    fn dec(val: &str) -> Decimal {
        Decimal::from_str(val).expect("Invalid decimal")
    }

    #[test]
    fn test_empty() {
        let correct = AccountNode::new("root");
        let to_test = build_tree(vec!());
        assert_eq!(to_test, correct);
    }
    
    #[test]
    fn test_single_account() {
        let correct = AccountNode {
            label: "root".to_string(),
            total_currency: dec("2.50"),
            sub_accounts: vec!(AccountNode {
                label: "Expenses".to_string(),
                total_currency: dec("2.50"),
                sub_accounts: vec!(),
            })
        };
        assert_eq!(correct, build_tree(vec!(
            ("Expenses", dec("2.50"))
        )));
    }

    #[test]
    fn test_single_path() {
        let correct = AccountNode {
            label: "root".to_string(),
            total_currency: dec("2.50"),
            sub_accounts: vec!(AccountNode {
                label: "Expenses".to_string(),
                total_currency: dec("2.50"),
                sub_accounts: vec!(AccountNode {
                    label: "Shopping".to_string(),
                    total_currency: dec("2.50"),
                    sub_accounts: vec!(AccountNode {
                        label: "Physical".to_string(),
                        total_currency: dec("2.50"),
                        sub_accounts: vec!()
                    })
                }),
            })
        };
        assert_eq!(correct, build_tree(vec!(
            ("Expenses:Shopping:Physical", dec("2.50"))
        )));
    }

    #[test]
    fn test_multiple_roots() {
        let correct = AccountNode {
            label: "root".to_string(),
            total_currency: Decimal::ZERO,
            sub_accounts: vec!(
                AccountNode {
                    label: "Expenses".to_string(),
                    total_currency: dec("2.50"),
                    sub_accounts: vec!(AccountNode {
                        label: "Shopping".to_string(),
                        total_currency: dec("2.50"),
                        sub_accounts: vec!(AccountNode {
                            label: "Physical".to_string(),
                            total_currency: dec("2.50"),
                            sub_accounts: vec!()
                        })
                    }),
                },
                AccountNode {
                    label: "Liabilities".to_string(),
                    total_currency: dec("-2.50"),
                    sub_accounts: vec!(AccountNode {
                        label: "CreditCard".to_string(),
                        total_currency: dec("-2.50"),
                        sub_accounts: vec!()
                    })
                }
            )
        };

        assert_eq!(correct, build_tree(vec!(
            ("Expenses:Shopping:Physical", dec("2.50")),
            ("Liabilities:CreditCard", dec("-2.50"))
        )));
    }

    #[test]
    fn test_overlap() {
        let correct = AccountNode {
            label: "root".to_string(),
            total_currency: dec("99.99"),
            sub_accounts: vec!(
                AccountNode {
                    label: "Expenses".to_string(),
                    total_currency: dec("6.53"),
                    sub_accounts: vec!(
                        AccountNode {
                            label: "Shopping".to_string(),
                            total_currency: dec("5.51"),
                            sub_accounts: vec!(
                                AccountNode {
                                    label: "Physical".to_string(),
                                    total_currency: dec("2.50"),
                                    sub_accounts: vec!()
                                },
                                AccountNode {
                                    label: "Online".to_string(),
                                    total_currency: dec("3.01"),
                                    sub_accounts: vec!(),
                                }
                            )
                        },
                        AccountNode {
                            label: "Entertainment".to_string(),
                            total_currency: dec("1.02"),
                            sub_accounts: vec!(),
                        }
                    ),
                },
                AccountNode {
                    label: "Liabilities".to_string(),
                    total_currency: dec("-6.53"),
                    sub_accounts: vec!(AccountNode {
                        label: "CreditCard".to_string(),
                        total_currency: dec("-6.53"),
                        sub_accounts: vec!(
                            AccountNode {
                                label: "Visa".to_string(),
                                total_currency: dec("-3.52"),
                                sub_accounts: vec!(),
                            },
                            AccountNode {
                                label: "MasterCard".to_string(),
                                total_currency: dec("-3.01"),
                                sub_accounts: vec!()
                            }
                        )
                    })
                },
                AccountNode {
                    label: "Offsetter".to_string(),
                    total_currency: dec("99.99"),
                    sub_accounts: vec!(),
                }
            )
        };

        assert_eq!(correct, build_tree(vec!(
            ("Expenses:Shopping:Physical", dec("2.50")),
            ("Expenses:Shopping:Online", dec("3.01")),
            ("Expenses:Entertainment", dec("1.02")),
            ("Liabilities:CreditCard:Visa", dec("-3.52")),
            ("Liabilities:CreditCard:MasterCard", dec("-3.01")),
            ("Offsetter", dec("99.99")),
        )));
    }
}
