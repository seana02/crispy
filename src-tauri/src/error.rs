use core::fmt;

#[derive(Debug)]
pub enum TransactionError {
    Rusqlite(rusqlite::Error),
    UnbalancedPostingError,
}

impl fmt::Display for TransactionError {
    fn fmt(&self, f: &mut fmt::Formatter) -> fmt::Result {
        write!(f, "transaction is unbalanced")
    }
}

impl From<rusqlite::Error> for TransactionError {
    fn from(error: rusqlite::Error) -> Self {
        TransactionError::Rusqlite(error)
    }
}
